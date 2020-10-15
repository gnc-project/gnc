package services

import (
	"fmt"
	"math/big"
	"testing"
)

func TestCalculateDifficulty(t *testing.T) {
	i:=big.NewInt(0)
	f:=i.Add(i,big.NewInt(2))
	fmt.Println(f)
	fmt.Println(i)
}
