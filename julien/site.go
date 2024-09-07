package julien

import (
	"fmt"
	"julien/driver"
	"os"
	"path"
	"strings"
)

type Site struct {
	meta map[string]interface{} `yaml:"meta"`
	body string                 `yaml:"content"`
	path string
}

func LoadSite(ppath string, s *Site) {
	yamler := driver.Yaml{}
	ppath = path.Clean(ppath)
	stats, err := os.Stat(ppath)
	if err != nil {
		panic(err)
	}
	if stats.IsDir() {
		err = fmt.Errorf(ppath + " is a directory not a file")
		panic(err)
	}
	cbytes, err := os.ReadFile(ppath)
	if err != nil {
		panic(err)
	}
	frontmatter, body, err := yamler.Parse(cbytes)
	if err != nil {
		panic(err)
	}
	s.path = ppath
	s.meta = *frontmatter
	s.body = body
}

func FindSite(ppath string) (*Site, error) {
	yamler := driver.Yaml{}
	ppath = path.Clean(ppath)
	stats, err := os.Stat(ppath)
	if err != nil {
		return nil, err
	}
	if stats.IsDir() {
		return nil, fmt.Errorf(ppath + " is a directory not a file")
	}
	cbytes, err := os.ReadFile(ppath)
	if err != nil {
		return nil, err
	}
	frontmatter, body, err := yamler.Parse(cbytes)
	if err != nil {
		return nil, err
	}
	return &Site{
		path: ppath,
		meta: *frontmatter,
		body: body,
	}, nil

}

func (site *Site) URL(upath ...string) string {
	root := site.GetString("url", "")
	if len(upath) == 0 {
		return root
	}
	ppath := path.Clean(upath[0])
	ppath = strings.TrimLeft(ppath, "/")
	ppath = strings.TrimRight(ppath, "/")
	return strings.TrimRight(root, "/") + "/" + ppath
}

func (site *Site) Get(key string, defaultValue ...interface{}) interface{} {
	value, ok := site.meta[key]
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
	}
	return value
}

func (site *Site) Has(key string) bool {
	_, ok := site.meta[key]
	return ok
}

func (site *Site) Body(body ...string) string {
	if len(body) > 0 {
		site.body = body[0]
	}
	return site.body
}

func (site *Site) GetString(key string, defaultValue ...string) string {
	value, ok := site.meta[key].(string)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
}

func (site *Site) GetInt(key string, defaultValue ...int) int {
	value, ok := site.meta[key].(int)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
}

func (site *Site) GetInt64(key string, defaultValue ...int64) int64 {
	value, ok := site.meta[key].(int64)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
}
