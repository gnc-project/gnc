package hashid

import "C"

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/big"
	"github.com/ethereum/go-ethereum/common/constants"
	"github.com/ethereum/go-ethereum/pocCrypto"
	"github.com/ethereum/go-ethereum/pocCrypto/shabal256"
	"time"
)

const (
	SCOOP_SIZE = 64
	NUM_SCOOPS = 4096
	NONCE_SIZE = NUM_SCOOPS * SCOOP_SIZE
	HASH_SIZE = 32
	HASH_CAP  = 4096
)
const (
	avx2Parallel    = 8
	sse4Parallel    = 4
	blockChainStart = 1407722400
	// GenesisBaseTarget is the base target of the first block
	GenesisBaseTarget = 18325193796
	// BlockChainStart is the timestamp of the first Block
	BlockChainStart = 1407722400
)


func CalculateScoop(genSig []byte,height *big.Int) uint64 {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, height.Uint64())
	hash:=shabal256.Sum256(append(genSig,bs...))
	bg:=big.NewInt(0).SetBytes(hash[:])
	return bg.Mod(bg,constants.SCOOPS_PER_PLOT_BIGINT).Uint64()
}


func CalcuHitAndScoop(pid,nonce *big.Int,genSig []byte,scoopNum uint64) uint64 {
	var gendata []byte
	ba := make([]byte, 8)
	binary.BigEndian.PutUint64(ba, pid.Uint64())
	bn := make([]byte,8)
	binary.BigEndian.PutUint64(bn, nonce.Uint64())
	accnon:=append(ba,bn...)
	for i:=NONCE_SIZE;i>0;i -= HASH_SIZE  {
		if len(gendata) >= HASH_CAP {
			break
		}
		gp:=shabal256.Sum256(append(gendata[:],accnon[:]...))
		gendata=append(gp[:],gendata[:]...)
	}
	for i:=NONCE_SIZE-128*HASH_SIZE; i>0; i -= HASH_SIZE {
		gl:=shabal256.Sum256(append(gendata[:HASH_CAP]))
		gendata=append(gl[:],gendata...)
	}
	final:=shabal256.Sum256(append(gendata,accnon...))

	for i:=0;i<NONCE_SIZE; i++ {
		gendata[i] ^= final[i % 32]
	}

	hashBuffer:=make([]byte,HASH_SIZE)
	revPos:=NONCE_SIZE-HASH_SIZE
	for pos:=32;pos<(NONCE_SIZE/2);pos+=64 {
		copyArr(gendata,pos,hashBuffer,0,HASH_SIZE)
		copyArr(gendata,revPos,gendata,pos,HASH_SIZE)
		copyArr(hashBuffer,0,gendata,revPos,HASH_SIZE)
		revPos -=64
	}
	scoop:=gendata[scoopNum*SCOOP_SIZE:scoopNum*SCOOP_SIZE+SCOOP_SIZE]
	finals2:=shabal256.Sum256(append(genSig,scoop...))
	f:=[]byte{finals2[7],finals2[6],finals2[5],finals2[4],finals2[3],finals2[2],finals2[1],finals2[0],}
	fin:=big.NewInt(0).SetBytes(f)
	return fin.Uint64()
}


func CalcuHit(pid,nonce *big.Int,genSig []byte) uint64 {
	var gendata []byte
	ba := make([]byte, 8)
	binary.BigEndian.PutUint64(ba, pid.Uint64())
	bn := make([]byte,8)
	binary.BigEndian.PutUint64(bn, nonce.Uint64())
	accnon:=append(ba,bn...)
	for i:=NONCE_SIZE;i>0;i -= HASH_SIZE  {
		if len(gendata) >= HASH_CAP {
			break
		}
		gp:=shabal256.Sum256(append(gendata[:],accnon[:]...))
		gendata=append(gp[:],gendata[:]...)
	}
	for i:=NONCE_SIZE-128*HASH_SIZE; i>0; i -= HASH_SIZE {
		gl:=shabal256.Sum256(append(gendata[:HASH_CAP]))
		gendata=append(gl[:],gendata...)
	}
	final:=shabal256.Sum256(append(gendata,accnon...))

	for i:=0;i<NONCE_SIZE; i++ {
		gendata[i] ^= final[i % 32]
	}

	hashBuffer:=make([]byte,HASH_SIZE)
	revPos:=NONCE_SIZE-HASH_SIZE
	for pos:=32;pos<(NONCE_SIZE/2);pos+=64 {
		copyArr(gendata,pos,hashBuffer,0,HASH_SIZE)
		copyArr(gendata,revPos,gendata,pos,HASH_SIZE)
		copyArr(hashBuffer,0,gendata,revPos,HASH_SIZE)
		revPos -=64
	}
	scn:=CalculateScoop(genSig,big.NewInt(500000))
	scoop:=gendata[scn*SCOOP_SIZE:scn*SCOOP_SIZE+SCOOP_SIZE]
	finals2:=shabal256.Sum256(append(genSig,scoop...))
	f:=[]byte{finals2[7],finals2[6],finals2[5],finals2[4],finals2[3],finals2[2],finals2[1],finals2[0],}
	fin:=big.NewInt(0).SetBytes(f)
	return fin.Uint64()
	//t:=0
	//for i:=0; i<len(gendata);i+=HASH_SIZE {
	//	g:=gendata[i:i+HASH_SIZE]
	//	ghx:=hex.EncodeToString(g)
	//	fmt.Println(ghx)
	//	t++
	//}
	//fmt.Println("t---->",t)
}


