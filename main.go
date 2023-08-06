package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

// NormalizeNewlines normalizes \r\n (windows) and \r (mac)
// into \n (unix)
func NormalizeNewlines(d []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)

	return d
}

var newline = "\n"

func SplitByNewLine(bytes []byte) []string {
	return strings.Split(string(NormalizeNewlines(bytes)), newline)
}

func GetDeviceId() string {
	// Get seial number, as reported by adb devices
	cmd := exec.Command("adb", "shell", "getprop", "ro.serialno")
	device_id_data, _ := cmd.Output()
	return strings.ReplaceAll(string(NormalizeNewlines(device_id_data)), newline, "")
}

type Config struct {
	dst string
}

func GetConfig() Config {
	configData, _ := os.ReadFile("config")
	config := SplitByNewLine(configData)

	dst := ""

	for _, v := range config {
		option := strings.Split(v, "=")
		key := option[0]
		value := option[1]

		// handle config options
		if key == "DST" {
			dst = value + string(os.PathSeparator)
		}
	}

	cfg := Config{dst}

	return cfg
}

func PullFiles(foldersToPull []string, dst string) {
	for _, v := range foldersToPull {
		if v == "" {
			continue
		}

		pull := exec.Command("adb", "pull", v, dst)

		// print command with src and dst paths
		fmt.Println(pull)

		stdout, _ := pull.Output()

		// print result
		fmt.Println(string(stdout))
	}
}

func MkdirP(dst string, isWindows bool) {
	// easy way to guesstimate the env we're running on, cuz
	// it can be: Windows, linux, Git Bash, WSL (run from .exe), WSL (run from "elf" linux binary)
	if isWindows {
		// cmd cuz os.MkdirAll is too stupid (maybe it's only a Windows problem?)
		cmd := exec.Command("cmd", "/c", "mkdir", dst)
		_, _ = cmd.Output()
	} else {
		cmd := exec.Command("mkdir", "-p", dst)
		_, _ = cmd.Output()
	}
}

func GetFoldersToPull() []string {
	// Dirs to be ignored when pulling
	pullignoreData, _ := os.ReadFile(".pullignore")
	pullignore := SplitByNewLine(pullignoreData)

	ls := exec.Command("adb", "shell", "ls", "-d", "/sdcard/*")

	stdout, _ := ls.Output()

	folders := SplitByNewLine(stdout)
	// exclude ignored folders
	foldersToPull := make([]string, len(folders))
	for _, v := range folders {
		if slices.Contains(pullignore, strings.Trim(v, " ")) == false {
			foldersToPull = append(foldersToPull, strings.Trim(v, " "))
		}
	}

	return foldersToPull
}

func IsWindows() bool {
	testCommand := exec.Command("adb", "shell", "echo", "test")
	stdout, _ := testCommand.Output()
	isWindows := strings.Contains(string(stdout), "\r\n")

	return isWindows
}

func PreparePull(foldersToPull []string, dst string) {
	fmt.Println("Pulling these dirs: " + strings.Trim(strings.Join(foldersToPull, " "), " "))
	fmt.Println("To...................")
	fmt.Println(dst)
	fmt.Println("Pulling.....................")

	MkdirP(dst, IsWindows())
}

func GetDestination(cfg Config) string {
	// Destination is a path from config + device_id folder
	return cfg.dst + GetDeviceId()
}

func pull(cfg Config) {
	dst := GetDestination(cfg)
	foldersToPull := GetFoldersToPull()

	PreparePull(foldersToPull, dst)
	PullFiles(foldersToPull, dst)
}

func PreparePush() {
	fmt.Println("Pushing.....................")
	fmt.Println("/sdcard/Pushed")
	fmt.Println("Pushing.....................")
}

func PushFiles(cfg Config) {
	dst := "/sdcard/Pushed"

	dirPath := GetDestination(cfg) + string(os.PathSeparator) + "Push"

	push := exec.Command("adb", "push", dirPath, dst)

	// print command
	fmt.Println(push)

	stdout, err := push.Output()

	fmt.Println(string(stdout))

	if err == nil {
		t := time.Now()
		date := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())

		os.Rename(dirPath, dirPath+date)
		MkdirP(dirPath, IsWindows())
	}
}

func push(cfg Config) {
	PreparePush()
	PushFiles(cfg)
}

func main() {
	// flags are easier than args, it just works (c)
	modePtr := flag.String("mode", "", "a string")

	flag.Parse()

	fmt.Println("Running in mode: " + *modePtr)

	cfg := GetConfig()

	if *modePtr == "pull" {
		pull(cfg)
	} else if *modePtr == "push" {
		push(cfg)
	} else {
		fmt.Println("no such mode, or no mode passed. try to run with '-mode pull'")
	}
}
