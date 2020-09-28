package cstruct

type FieldType string

const (
	// C语言中的数据类型
	UInt8   = FieldType("uint8")
	UInt16  = FieldType("uint16")
	UInt32  = FieldType("uint32")
	UInt64  = FieldType("uint64")
	Int8    = FieldType("int8")
	Int16   = FieldType("int16")
	Int32   = FieldType("int32")
	Int64   = FieldType("int64")
	Float32 = FieldType("float32")
	Float64 = FieldType("float64")
	// string和bytes的区别是，bytes解析出来后是base64编码，string是asc/utf-8编码
	Bytes  = FieldType("bytes")
	String = FieldType("string")
	Hex    = FieldType("hex")
	// 内部结构体
	Struct = "struct"

	LEN1 = 1
	LEN2 = 2
	LEN4 = 4
	LEN8 = 8
)

type CStruct struct {
	Name     string
	FieldSet FieldSet `json:"field_set" bson:"field_set"`
}

type FieldSet []Field

func (fs FieldSet) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fs FieldSet) Len() int {
	return len(fs)
}

func (fs FieldSet) Less(i, j int) bool {
	return fs[i].Index < fs[j].Index
}

type Field struct {
	Index       uint       `json:"index" bson:"index"`                             // 字段序号
	Name        string     `json:"name" bson:"name"`                               // 字段名
	DisplayName string     `json:"display_name" bson:"display_name"`               // 显示名称
	Display     bool       `json:"display" bson:"display"`                         // 是否显示
	Type        FieldType  `json:"type" bson:"type"`                               // 字段类型
	ByteOrder   string     `json:"byte_order" bson:"byte_order"`                   // 字节序
	Size        *uint      `json:"size,omitempty" bson:"size,omitempty"`           // 字段的大小
	ElemSize    *FieldType `json:"elem_size,omitempty" bson:"elem_size,omitempty"` // elem的大小
	Elem        []Field    `json:"elem,omitempty" bson:"elem,omitempty"`           // elem的内容
}

func (c *CStruct) ToStatement() ([]Field, error) {
	return c.FieldSet, nil
}
