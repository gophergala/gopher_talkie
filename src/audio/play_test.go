package audio

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
	"time"
)

func TestPlayAIFF(t *testing.T) {
	p := path.Join(os.TempDir(), "test_record.aiff")
	os.RemoveAll(p)
	fmt.Printf("%s\n", p)

	fmt.Printf("Recording...\n")
	options := RecordOptions{
		MaxDuration: time.Duration(2) * time.Second,
		FilePath:    p,
	}
	err := RecordAIFF(options)
	assert.Nil(t, err)

	stat, err := os.Stat(p)
	assert.Nil(t, err)
	assert.True(t, stat.Size() > 10000)

	// test play it
	fmt.Printf("Playing...\n")
	err = PlayAIFF(p, nil)
	assert.Nil(t, err)

	fmt.Println("DONE")
}
