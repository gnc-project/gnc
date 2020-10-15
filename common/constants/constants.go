package constants

import (
	"math/big"
	"time"
)

var GenerationSignature [32]byte

var EPOCH_BEGINNING int64
var TWO64 *big.Int
var SCOOPS_PER_PLOT_BIGINT = big.NewInt(4096)
func init(){

	 EPOCH_BEGINNING=time.Date(2014,time.August,11,2,0,0,0,time.UTC).UnixNano() / 1e6
	 TWO64 = big.NewInt(0).Exp(big.NewInt(2),big.NewInt(64),nil)
}

func pow(x, n uint64) uint64 {
	var ret uint64 = 1 
	for n != 0 {
		if n%2 != 0 {
			ret = ret * x
		}
		n /= 2
		x = x * x
	}
	return ret
}

const (


	RECIPIENT_ASSIGNMENT_WAIT_TIME = 4 
	GENESIS_BLOCK_HASH = "asdfasdasf" 
	MINING_RATIO = 3 


	//// 4398046511104 / 240 = 18325193796
	//const uint64_t BHD_BASE_TARGET_240 = 18325193796ull;
	//
	//// 4398046511104 / 300 = 14660155037
	//const uint64_t BHD_BASE_TARGET_300 = 14660155037ull;
	//
	//// 4398046511104 / 180 = 24433591728
	//const uint64_t BHD_BASE_TARGET_180 = 24433591728ull;

	// 4398046511104 / 180 = 24433591728  
	// coin = INITIAL_BASE_TARGET / block.Difficulty * MINING_RATIO
	//T=(BASE_TARGET_180 or INITIAL_BASE_TARGET) / block.Difficulty / 1024
	INITIAL_BASE_TARGET int64  = 48867183456
	MAX_BASE_TARGET int64 = 48867183456
)



