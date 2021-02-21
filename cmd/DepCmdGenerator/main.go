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
