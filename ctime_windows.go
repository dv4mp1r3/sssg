package main

import (
	"io/fs"
	"syscall"
	"time"
)

func getCtime(fInfo fs.FileInfo) time.Time {

	d := fInfo.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, getCtimeNSec(d))
}

func getCtimeNSec(stat *syscall.Win32FileAttributeData) int64 {
	return stat.CreationTime.Nanoseconds()
}
