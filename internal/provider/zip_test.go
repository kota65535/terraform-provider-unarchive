package provider

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnzip(t *testing.T) {
	files, err := UnzipSource("commons-text-1.9.jar", "**/*.txt", ".")
	assert.Nil(t, err)
	assert.Equal(t, []string{
		"META-INF/LICENSE.txt",
		"META-INF/NOTICE.txt",
	}, files)
}
