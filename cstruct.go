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

	LEN0 = 0
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
	Index uint      // 字段序号
	Name  string    // 字段名
	Type  FieldType // 字段类型
	Size  uint      // 字段长度
}

func (c *CStruct) ToStatement() ([]Field, error) {
	return c.FieldSet, nil
}
