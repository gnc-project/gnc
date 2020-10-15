// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethash

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	// "runtime"
	"time"
	// mapset "github.com/deckarep/golang-set"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/common/constants"
	"github.com/ethereum/go-ethereum/pocCrypto/hashid"
)

// Ethash proof-of-work protocol constants.
var (
	FrontierBlockReward       = big.NewInt(5e+18) // Block reward in wei for successfully mining a block
	ByzantiumBlockReward      = big.NewInt(3e+18) // Block reward in wei for successfully mining a block upward from Byzantium
	ConstantinopleBlockReward = big.NewInt(2e+18) // Block reward in wei for successfully mining a block upward from Constantinople
	allowedFutureBlockTime    = 15 * time.Second  // Max time from current time allowed for blocks, before they're considered future blocks

	// calcDifficultyConstantinople is the difficulty adjustment algorithm for Constantinople.
	// It returns the difficulty that a new block should have when created at time given the
	// parent block's time and difficulty. The calculation uses the Byzantium rules, but with
	// bomb offset 5M.
	// Specification EIP-1234: https://eips.ethereum.org/EIPS/eip-1234
	calcDifficultyConstantinople = makeDifficultyCalculator(big.NewInt(5000000))

	// calcDifficultyByzantium is the difficulty adjustment algorithm. It returns
	// the difficulty that a new block should have when created at time given the
	// parent block's time and difficulty. The calculation uses the Byzantium rules.
	// Specification EIP-649: https://eips.ethereum.org/EIPS/eip-649
	calcDifficultyByzantium = makeDifficultyCalculator(big.NewInt(3000000))

)

var (
	POCreward uint64 =8000


	blcokToatal uint64 = 365*24*60 / 4


	subsidyHalvingInterval uint64= blcokToatal*2;
)
var (
	InitDifficulty =big.NewInt(48867183456)
	mortgageToken=big.NewInt(3e+18)
)

func GetReward(height uint64) uint64 {

	halvings := height/subsidyHalvingInterval
	subsidy :=POCreward
	subsidy >>= halvings

	return subsidy
}



// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	errZeroBlockTime     = errors.New("timestamp equals parent's")
	errInvalidDifficulty = errors.New("non-positive difficulty")
	errInvalidMixDigest  = errors.New("invalid mix digest")
	errInvalidPoW        = errors.New("invalid proof-of-work")
)

// Author implements consensus.Engine, returning the header's coinbase as the
// proof-of-work verified author of the block.
func (ethash *Ethash) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of the
// stock Ethereum ethash engine.
func (ethash *Ethash) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	// If we're running a full engine faking, accept any input as valid
	if ethash.config.PowMode == ModeFullFake {
		return nil
	}
	// Short circuit if the header is known, or its parent not
	number := header.Number.Uint64()
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// Sanity checks passed, do a proper verification
	return ethash.verifyHeader(chain, header, parent,seal)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications.
func (ethash *Ethash) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	// If we're running a full engine faking, accept any input as valid
	if ethash.config.PowMode == ModeFullFake || len(headers) == 0 {
		abort, results := make(chan struct{}), make(chan error, len(headers))
		for i := 0; i < len(headers); i++ {
			results <- nil
		}
		return abort, results
	}

	// Spawn as many workers as allowed threads
	// workers := runtime.GOMAXPROCS(0)
	// if len(headers) < workers {
	// 	workers = len(headers)
	// }
	// Create a task channel and spawn the verifiers
	var (
		inputs = make(chan int)
		done   = make(chan int)
		errors = make([]error, len(headers))
		abort  = make(chan struct{})
	)
	// for i := 0; i < workers; i++ {
		go func() {
			for index := range inputs {
				errors[index] = ethash.verifyHeaderWorker(chain, headers, seals, index)
				done <- index
			}
		}()
	// }

	errorsOut := make(chan error, len(headers))
	go func() {
		defer close(inputs)
		var (
			in, out = 0, 0
			checked = make([]bool, len(headers))
			inputs  = inputs
		)
		for {
			select {
			case inputs <- in:
				if in++; in == len(headers) {
					// Reached end of headers. Stop sending to workers.
					inputs = nil
				}
			case index := <-done:
				for checked[index] = true; checked[out]; out++ {
					errorsOut <- errors[out]
					if out == len(headers)-1 {
						return
					}
				}
			case <-abort:
				return
			}
		}
	}()
	return abort, errorsOut
}

