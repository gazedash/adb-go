package main

import (
	"fmt"
	"os"
	"strings"

	// "os"
	"os/exec"

	"golang.org/x/exp/slices"
	// "os/exec"
	// "strings"
)
	// app := "adb// app := "adb"
// arg0 := "pull"
// arg1 := "sdcard/"

func main() {
pullignoreData, _ := os.ReadFile(".pullignore")
pullignore := strings.Split(string(pullignoreData), "\n")

ls := exec.Command("adb", "shell", "ls", "-d", "/sdcard/*")

// cmd = exec.Command(app, arg0, arg1)
stdout, _ := ls.Output()

folders := strings.Split(string(stdout), "\n")

// fmt.Println(string(stdout))

// Print the output

foldersToPull := ""
for _, v := range folders {
	if (slices.Contains(pullignore, v) == false) {
		foldersToPull += " " + v + " "
	}
}

// fmt.Println(folders)

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
