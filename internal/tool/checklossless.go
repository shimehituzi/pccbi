package tool

import (
	"fmt"
	"os"
	"os/exec"
)

func TestLossless(srcPath, recPath, srcPlyPath, recPlyPath string) bool {
	out1, err := exec.Command("cmp", srcPath, recPath).Output()
	if err != nil {
		fmt.Println("ERROR", err)
		fmt.Printf("%s\n", out1)
		return false
	}

	out2, err := exec.Command("cmp", srcPlyPath, recPlyPath).Output()
	if err != nil {
		fmt.Println("ERROR", err)
		fmt.Printf("%s\n", out2)
		return false
	}

	fmt.Println("OK")
	return true
}

func DeleteTmpFile(isLossless bool, srcPath, recPath, srcPlyPath, recPlyPath, etcPath string) {
	if !isLossless {
		return
	}
	if err := os.Remove(srcPath); err != nil {
		panic(err)
	}
	if err := os.Remove(recPath); err != nil {
		panic(err)
	}
	if err := os.Remove(srcPlyPath); err != nil {
		panic(err)
	}
	if err := os.Remove(recPlyPath); err != nil {
		panic(err)
	}
	if err := os.Remove(etcPath); err != nil {
		panic(err)
	}
}
