package shell

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	var b strings.Builder
	Run([]byte("echo testing"), &b)
	assert.Equal(t, "testing", b.String())
}
