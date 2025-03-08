package tus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigMissingStore(t *testing.T) {
	c := Config{
		ChunkSize:           1048576 * 15, // 15 MB
		Resume:              true,
		OverridePatchMethod: false,
		Store:               nil,
		Header:              nil,
	}

	assert.NotNil(t, c.Validate())
}

func TestConfigChunkSizeZero(t *testing.T) {
	c := Config{
		ChunkSize:           0,
		Resume:              false,
		OverridePatchMethod: false,
		Store:               nil,
		Header:              nil,
	}

	assert.NotNil(t, c.Validate())
}

func TestConfigValid(t *testing.T) {
	c := DefaultConfig()
	assert.Nil(t, c.Validate())
}
