package cstruct

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestScanner_Scan(t *testing.T) {
	scanner := NewScannerString(text)
	tokens, err := scanner.Scan()
	assert.Equal(t, err, nil)
	for _, token := range tokens {
		fmt.Println(TokenTypeMap[token.Type()], token.Value())
	}
}

var text = `typedef struct
{
	s32 Longitude;
	s32 Latitude;
	s16 Height;
	u8  Speed;
	u8  Type;
	u8  HeadType;
	u8  HeightType;
	u16 Reserved2;
    u8  Param0;
    u8  Param1;
    u8  Param2;
    u8  Param3;
    u32 Reserved3;
} MSG_Waypoint_t;`

var nestedText = `typedef struct{
    u32 ITOW;
    u16 Week;
    u16 VITOW;
}UTC_t;

typedef struct
{
    u8  CustomData[16];
    UTC_t StartTime;
    UTC_t EndTime;
    u16 StartVoltage;
    u16 EndVoltage;
    u8  StartPoint;
    u8  EndPoint;
    u8  CurrentPoint;
    u8  Status;
} MSG_RouteData_t;`