func (ethash *Ethash) verifyHeaderWorker(chain consensus.ChainReader, headers []*types.Header, seals []bool, index int) error {
	var parent *types.Header
	if index == 0 {
		parent = chain.GetHeader(headers[0].ParentHash, headers[0].Number.Uint64()-1)
	} else if headers[index-1].Hash() == headers[index].ParentHash {
		parent = headers[index-1]
	}
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	if chain.GetHeader(headers[index].Hash(), headers[index].Number.Uint64()) != nil {
		return nil // known block
	}

	if err:=ethash.verifyGenerationSignature(chain,headers,index);err!=nil{
		return err
	}

	expected,err:=ethash.verifyDifficulty(chain,headers,index)
	if err!=nil{
		return err
	}
	if expected.Cmp(headers[index].Difficulty) != 0 {
		return fmt.Errorf("invalid difficulty: have %v, want %v", headers[index].Difficulty, expected)
	}

	return ethash.verifyHeader(chain, headers[index], parent,seals[index])
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of the stock Ethereum ethash engine.

// verifyHeader checks whether a header conforms to the consensus rules of the
// stock Ethereum ethash engine.
// See YP section 4.3.4. "Block Header Validity"
func (ethash *Ethash) verifyHeader(chain consensus.ChainReader, header, parent *types.Header, seal bool) error {
//verifyTheFates
	if header.TheFates!=parent.TheFates{
	   return fmt.Errorf("invalid TheFates:have: %d ,want %d", header.TheFates.Hex(), parent.TheFates.Hex())
	}
	// Ensure that the header's extra-data section is of a reasonable size
	if uint64(len(header.Extra)) > params.MaximumExtraDataSize {
		return fmt.Errorf("extra-data too long: %d > %d", len(header.Extra), params.MaximumExtraDataSize)
	}
	// Verify the header's timestamp
	if header.Time > uint64(time.Now().Add(allowedFutureBlockTime).Unix()) {
			return consensus.ErrFutureBlock
		}

	if header.Time <= parent.Time {
		return errZeroBlockTime
	}
	// Verify the block's difficulty based in its timestamp and parent's difficulty
	// Verify that the gas limit is <= 2^63-1
	cap := uint64(0x7fffffffffffffff)
	if header.GasLimit > cap {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, cap)
	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
	}

	// Verify that the gas limit remains within allowed bounds
	diff := int64(parent.GasLimit) - int64(header.GasLimit)
	if diff < 0 {
		diff *= -1
	}
	limit := parent.GasLimit / params.GasLimitBoundDivisor

	if uint64(diff) >= limit || header.GasLimit < params.MinGasLimit {
		return fmt.Errorf("invalid gas limit: have %d, want %d += %d", header.GasLimit, parent.GasLimit, limit)
	}
	// Verify that the block number is parent's +1
	if diff := new(big.Int).Sub(header.Number, parent.Number); diff.Cmp(big.NewInt(1)) != 0 {
		return consensus.ErrInvalidNumber
	}
	// Verify the engine specific seal securing the block
	if seal {
		if err := ethash.VerifySeal(chain, header); err != nil {
			return err
		}
	}
	// If all checks passed, validate any special fields for hard forks
	if err := misc.VerifyDAOHeaderExtraData(chain.Config(), header); err != nil {
		return err
	}
	if err := misc.VerifyForkHashes(chain.Config(), header); err != nil {
		return err
	}


	return nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func (ethash *Ethash) CalcDifficulty(chain consensus.ChainReader, header *types.Header, parent *types.Header) *big.Int {
	return CalcDifficulty(chain, header, parent)
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func CalcDifficulty(chain consensus.ChainReader, header *types.Header, parent *types.Header) *big.Int {
	// next := new(big.Int).Add(parent.Number, big1)
	// switch {
	// case config.IsConstantinople(next):
	// 	return calcDifficultyConstantinople(time, parent)
	// case config.IsByzantium(next):
	// 	return calcDifficultyByzantium(time, parent)
	// case config.IsHomestead(next):
	// 	return calcDifficultyHomestead(time, parent)
	// default:
	// 	return calcDifficultyFrontier(time, parent)
	// }
	difficulty:=calculateDifficulty(chain,header,parent)
	return difficulty
}

// Some weird constants to avoid constant memory allocs for them.
var (
	expDiffPeriod = big.NewInt(100000)
	big1          = big.NewInt(1)
	big2          = big.NewInt(2)
	big9          = big.NewInt(9)
	big10         = big.NewInt(10)
	bigMinus99    = big.NewInt(-99)
)

// makeDifficultyCalculator creates a difficultyCalculator with the given bomb-delay.
// the difficulty is calculated with Byzantium rules, which differs from Homestead in
// how uncles affect the calculation
func makeDifficultyCalculator(bombDelay *big.Int) func(time uint64, parent *types.Header) *big.Int {
	// Note, the calculations below looks at the parent number, which is 1 below
	// the block number. Thus we remove one from the delay given
	bombDelayFromParent := new(big.Int).Sub(bombDelay, big1)
	return func(time uint64, parent *types.Header) *big.Int {
		// https://github.com/ethereum/EIPs/issues/100.
		// algorithm:
		// diff = (parent_diff +
		//         (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
		//        ) + 2^(periodCount - 2)

		bigTime := new(big.Int).SetUint64(time)
		bigParentTime := new(big.Int).SetUint64(parent.Time)

		// holds intermediate values to make the algo easier to read & audit
		x := new(big.Int)
		y := new(big.Int)

		// (2 if len(parent_uncles) else 1) - (block_timestamp - parent_timestamp) // 9
		x.Sub(bigTime, bigParentTime)
		x.Div(x, big9)
		x.Sub(big1, x)
		// max((2 if len(parent_uncles) else 1) - (block_timestamp - parent_timestamp) // 9, -99)
		if x.Cmp(bigMinus99) < 0 {
			x.Set(bigMinus99)
		}
		// parent_diff + (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
		y.Div(parent.Difficulty, params.DifficultyBoundDivisor)
		x.Mul(y, x)
		x.Add(parent.Difficulty, x)

		// minimum difficulty can ever be (before exponential factor)
		if x.Cmp(params.MinimumDifficulty) < 0 {
			x.Set(params.MinimumDifficulty)
		}
		// calculate a fake block number for the ice-age delay
		// Specification: https://eips.ethereum.org/EIPS/eip-1234
		fakeBlockNumber := new(big.Int)
		if parent.Number.Cmp(bombDelayFromParent) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, bombDelayFromParent)
		}
		// for the exponential factor
		periodCount := fakeBlockNumber
		periodCount.Div(periodCount, expDiffPeriod)

		// the exponential factor, commonly referred to as "the bomb"
		// diff = diff + 2^(periodCount - 2)
		if periodCount.Cmp(big1) > 0 {
			y.Sub(periodCount, big2)
			y.Exp(big2, y, nil)
			x.Add(x, y)
		}
		return x
	}
}

