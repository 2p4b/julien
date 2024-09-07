package form

import (
	"julien/fs"
	"strings"
	"time"
)

type Doc struct {
	meta  map[string]interface{}
	body  string
	entry *fs.Entry
	form  *Form
}

func (doc *Doc) Name() string {
	return strings.Trim(doc.entry.Name(), doc.Ext())
}

func (doc *Doc) Path() string {
	return doc.entry.Path()
}

func (doc *Doc) Size() int64 {
	return doc.entry.Size()
}

func (doc *Doc) Ext() string {
	return doc.entry.Ext()
}

func (doc *Doc) Has(key string) bool {
	_, ok := doc.meta[key]
	return ok
}

func (doc *Doc) Get(key string, defaultValue ...interface{}) interface{} {
	value, ok := doc.meta[key]
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
	}
	return value
}

func (doc *Doc) Body(body ...string) string {
	if len(body) > 0 {
		doc.body = body[0]
	}
	return doc.body
}

func (doc *Doc) GetString(key string, defaultValue ...string) string {
	value, ok := doc.meta[key].(string)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
}

func (doc *Doc) GetInt(key string, defaultValue ...int) int {
	value, ok := doc.meta[key].(int)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
}

func (doc *Doc) GetInt64(key string, defaultValue ...int64) int64 {
	value, ok := doc.meta[key].(int64)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
}

func (doc *Doc) Set(key string, value interface{}) {
	doc.meta[key] = value
}

func (doc *Doc) Fill(frontmatter map[string]interface{}) {
	for key, value := range frontmatter {
		doc.meta[key] = value
	}
}

func (doc *Doc) Dump() ([]byte, error) {
	bytes, err := doc.form.root.driver.Dump(&doc.meta, doc.body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (doc *Doc) Timestamp() time.Time {
	return doc.entry.Timestamp()
}

func (doc *Doc) HasKey(string) bool {
	return false
}

func (doc *Doc) Save() error {
	bytes, err := doc.Dump()
	if err != nil {
		return err
	}
	return doc.entry.Write(bytes)
}
