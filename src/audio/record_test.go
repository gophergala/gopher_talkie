package audio

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
	"time"
)

func TestRecordAIFF(t *testing.T) {
	p := path.Join(os.TempDir(), "test_record.aiff")
	os.RemoveAll(p)
	fmt.Printf("%s\n", p)

	options := RecordOptions{
		MaxDuration: time.Duration(2) * time.Second,
		FilePath:    p,
	}
	err := RecordAIFF(options)
	assert.Nil(t, err)

	stat, err := os.Stat(p)
	assert.Nil(t, err)
	assert.True(t, stat.Size() > 10000)

	// test cancel recording using signal
	p2 := path.Join(os.TempDir(), "test_record2.aiff")
	os.RemoveAll(p2)
	fmt.Printf("%s\n", p2)

	dur := time.Duration(10) * time.Second
	sig := make(chan int)
	options2 := RecordOptions{
		MaxDuration: dur,
		FilePath:    p2,
		StopSignal:  sig,
	}
	go RecordAIFF(options2)

	// cancel recording immediately
	sig <- 1

	stat2, err := os.Stat(p2)
	assert.Nil(t, err)
	assert.True(t, stat2.Size() < 10000)
}
