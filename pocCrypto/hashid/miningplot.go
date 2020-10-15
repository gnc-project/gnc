package hashid
//
//import (
//	"encoding/binary"
//	"math/big"
//)
//
//const (
//	HASH_SIZE int64 = 32
//	HASHES_PER_SCOOP int64 = 2
//	SCOOP_SIZE int64 = HASHES_PER_SCOOP*HASH_SIZE//64
//	SCOOPS_PER_PLOT int64 = 4096
//	PLOT_SIZE int64 = SCOOPS_PER_PLOT * SCOOP_SIZE//262144
//	BASE_LENGTH int64 = 16
//	HASH_CAP int64 = 4096
//)
//
//type MiningPlot struct {
//	HASH_SIZE int64
//	HASHES_PER_SCOOP int64
//	SCOOP_SIZE int64
//	SCOOPS_PER_PLOT int64
//	SCOOPS_PER_PLOT_BIGINT big.Int
//	PLOT_SIZE int64
//	BASE_LENGTH int64
//
//	HASH_CAP int64
//	Data []byte
//	DataOffset int64
//}
//
//func NewDefaultMininPlot(addr big.Int,  nonce big.Int,  pocVersion int64,
//	 buffer []byte,  bufferOffset int64) *MiningPlot {
//	mp:=&MiningPlot{
//		HASH_SIZE:HASH_SIZE,
//		HASHES_PER_SCOOP:HASHES_PER_SCOOP,
//		SCOOP_SIZE:SCOOP_SIZE,
//		SCOOPS_PER_PLOT:SCOOPS_PER_PLOT,
//		SCOOPS_PER_PLOT_BIGINT:*big.NewInt(SCOOPS_PER_PLOT),
//		PLOT_SIZE:PLOT_SIZE,
//		BASE_LENGTH:BASE_LENGTH,
//		HASH_CAP:HASH_CAP,
//	}
//	mp.Data = buffer
//	mp.DataOffset= bufferOffset
//
//	baseData := make([]byte, BASE_LENGTH)
//	binary.BigEndian.PutUint64(baseData,addr.Uint64())
//	binary.BigEndian.PutUint64(baseData,nonce.Uint64())
//
//	limit:=PLOT_SIZE - 128 * HASH_SIZE
//	//buf:=make([]byte,64)
//	for i:=PLOT_SIZE; i> limit; i -= HASH_SIZE {
//		//blockCopy(mp.Data,buf,i+bufferOffset,PLOT_SIZE - i);
//	}
//	return mp
//}