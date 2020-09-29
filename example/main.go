package main

import (
	"bytes"
	"fmt"
	"github.com/geange/cstruct-go"
)

func main() {
	buf := bytes.NewBufferString(nestedText)
	lex, err := cstruct.NewLexer(buf)
	if err != nil {
		panic(err)
	}

	list, err := lex.AllStructure()
	if err != nil {
		panic(err)
	}

	for _, item := range list {
		fmt.Println(item.Name)
		for _, v := range item.FieldSet {
			fmt.Println(v.Name, v.Type, v.Size)
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
