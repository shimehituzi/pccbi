package tool

import (
	"fmt"
	"os"
	"os/exec"
)

func TestLossless(sortedPath, recPath string) bool {
	out, err := exec.Command("cmp", sortedPath, recPath).Output()
	if err != nil {
		fmt.Println("ERROR", err)
		fmt.Printf("%s\n", out)
		return false
	} else {
		fmt.Println("OK")
		return true
	}
}

func DeleteTmpFile(isLossless bool, sortedPath, dstPath, etcPath string) {
	if !isLossless {
		return
	}
	if err := os.Remove(sortedPath); err != nil {
		panic(err)
	}
	if err := os.Remove(dstPath); err != nil {
		panic(err)
	}
	if err := os.Remove(etcPath); err != nil {
		panic(err)
	}
}
