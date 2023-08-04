package main

import (
	"fmt"
	"os"
	"strings"
	"os/exec"
	"golang.org/x/exp/slices"
)

func main() {
pullignoreData, _ := os.ReadFile(".pullignore")
pullignore := strings.Split(string(pullignoreData), "\r\n")

ls := exec.Command("adb", "shell", "ls", "-d", "/sdcard/*")

stdout, _ := ls.Output()

folders := strings.Split(string(stdout), "\r\n")

foldersToPull := ""
for _, v := range folders {
	if (slices.Contains(pullignore, strings.Trim(v, " ")) == false) {
		foldersToPull += v + " "
	}
}

fmt.Println(foldersToPull)

// cmd = exec.Command(app, arg0, arg1)
// stdout, err := cmd.Output()

// if err != nil {
// 	fmt.Println(err.Error())
// 	return
// }

// // Print the output
// fmt.Println(string(stdout))
// utput
// fmt.Println(string(stdout))
}
