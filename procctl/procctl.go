package procctl

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strconv"

	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func RunAsUser(cmd_str string) {
	if pid := findpid("explorer.exe"); pid != 0 {
		handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, pid)
		if err != nil {
			log.Fatalf("error OpenProcess: %v", err)
		}

		var token windows.Token
		err = windows.OpenProcessToken(handle, windows.TOKEN_ALL_ACCESS, &token)
		if err != nil {
			log.Fatalf("error OpenProcessToken: %v", err)
		}

		var outProcInfo windows.ProcessInformation
		var startupInfo windows.StartupInfo
		appName, err := windows.UTF16PtrFromString(cmd_str)
		if err != nil {
			log.Fatalf("error UTF16PtrFromString: %v", err)
		}
		err = windows.CreateProcessAsUser(token, appName, nil, nil, nil, true, windows.NORMAL_PRIORITY_CLASS, nil, nil, &startupInfo, &outProcInfo)
		if err != nil {
			log.Fatalf("error CreateProcessAsUser: %v", err)
		}
	}
}

// start a subprocess and collect it's stdout and stderr real-timely
func Run(cmd []string) (*exec.Cmd, chan string) {
	p := exec.Command(cmd[0], cmd[1:]...)
	stderr, _ := p.StderrPipe()
	stdout, _ := p.StdoutPipe()

	err := p.Start()

	if err == nil {
		ch := make(chan string, 10)
		closed := false
		reader := func(r io.ReadCloser) {
			defer func() {
				if p := recover(); p != nil {
					// we DO have dealt with double-close(see below), but there might be other panics...
					fmt.Printf("Panic: %v\n", p)
				}
			}()
			buf := make([]byte, 256)
			nread, err := r.Read(buf)
			for err == nil || err == io.EOF {
				if nread > 0 {
					ch <- string(buf[:nread])
				}
				if err == io.EOF {
					break
				}
				nread, err = r.Read(buf)
			}

			// double-close cause panic
			if !closed {
				closed = true
				close(ch)
			}
		}

		go reader(stdout)
		go reader(stderr)
		return p, ch
	}

	return nil, nil
}

func ProcessExists(img_name string) bool {
	return findpid(img_name) != 0
}

// find the pid of process with given image name
// currently, the pid of the first process is returned, if multiple processes exist.
// return 0 represents failure (althrough System Idle Process has pid 0)
func findpid(img_name string) uint32 {
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