func CalcuDeadline(pid *big.Int,nonce *big.Int,baseTarget *big.Int,scn uint64,gensig []byte) int64 {
	f:=CalcuHitAndScoop(pid,nonce,gensig,scn)
	fin:=big.NewInt(0).SetUint64(f)
	fb:=big.NewInt(0).Div(fin,baseTarget)
	return fb.Int64()
}


// DateToTimeStamp yields a timestamp counted since block chain start
func DateToTimestamp(date time.Time) uint32 {
	ts := date.Unix() - blockChainStart
	if ts < 0 {
		return 0
	}
	return uint32(ts)
}


func DecodeGeneratorSignature(genSigStr string) ([]byte, error) {
	if len(genSigStr) != 64 {
		return nil, errors.New("Generation signature's length differs from 64")
	}
	genSig, err := hex.DecodeString(genSigStr)
	if err != nil {
		return nil, err
	}
	return genSig, nil
}



func CalculateGenSinGen(lastSig []byte,Generator *big.Int) []byte {
	bs:=make([]byte,8)
	binary.BigEndian.PutUint64(bs, Generator.Uint64())
	hash:=shabal256.Sum256(append(lastSig,bs...))
	return hash[:]
}


func CalculateGenerationSignature(lastSig []byte,pub []byte) []byte {
	_, id := crypto.BytesToHashAndID(pub)
	return CalculateGenSinGen(lastSig,id)
}


func CalculateGenerFSP(secretPhrase string)( PID *big.Int){
	pub:=crypto.SecretPhraseToPublicKey(secretPhrase)
	_,id:=crypto.BytesToHashAndID(pub)
	return id
}


func copyArr(src []byte,srcOffset int, dst []byte,dstOffset,count int)(bool,error){
	srcLen := len(src)
	if srcOffset >srcLen || count > srcLen || srcOffset+count > srcLen{
		return false,errors.New("index out of array")
	}
	dstLen := len(dst)
	if dstOffset >dstLen || count > dstLen || dstOffset+count > dstLen {
		return false, errors.New("index out of array")
	}
	index:=0
	for i:=srcOffset; i<srcOffset+count; i++{
		dst[dstOffset+index] = src[srcOffset+index]
		index++
	}
	return true,nil
}












//
//func CalHit(aid,nonce,height *big.Int,genSig []byte)  uint64 {
//	scoop:=CalculateScoop(genSig,height)
//	req:=NewCalcDeadlineRequest(aid.Uint64(),
//		nonce.Uint64(),0,uint32(scoop),height.Uint64(),genSig).native
//	calculateHit(req,genSig)
//	return uint64(req.hit)
//}
//
//func calculateHit(req *C.CalcDeadlineRequest,gen []byte)  {
//	g:=(*C.uint8_t)(unsafe.Pointer(&gen[0]))
//	C.calculate_hit(req,g)
//}
//
//// BurstToPlanck converts an amount in burst to an amount in burst
//func BurstToPlanck(n float64) int64 {
//	return int64(n * 100000000)
//}
//
//// PlanckToBurst converts an amount in placnk to an amount in burst
//func PlanckToBurst(n int64) float64 {
//	return float64(n) / 100000000.0
//}



