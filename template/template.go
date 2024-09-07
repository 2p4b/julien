package template

import (
	"fmt"
	"julien/driver"
	"os"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/template/django/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/template/jet/v2"
	"github.com/gofiber/template/mustache/v2"
)

type Template struct {
	Meta    map[string]interface{}
	Content string
	Path    string
}

func Find(ppath string, index ...string) (*Template, error) {
	yamler := driver.Yaml{}
	stats, err := os.Stat(ppath)
	if len(index) == 0 {
		index = append(index, "index.md")
	}
	if err != nil {
		return nil, err
	}
	if !stats.IsDir() {
		return nil, fmt.Errorf(ppath + " is not a directory")
	}
	cbytes, err := os.ReadFile(path.Join(ppath, index[0]))
	if err != nil {
		return nil, err
	}
	frontmatter, content, err := yamler.Parse(cbytes)
	if err != nil {
		return nil, err
	}
	return &Template{
		Path:    ppath,
		Meta:    *frontmatter,
		Content: content,
	}, nil

}

func (tmpl *Template) Engine(reload bool) fiber.Views {
	name := tmpl.GetString("type", "html")
	switch name {
	case "pongo":
		PogoInit()
		vengine := django.New(tmpl.Path, "."+tmpl.GetString("ext", "html"))
		vengine.Reload(reload)
		return vengine

	case "html":
		vengine := html.New(tmpl.Path, "."+tmpl.GetString("ext", "html"))
		vengine.Reload(reload)
		return vengine

	case "mustache":
		vengine := mustache.New(tmpl.Path, "."+tmpl.GetString("ext", "mustache"))
		vengine.Reload(reload)
		return vengine

	case "jet":
		vengine := jet.New(tmpl.Path, "."+tmpl.GetString("ext", "jet"))
		vengine.Reload(reload)
		return vengine

	default:
		log.Error("template " + name + " engine not found default to html")
		vengine := html.New(tmpl.Path, "."+tmpl.GetString("ext", "html"))
		vengine.Reload(reload)
		return vengine
	}
}

func (tmpl *Template) Get(key string, defaultValue ...interface{}) interface{} {
	value, ok := tmpl.Meta[key]
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
	}
	return value
}

func (tmpl *Template) GetString(key string, defaultValue ...string) string {
	value, ok := tmpl.Meta[key].(string)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
}

func (tmpl *Template) GetInt(key string, defaultValue ...int) int {
	value, ok := tmpl.Meta[key].(int)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
}

func (tmpl *Template) GetInt64(key string, defaultValue ...int64) int64 {
	value, ok := tmpl.Meta[key].(int64)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
}
