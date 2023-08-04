package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/exp/slices"
)

func pull() {
	// Dirs to be ignored when pulling
	pullignoreData, _ := os.ReadFile(".pullignore")
	pullignore := strings.Split(string(pullignoreData), "\r\n")
	configData, _ := os.ReadFile("config")
	config := strings.Split(string(configData), "\r\n")
	dst := ""

	// Get seial number, as reported by adb devices
	cmd := exec.Command("adb", "shell", "getprop", "ro.serialno")
	device_id, _ := cmd.Output()

	for _, v := range config {
		option := strings.Split(v, "=")
		key := option[0]
		value := option[1]

		// handle config options
		if key == "DST" {
			dst = value + "\\" + string(device_id)
		}
	}

	ls := exec.Command("adb", "shell", "ls", "-d", "/sdcard/*")

	stdout, _ := ls.Output()

	folders := strings.Split(string(stdout), "\r\n")

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

	// cmd cuz os.MkdirAll is too stupid (maybe it's only a Windows problem)
	cmd = exec.Command("cmd", "/c", "mkdir", dst)

	_, _ = cmd.Output()

	for _, v := range foldersToPull {

		if v == "" {
			continue
		}

		pull := exec.Command("adb", "pull", v, strings.ReplaceAll(dst, "\r\n", ""))

		// print command with src and dst paths
		fmt.Println(pull)

		stdout, _ := pull.Output()

		// print result
		fmt.Println(string(stdout))
	}
}

func main() {
	modePtr := flag.String("mode", "", "a string")

	flag.Parse()

	fmt.Println("Running in mode: " + *modePtr)

	if *modePtr == "pull" {
		pull()
	} else {
		fmt.Println("no such mode, or no mode passed. try to run with '-mode pull'")
	}
}
