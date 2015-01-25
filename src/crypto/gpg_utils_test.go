package crypto

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	// "os"
	"testing"
)

func TestGPGListSecretKeys(t *testing.T) {
	list, err := GPGListSecretKeys("")
	assert.Nil(t, err)
	assert.True(t, len(list) > 0)

	for i := range list {
		fmt.Printf("%v\n", list[i])
		assert.NotEmpty(t, list[i].Name)
		assert.NotEmpty(t, list[i].Email)
	}
}

func TestGPGListPublicKeys(t *testing.T) {
	list, err := GPGListPublicKeys("")
	assert.Nil(t, err)
	assert.True(t, len(list) > 0)

	for i := range list {
		fmt.Printf("%v\n", list[i])
		assert.NotEmpty(t, list[i].Name)
		assert.NotEmpty(t, list[i].Email)
	}
}

func TestGPGEncrypt(t *testing.T) {
	src := []byte("hello")
	uid := "B44966D6"
	recipient := "B44966D6"
	dst, err := GPGEncrypt(uid, recipient, bytes.NewBuffer(src))
	assert.Nil(t, err)
	assert.NotEmpty(t, dst)

	src2, err := GPGDecrypt(uid, bytes.NewBuffer(dst))
	assert.Nil(t, err)
	assert.Equal(t, src, src2)
}
