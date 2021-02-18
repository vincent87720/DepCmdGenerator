package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type Settings struct {
	StartHostInfo string
	StartIPAddr   string
	Mask          string
	Gateway       string
	DNSServer     []string
	Qty           int
	Interface     string
}

func LoadFromFile() (s Settings, err error) {
	content, err := ioutil.ReadFile("settings.json")
	if err != nil {
		err = errors.WithStack(fmt.Errorf("settings.json file not found"))
		return
	}

	err = json.Unmarshal(content, &s)
	if err != nil {
		err = errors.WithStack(fmt.Errorf("Unable to unmarshal the settings.json file"))
		return
	}

	return s, nil
}

func CreateFile(path string) *os.File {
	fo, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}
	return fo
}

func WriteCmd(w *bufio.Writer, command string) {

	utf8ToBig5 := traditionalchinese.Big5.NewEncoder()
	big5, _, _ := transform.String(utf8ToBig5, command)
	_, err := w.WriteString(big5)
	if err != nil {
		fmt.Println(err)
	}
}
