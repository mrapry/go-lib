package goson

/*
	Support deeply nested struct, with time complexity is total fields (include total fields in each slice/map)
*/

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	whiteListType = map[string]bool{
		"ptr": true, "slice": true, "map": true,
	}
)

func makeZeroField(obj reflect.Value) {
	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		if !field.CanSet() { // break if field cannot set a value (usually an unexported field in struct), to avoid a panic
			return
		}

		// if field is struct or types (nested struct or slice), process with recursive
		processTypeOfValue(field)

		jsonTag := obj.Type().Field(i).Tag.Get("json")
		jsonTags := strings.Split(jsonTag, ",")
		if len(jsonTags) > 1 && jsonTags[1] == "omitempty" {
			field.Set(reflect.Zero(reflect.TypeOf(field.Interface())))
		}
	}
}

func processPointer(ptr reflect.Value) {
	val := ptr.Interface()
	if ptr.IsNil() {
		ptr = reflect.ValueOf(&val).Elem()
		ptr = reflect.New(ptr.Elem().Type().Elem()) // create from new domain model type of field
	}

	obj := ptr.Elem()
	processTypeOfValue(obj)
}

func fetchSliceType(slice reflect.Value) {
	for i := 0; i < slice.Len(); i++ {
		obj := slice.Index(i)
		processTypeOfValue(obj)
	}
}

func fetchMapType(mapper reflect.Value) {
	for _, idx := range mapper.MapKeys() {
		obj := mapper.MapIndex(idx)
		processTypeOfValue(obj)
	}
}

func processTypeOfValue(obj reflect.Value) {
	switch obj.Kind() {
	case reflect.Struct:
		makeZeroField(obj)
	case reflect.Slice:
		fetchSliceType(obj)
	case reflect.Map:
		fetchMapType(obj)
	case reflect.Ptr:
		processPointer(obj)
	}
}

// MakeZeroOmitempty tools for make Zero value in field contains `json: omitempty` tag
func MakeZeroOmitempty(obj interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	refValue := reflect.ValueOf(obj)
	if ok := whiteListType[refValue.Kind().String()]; !ok {
		return fmt.Errorf("invalid type %v: accept pointer, slice, and map", refValue.Kind())
	}

	processTypeOfValue(refValue)
	return
}
