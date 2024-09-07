package driver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYaml_Init(t *testing.T) {
	y := &Yaml{}
	assert.NoError(t, y.Init(nil))
}

func TestYaml_Dump(t *testing.T) {
	driver := &Yaml{}

	frontmatter := &map[string]interface{}{
		"title": "My Title",
		"tags":  []string{"tag1", "tag2"},
	}
	content := "This is the content."

	expected := "---\ntags:\n- tag1\n- tag2\ntitle: My Title\n---\nThis is the content."

	result, err := driver.Dump(frontmatter, content)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(result))
}

// Add more tests
