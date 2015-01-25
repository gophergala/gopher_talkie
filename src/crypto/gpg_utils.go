package crypto

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"
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
	CreatedAt time.Time
}

// func ListSecretKeysOld() ([]Key, error) {

// 	var buf bytes.Buffer

// 	gpg := exec.Command(GPGPath, "--list-secret-keys")

// 	grep := exec.Command("grep", "uid")
// 	grep.Stdin, _ = gpg.StdoutPipe()

// 	sed := exec.Command("sed", "-E", "s/uid +//g")
// 	sed.Stdin, _ = grep.StdoutPipe()
// 	sed.Stdout = &buf

// 	grep.Start()
// 	sed.Start()

// 	gpg.Run()
// 	grep.Wait()
// 	sed.Wait()

// 	var list []Key
// 	re := regexp.MustCompile("(.*)<(.*@.*)>")
// 	for {
// 		line, err := buf.ReadString('\n')
// 		fmt.Println(line)
// 		if err == io.EOF {
// 			break
// 		}

// 		// parse line
// 		matches := re.FindAllStringSubmatch(line, -1)
// 		if len(matches) > 0 && len(matches[0]) > 2 {
// 			key := Key{
// 				Name:  matches[0][1],
// 				Email: matches[0][2],
// 			}
// 			// fmt.Printf("key: %v\n", key)
// 			list = append(list, key)
// 		}
// 	}

// 	return list, nil
// }

func GPGListSecretKeys() ([]Key, error) {

	var buf bytes.Buffer

	gpg := exec.Command(GPGPath, "--list-secret-keys", "--display-charset", "utf-8")
	gpg.Stdout = &buf
	gpg.Run()

	var list []Key
	re1 := regexp.MustCompile("sec\\s+(.*)/(.*) (\\d{4}-\\d{2}-\\d{2})")
	re2 := regexp.MustCompile("uid\\s+(.*)<(.*@.*)>")
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}

		var key Key

		// match sec
		matches1 := re1.FindAllStringSubmatch(line, -1)
		if len(matches1) > 0 && len(matches1[0]) > 3 {
			key.PublicKey = strings.TrimSpace(matches1[0][2])
			key.CreatedAt, _ = time.Parse("2006-01-02", strings.TrimSpace(matches1[0][3]))

			// read next line
			line, err := buf.ReadString('\n')
			if err != nil {
				break
			}
			// match uid
			matches2 := re2.FindAllStringSubmatch(line, -1)
			if len(matches2) > 0 && len(matches2[0]) > 2 {
				key.Name = strings.TrimSpace(matches2[0][1])
				key.Email = strings.TrimSpace(matches2[0][2])

				list = append(list, key)
			}
		}
	}
	return list, nil
}

func GPGEncrypt(uid, recipient string, src io.Reader) ([]byte, error) {
	srcfile := path.Join(os.TempDir(), uuid.NewUUID().String())
	dstfile := fmt.Sprintf("%s.gpg", srcfile)
	defer func() {
		os.RemoveAll(srcfile)
		os.RemoveAll(dstfile)
	}()

	// save data to srcfile
	wd, err := os.Create(srcfile)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(wd, src)
	if err != nil {
		return nil, err
	}
	wd.Close()

	// encrypt using gpg
	gpg := exec.Command(GPGPath, "-u", uid, "-se", "-r", recipient, "-o", dstfile, srcfile)
	gpg.Stdin = os.Stdin
	if err = gpg.Run(); err != nil {
		return nil, err
	}

	// read dstfile
	rd, err := os.Open(dstfile)
	if err != nil {
		return nil, err
	}
	defer rd.Close()
	return ioutil.ReadAll(rd)
}

func GPGDecrypt(uid string, src io.Reader) ([]byte, error) {
	dstfile := path.Join(os.TempDir(), uuid.NewUUID().String())
	srcfile := fmt.Sprintf("%s.gpg", dstfile)
	defer func() {
		os.RemoveAll(srcfile)
		os.RemoveAll(dstfile)
	}()

	// save data to srcfile
	wd, err := os.Create(srcfile)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(wd, src)
	if err != nil {
		return nil, err
	}
	wd.Close()

	// decrypt using gpg
	gpg := exec.Command(GPGPath, "-u", uid, "-o", dstfile, srcfile)
	err = gpg.Run()
	if err != nil {
		return nil, err
	}

	// read dstfile
	rd, err := os.Open(dstfile)
	if err != nil {
		return nil, err
	}
	defer rd.Close()
	return ioutil.ReadAll(rd)
}
