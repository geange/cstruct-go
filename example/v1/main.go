package main

import (
	"bytes"
	"fmt"
	"github.com/geange/cstruct-go"
)

func main() {
	buf := bytes.NewBufferString(nestedText)
	lex, err := cstruct.NewLexerV1(buf)
	if err != nil {
		panic(err)
	}

	items, err := lex.AllStructureV1()
	if err != nil {
		panic(err)
	}

	for _, item := range items {
		fmt.Println(item.FieldName)

		fields, err := item.Flat()
		if err != nil {
			panic(err)
		}
		for _, v := range fields {
			fmt.Println(v.FieldName(), v.Type())
		}
	}
}

var (
	nestedText = `typedef struct{
    u32 ITOW;
    u16 Week;
    u16 VITOW;
}UTC_t;

typedef struct
{
    byte CustomData[16];
    UTC_t StartTime;
    UTC_t EndTime;
    u16 StartVoltage;
    u16 EndVoltage;
    u8  StartPoint;
    u8  EndPoint;
    u8  CurrentPoint;
    u8  Status;
} MSG_RouteData_t;`
)
