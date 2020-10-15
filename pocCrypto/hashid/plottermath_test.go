package hashid

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"plotterid/crypto"
	"testing"
	"time"
)

func TestCalHit(t *testing.T) {
	genSig, _ := DecodeGeneratorSignature("6ec823b5fd86c4aee9f7c3453cacaf4a43296f48ede77e70060ca8225c2855d0")

	hit:=CalcuHit(big.NewInt(7009665667967103287),
		big.NewInt(0),
		genSig)
	if hit!=13406969278176040930 {
		t.Errorf("hit err")
	}
	fmt.Println("dfassssss ",hit)//13406969278176040930
}

func TestCalculateDeadline(t *testing.T)  {
	genSig, _ := DecodeGeneratorSignature("6ec823b5fd86c4aee9f7c3453cacaf4a43296f48ede77e70060ca8225c2855d0")
	scop:=CalculateScoop(genSig,big.NewInt(500000))
	fmt.Println(uint32(scop))
	u:=CalcuDeadline(big.NewInt(7009665667967103287),
		big.NewInt(0),
		scop,
		big.NewInt(70312),
		genSig)
	fmt.Println(u)
	if u!=190678252334964 {
		t.Errorf("err")
	}
}



func TestCalculateGenerationSignature(t *testing.T)  {
	sig,e:=hex.DecodeString("305a98571a8b96f699449dd71eff051fc10a3475bce18c7dac81b3d9316a9780")
	fmt.Println(e)
	p,es:=hex.DecodeString("a44e4299354f59919329a0bfbac7d6858873ef06c8db3a6a90158f581478bd38")
	fmt.Println(es)
	_,g:=crypto.BytesToHashAndID(p)
	fmt.Println(g)
	s:=CalculateGenSinGen(sig,big.NewInt(0).SetUint64(10446462338210047360))
	fmt.Println(hex.EncodeToString(s))
}


func TestDateToTimestamp(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	assert.Equal(t, uint32(0), DateToTimestamp(time.Date(1995, time.August, 2, 2, 2, 0, 0, loc)))
	assert.Equal(t, uint32(62380920), DateToTimestamp(time.Date(2016, time.August, 2, 2, 2, 0, 0, loc)))
}