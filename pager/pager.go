package pager

import (
	"julien/contract"
	"julien/fs"
	"path"
	"strings"
)

var EMPTY_PAGES = make([]*Page, 0)

type Page struct {
	meta     map[string]interface{}
	body     string
	entry    *fs.Entry
	root     *Root
	extended map[interface{}]interface{}
}

type Root struct {
	disk   *fs.Disk
	driver contract.Driver
}

func Init(disk *fs.Disk, driver contract.Driver) Root {
	return Root{
		disk:   disk,
		driver: driver,
	}
}

func (root *Root) create_entry_page(entry *fs.Entry) (*Page, error) {
	raw, err := entry.Read()
	if err != nil {
		return nil, err
	}

	frontmatter, body, err := root.driver.Parse(raw)
	if err != nil {
		return nil, err
	}

	page := Page{
		meta:     *frontmatter,
		body:     body,
		entry:    entry,
		root:     root,
		extended: make(map[interface{}]interface{}),
	}

	if page.IsFile() && !page.IsRootIndex() {
		// Attempt to parse pareent dir index entry
		// if current entry is a file
		dirpage, err := root.Find(path.Dir(page.Path()))
		if err == nil {
			spec, ok := dirpage.Get("page").(map[interface{}]interface{})
			if ok {
				page.extended = spec
			}
		}
	}

	return &page, nil
}

func (root *Root) List(ppath string) ([]*Page, error) {
	pages := make([]*Page, 0)
	entries, err := root.disk.List(ppath)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsIndex() {
			var err error = nil
			var page *Page = nil

			if entry.IsFile() {
				page, err = root.create_entry_page(entry)
			} else if entry.IsDir() {
				page, err = root.Find(entry.Path())
			}

			if err != nil {
				// Ignore directories that can't be index
				if entry.IsDir() {
					continue
				}
				return nil, err
			}

			if page == nil && err == nil {
				continue
			}
			pages = append(pages, page)
		}
	}
	return pages, nil
}

func (root *Root) Find(ppath string) (*Page, error) {
	entry, err := root.disk.Find(path.Clean(ppath))
	if err != nil {
		return nil, err
	}
	return root.create_entry_page(entry)
}

func (root *Root) Open(ppath string) *Page {
	page, err := root.Find(ppath)
	if err != nil {
		return nil
	}
	return page
}

func (root *Root) Dump(ppath string, content []byte) (*Page, error) {
	err := root.disk.Dump(ppath, content)
	if err != nil {
		return nil, err
	}
	return root.Find(ppath)
}

func (page *Page) Name() string {
	return strings.Trim(page.entry.Name(), page.Ext())
}

func (page *Page) Body(body ...string) string {
	if len(body) > 0 {
		page.body = body[0]
	}
	return page.body
}

func (page *Page) Entries() []*Page {
	if !page.IsDir() {
		return EMPTY_PAGES
	}
	pages, err := page.root.List(page.EPath())
	if err != nil {
		return EMPTY_PAGES
	}
	return pages
}
func (page *Page) Collection() *Collection {
	return NewPageCollection(page.Entries())
}

func (page *Page) EPath() string {
	return page.entry.Path()
}

func (page *Page) IsIndex() bool {
	if page.IsDir() {
		return false
	}
	disk := page.root.disk
	ifname := disk.Index() + "." + disk.Ext()
	filename := path.Base(page.EPath())
	if filename == disk.Index() || filename == (ifname) {
		return true
	}
	return false
}

func (page *Page) Path() string {
	epath := page.EPath()
	if page.entry.IsDir() {
		return epath
	}

	filename := path.Base(epath)
	if filename == "." {
		return "/"
	}
	if page.IsIndex() {
		ppath := path.Dir(epath)
		if ppath == "/" || ppath == "." {
			return ""
		}
		return ppath
	}
	return strings.TrimRight(epath, "."+page.Ext())
}

func (page *Page) APath() string {
	ppath := page.Path()
	ppath = strings.TrimLeft(ppath, "/")
	ppath = strings.TrimRight(ppath, "/")
	return "/" + ppath
}

func (page *Page) Size() int64 {
	return page.entry.Size()
}

func (page *Page) Ext() string {
	return page.entry.Ext()
}

func (page *Page) Get(key string, defaultValue ...interface{}) interface{} {
	value, ok := page.meta[key]
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
	}
	return value
}

func (page *Page) GetString(key string, defaultValue ...string) string {
	value, ok := page.meta[key].(string)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
}

func (page *Page) GetInt(key string, defaultValue ...int) int {
	value, ok := page.meta[key].(int)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
}

func (page *Page) GetInt64(key string, defaultValue ...int64) int64 {
	value, ok := page.meta[key].(int64)
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
}

func (page *Page) Has(key string) bool {
	_, ok := page.meta[key]
	return ok
}

func (page *Page) Set(key string, value interface{}) {
	page.meta[key] = value
}

func (page *Page) Fill(frontmatter map[string]interface{}) {
	for key, value := range frontmatter {
		page.meta[key] = value
	}
}

func (page *Page) Dump() ([]byte, error) {
	bytes, err := page.root.driver.Dump(&page.meta, page.body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (page *Page) Save() error {
	bytes, err := page.Dump()
	if err != nil {
		return err
	}
	return page.entry.Write(bytes)
}

func (page *Page) Metadata() map[string]interface{} {
	return page.meta
}

func (page *Page) IsDir() bool {
	return page.entry.IsDir()
}

func (page *Page) IsFile() bool {
	return page.entry.IsFile()
}

func (page *Page) IsRootIndex() bool {
	parts := strings.Split(page.Path(), "/")
	plen := len(parts)
	if page.IsFile() && plen == 1 && parts[0] == "index" {
		return true
	}
	return false

}

func (page *Page) GetStringValueSpecOrNameRecusive(key string) string {
	value, ok := page.Get(key).(string)
	if ok {
		return value
	}

	if page.IsFile() {
		value, ok = page.extended[key].(string)
		if ok {
			return value
		}
	}

	if !page.IsRootIndex() {
		dirpath := path.Dir(page.Path())
		dirpage, err := page.root.Find(dirpath)
		if err != nil {
			return page.Name()
		}

		spec, ok := dirpage.Get("page").(map[interface{}]interface{})
		if ok {
			specvalue, ok := spec[key].(string)
			if ok {
				return specvalue
			}
		}
		keval, ok := dirpage.Get(key).(string)
		if ok {
			return keval
		}
	}

	return page.Name()

}

func (page *Page) View() string {
	return page.GetStringValueSpecOrNameRecusive("view")
}

func (page *Page) Layout() string {
	return page.GetStringValueSpecOrNameRecusive("layout")
}