// calcDifficultyHomestead is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time given the
// parent block's time and difficulty. The calculation uses the Homestead rules.
func calcDifficultyHomestead(time uint64, parent *types.Header) *big.Int {
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2.md
	// algorithm:
	// diff = (parent_diff +
	//         (parent_diff / 2048 * max(1 - (block_timestamp - parent_timestamp) // 10, -99))
	//        ) + 2^(periodCount - 2)

	bigTime := new(big.Int).SetUint64(time)
	bigParentTime := new(big.Int).SetUint64(parent.Time)

	// holds intermediate values to make the algo easier to read & audit
	x := new(big.Int)
	y := new(big.Int)

	// 1 - (block_timestamp - parent_timestamp) // 10
	x.Sub(bigTime, bigParentTime)
	x.Div(x, big10)
	x.Sub(big1, x)

	// max(1 - (block_timestamp - parent_timestamp) // 10, -99)
	if x.Cmp(bigMinus99) < 0 {
		x.Set(bigMinus99)
	}
	// (parent_diff + parent_diff // 2048 * max(1 - (block_timestamp - parent_timestamp) // 10, -99))
	y.Div(parent.Difficulty, params.DifficultyBoundDivisor)
	x.Mul(y, x)
	x.Add(parent.Difficulty, x)

	// minimum difficulty can ever be (before exponential factor)
	if x.Cmp(params.MinimumDifficulty) < 0 {
		x.Set(params.MinimumDifficulty)
	}
	// for the exponential factor
	periodCount := new(big.Int).Add(parent.Number, big1)
	periodCount.Div(periodCount, expDiffPeriod)

	// the exponential factor, commonly referred to as "the bomb"
	// diff = diff + 2^(periodCount - 2)
	if periodCount.Cmp(big1) > 0 {
		y.Sub(periodCount, big2)
		y.Exp(big2, y, nil)
		x.Add(x, y)
	}
	return x
}

