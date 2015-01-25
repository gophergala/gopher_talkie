package crypto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListSecretKeys(t *testing.T) {
	list, err := ListSecretKeys()
	assert.Nil(t, err)
	assert.True(t, len(list) > 0)

	for i := range list {
		fmt.Printf("%v\n", list[i])
		assert.NotEmpty(t, list[i].Name)
		assert.NotEmpty(t, list[i].Email)
	}

}
