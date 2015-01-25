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
	"runtime"
	"strings"
	"time"
)

var (
	GPGPath = "gpg"
)

func init() {
	var p string
	switch runtime.GOOS {
	case "darwin":
		p = "/usr/local/bin/gpg"
	case "linux":
		p = "/usr/bin/gpg"
	}
	_, err := os.Stat(p)
	if err == nil {
		GPGPath = p
	}
}

type Key struct {
	Name      string
	Email     string
	PublicKey string
	SecretKey string
	CreatedAt time.Time
}

func parseBuffer(buf *bytes.Buffer) ([]Key, error) {

	var list []Key
	re1 := regexp.MustCompile("(sec|pub)\\s+(.*)/(.*) (\\d{4}-\\d{2}-\\d{2})")
	re2 := regexp.MustCompile("uid\\s+(.*)<(.*@.*)>")
	for {
		line, err := buf.ReadString('\n')
		// fmt.Printf("%s\n", line)
		if err == io.EOF {
			break
		}

		var key Key

		// match sec
		matches1 := re1.FindAllStringSubmatch(line, -1)
		if len(matches1) > 0 && len(matches1[0]) > 4 {
			key.PublicKey = strings.TrimSpace(matches1[0][3])
			// fmt.Printf("d: %s\n", matches1[0][4])
			key.CreatedAt, _ = time.Parse("2006-01-02", strings.TrimSpace(matches1[0][4]))

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

func GPGListPublicKeys(search string) ([]Key, error) {
	var buf bytes.Buffer
	var gpg *exec.Cmd

	search = strings.TrimSpace(search)
	if len(search) > 0 {
		gpg = exec.Command(GPGPath, "--list-public-keys", "--display-charset", "utf-8", search)
	} else {
		gpg = exec.Command(GPGPath, "--list-public-keys", "--display-charset", "utf-8")
	}
	gpg.Stdout = &buf
	gpg.Run()

	return parseBuffer(&buf)
}

func GPGRecvKey(key string) error {
	gpg := exec.Command(GPGPath, "--batch", "--yes", "--recv-keys", "--display-charset", "utf-8", key)
	return gpg.Run()
}

func GPGListSecretKeys(search string) ([]Key, error) {
	var buf bytes.Buffer
	var gpg *exec.Cmd

	search = strings.TrimSpace(search)
	if len(search) > 0 {
		gpg = exec.Command(GPGPath, "--list-secret-keys", "--display-charset", "utf-8", search)
	} else {
		gpg = exec.Command(GPGPath, "--list-secret-keys", "--display-charset", "utf-8")
	}
	gpg.Stdout = &buf
	gpg.Run()

	return parseBuffer(&buf)
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
	gpg := exec.Command(GPGPath, "--trust-model", "always", "-u", uid, "-se", "-r", recipient, "-o", dstfile, srcfile)
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

func GPGSearch(key string) ([]Key, error) {
	var buf bytes.Buffer

	// search key from gpg keyserver
	gpg := exec.Command(GPGPath, "--batch", "--display-charset", "utf-8", "--search", key)
	gpg.Stdout = &buf
	gpg.Run()

	return parseBuffer(&buf)
}