// calcDifficultyFrontier is the difficulty adjustment algorithm. It returns the
// difficulty that a new block should have when created at time given the parent
// block's time and difficulty. The calculation uses the Frontier rules.
func calcDifficultyFrontier(time uint64, parent *types.Header) *big.Int {
	diff := new(big.Int)
	adjust := new(big.Int).Div(parent.Difficulty, params.DifficultyBoundDivisor)
	bigTime := new(big.Int)
	bigParentTime := new(big.Int)

	bigTime.SetUint64(time)
	bigParentTime.SetUint64(parent.Time)

	if bigTime.Sub(bigTime, bigParentTime).Cmp(params.DurationLimit) < 0 {
		diff.Add(parent.Difficulty, adjust)
	} else {
		diff.Sub(parent.Difficulty, adjust)
	}
	if diff.Cmp(params.MinimumDifficulty) < 0 {
		diff.Set(params.MinimumDifficulty)
	}

	periodCount := new(big.Int).Add(parent.Number, big1)
	periodCount.Div(periodCount, expDiffPeriod)
	if periodCount.Cmp(big1) > 0 {
		// diff = diff + 2^(periodCount - 2)
		expDiff := periodCount.Sub(periodCount, big2)
		expDiff.Exp(big2, expDiff, nil)
		diff.Add(diff, expDiff)
		diff = math.BigMax(diff, params.MinimumDifficulty)
	}
	return diff
}

// VerifySeal implements consensus.Engine, checking whether the given block satisfies
// the PoW difficulty requirements.
func (ethash *Ethash) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	return ethash.verifySeal(chain, header, false)
}

// verifySeal checks whether a block satisfies the PoW difficulty requirements,
// either using the usual ethash cache for it, or alternatively using a full DAG
// to make remote mining fast.
func (ethash *Ethash) verifySeal(chain consensus.ChainReader, header *types.Header, fulldag bool) error {
	// If we're running a fake PoW, accept any seal as valid
	if ethash.config.PowMode == ModeFake || ethash.config.PowMode == ModeFullFake {
		time.Sleep(ethash.fakeDelay)
		if ethash.fakeFail == header.Number.Uint64() {
			return errInvalidPoW
		}
		return nil
	}
	// If we're running a shared PoW, delegate verification to it
	if ethash.shared != nil {
		return ethash.shared.verifySeal(chain, header, fulldag)
	}
	// Ensure that we have a valid difficulty for the block
	if header.Difficulty.Sign() <= 0 {
		return errInvalidDifficulty
	}
	// Recompute the digest and PoW values
	// number := header.Number.Uint64()

	// var (
	// 	digest []byte
	// 	result []byte
	// )
	// // If fast-but-heavy PoW verification was requested, use an ethash dataset
	// if fulldag {
	// 	dataset := ethash.dataset(number, true)
	// 	if dataset.generated() {
	// 		digest, result = hashimotoFull(dataset.dataset, ethash.SealHash(header).Bytes(), header.Nonce.Uint64())

	// 		// Datasets are unmapped in a finalizer. Ensure that the dataset stays alive
	// 		// until after the call to hashimotoFull so it's not unmapped while being used.
	// 		runtime.KeepAlive(dataset)
	// 	} else {
	// 		// Dataset not yet generated, don't hang, use a cache instead
	// 		fulldag = false
	// 	}
	// }
	// // If slow-but-light PoW verification was requested (or DAG not yet ready), use an ethash cache
	// if !fulldag {
	// 	cache := ethash.cache(number)

	// 	size := datasetSize(number)
	// 	if ethash.config.PowMode == ModeTest {
	// 		size = 32 * 1024
	// 	}
	// 	digest, result = hashimotoLight(size, cache.cache, ethash.SealHash(header).Bytes(), header.Nonce.Uint64())

	// 	// Caches are unmapped in a finalizer. Ensure that the cache stays alive
	// 	// until after the call to hashimotoLight so it's not unmapped while being used.
	// 	runtime.KeepAlive(cache)
	// }
	// Verify the calculated values against the ones provided in the header
	// if !bytes.Equal(header.MixDigest[:], digest) {
	// 	return errInvalidMixDigest
	// }
	// target := new(big.Int).Div(two256, header.Difficulty)
	// if new(big.Int).SetBytes(result).Cmp(target) > 0 {
	// 	return errInvalidPoW
	// }


	return nil
}

