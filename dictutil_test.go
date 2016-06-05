package dictutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsGreetingMsg(t *testing.T) {
	testcases := map[string]bool{
		"hello": true,
		"hi":    true,
		"helu":  false,
	}

	for key, val := range testcases {
		_, v := IsGreetingMsg(key)
		assert.Equal(t, val, v)
	}
}
