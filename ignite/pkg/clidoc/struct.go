package clidoc

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	// Docs represents a slice of Doc.
	Docs []Doc
	// Doc represents the struct documentation with tag comments.
	Doc struct {
		Key     string
		Value   Docs
		Comment string
	}
)

// Strings convert Docs to string.
func (d Docs) String() string {
	var sb strings.Builder
	d.writeString(&sb, 0)
	return sb.String()
}

// writeString appends the contents of Docs to b's buffer at level.
func (d Docs) writeString(sb *strings.Builder, level int) {
	indent := strings.Repeat("    ", level)
	for _, doc := range d {
		sb.WriteString(fmt.Sprintf("%s- %s: %s\n", indent, doc.Key, doc.Comment))
		if len(doc.Value) > 0 {
			doc.Value.writeString(sb, level+1)
		}
	}
}

// GenDoc to generate documentation from a struct.
func GenDoc(v interface{}) (fields Docs, err error) {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct && t.Kind() != reflect.Ptr {
		return fields, nil
	}
	for i := 0; i < t.NumField(); i++ {
		var (
			field = t.Field(i)
			doc   = field.Tag.Get("doc")
			yaml  = field.Tag.Get("yaml")
		)

		tags := strings.Split(yaml, ",")
		name := tags[0]
		if name == "" {
			name = strings.ToLower(field.Name)
		}
		if len(tags) > 1 && strings.Contains(tags[1], "inline") {
			elemFields, err := GenDoc(reflect.New(field.Type).Elem().Interface())
			if err != nil {
				return nil, err
			}
			fields = append(fields, elemFields...)
			continue
		}

		var elemFields Docs
		switch field.Type.Kind() {
		case reflect.Struct:
			elemFields, err = GenDoc(reflect.New(field.Type).Elem().Interface())
			if err != nil {
				return nil, err
			}
		case reflect.Ptr:
			elemFields, err = GenDoc(reflect.New(field.Type.Elem()).Elem().Interface())
			if err != nil {
				return nil, err
			}
		case reflect.Slice:
			name = fmt.Sprintf("%s [array]", name)
			elemFields, err = GenDoc(reflect.New(field.Type.Elem()).Elem().Interface())
			if err != nil {
				return nil, err
			}
		default:
		}
		fields = append(fields, Doc{
			Key:     name,
			Comment: doc,
			Value:   elemFields,
		})
	}

	return fields, nil
}
