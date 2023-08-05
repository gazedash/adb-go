package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

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

func pull() {
	// Dirs to be ignored when pulling
	pullignoreData, _ := os.ReadFile(".pullignore")
	pullignore := strings.Split(string(NormalizeNewlines(pullignoreData)), newline)
	configData, _ := os.ReadFile("config")
	config := strings.Split(string(NormalizeNewlines(configData)), newline)
	dst := ""

	// Get seial number, as reported by adb devices
	cmd := exec.Command("adb", "shell", "getprop", "ro.serialno")
	device_id_data, _ := cmd.Output()
	device_id := strings.ReplaceAll(string(NormalizeNewlines(device_id_data)), newline, "")

	for _, v := range config {
		option := strings.Split(v, "=")
		key := option[0]
		value := option[1]

		// handle config options
		if key == "DST" {
			dst = value + string(os.PathSeparator) + device_id
		}
	}

	ls := exec.Command("adb", "shell", "ls", "-d", "/sdcard/*")

	stdout, _ := ls.Output()

	folders := strings.Split(string(NormalizeNewlines(stdout)), newline)

	// exclude ignored folders
	foldersToPull := make([]string, len(folders))
	for _, v := range folders {
		if slices.Contains(pullignore, strings.Trim(v, " ")) == false {
			foldersToPull = append(foldersToPull, strings.Trim(v, " "))
		}
	}

	fmt.Println("Pulling these dirs: " + strings.Trim(strings.Join(foldersToPull, " "), " "))
	fmt.Println(".....................")
	fmt.Println("to...................")
	fmt.Println(dst)
	fmt.Println(".....................")

	// easy way to guesstimate the env we're running on, cuz
	// it can be: Windows, linux, Git Bash, WSL (run from .exe), WSL (run from "elf" linux binary)
	if strings.Contains(string(stdout), "\r\n") {
		// cmd cuz os.MkdirAll is too stupid (maybe it's only a Windows problem?)
		cmd = exec.Command("cmd", "/c", "mkdir", dst)
		_, _ = cmd.Output()
	} else {
		cmd = exec.Command("mkdir", "-p", dst)
		_, _ = cmd.Output()
	}

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

func main() {
	// flags are easier than args, it just works (c)
	modePtr := flag.String("mode", "", "a string")

	flag.Parse()

	fmt.Println("Running in mode: " + *modePtr)

	if *modePtr == "pull" {
		pull()
	} else {
		fmt.Println("no such mode, or no mode passed. try to run with '-mode pull'")
	}
}
