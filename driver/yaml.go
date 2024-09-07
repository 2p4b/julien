package driver

import (
	"fmt"
	"regexp"

	"github.com/go-yaml/yaml"
)

type Yaml struct{}

func (y *Yaml) Init(_ interface{}) error {
	return nil
}

func (driver *Yaml) Dump(frontmatter *map[string]interface{}, content string) ([]byte, error) {
	yaml, err := yaml.Marshal(frontmatter)
	if err != nil {
		return nil, err
	}
	return []byte("---\n" + string(yaml) + "---\n" + content), nil
}

func (driver *Yaml) Parse(raw []byte) (*map[string]interface{}, string, error) {
	head, content, err := driver.Parts(string(raw))
	if err != nil {
		return nil, "", err
	}

	frontmatter := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(head), &frontmatter)
	if err != nil {
		return nil, "", err
	}

	return &frontmatter, content, nil

}

func (driver *Yaml) Parts(content string) (string, string, error) {
	if content == "" {
		return "", "", nil
	}

	re := regexp.MustCompile(`(?m)^\s*---\n([\s\S]+?)\n---(\n([\S\s]+))?$`)
	matches := re.FindStringSubmatch(content)
	matched_len := len(matches)

	if matched_len == 2 {
		return matches[1], "", nil
	}

	if matched_len == 4 {
		return matches[1], matches[3], nil
	}

	return "", "", fmt.Errorf("no match found")
}
