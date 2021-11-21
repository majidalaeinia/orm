package orm

import (
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)


type FieldMetadata struct {
	Name             string
	IsPK    		 bool
	Virtual 		 bool
	Type    reflect.Type
}

type FieldTag struct {
	Name  string
	Virtual bool
	PK    bool
}

type HasFields interface {
	Fields() []*FieldMetadata
}

func fieldMetadataFromTag(t string) FieldTag {
	if t == "" {
		return FieldTag{}
	}
	tuples := strings.Split(t, " ")
	var tag FieldTag
	kv := map[string]string{}
	for _, tuple := range tuples {
		parts := strings.Split(tuple, "=")
		key := parts[0]
		value := parts[1]
		kv[key] = value
		if key == "name" {
			tag.Name = value
		} else if key == "pk" {
			tag.PK = true
		} else if key == "virtual" {
			tag.Virtual = true
		}
	}
	return tag
}

func fieldsOf(obj interface{}, dialect *Dialect) []*FieldMetadata {
	hasFields, is := obj.(HasFields)
	if is {
		return hasFields.Fields()
	}
	t := reflect.TypeOf(obj)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()

	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	var fms []*FieldMetadata
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		tagParsed := fieldMetadataFromTag(ft.Tag.Get("orm"))
		fm := &FieldMetadata{}
		fm.Type = ft.Type
		if tagParsed.Name != "" {
			fm.Name = tagParsed.Name
		} else {
			fm.Name = strcase.ToSnake(ft.Name)
		}
		if tagParsed.PK || strings.ToLower(ft.Name) == "id" {
			fm.IsPK = true
		}
		if tagParsed.Virtual || ft.Type.Kind() == reflect.Struct || ft.Type.Kind() == reflect.Slice || ft.Type.Kind() == reflect.Ptr {
			fm.Virtual = true
		}
		fms = append(fms, fm)
	}
	return fms
}