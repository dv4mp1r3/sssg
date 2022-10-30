package main

import (
	"syscall"
)

func getCtimeSec(stat *syscall.Stat_t) int64 {
	return stat.Ctimespec.Sec
}

func getCtimeNSec(stat *syscall.Stat_t) int64 {
	return stat.Ctimespec.Nsec
}
