package cstruct

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestToUnderLine(t *testing.T) {
	res := ToUnderLine("EndTime_StartTime_Itow")
	assert.Equal(t, res, "end_time_start_time_itow")
}
