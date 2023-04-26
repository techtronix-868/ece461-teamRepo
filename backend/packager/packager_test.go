package packager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackageJSONReader(t *testing.T) {
	dir := "../../frontend"

	metadata, err := ReadPackageJson(dir)
	assert.Nil(t, err, "Unexpected item in bagging area")
	assert.NotEmpty(t, metadata)
	assert.NotEmpty(t, metadata.Name)
	assert.NotEmpty(t, metadata.Version)
}
