package form

import (
	"reflect"
	"sort"
	"strings"
)

type Collection struct {
	entries []*Doc
}

func NewDocCollection(entries []*Doc) *Collection {
	return &Collection{
		entries: entries,
	}
}

func (c *Collection) Entries() []*Doc {
	return c.entries
}

func (c *Collection) Count() int {
	return len(c.entries)
}

func (c *Collection) First() *Doc {
	if c.Count() == 0 {
		return nil
	}
	return c.Entries()[0]
}

func (c *Collection) Last() *Doc {
	entries := c.Entries()
	size := len(entries)
	if size == 0 {
		return nil
	}
	return entries[size-1]
}

func (c *Collection) Get(index int) *Doc {
	if index > -1 && index < c.Count() {
		return c.Entries()[index]
	}
	return nil
}

func (c *Collection) Slice(start, end int) *Collection {
	return NewDocCollection(c.Entries()[start:end])
}

func (c *Collection) Having(key string) *Collection {
	entries := make([]*Doc, 0)
	for _, entry := range c.Entries() {
		if entry != nil {
			if entry.Has(key) {
				entries = append(entries, entry)
			}
		}
	}
	return NewDocCollection(entries)
}

func (c *Collection) Without(key string) *Collection {
	entries := make([]*Doc, 0)
	for _, entry := range c.Entries() {
		if entry != nil {
			if !entry.Has(key) {
				entries = append(entries, entry)
			}
		}
	}
	return NewDocCollection(entries)
}

func (c *Collection) Where(key string, value interface{}) *Collection {
	entries := make([]*Doc, 0)
	for _, entry := range c.Entries() {
		if entry != nil {
			if entry.Has(key) && entry.Get(key) == value {
				entries = append(entries, entry)
			}
		}
	}
	return NewDocCollection(entries)
}

func (c *Collection) WhereNot(key string, value interface{}) *Collection {
	entries := make([]*Doc, 0)
	for _, entry := range c.Entries() {
		if entry != nil {
			if !entry.Has(key) || entry.Get(key) != value {
				entries = append(entries, entry)
			}
		}
	}
	return NewDocCollection(entries)
}

func (c *Collection) Search(token string) *Collection {
	entries := make([]*Doc, 0)
	for _, entry := range c.Entries() {
		if entry != nil {
			if strings.Contains(entry.Body(), token) {
				entries = append(entries, entry)
			}
		}
	}
	return NewDocCollection(entries)
}

func (c *Collection) SortBy(key string, order ...string) *Collection {
	entries := make([]*Doc, 0)
	copy(entries, c.Entries())
	sort.SliceStable(entries, func(i, j int) bool {
		aval := entries[i]
		bval := entries[j]
		avalue := aval.Get(key)
		bvalue := bval.Get(key)

		atype := reflect.TypeOf(avalue).String()
		btype := reflect.TypeOf(bvalue).String()

		if atype != btype {
			return false
		}
		sorder := "asc"
		if len(order) > 0 {
			sorder = order[0]
		}
		if sorder == "desc" {
			switch avalue.(type) {
			case string:
				return avalue.(string) > bvalue.(string)
			case int:
				return avalue.(int) > bvalue.(int)
			case int64:
				return avalue.(int64) > bvalue.(int64)
			case float32:
				return avalue.(float32) > bvalue.(float32)
			case float64:
				return avalue.(float64) > bvalue.(float64)
			default:
				return false
			}

		} else {

			switch avalue.(type) {
			case string:
				return avalue.(string) < bvalue.(string)
			case int:
				return avalue.(int) < bvalue.(int)
			case int64:
				return avalue.(int64) < bvalue.(int64)
			case float32:
				return avalue.(float32) < bvalue.(float32)
			case float64:
				return avalue.(float64) < bvalue.(float64)
			default:
				return false
			}
		}
	})
	return NewDocCollection(entries)
}
