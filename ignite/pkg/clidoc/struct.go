package clidoc

import (
	"fmt"
	"reflect"
	"strings"
)

type Doc struct {
	Key     string
	Value   []Doc
	Comment string
}

// GenDoc to generate documentation
func GenDoc(v interface{}) (fields []Doc, err error) {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return fields, nil
	}
	for i := 0; i < t.NumField(); i++ {
		var (
			field = t.Field(i)
			doc   = field.Tag.Get("doc")
			tag   = field.Tag.Get("yaml")
		)

		tags := strings.Split(tag, ",")
		tag = tags[0]
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}
		if len(tags) > 1 && strings.Contains(tags[1], "inline") {
			elemFields, err := GenDoc(reflect.New(field.Type).Elem().Interface())
			if err != nil {
				return nil, err
			}
			fields = append(fields, elemFields...)
		}

		var elemFields []Doc
		if field.Type.Kind() == reflect.Slice {
			tag = fmt.Sprintf("%s [array]", tag)
			elemFields, err = GenDoc(reflect.New(field.Type.Elem()).Elem().Interface())
			if err != nil {
				return nil, err
			}
		}
		if field.Type.Kind() == reflect.Struct {
			elemFields, err = GenDoc(reflect.New(field.Type).Elem().Interface())
			if err != nil {
				return nil, err
			}
		}
		fields = append(fields, Doc{
			Key:     tag,
			Comment: doc,
			Value:   elemFields,
		})
	}

	return fields, nil
}
