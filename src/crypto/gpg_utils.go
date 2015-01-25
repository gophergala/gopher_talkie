package crypto

import (
	"bytes"
	"fmt"
	"io"
	_ "io/ioutil"
	_ "os"
	"os/exec"
	"regexp"
	_ "strings"
)

var (
	GPGPath = "/usr/local/bin/gpg"
)

func init() {
	output, err := exec.Command("which gpg").CombinedOutput()
	if err == nil {
		GPGPath = string(output)
	}
}

type Key struct {
	Name      string
	Email     string
	PublicKey string
	SecretKey string
}

func ListSecretKeys() ([]Key, error) {

	var buf bytes.Buffer

	gpg := exec.Command(GPGPath, "--list-secret-keys")

	grep := exec.Command("grep", "uid")
	grep.Stdin, _ = gpg.StdoutPipe()

	sed := exec.Command("sed", "-E", "s/uid +//g")
	sed.Stdin, _ = grep.StdoutPipe()
	sed.Stdout = &buf

	grep.Start()
	sed.Start()

	gpg.Run()
	grep.Wait()
	sed.Wait()

	var list []Key
	re := regexp.MustCompile("(.*)<(.*@.*)>")
	for {
		line, err := buf.ReadString('\n')
		fmt.Println(line)
		if err == io.EOF {
			break
		}

		// parse line
		matches := re.FindAllStringSubmatch(line, -1)
		if len(matches) > 0 && len(matches[0]) > 2 {
			key := Key{
				Name:  matches[0][1],
				Email: matches[0][2],
			}
			// fmt.Printf("key: %v\n", key)
			list = append(list, key)
		}
	}
	return list, nil
}
