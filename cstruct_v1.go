package cstruct

import (
	"errors"
	"fmt"
)

type StatementV1 interface {
	Type() FieldType
	Name() string
	SetFieldName(name string)
	Copy() StatementV1
}

type BaseV1 struct {
	index int
	name  string
	fType FieldType
	size  int
}

func (b *BaseV1) Type() FieldType {
	return b.fType
}

func (b *BaseV1) Copy() StatementV1 {
	return &BaseV1{
		index: b.index,
		name:  b.name,
		fType: b.fType,
		size:  b.size,
	}
}

func (b *BaseV1) Flat() ([]StatementV1, error) {
	return []StatementV1{
		&BaseV1{
			index: b.index,
			name:  b.name,
			fType: b.fType,
			size:  b.size,
		},
	}, nil
}

func (b *BaseV1) Name() string {
	return b.name
}

func (b *BaseV1) SetFieldName(name string) {
	b.name = name
}

type U8 struct {
	BaseV1
}

type U16 struct {
	BaseV1
}

type U32 struct {
	BaseV1
}

type U64 struct {
	BaseV1
}

type S8 struct {
	BaseV1
}

type S16 struct {
	BaseV1
}

type S32 struct {
	BaseV1
}

type S64 struct {
	BaseV1
}

type F32 struct {
	BaseV1
}

type F64 struct {
	BaseV1
}

type Byte struct {
	BaseV1
}

type BaseArrayV1 struct {
	index int
	name  string
	fType StatementV1
	size  int
}

func (b *BaseArrayV1) Type() FieldType {
	return Array
}

func (b *BaseArrayV1) Name() string {
	return b.name
}

func (b *BaseArrayV1) SetFieldName(name string) {
	b.name = name
}

func (b *BaseArrayV1) Copy() StatementV1 {
	return &BaseArrayV1{
		index: b.index,
		name:  b.name,
		fType: b.fType.Copy(),
		size:  b.size,
	}
}

func (b *BaseArrayV1) Flat() ([]StatementV1, error) {
	return flatDeep(b, "")
}

type BaseStructV1 struct {
	index  int
	sName  string
	fName  string
	fields []StatementV1
}

func (b *BaseStructV1) ClassName() string {
	return b.sName
}

func (b *BaseStructV1) Type() FieldType {
	return Struct
}

func (b *BaseStructV1) Name() string {
	return b.fName
}

func (b *BaseStructV1) SetFieldName(name string) {
	b.fName = name
}

func (b *BaseStructV1) Copy() StatementV1 {
	r := make([]StatementV1, 0, len(b.fields))
	for _, v := range b.fields {
		r = append(r, v.Copy())
	}

	return &BaseStructV1{
		index:  b.index,
		sName:  b.sName,
		fName:  b.fName,
		fields: r,
	}
}

func (b *BaseStructV1) Flat() ([]StatementV1, error) {
	return flatDeep(b, "")
}

func flatDeep(item StatementV1, parentName string) ([]StatementV1, error) {
	r := make([]StatementV1, 0)
	switch item.(type) {
	case *BaseStructV1:
		obj := item.(*BaseStructV1)
		for _, v := range obj.fields {
			// 如果是普通类型的数据直接返回
			if !(v.Type() == Struct || v.Type() == Array) {
				if parentName == "" {
					r = append(r, v)
				} else {
					tObj := v.Copy()
					fName := fmt.Sprintf("%s_%s", parentName, tObj.Name())
					tObj.SetFieldName(fName)
					r = append(r, tObj)
				}
				continue
			}

			// 如果不是继续深度递归下去
			fs, err := flatDeep(v, v.Name())
			if err != nil {
				return nil, err
			}
			for _, f := range fs {
				r = append(r, f)
			}
		}
	case *BaseArrayV1:
		obj := item.(*BaseArrayV1)
		switch obj.fType.(type) {
		case *BaseStructV1:
			fields, err := flatDeep(obj.fType, obj.name)
			if err != nil {
				return nil, err
			}
			for i := 0; i < obj.size; i++ {
				for _, v := range fields {
					tObj := v.Copy()
					name := fmt.Sprintf("%s_%d", v.Name(), i)
					tObj.SetFieldName(name)
					r = append(r, tObj)
				}
			}
		case *BaseArrayV1:
			return nil, errors.New("unsupported mul array")
		default:
			for i := 0; i < obj.size; i++ {
				v := obj.fType.Copy()
				name := fmt.Sprintf("%s_%d", parentName, i)
				v.SetFieldName(name)
				r = append(r, v)
			}
		}
	default:
		return nil, fmt.Errorf("base type not allow deep scan")
	}
	return r, nil
}
