package web

import (
	jutils "julien/utils"
)

type FormData struct {
	Name      string                 `json:"name"`
	Data      map[string]interface{} `json:"data"`
	Errors    map[string][]string    `json:"errors"`
	Timestamp int64                  `json:"timestamp"`
}

func (f *FormData) Exist() bool {
	return f.Name == ""
}

func (f *FormData) Get(key string, defaultValue ...interface{}) interface{} {
	value, ok := f.Data[key]
	if ok {
		return value
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
	}
	return value
}

func (f *FormData) GetErrors(key string) interface{} {
	return f.Errors[key]
}

func (f *FormData) HasErrors(key string, tags ...string) bool {
	errors := f.Errors[key]
	if len(errors) == 0 {
		return false
	}
	if len(tags) == 0 {
		return true
	}
	for _, tag := range tags {
		if jutils.ArrayIncludes(errors, tag) {
			return true
		}
	}
	return false
}
