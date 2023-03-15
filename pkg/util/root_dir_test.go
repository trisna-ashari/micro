package util_test

import (
	"micro/pkg/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtilRootDir(t *testing.T) {
	rootDir := util.RootDir()

	assert.NotEmpty(t, rootDir)
}
