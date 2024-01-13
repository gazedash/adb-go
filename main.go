package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/websocket"
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
	configData, err := os.ReadFile("config")

	if err != nil {
		file, _ := os.Create("config")
		winPath := "DST=C:\\AndroidBackup"
		linuxPath := "~"

		path := winPath

		if runtime.GOOS == "windows" {
			path = winPath
		} else {
			path = linuxPath
		}

		file.Write([]byte(path))

		configData, _ = os.ReadFile("config")
	}

	configData, _ = os.ReadFile("config")

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
		print(pull.String())

		stdout, _ := pull.Output()

		// print result
		print(string(stdout))
	}

	print("pull finished")
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

func print(str string) {
	fmt.Println(str)

	if msgChan != nil {
		msgChan <- str
	}
}

type Target struct {
	separator string
	path      string
}

var target = Target{
	separator: "/",
	path:      "sdcard",
}

func (t *Target) getPath() string {
	return t.separator + t.path
}

func (t *Target) getDir(dir string) string {
	return t.getPath() + t.separator + dir
}

func (t *Target) allFilesGlob() string {
	return t.getDir("*")
}

func GetFoldersToPull() []string {
	// Dirs to be ignored when pulling
	pullignoreData, err := os.ReadFile(".pullignore")

	if (err != nil) {
		file, _ := os.Create(".pullignore")
		file.Write([]byte("/sdcard/Aboba\nArt"))
	}

	pullignoreData, _ = os.ReadFile(".pullignore")

	pullignore := strings.ReplaceAll(string(NormalizeNewlines(pullignoreData)), newline, " ")

	ls := exec.Command("adb", "shell", "ls", "-d", target.allFilesGlob())

	stdout, _ := ls.Output()

	folders := SplitByNewLine(stdout)
	// exclude ignored folders
	foldersToPull := make([]string, len(folders))
	for _, v := range folders {
		folder := strings.Trim(v, " ")
		if strings.Contains(pullignore, folder) == false {
			foldersToPull = append(foldersToPull, folder)
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
	print("Pulling these dirs: " + strings.Trim(strings.Join(foldersToPull, " "), " "))
	print("To...................")
	print(dst)
	print("Pulling.....................")

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

var pushDirName = "Push"

func PreparePush() {
	print("Pushing.....................")
	print(target.getDir(pushDirName))
	print("Pushing.....................")
}

func IsEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return true
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false
}

func PushFiles(cfg Config) {
	dirPath := GetDestination(cfg) + string(os.PathSeparator) + pushDirName

	MkdirP(dirPath, IsWindows())

	push := exec.Command("adb", "push", dirPath, target.getPath())

	// print command
	print(push.String())

	stdout, err := push.Output()

	if err != nil {
		print(err.Error())
	}
	print(string(stdout))

	if IsEmpty(dirPath) {
		print("push dir empty, fin")
		return
	}

	if err == nil {
		t := time.Now()
		date := fmt.Sprintf("%d-%02d-%02dT%02d_%02d_%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())

		// fmt.Println(date)
		// fmt.Println(dirPath)
		// fmt.Println(dirPath + date)

		err := os.Rename(dirPath, dirPath+"_"+date)

		print(err.Error())

		MkdirP(dirPath, IsWindows())
	}

	print("push finished")
}

func PushPathFiles(cfg Config, path string) {
	push := exec.Command("adb", "push", path, target.getPath())

	// print command
	print(push.String())

	stdout, err := push.Output()

	if err != nil {
		print(err.Error())
	}

	print(string(stdout))

	print("push path finished")
}

func push(cfg Config) {
	PreparePush()
	PushFiles(cfg)
}

func PreparePushPath(path string) {
	print("@ Pushing path.....................")
	print(path)
	print("@ Pushing path.....................")
}

func pushPath(cfg Config, path string) {
	PreparePushPath(path)
	PushPathFiles(cfg, path)
}

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		panic(err)
	}
}

var msgChan chan string = make(chan string)

func wsHandler(mux *http.ServeMux) {
	mux.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		for {
			msg := <-msgChan
			websocket.JSON.Send(ws, msg)
		}
	}))
}

func serveHtml(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := os.ReadFile("./index.html")
		w.Write(bytes)
	})
}

func pullHandler(mux *http.ServeMux, cfg Config) {
	mux.HandleFunc("/pull", func(w http.ResponseWriter, r *http.Request) {
		pull(cfg)
		msgChan <- "FinishedEvent"
	})
}

func pushHandler(mux *http.ServeMux, cfg Config) {
	mux.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		push(cfg)
		msgChan <- "FinishedEvent"
	})
}

func pushPathHandler(mux *http.ServeMux, cfg Config) {
	mux.HandleFunc("/pushPath", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		print(string(body))
		pushPath(cfg, string(body))
		msgChan <- "FinishedEvent"
	})
}

func syncHandler(mux *http.ServeMux, cfg Config) {
	mux.HandleFunc("/sync", func(w http.ResponseWriter, r *http.Request) {
		doSync(cfg)
		msgChan <- "FinishedEvent"
		print("server sync end")
	})
}

func server(cfg Config) {
	port := "5151"
	addr := "http://localhost:" + port

	OpenBrowser(addr)

	mux := http.NewServeMux()

	serveHtml(mux)
	wsHandler(mux)
	syncHandler(mux, cfg)
	pullHandler(mux, cfg)
	pushHandler(mux, cfg)
	pushPathHandler(mux, cfg)

	s := http.Server{Addr: ":" + port, Handler: mux}

	s.ListenAndServe()
}

var serverMode = "server"

func doSync(cfg Config) {
	push(cfg)
	pull(cfg)
	print("doSync end")
}

func main() {
	cfg := GetConfig()

	// flags are easier than args, it just works (c)
	modePtr := flag.String("mode", "", "a string")

	if *modePtr == "" {
		modePtr = &serverMode
	}

	flag.Parse()

	fmt.Println("Running in mode: " + *modePtr)

	if *modePtr == "pull" {
		pull(cfg)
	} else if *modePtr == "push" {
		push(cfg)
	} else if *modePtr == "sync" {
		doSync(cfg)

		print("cmd sync end")
	} else if *modePtr == "server" {
		server(cfg)

		fmt.Println("starting in server mode by default...")
	} else {
		fmt.Println("no such mode, or no mode passed. try to run with '-mode pull'")
		fmt.Println("available modes: pull, push, sync, server")
	}
}