//func CalcuDeadline(pid *big.Int,nonce *big.Int,scn int,
//	baseTarget *big.Int,gensig []byte) uint64 {
//	var gendata []byte
//	ba := make([]byte, 8)
//	binary.BigEndian.PutUint64(ba, pid.Uint64())
//	bn := make([]byte,8)
//	binary.BigEndian.PutUint64(bn, nonce.Uint64())
//	accnon:=append(ba,bn...)
//	for i:=NONCE_SIZE;i>0;i -= HASH_SIZE  {
//		if len(gendata) >= HASH_CAP {
//			break
//		}
//		gp:=shabal256.Sum256(append(gendata[:],accnon[:]...))
//		gendata=append(gp[:],gendata[:]...)
//	}
//	for i:=NONCE_SIZE-128*HASH_SIZE; i>0; i -= HASH_SIZE {
//		gl:=shabal256.Sum256(append(gendata[:HASH_CAP]))
//		gendata=append(gl[:],gendata...)
//	}
//
//	final:=shabal256.Sum256(append(gendata,accnon...))
//
//	for i:=0;i<NONCE_SIZE; i++ {
//		gendata[i] ^= final[i % 32]
//	}
//
//	hashBuffer:=make([]byte,HASH_SIZE)
//	revPos:=NONCE_SIZE-HASH_SIZE
//	for pos:=32;pos<(NONCE_SIZE/2);pos+=64 {
//		copyArr(gendata,pos,hashBuffer,0,HASH_SIZE)
//		copyArr(gendata,revPos,gendata,pos,HASH_SIZE)
//		copyArr(hashBuffer,0,gendata,revPos,HASH_SIZE)
//		revPos -=64
//	}
//
//	//t:=0
//	//for i:=0; i<len(gendata);i+=HASH_SIZE {
//	//	g:=gendata[i:i+HASH_SIZE]
//	//	ghx:=hex.EncodeToString(g)
//	//	fmt.Println(ghx)
//	//	t++
//	//}
//	//fmt.Println("t---->",t)
//	scoop:=gendata[scn*SCOOP_SIZE:scn*SCOOP_SIZE+SCOOP_SIZE]
//	finals2:=shabal256.Sum256(append(gensig,scoop...))
//	f:=[]byte{finals2[7],finals2[6],finals2[5],finals2[4],finals2[3],finals2[2],finals2[1],finals2[0],}
//	fin:=big.NewInt(0).SetBytes(f)
//	fb:=big.NewInt(0).Div(fin,baseTarget)
//	return fb.Uint64()
//}

//func CalcuDeadline(accountID,nonce,baseTarget,height *big.Int,
//	scoop int64,genSig []byte) *big.Int {
//	req:=NewCalcDeadlineRequest(accountID.Uint64(),
//		nonce.Uint64(),
//		baseTarget.Uint64(),
//		uint32(scoop),
//		height.Uint64(),
//		genSig).native
//	calculateDeadline(req,genSig)
//	return big.NewInt(0).SetUint64(uint64(req.deadline))
//}

// CalcScoop calculated the scoop for a given height and generation signature
//func CalcScoop(height int32, genSig []byte) uint32 {
//	return uint32(C.calculate_scoop(C.uint64_t(height), (*C.uint8_t)(&genSig[0])))
//}

// NewCalcDeadlineRequest allocated paramters neeeded for deadline
// calculation native so that C can deal with it
//func NewCalcDeadlineRequest(accountID, nonce, baseTarget uint64, scoop uint32, height uint64, genSig []byte) *CalcDeadlineRequest {
//	var deadline C.uint64_t
//	return &CalcDeadlineRequest{
//		native: &C.CalcDeadlineRequest{
//			account_id:  C.uint64_t(accountID),
//			nonce:       C.uint64_t(nonce),
//			base_target: C.uint64_t(baseTarget),
//			scoop_nr:    C.uint32_t(scoop),
//			//gen_sig:     (*C.uint8_t)(unsafe.Pointer(&genSig[0])),
//			deadline:    deadline,
//			with_scoop: false,
//			height:      C.uint64_t(height)},
//		//done: make(chan struct{})
//	}
//}


// CalculateDeadline calculates a single deadline
//  TODO: might fail on some go versions(cgo argument has Go pointer to Go pointer): GODEBUG=cgocheck=0
//func calculateDeadline(req *C.CalcDeadlineRequest,gen []byte) {
//	g:=(*C.uint8_t)(unsafe.Pointer(&gen[0]))
//	C.calculate_deadline(req,g)
//}

//// #cgo LDFLAGS: -L../c -lmcmath
///*
//#include "../c/mcmath.h"
//#include "stdlib.h"
//
//CalcDeadlineRequest** alloc_reqs_avx2() {
//  return (CalcDeadlineRequest**) malloc(8 * sizeof(CalcDeadlineRequest*));
//}
//*/
// CalcDeadlineRequest stores paramters native that are
// needed for deadline calculation
//type CalcDeadlineRequest struct {
//	native    *C.CalcDeadlineRequest
//	withScoop bool
//	//done      chan struct{}
//}