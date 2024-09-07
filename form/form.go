package form

import (
	"julien/contract"
	"julien/fs"
	"path"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

var EMPTY_DOC_ARRAY = make([]*Doc, 0)

type Root struct {
	disk   *fs.Disk
	data   *fs.Disk
	driver contract.Driver
}

type Form struct {
	meta  map[string]interface{}
	body  string
	entry *fs.Entry
	root  *Root
}

func Init(fdisk *fs.Disk, ddisk *fs.Disk, driver contract.Driver) Root {
	return Root{
		data:   ddisk,
		disk:   fdisk,
		driver: driver,
	}
}

func (root *Root) Find(ppath string) (*Form, error) {
	entry, err := root.disk.Find(path.Clean(ppath))
	if err != nil {
		return nil, err
	}
	raw, err := entry.Read()
	if err != nil {
		return nil, err
	}

	frontmatter, body, err := root.driver.Parse(raw)
	if err != nil {
		return nil, err
	}

	return &Form{
		meta:  *frontmatter,
		body:  body,
		entry: entry,
		root:  root,
	}, nil

}

func (root *Root) Open(ppath string) *Form {
	fm, err := root.Find(ppath)
	if err != nil {
		return nil
	}
	return fm
}

func (root *Root) Dump(ppath string, content []byte) (*Form, error) {
	err := root.disk.Dump(ppath, content)
	if err != nil {
		return nil, err
	}
	return root.Find(ppath)
}

func (fm *Form) crate_entry_doc(entry *fs.Entry) (*Doc, error) {
	raw, err := entry.Read()
	if err != nil {
		return nil, err
	}

	frontmatter, body, err := fm.root.driver.Parse(raw)
	if err != nil {
		return nil, err
	}

	return &Doc{
		form:  fm,
		meta:  *frontmatter,
		body:  body,
		entry: entry,
	}, nil
}

func (fm *Form) Name() string {
	return strings.Trim(fm.entry.Name(), fm.Ext())
}

func (form *Form) Path() string {
	return form.entry.Path()
}

func (form *Form) Size() int64 {
	return form.entry.Size()
}

func (form *Form) Ext() string {
	return form.entry.Ext()
}

func (form *Form) Get(key string, defaultValue ...interface{}) interface{} {
	value, ok := form.meta[key]
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
	}
	return value
}

func (form *Form) Body(body ...string) string {
	if len(body) > 0 {
		form.body = body[0]
	}
	return form.body
}

func (form *Form) GetString(key string, defaultValue ...string) string {
	value, ok := form.meta[key].(string)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
}

func (form *Form) GetInt(key string, defaultValue ...int) int {
	value, ok := form.meta[key].(int)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
}

func (form *Form) GetInt64(key string, defaultValue ...int64) int64 {
	value, ok := form.meta[key].(int64)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
}

func (form *Form) Set(key string, value interface{}) {
	form.meta[key] = value
}

func (form *Form) Fill(frontmatter map[string]interface{}) {
	for key, value := range frontmatter {
		form.meta[key] = value
	}
}

func (fm *Form) Dump() ([]byte, error) {
	bytes, err := fm.root.driver.Dump(&fm.meta, fm.body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (fm *Form) List() ([]*Doc, error) {
	doclist := make([]*Doc, 0)
	entries, err := fm.root.data.List(fm.Name())
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		doc, err := fm.crate_entry_doc(entry)
		if err != nil {
			continue
		}
		doclist = append(doclist, doc)
	}
	return doclist, nil
}

func (fm *Form) Collection() *Collection {
	docs, err := fm.List()
	if err != nil {
		log.Error(err)
		return NewDocCollection(EMPTY_DOC_ARRAY)
	}
	return NewDocCollection(docs)
}

func (fm *Form) Find(name string) (*Doc, error) {
	entry, err := fm.root.data.Find(path.Join(fm.Name(), name))
	if err != nil {
		return nil, err
	}
	return fm.crate_entry_doc(entry)
}

func (fm *Form) Open(name string) *Doc {
	doc, err := fm.Find(name)
	if err != nil {
		return nil
	}
	return doc
}

func (fm *Form) Compose(name string, content []byte) (*Doc, error) {
	err := fm.root.data.Dump(path.Join(fm.Name(), name), content)
	if err != nil {
		return nil, err
	}
	return fm.Find(name)
}

func (fm *Form) Timestamp() time.Time {
	return fm.entry.Timestamp()
}