// Prepare implements consensus.Engine, initializing the difficulty field of a
// header to conform to the ethash protocol. The changes are done inline.
func (ethash *Ethash) Prepare(chain consensus.ChainReader, header *types.Header) error {
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	header.Difficulty = ethash.CalcDifficulty(chain, header, parent)
	return nil
}

// Finalize implements consensus.Engine, accumulating the block and uncle rewards,
// setting the final state on the header
func (ethash *Ethash) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction) {
	// Accumulate any block and uncle rewards and commit the final state root
	accumulateRewards(chain,chain.Config(), state, header)
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
}

// FinalizeAndAssemble implements consensus.Engine, accumulating the block and
// uncle rewards, setting the final state and assembling the block.
func (ethash *Ethash) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) (*types.Block, error) {
	// Accumulate any block and uncle rewards and commit the final state root

	accumulateRewards(chain,chain.Config(), state, header)
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))

	// Header seems complete, assemble into a block and return
	return types.NewBlock(header, txs, receipts), nil
}

// SealHash returns the hash of a block prior to it being sealed.
func (ethash *Ethash) SealHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()

	rlp.Encode(hasher, []interface{}{
		header.ParentHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		// header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		// header.Time,
		header.Extra,
//linzhaojie header code 
		header.TheFates,
		
	})
	hasher.Sum(hash[:0])
	return hash
}

// Some weird constants to avoid constant memory allocs for them.
var (
	big8  = big.NewInt(8)
	big32 = big.NewInt(32)
)

// AccumulateRewards credits the coinbase of the given block with the mining
// reward. The total reward consists of the static block reward and rewards for
// included uncles. The coinbase of each uncle block is also rewarded.
func accumulateRewards(chain consensus.ChainReader,config *params.ChainConfig, state *state.StateDB, header *types.Header) {
	// Select the correct block reward based on chain progression
	blockReward:=GetReward(header.Number.Uint64())
	reward:= new(big.Int).Mul(big.NewInt(int64(blockReward)),big.NewInt(1e+18))

	minerCapacityTB:=setMinerCapacityTB(chain,header,state)
	requiredAmount:=new(big.Int).Mul(minerCapacityTB,mortgageToken)
	header.RequiredAmount=requiredAmount
	enough:=state.GetBalance(header.Coinbase).Cmp(requiredAmount)
	if enough==1||enough==0{
	state.AddBalance(header.Coinbase, reward)
	header.RewardAmount=reward
	state.AddTotalMine(header.Coinbase)
	}else{
	state.AddBalance(header.Coinbase, reward)
	header.RewardAmount=reward
	state.AddTotalMine(header.Coinbase)
	}	
}

