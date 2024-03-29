# adb-go

![alt text](https://github.com/gazedash/adb-go/blob/master/Screenshot%202023-10-03%20at%2022-06-47%20adbgo.png)

NOTE: if you want to pull everything, and I mean EVERYTHING from your phone, just do `adb pull sdcard/`, otherwise
read on:

All of this can be done with a few lines of bash, however:

I want to be able to plug any device and with one simple command pull all data from it,
exluding some folders

And I find any programming language more readable and maintainable than Bash/Zsh etc

Why not just manually copy through MTP?

* Need to enable it every time the device is plugged in
* It's platform dependent, and device vendor dependent
* It hangs on big files or if there are too many files
* It's manual, because...
* It's hard to automate, even though libs for MTP exist

So my target workflow is:

Plug in the device, run this tool, PROFIT

Plug any other device, run this tool, PROFIT

Pulled files are grouped by device id

### Running (Windows)

Download latest release, unzip
Make sure to tweak config and .pullignore for your needs!

Click run_me OR
Open terminal, navigate to the folder, and

```
adbgo.exe -mode pull
```

### Running (other OS)

```
go run main.go -mode pull
```

### Building from source

If a binary is more convenient for you, you can build it yourself:

```
# Build Linux x64
GOOS=linux GOARCH=amd64 go build -o adbgo_linux_x64
# Build MacOS ARM
GOOS=darwin GOARCH=arm64 go build -o adbgo_darwin_arm64
```

Note: this is not tested.

TODO:

- [ ] One button to rule them all?

- [ ] Ability to pass .pullignore from cmd OR

- [ ] Implement per-device .pullignore based on device id

- [X] Implement push

- [ ] Think about automatic sync when device is plugged in (autorun, service)

- [ ] Make it cross platform - support unix-style paths too, etc (implemented, not tested)

- [ ] Ignore already pulled files? Does adb check if files exist and doesn't copy same files or not? (apparently adb skips already pulled/pushed files)

- [X] GUI - basic GUI in HTML done

- [ ] Pass logs from server to client

- [ ] Show connected device id

PROFIT

Build w/o terminal window - `go build -ldflags -H=windowsgui`