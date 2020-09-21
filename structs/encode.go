package structs

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
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
		name, tagOptions := parseTag(tag, e.separated)
		if !isValidTag(name) {
			name = sf.Name
		}
		vf := v.Field(i)
		if vf.Type().Kind() == reflect.Ptr {
			vf = vf.Elem()
		}
		// Type value
		switch vf.Type().Kind() {
		case reflect.Struct, reflect.Ptr, reflect.Interface:
			if tagOptions.Contains("json", e.separated) {
				if jsonMarshaler, ok := v.Field(i).Interface().(json.Marshaler); ok {
					data, err := jsonMarshaler.MarshalJSON()
					if err != nil {
						return nil, err
					}
					ret[name] = data
					continue
				} else {
					data, err := json.Marshal(v.Field(i).Interface())
					if err != nil {
						return nil, err
					}
					ret[name] = data
					continue
				}
			} else if tagOptions.Contains("string", e.separated) {
				// string tag option
				if textMarshaler, ok := v.Field(i).Interface().(encoding.TextMarshaler); ok {
					data, err := textMarshaler.MarshalText()
					if err != nil {
						return nil, err
					}
					ret[name] = data
					continue
				} else if jsonMarshaler, ok := v.Field(i).Interface().(json.Marshaler); ok {
					data, err := jsonMarshaler.MarshalJSON()
					if err != nil {
						return nil, err
					}
					ret[name] = data
					continue
				} else if binaryMarshaler, ok := v.Field(i).Interface().(encoding.BinaryMarshaler); ok {
					data, err := binaryMarshaler.MarshalBinary()
					if err != nil {
						return nil, err
					}
					ret[name] = data
					continue
				} else {
					return nil, fmt.Errorf("filed: %s with 'string' tag have not implment TextMarshaler/json.Marshaler/BinaryMarshaler", sf.Name)
				}
			} else if tagOptions.Contains("bytes", e.separated) {
				// bytes tag option
				if binaryMarshaler, ok := v.Field(i).Interface().(encoding.BinaryMarshaler); ok {
					data, err := binaryMarshaler.MarshalBinary()
					if err != nil {
						return nil, err
					}
					ret[name] = data
					continue
				} else if jsonMarshaler, ok := v.Field(i).Interface().(json.Marshaler); ok {
					data, err := jsonMarshaler.MarshalJSON()
					if err != nil {
						return nil, err
					}
					ret[name] = data
					continue
				} else if textMarshaler, ok := v.Field(i).Interface().(encoding.TextMarshaler); ok {
					data, err := textMarshaler.MarshalText()
					if err != nil {
						return nil, err
					}
					ret[name] = data
					continue
				} else {
					return nil, fmt.Errorf("filed: %s with 'bytes' tag have not implment BinaryMarshaler/json.Marshaler/TextMarshaler", sf.Name)
				}
			} else {
				// Recursion call anonymous filed
				sr, err := e.MarshalMap(v.Field(i).Interface())
				if err != nil {
					return nil, err
				}
				ret[name] = sr
			}
		default:
			ret[name] = v.Field(i).Interface()
		}
	}
	return ret, nil
}