func calculateDifficulty(chain consensus.ChainReader,header *types.Header,parent *types.Header) *big.Int {
	if header.Number.Uint64()<=24{
		return big.NewInt(constants.INITIAL_BASE_TARGET)
	}else {
		 itHeader := parent
		 avgDifficulty := itHeader.Difficulty

		 var blockCounter int64 = 1


		 for i:=0; i<24; i++  {
			parentNumber := itHeader.Number

			itHeader = chain.GetHeader(itHeader.ParentHash,parentNumber.Uint64()-1)
			// if itHeader == nil {
			// 	return nil,errors.New(fmt.Sprintf(
			// 		"ParentBlock does no longer exist for block height %d",parentNumber))
			// }
			blockCounter++
			MU:=new(big.Int).Mul(avgDifficulty,big.NewInt(blockCounter))
			ADD:=new(big.Int).Add(MU,itHeader.Difficulty)
			DIV:=new(big.Int).Div(ADD,big.NewInt(blockCounter + 1))
			avgDifficulty=DIV
			

		 }
		//  fmt.Println("------------------------------------")
        //  fmt.Println("header.Number",header.Number)
        //  fmt.Println("avgDifficulty",avgDifficulty)

		 difTime := header.Time - itHeader.Time
        //  fmt.Println("header.Time",header.Time)
        //  fmt.Println("itHeader.Time",itHeader.Time)
        //  fmt.Println("difTime",difTime)

		 targetTimespan := 24 * constants.RECIPIENT_ASSIGNMENT_WAIT_TIME *60



		if difTime < uint64(targetTimespan) / 2 {
			// fmt.Println("11111")
			difTime = uint64(targetTimespan) / 2
		}


		if difTime > uint64(targetTimespan) * 2{
			// fmt.Println("222222")
			difTime = uint64(targetTimespan) * 2
		}



		curDifficulty:=parent.Difficulty
        // fmt.Println("curDifficulty",curDifficulty)

		MU:=new(big.Int).Mul(avgDifficulty,new(big.Int).SetUint64(difTime))
		DIV:=new(big.Int).Div(MU,big.NewInt(int64(targetTimespan)))
		newDifficulty:=DIV.Int64()
		// fmt.Println("newDifficulty",newDifficulty)
		if newDifficulty < 0 || newDifficulty > constants.MAX_BASE_TARGET{
			// fmt.Println("reset")
			newDifficulty = constants.MAX_BASE_TARGET
		}
		if newDifficulty == 0 {
			newDifficulty = 1
		}

		if newDifficulty < curDifficulty.Int64()*8/10{
			newDifficulty = curDifficulty.Int64()*8/10
		}
		if newDifficulty> curDifficulty.Int64()*12 / 10{
			newDifficulty = curDifficulty.Int64()*12 / 10
		}
		// fmt.Println("newDifficulty",newDifficulty)
		return big.NewInt(int64(newDifficulty))
		
	}
	// return nil
}
func (ethash *Ethash)verifyDifficulty(chain consensus.ChainReader,headers []*types.Header,index int) (*big.Int,error) {
	var parent *types.Header
	if index == 0 {
		parent = chain.GetHeader(headers[0].ParentHash, headers[0].Number.Uint64()-1)
	} else if headers[index-1].Hash() == headers[index].ParentHash {
		
		parent = headers[index-1]

	}
	if parent == nil {
		return nil,consensus.ErrUnknownAncestor
	}

	if headers[index].Number.Uint64()<=24{

		return big.NewInt(constants.INITIAL_BASE_TARGET),nil
   
	}else {
		 itHeader := parent
		 avgDifficulty := itHeader.Difficulty
		 var blockCounter int64 = 1
		 Index:=index-1

		 for i:=0; i<24; i++  {
			parentNumber := itHeader.Number
			if Index<= 0 {
				itHeader = chain.GetHeader(itHeader.ParentHash,parentNumber.Uint64()-1)
			} else if headers[Index-1].Hash() == headers[Index].ParentHash {
				itHeader = headers[Index-1]
				Index--

			}
			
			blockCounter++

			MU:=new(big.Int).Mul(avgDifficulty,big.NewInt(blockCounter))
			ADD:=new(big.Int).Add(MU,itHeader.Difficulty)
	
			DIV:=new(big.Int).Div(ADD,big.NewInt(blockCounter + 1))
			avgDifficulty=DIV

		 }
	//   fmt.Println("------------------------------------")
    //   fmt.Println("header.Number",headers[index].Number)
	//   fmt.Println("avgDifficulty",avgDifficulty)

	  difTime := headers[index].Time - itHeader.Time
	//   fmt.Println("header.Time",headers[index].Time)
	//   fmt.Println("itHeader.Time",itHeader.Time)
	//   fmt.Println("difTime",difTime)

	  targetTimespan := 24 * constants.RECIPIENT_ASSIGNMENT_WAIT_TIME *60



	 if difTime < uint64(targetTimespan) / 2 {
		//  fmt.Println("11111")
		 difTime = uint64(targetTimespan) / 2
	 }


	 if difTime > uint64(targetTimespan) * 2{
		//  fmt.Println("222222")
		 difTime = uint64(targetTimespan) * 2
	 }


	 curDifficulty:=parent.Difficulty
	//  fmt.Println("curDifficulty",curDifficulty)

	 MU:=new(big.Int).Mul(avgDifficulty,new(big.Int).SetUint64(difTime))
	 DIV:=new(big.Int).Div(MU,big.NewInt(int64(targetTimespan)))
	 newDifficulty:=DIV.Int64()
	//  fmt.Println("newDifficulty",newDifficulty)
	 if newDifficulty < 0 || newDifficulty > constants.MAX_BASE_TARGET{
		//  fmt.Println("reset")
		 newDifficulty = constants.MAX_BASE_TARGET
	 }
	 if newDifficulty == 0 {
		 newDifficulty = 1
	 }

	 if newDifficulty < curDifficulty.Int64()*8/10{
		 newDifficulty = curDifficulty.Int64()*8/10
	 }
	 if newDifficulty> curDifficulty.Int64()*12 / 10{
		 newDifficulty = curDifficulty.Int64()*12 / 10
	 }
	//  fmt.Println("newDifficulty",newDifficulty)
	 return big.NewInt(int64(newDifficulty)),nil
		
	}
	return nil,nil
}
func setMinerCapacityTB(chain consensus.ChainReader,header *types.Header,state *state.StateDB)*big.Int{
    parentBlock:=chain.CurrentBlock()
	networkCapacityTB:=	new(big.Int).Div(InitDifficulty,parentBlock.Difficulty())
	header.NetCapacity=networkCapacityTB.Uint64()*10
	// fmt.Println("blockNumber",header.Number)
    // fmt.Println("InitDifficulty:",InitDifficulty)
    // fmt.Println("Difficulty:",header.Difficulty)
	// fmt.Println("networkCapacityTB:",networkCapacityTB)
	TotalMine:=state.GetTotalMine(header.Coinbase)
	if TotalMine==0{
		TotalMine=1
	}
	// fmt.Println("TotalMine",TotalMine)
	// fmt.Println("CoinBase:",header.Coinbase.Hex())
	minerCapacityTB:=new(big.Int).Div(new(big.Int).Mul(networkCapacityTB,big.NewInt(int64(TotalMine))),header.Number)

	// fmt.Println("minerCapacityTB:",minerCapacityTB)
	// fmt.Println("-----------------------------------")
  
	if minerCapacityTB.Cmp(big.NewInt(1))==-1{
		minerCapacityTB=big.NewInt(1)
	}
	    header.MinerCapacity=minerCapacityTB.Uint64()*10
		return minerCapacityTB
}



