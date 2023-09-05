//go:build windows
// +build windows

package paniclog

import (
	"log"
	"os"
	"syscall"
)

// redirectStderr to the file passed in
func RedirectStderr(f *os.File) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setStdHandle := kernel32.NewProc("SetStdHandle")
	sh := syscall.STD_ERROR_HANDLE
	v, _, err := setStdHandle.Call(uintptr(sh), uintptr(f.Fd()))
	if v == 0 {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}
	os.Stderr = f
}
