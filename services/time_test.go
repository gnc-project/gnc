package services

import (
	"fmt"
	"testing"
)

func TestEpochTime_GetTime(t *testing.T) {
	d:=EpochTime{}
	//1407722400000
	fmt.Println(d.GetTimeInMillis())
}