func (ethash *Ethash) verifyGenerationSignature(chain consensus.ChainReader,headers []*types.Header,index int)error{
	var parent *types.Header
	if index == 0 {
		parent = chain.GetHeader(headers[0].ParentHash, headers[0].Number.Uint64()-1)
	} else if headers[index-1].Hash() == headers[index].ParentHash {
		parent = headers[index-1]
	}
	if parent==nil{
		log.Error("Con't verify generation signature because pre block is missing","err")
	}


	bspr:=hashid.CalculateGenSinGen(parent.GenerationSignature,parent.Generator)
	if !bytes.Equal(headers[index].GenerationSignature,bspr) {
		return fmt.Errorf("invalid GenerationSignature (remote:%x,local:%x)",headers[index].GenerationSignature,bspr)
	}

	scoopnum:=hashid.CalculateScoop(headers[index].GenerationSignature,headers[index].Number)
	deadline:=hashid.CalcuDeadline(headers[index].Generator,big.NewInt(int64(headers[index].Nonce.Uint64())),parent.Difficulty,scoopnum,headers[index].GenerationSignature)
	elapsedTime:=headers[index].Time-parent.Time

	if big.NewInt(int64(elapsedTime)).Cmp(big.NewInt(deadline))<0{
		return fmt.Errorf("invalid deadline (elapsedTime: %v,deadline: %v)",elapsedTime,deadline)
	}
     return nil
}