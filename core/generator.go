package core

import (
	"github.com/orcaman/concurrent-map"
	"math/big"
	"github.com/ethereum/go-ethereum/pocCrypto/hashid"
	// "time"
	"github.com/ethereum/go-ethereum/services"
	// "github.com/ethereum/go-ethereum/consensus"
	// "encoding/hex"
	"github.com/ethereum/go-ethereum/core/types"

)

const (
	NOCE_SUBMITTED = iota
	GENERATION_DEADLINE
)

type Generator interface {

	AddNonce(pid *big.Int,nonce *big.Int)GeneratorState


	GetAllGenerators()cmap.ConcurrentMap


	GenerateForBlockchainProcessor()
}


type GeneratorImpl struct {
	Generators cmap.ConcurrentMap
	TimeService services.TimeService
}
type GeneratorState interface {
	GetAccountID() *big.Int
	GetDeadline() *big.Int
	GetBlockHeight() *big.Int
	GetGenerationSignature() []byte
	GetNonce() *big.Int
}

func NewGenerator() *GeneratorImpl {
    timeServiceImpl:=services.NewTimeServiceImpl()
	return &GeneratorImpl{Generators:cmap.New(),TimeService:&timeServiceImpl}
}


func (ge *GeneratorImpl)AddNonce(pid *big.Int,nonce *big.Int,number *big.Int,header *types.Header,lastBlock *types.Block)GeneratorState {
	gsinew:=ge.NewGeneratorStateImpl(pid,nonce,number,header,lastBlock)

	// if g,b:=ge.Generators.Get(pid.String());b{
	// 	gsi=g.(GeneratorStateImpl)
	// }

	// if gsinew==nil || gsinew.blockHeight.Cmp(gsi.blockHeight)>0 || gsinew.GetDeadline().Cmp(gsi.GetDeadline())<0{
		ge.Generators.Set(pid.String(),gsinew)
	// }
	return gsinew
}


func (ge *GeneratorImpl)GetAllGenerators()cmap.ConcurrentMap {
	return ge.Generators
}


type GeneratorStateImpl struct {

	accountID *big.Int `json:"accountID"`

	deadline	*big.Int  `json:"deadline"`
	nonce		*big.Int  `json:"nonce"`
	blockHeight	*big.Int   `json:"blockHeight"`

	generationSignature []byte  `json:"generationSignature"`
}

func (g *GeneratorStateImpl) GetNonce() *big.Int {
	return g.nonce
}

func (g *GeneratorStateImpl) GetAccountID() *big.Int {
	return g.accountID
}

func (g *GeneratorStateImpl) GetDeadline() *big.Int {
	return g.deadline
}

func (g *GeneratorStateImpl) GetBlockHeight() *big.Int {
	return g.blockHeight
}
func (g *GeneratorStateImpl) GetGenerationSignature() []byte {
	return g.generationSignature
}

func (ge *GeneratorImpl)NewGeneratorStateImpl(accountid *big.Int,nonce *big.Int,number *big.Int,header *types.Header,lastBlock *types.Block) *GeneratorStateImpl {

	gsi:=&GeneratorStateImpl{
		nonce:nonce,
		accountID:accountid,
		blockHeight:number,
		generationSignature:header.GenerationSignature,
	}
	scoopnum:=hashid.CalculateScoop(header.GenerationSignature,gsi.blockHeight)
	deadline:=hashid.CalcuDeadline(accountid,nonce,lastBlock.Difficulty(),scoopnum,gsi.generationSignature)

	gsi.deadline=big.NewInt(deadline)

	return gsi
}