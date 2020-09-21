package structs

import (
	"bytes"
	"math"
	"reflect"
	"testing"
)

type AnonymousT struct {
	Anonymous string
}

type T struct {
	Name string
}

type StringTag struct {
	Data interface{} `map:"data,string"`
}

type BytesTag struct {
	Data interface{} `map:"data,bytes"`
}

type JSONTag struct {
	Data interface{} `map:"data,json"`
}

type Text struct {
	Name string
}

func (t Text) MarshalText() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteString(t.Name)
	return buf.Bytes(), nil
}

type Binary struct {
	Name string
}

func (t Binary) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteString(t.Name)
	return buf.Bytes(), nil
}

type JSON struct {
	Name string
}

func (t JSON) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteByte('{')
	buf.WriteString(`"name":"`)
	buf.WriteString(t.Name)
	buf.WriteByte('"')
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

type Pointer struct {
	Int    *int64
	Float  *float64
	String *string
	Bytes  *[]byte
}

func Test_encodeState_MarshalMap(t *testing.T) {
	type fields struct {
		tag       string
		separated string
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "nil",
			fields:  fields{},
			args:    args{},
			want:    nil,
			wantErr: false,
		},
		{
			name: "string",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: struct {
					String string
				}{
					String: "I'am a string",
				},
			},
			want: map[string]interface{}{
				"String": "I'am a string",
			},
			wantErr: false,
		},
		{
			name: "bytes",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: struct {
					Bytes []byte
				}{
					Bytes: []byte{1, 2, 4, 5, 6, 7, 8},
				},
			},
			want: map[string]interface{}{
				"Bytes": []byte{1, 2, 4, 5, 6, 7, 8},
			},
			wantErr: false,
		},
		{
			name: "float",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: struct {
					Float32 float32
					Float64 float64
				}{
					Float32: math.MaxFloat32,
					Float64: math.MaxFloat64,
				},
			},
			want: map[string]interface{}{
				"Float32": float32(math.MaxFloat32),
				"Float64": float64(math.MaxFloat64),
			},
			wantErr: false,
		},
		{
			name: "int",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: struct {
					Int   int
					Int8  int8
					Int16 int16
					Int32 int32
					Int64 int64
				}{
					Int:   -1234567,
					Int8:  math.MaxInt8,
					Int16: math.MaxInt16,
					Int32: math.MaxInt32,
					Int64: int64(math.MaxInt64),
				},
			},
			want: map[string]interface{}{
				"Int":   -1234567,
				"Int8":  int8(math.MaxInt8),
				"Int16": int16(math.MaxInt16),
				"Int32": int32(math.MaxInt32),
				"Int64": int64(math.MaxInt64),
			},
			wantErr: false,
		},
		{
			name: "uint",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: struct {
					Uint   uint
					Uint8  uint8
					Uint16 uint16
					Uint32 uint32
					Uint64 uint64
				}{
					Uint:   1234567,
					Uint8:  math.MaxUint8,
					Uint16: math.MaxUint16,
					Uint32: math.MaxUint32,
					Uint64: uint64(math.MaxUint64),
				},
			},
			want: map[string]interface{}{
				"Uint":   uint(1234567),
				"Uint8":  uint8(math.MaxUint8),
				"Uint16": uint16(math.MaxUint16),
				"Uint32": uint32(math.MaxUint32),
				"Uint64": uint64(math.MaxUint64),
			},
			wantErr: false,
		},
		{
			name: "tag",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: struct {
					Int    int    `map:"int"`
					Int8   int8   `map:"int_8"`
					Int16  int16  `map:"int_16"`
					Int32  int32  `map:"int_32"`
					Int64  int64  `map:"int_64"`
					Ignore string `map:"-"`
					JSON   string `json:"json"`
				}{
					Int:    -1234567,
					Int8:   math.MaxInt8,
					Int16:  math.MaxInt16,
					Int32:  math.MaxInt32,
					Int64:  int64(math.MaxInt64),
					Ignore: "Ignore",
					JSON:   "I'am JSON",
				},
			},
			want: map[string]interface{}{
				"int":    -1234567,
				"int_8":  int8(math.MaxInt8),
				"int_16": int16(math.MaxInt16),
				"int_32": int32(math.MaxInt32),
				"int_64": int64(math.MaxInt64),
				"JSON":   "I'am JSON",
			},
			wantErr: false,
		},
		{
			name: "anonymous",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: struct {
					Name string
					AnonymousT
				}{
					Name: "Name",
					AnonymousT: AnonymousT{
						Anonymous: "Anonymous",
					},
				},
			},
			want: map[string]interface{}{
				"Name":      "Name",
				"Anonymous": "Anonymous",
			},
			wantErr: false,
		},
		{
			name: "anonymous pointer",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: struct {
					String string
					*AnonymousT
				}{
					String: "Name",
					AnonymousT: &AnonymousT{
						Anonymous: "Anonymous",
					},
				},
			},
			want: map[string]interface{}{
				"String":    "Name",
				"Anonymous": "Anonymous",
			},
			wantErr: false,
		},
		{
			name: "multiple struct",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: struct {
					T
					T1 T
					T2 *T
				}{
					T: T{
						Name: "T",
					},
					T1: T{
						Name: "T1",
					},
					T2: &T{
						Name: "T2",
					},
				},
			},
			want: map[string]interface{}{
				"Name": "T",
				"T1": map[string]interface{}{
					"Name": "T1",
				},
				"T2": map[string]interface{}{
					"Name": "T2",
				},
			},
			wantErr: false,
		},
		{
			name: "interface",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: struct {
					Any interface{}
				}{
					Any: T{
						Name: "Interface",
					},
				},
			},
			want: map[string]interface{}{
				"Any": map[string]interface{}{
					"Name": "Interface",
				},
			},
			wantErr: false,
		},
		{
			name: "string tag option error",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: struct {
					Any interface{} `map:"any,string"`
				}{
					Any: T{
						Name: "Interface",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "bytes tag option error",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: struct {
					Any interface{} `map:"any,bytes"`
				}{
					Any: T{
						Name: "Interface",
					},
				},
			},
			wantErr: true,
		},
		// string tag
		{
			name: "string tag option with Text",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: StringTag{
					Data: Text{
						Name: "text",
					},
				},
			},
			want: map[string]interface{}{
				"data": []byte{116, 101, 120, 116},
			},
			wantErr: false,
		},
		{
			name: "string tag option with Text pointer",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: StringTag{
					Data: &Text{
						Name: "text",
					},
				},
			},
			want: map[string]interface{}{
				"data": []byte{116, 101, 120, 116},
			},
			wantErr: false,
		},
		{
			name: "string tag option with Binary",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: StringTag{
					Data: Binary{
						Name: "binary",
					},
				},
			},
			want: map[string]interface{}{
				"data": []byte{98, 105, 110, 97, 114, 121},
			},
			wantErr: false,
		},
		{
			name: "string tag option with Binary pointer",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: StringTag{
					Data: &Binary{
						Name: "binary",
					},
				},
			},
			want: map[string]interface{}{
				"data": []byte{98, 105, 110, 97, 114, 121},
			},
			wantErr: false,
		},
		{
			name: "string tag option with JSON",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: StringTag{
					Data: JSON{
						Name: "json",
					},
				},
			},
			want: map[string]interface{}{
				"data": []byte(`{"name":"json"}`),
			},
			wantErr: false,
		},
		{
			name: "string tag option with JSON pointer",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: StringTag{
					Data: &JSON{
						Name: "json",
					},
				},
			},
			want: map[string]interface{}{
				"data": []byte(`{"name":"json"}`),
			},
			wantErr: false,
		},
		// bytes tag
		{
			name: "bytes tag option with Text",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: BytesTag{
					Data: Text{
						Name: "text",
					},
				},
			},
			want: map[string]interface{}{
				"data": []byte{116, 101, 120, 116},
			},
			wantErr: false,
		},
		{
			name: "bytes tag option with Text pointer",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: BytesTag{
					Data: &Text{
						Name: "text",
					},
				},
			},
			want: map[string]interface{}{
				"data": []byte{116, 101, 120, 116},
			},
			wantErr: false,
		},
		{
			name: "bytes tag option with Binary",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: BytesTag{
					Data: Binary{
						Name: "binary",
					},
				},
			},
			want: map[string]interface{}{
				"data": []byte{98, 105, 110, 97, 114, 121},
			},
			wantErr: false,
		},
		{
			name: "bytes tag option with Binary pointer",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: BytesTag{
					Data: &Binary{
						Name: "binary",
					},
				},
			},
			want: map[string]interface{}{
				"data": []byte{98, 105, 110, 97, 114, 121},
			},
			wantErr: false,
		},
		{
			name: "bytes tag option with JSON",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: BytesTag{
					Data: JSON{
						Name: "json",
					},
				},
			},
			want: map[string]interface{}{
				"data": []byte(`{"name":"json"}`),
			},
			wantErr: false,
		},
		{
			name: "bytes tag option with JSON pointer",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: BytesTag{
					Data: &JSON{
						Name: "json",
					},
				},
			},
			want: map[string]interface{}{
				"data": []byte(`{"name":"json"}`),
			},
			wantErr: false,
		},
		// json tag
		{
			name: "json tag option with JSON",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: JSONTag{
					Data: JSON{
						Name: "json",
					},
				},
			},
			want: map[string]interface{}{
				"data": []byte(`{"name":"json"}`),
			},
			wantErr: false,
		},
		{
			name: "json tag option with JSON pointer",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: JSONTag{
					Data: &JSON{
						Name: "json",
					},
				},
			},
			want: map[string]interface{}{
				"data": []byte(`{"name":"json"}`),
			},
			wantErr: false,
		},
		{
			name: "pointer",
			fields: fields{
				tag:       Tag,
				separated: Separated,
			},
			args: args{
				data: Pointer{
					Int:    new(int64),
					Float:  new(float64),
					String: new(string),
					Bytes:  &[]byte{1, 2, 3, 4, 5, 6},
				},
			},
			want: map[string]interface{}{
				"Int":    new(int64),
				"Float":  new(float64),
				"String": new(string),
				"Bytes":  &[]byte{1, 2, 3, 4, 5, 6},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := MapEncode{
				tag:       tt.fields.tag,
				separated: tt.fields.separated,
			}
			got, err := e.MarshalMap(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalMap()\n\tgot  %v\n\twant %v", got, tt.want)
			}
		})
	}
}
