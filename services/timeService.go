package services

import (
	"errors"
	"sync/atomic"
	"time"
)

type Times interface{
	getTime() int64
	getTimeInMillis() int64
}

type EpochTime struct{}

func (e *EpochTime) getTime() int64{
	return ((time.Now().UnixNano() / 1e6)) / 1000
}
func (e *EpochTime) getTimeInMillis() int64{
	return (time.Now().UnixNano() / 1e6)
}

type FasterTime struct {
	multiplier			int64
	systemStartTime		int64
	time				int64
}

func (f FasterTime) getTime() int64 {
	return f.time+((time.Now().UnixNano()/1e6 - f.systemStartTime) / (1000 / f.multiplier))
}

func (f FasterTime) getTimeInMillis() int64 {
	return f.time + ((time.Now().UnixNano()/1e6 - f.systemStartTime) / f.multiplier)
}

func NewFasterTime(tim int64, multiplier int64) (Times,error) {
	if multiplier > 1000 || multiplier <= 0 {
		return nil,errors.New("Time multiplier must be between 1 and 1000")
	}
	return FasterTime{
		multiplier:      multiplier,
		systemStartTime: time.Now().UnixNano()/1e6,
		time:            tim,
	},nil
}


type TimeService interface {

	GetEpochTime() uint64

	GetEpochTimeMillis()uint64

	SetTime(fasterTime *FasterTime)
}

type TimeServiceImpl struct {
	at *atomic.Value
}

func NewTimeServiceImpl() TimeServiceImpl {
	a:=new(atomic.Value)
	ts:=TimeServiceImpl{a}
	ts.at.Store(&EpochTime{})
	return ts
}

func (ti *TimeServiceImpl) GetEpochTime()uint64 {
	return uint64(ti.at.Load().(Times).getTime())
}
func (ti *TimeServiceImpl) GetEpochTimeMillis()uint64 {
	return uint64(ti.at.Load().(Times).getTimeInMillis())
}
func (ti *TimeServiceImpl) SetTime(fasterTime *FasterTime) {
	ti.at.Store(fasterTime)
}