package main

import (
	"fmt"
	"os"
	"path/filepath"

	procctl "wisersoft.com.cn/wsrefresher/test"
)

func main() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	file := filepath.Base(os.Args[0])
	r := filepath.Join(dir, os.Args[0])
	fmt.Println(dir)
	fmt.Println(file)
	fmt.Println(r)

	procctl.RunAsUser(`C:\Program Files\Notepad++\notepad++.exe`)
}
