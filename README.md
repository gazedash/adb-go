# adb-go

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

```
go run main.go -mode pull
```

TODO:

- [ ] Ability to pass .pullignore from cmd OR

- [ ] Implement per-device .pullignore based on device id

- [ ] Implement push

- [ ] Think about automatic sync when device is plugged in

PROFIT
