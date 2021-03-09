package main

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	"github.com/vincent87720/DepCmdGenerator/internal/file"
)

func main() {
	settings, err := file.LoadFromFile()
	if err != nil {
		fmt.Println(err)
	}

	hostFront, hostCounter, err := parse(settings.StartHostInfo)
	if err != nil {
		fmt.Println(err)
	}

	ipFront, ipCounter, err := parse(settings.StartIPAddr)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < settings.Qty; i++ {
		var s string
		hostBuf := bytes.NewBufferString(s)
		fmt.Fprintf(hostBuf, "%s%03d", hostFront, hostCounter)

		f := file.CreateFile(hostBuf.String() + ".bat")
		defer f.Close()

		w := bufio.NewWriter(f)

		batchGotAdminCmd(w)

		file.WriteCmd(w, "netsh interface ip set address name=\""+settings.Interface+"\" source=static address="+ipFront+strconv.Itoa(ipCounter)+" mask="+settings.Mask+" gateway="+settings.Gateway+" gwmetric=1\r\n")
		for idx, val := range settings.DNSServer {
			if idx == 0 {
				file.WriteCmd(w, "netsh interface ip set dnsserver name=\""+settings.Interface+"\" source=static address="+val+"\r\n")
			} else {
				file.WriteCmd(w, "netsh interface ip add dnsserver name=\""+settings.Interface+"\" address="+val+"\r\n")
			}
		}

		file.WriteCmd(w, "REG ADD HKLM\\SYSTEM\\ControlSet001\\Services\\Tcpip\\Parameters /v \"NV Hostname\" /t REG_SZ /d "+hostBuf.String()+" /f")

		w.Flush()

		hostCounter++
		ipCounter++
	}
}

func parse(s string) (front string, startID int, err error) {
	expFront, err := regexp.Compile(`.*[.-]`)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	expEnd, err := regexp.Compile(`[0-9]{1,3}$`)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	front = expFront.FindString(s)
	id := expEnd.FindString(s)

	int64id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	startID = int(int64id)

	return front, startID, nil
}

func batchGotAdminCmd(w *bufio.Writer) {
	file.WriteCmd(w, "@echo off\r\n")
	file.WriteCmd(w, "\r\n")
	file.WriteCmd(w, ":: BatchGotAdmin\r\n")
	file.WriteCmd(w, ":-------------------------------------\r\n")
	file.WriteCmd(w, "REM  --> Check for permissions\r\n")
	file.WriteCmd(w, ">nul 2>&1 \"%SYSTEMROOT%\\system32\\cacls.exe\" \"%SYSTEMROOT%\\system32\\config\\system\"\r\n")
	file.WriteCmd(w, "\r\n")
	file.WriteCmd(w, "REM --> If error flag set, we do not have admin.\r\n")
	file.WriteCmd(w, "if '%errorlevel%' NEQ '0' (\r\n")
	file.WriteCmd(w, "    echo Requesting administrative privileges...\r\n")
	file.WriteCmd(w, "    goto UACPrompt\r\n")
	file.WriteCmd(w, ") else ( goto gotAdmin )\r\n")
	file.WriteCmd(w, "\r\n")
	file.WriteCmd(w, ":UACPrompt\r\n")
	file.WriteCmd(w, "    echo Set UAC = CreateObject^(\"Shell.Application\"^) > \"%temp%\\getadmin.vbs\"\r\n")
	file.WriteCmd(w, "    set params = %*:\"=\"\"\r\n")
	file.WriteCmd(w, "    echo UAC.ShellExecute \"cmd.exe\", \"/c %~s0 %params%\", \"\", \"runas\", 1 >> \"%temp%\\getadmin.vbs\"\r\n")
	file.WriteCmd(w, "\r\n")
	file.WriteCmd(w, "    \"%temp%\\getadmin.vbs\"\r\n")
	file.WriteCmd(w, "    del \"%temp%\\getadmin.vbs\"\r\n")
	file.WriteCmd(w, "    exit /B\r\n")
	file.WriteCmd(w, "\r\n")
	file.WriteCmd(w, ":gotAdmin\r\n")
	file.WriteCmd(w, "    pushd \"%CD%\"\r\n")
	file.WriteCmd(w, "    CD /D \"%~dp0\"\r\n")
	file.WriteCmd(w, ":--------------------------------------\r\n")
	file.WriteCmd(w, "@echo on\r\n")
}
