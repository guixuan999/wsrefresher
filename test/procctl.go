package procctl

import (
	"bytes"
	"log"
	"os/exec"
	"regexp"
	"strconv"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// find the pid of process with given image name
// currently, the pid of the first process is returned, if multiple processes exist.
// return 0 represents failure (althrough System Idle Process has pid 0)
func Findpid(img_name string) uint32 {
	if out, err := exec.Command("tasklist").Output(); err != nil {
		log.Fatalf("error executing tasklist: %v", err)
	} else {
		gbk_decoder := simplifiedchinese.GBK.NewDecoder()
		decoded, _ := gbk_decoder.Bytes(out)

		lines := bytes.Split([]byte(decoded), []byte{0xa})
		for i, line := range lines {
			if string(line) == "" || i == 0 {
				continue
			}
			re, _ := regexp.Compile(`\s+`)
			fileds := re.Split(string(line), -1)
			name := fileds[0]
			pid, _ := strconv.Atoi(fileds[1])
			if name == img_name {
				return uint32(pid)
			}
		}
	}
	return 0
}
