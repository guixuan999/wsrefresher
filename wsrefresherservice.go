package main

import (
	_LOG "log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kardianos/service"
	"wisersoft.com.cn/wsrefresher/procctl"
	"wisersoft.com.cn/wsrefresher/utils"
)

var base = filepath.Base(os.Args[0])
var dir = filepath.Dir(os.Args[0])
var fn_no_ext = strings.TrimSuffix(base, filepath.Ext(base))
var fn_log = filepath.Join(dir, fn_no_ext+".log")
var log = utils.GetLogger(fn_log)

func main() {
	svcConfig := &service.Config{
		Name:        "WSRefresher",                 // service name
		DisplayName: "Wisersoft refresher service", // service name for display
		Description: "Wisersoft refresher service", // service description
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		_LOG.Fatal(err)
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			s.Install()
			_LOG.Println(
				`Service installed.
-----------------------------------------
start service:  sc start wsrefresher
stop service:   sc stop wsrefresher
remove service: sc delete wsrefresher`)
			return
		}

		if os.Args[1] == "remove" {
			s.Uninstall()
			_LOG.Println("Service removed")
			return
		}
	}

	// run in non-service mode provided no command parameter
	if err = s.Run(); err != nil {
		_LOG.Println(err)
	}
}

type program struct{}

func (p *program) Start(s service.Service) error {
	log.Info("service started!\n")
	go p.run()
	return nil
}

func (p *program) run() {
	for {
		if !procctl.ProcessExists("wsrefresher.exe") {
			procctl.RunAsUser(filepath.Join(dir, "wsrefresher", "wsrefresher.exe"))
		}
		time.Sleep(time.Second * 5)
	}
}

func (p *program) Stop(s service.Service) error {
	cmd, _ := procctl.Run(strings.Split("taskkill /F /IM wsrefresher.exe", " "))
	if cmd != nil && cmd.Process != nil {
		cmd.Wait()
		log.Info("kill wsrefresher.exe process(es) ok.\n")
	} else {
		log.Info("kill wsrefresher.exe process(es) failed, can't start process taskkill!\n")
	}
	log.Info("service stopped!\n")
	return nil
}
