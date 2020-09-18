package structs

import (
	"errors"
	"reflect"
)

type MapMarshaler interface {
	MarshalMap() (map[string]interface{}, error)
}

var (
	Tag       = "map"
	Separated = ","
)

var (
	mapEncode = MapEncode{
		tag:       Tag,
		separated: Separated,
	}
)

func MarshalMap(data interface{}) (map[string]interface{}, error) {
	return mapEncode.MarshalMap(data)
}

type MapEncode struct {
	tag       string
	separated string
}

func (e MapEncode) MarshalMap(data interface{}) (map[string]interface{}, error) {
	if data == nil {
		return nil, nil
	}
	if e.separated == "" {
		e.separated = Separated
	}
	if e.tag == "" {
		e.tag = Tag
	}
	if f, ok := data.(MapMarshaler); ok {
		return f.MarshalMap()
	}
	var v = reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("unsupported type: " + v.Type().String())
	}
	var (
		t   = v.Type()
		ret = make(map[string]interface{})
	)
	for i := 0; i < t.NumField(); i++ {
		var sf = t.Field(i)
		isUnexported := sf.PkgPath != ""
		// Anonymous
		if sf.Anonymous {
			at := sf.Type
			if at.Kind() == reflect.Ptr {
				at = at.Elem()
			}
			if isUnexported && at.Kind() != reflect.Struct {
				// Ignore embedded fields of unexported non-struct types.
				continue
			}
			// Recursion call anonymous filed
			ar, err := e.MarshalMap(v.Field(i).Interface())
			if err != nil {
				return nil, err
			}
			for key, value := range ar {
				ret[key] = value
			}
			continue
		} else if isUnexported {
			// Ignore unexported non-embedded fields.
			continue
		}
		// Tag
		tag := sf.Tag.Get(e.tag)
		if tag == "-" {
			continue
		}
		name, _ := parseTag(tag, e.separated)
		if !isValidTag(name) {
			name = sf.Name
		}
		// Type value
		switch sf.Type.Kind() {
		case reflect.Struct, reflect.Ptr, reflect.Interface:
			sr, err := e.MarshalMap(v.Field(i).Interface())
			if err != nil {
				return nil, err
			}
			ret[name] = sr
		default:
			ret[name] = v.Field(i).Interface()
		}
	}
	return ret, nil
}
