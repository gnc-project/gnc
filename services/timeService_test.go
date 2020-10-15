package services

import (
	"fmt"
	"testing"
)

func TestEpochTime_GetTime(t *testing.T) {
	nt:=NewTimeServiceImpl()
	fmt.Println(nt.GetEpochTime())
}
