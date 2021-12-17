package tool

import (
	"fmt"
	"os"
	"os/exec"
)

func TestLossless(srcPlyPath, recPlyPath string) bool {
	out, err := exec.Command("cmp", srcPlyPath, recPlyPath).Output()
	if err != nil {
		fmt.Println("ERROR", err)
		fmt.Printf("%s\n", out)
		return false
	} else {
		fmt.Println("OK")
		return true
	}
}

func DeleteTmpFile(isLossless bool, srcPlyPath, recPath, etcPath string) {
	if !isLossless {
		return
	}
	if err := os.Remove(srcPlyPath); err != nil {
		panic(err)
	}
	if err := os.Remove(recPath); err != nil {
		panic(err)
	}
	if err := os.Remove(etcPath); err != nil {
		panic(err)
	}
}
