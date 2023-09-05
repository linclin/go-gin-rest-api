//go:build unix
// +build unix

package paniclog

import (
	"log"
	"os"
	"syscall"
)

// redirectStderr to the file passed in
func RedirectStderr(f *os.File) {
	err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}
}
