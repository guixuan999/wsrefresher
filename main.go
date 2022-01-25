package main

import (
	"fmt"

	procctl "wisersoft.com.cn/wsrefresher/test"
)

func main() {
	fmt.Printf("pid: %d", procctl.Findpid("explorer.exe"))
}
