package cstruct

import "fmt"

type CStructureV1 interface {
	Flat() ([]CStructureV1, error)
	FieldName() string
	SetFieldName(name string)
	Type() FieldType
	Index() uint
	SetIndex(index uint)
}

type BaseStruct struct {
	index uint
	name  string
	fType FieldType
	size  uint
}

func (b *BaseStruct) Flat() ([]CStructureV1, error) {
	field := BaseStruct{
		index: b.index,
		name:  b.name,
		fType: b.fType,
		size:  b.size,
	}
	return []CStructureV1{&field}, nil
}

func (b *BaseStruct) FieldName() string {
	return b.name
}

func (b *BaseStruct) SetFieldName(name string) {
	b.name = name
}

func (b *BaseStruct) Type() FieldType {
	return b.fType
}

func (b *BaseStruct) Index() uint {
	return b.index
}

func (b *BaseStruct) SetIndex(index uint) {
	b.index = index
}

type U8 struct {
	BaseStruct
}

type S8 struct {
	BaseStruct
}

type U16 struct {
	BaseStruct
}

type S16 struct {
	BaseStruct
}

type U32 struct {
	BaseStruct
}

type S32 struct {
	BaseStruct
}

type U64 struct {
	BaseStruct
}

type S64 struct {
	BaseStruct
}

type F32 struct {
	BaseStruct
}

type F64 struct {
	BaseStruct
}

type ByteBS struct {
	BaseStruct
}

type ByteHex struct {
	BaseStruct
}

type ByteString struct {
	BaseStruct
}

type ArrayV1 struct {
	index uint
	name  string
	size  uint
	field CStructureV1
}

func (a *ArrayV1) Flat() ([]CStructureV1, error) {
	r := make([]CStructureV1, 0)
	index := uint(0)
	for i := 0; i < int(a.size); i++ {
		fields, err := a.field.Flat()
		if err != nil {
			return nil, err
		}

		for _, f := range fields {
			switch a.field.Type() {
			case Struct:
				f.SetFieldName(fmt.Sprintf("%s_%d_%s", a.name, i, f.FieldName()))
			default:
				f.SetFieldName(fmt.Sprintf("%s_%d", a.name, i))
			}

			f.SetIndex(index)
			index++
			r = append(r, f)
		}
	}
	return r, nil
}

func (a *ArrayV1) FieldName() string {
	return a.name
}

func (a *ArrayV1) SetFieldName(name string) {
	a.name = name
}

func (a *ArrayV1) Type() FieldType {
	return Array
}

func (a *ArrayV1) Index() uint {
	return a.index
}

func (a *ArrayV1) SetIndex(index uint) {
	a.index = index
}

type CStructV1 struct {
	index  uint
	gName  string
	fName  string
	fields []CStructureV1
}

func (c *CStructV1) Copy() *CStructV1 {
	return &CStructV1{
		index:  c.index,
		gName:  c.gName,
		fName:  c.fName,
		fields: c.fields,
	}
}

func (c *CStructV1) SetFieldName(name string) {
	c.fName = name
}

func (c *CStructV1) Index() uint {
	return c.index
}

func (c *CStructV1) SetIndex(index uint) {
	c.index = index
}

func (c *CStructV1) Flat() ([]CStructureV1, error) {
	r := make([]CStructureV1, 0, len(c.fields))
	index := uint(0)
	for _, f := range c.fields {
		fields, err := FlatDeep(f)
		if err != nil {
			return nil, err
		}

		parentName := f.FieldName()

		//fmt.Println("----", parentName)

		for _, field := range fields {
			field.SetIndex(index)

			name := fmt.Sprintf("%s_%s", parentName, field.FieldName())
			field.SetFieldName(name)
			index++
			r = append(r, field)
		}
	}
	return r, nil
}

func (c *CStructV1) FieldName() string {
	return c.fName
}

func (c *CStructV1) Type() FieldType {
	return Struct
}

func FlatDeep(item CStructureV1) ([]CStructureV1, error) {
	r := make([]CStructureV1, 0)
	fields, err := item.Flat()
	if err != nil {
		return nil, err
	}
	for _, v := range fields {
		switch v.Type() {
		case Struct, Array:
			tFields, err := FlatDeep(v)
			if err != nil {
				return nil, err
			}
			r = append(r, tFields...)
		default:
			r = append(r, v)
		}
	}
	return r, nil
}
