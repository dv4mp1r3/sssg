package main

import (
	"io/fs"
	"syscall"
	"time"
)

func getCtime(fInfo fs.FileInfo) time.Time {
	stat := fInfo.Sys().(*syscall.Stat_t)
	return time.Unix(int64(getCtimeSec(stat)), int64(getCtimeNSec(stat)))
}

func getCtimeSec(stat *syscall.Stat_t) int64 {
	return stat.Ctim.Sec
}

func getCtimeNSec(stat *syscall.Stat_t) int64 {
	return stat.Ctim.Nsec
}
