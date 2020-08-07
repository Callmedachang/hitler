package strorage

import (
	"hitler/id_gen"
	"time"
)

const (
	defaultSequenceCap = 20
	defaultMachineCap  = 15
	defaultTimeCap     = 28
)

type RBuffer struct {
	g          *id_gen.Generator
	iDsCursor  int64
	iDsTail    int64
	iDs        []int64
	flagCursor int64
	flagTail   int64
	flags      []bool
	size       int64
	mid        int64 //当钱的机器ID
	//以下三个参数想起来为63！
	sequenceCap uint //自增序列号的占用位数
	machineCap  uint //机器ID的占用位数
	timeCap     uint //时间Seconds的占用位数
}

type RBufferConfig struct {
	DbUrl       string
	SequenceCap uint //自增序列号的占用位数
	MachineCap  uint //机器ID的占用位数
	TimeCap     uint //时间Seconds的占用位数
}

func NewRBuffer(conf *RBufferConfig) *RBuffer {
	if conf.SequenceCap == 0 && conf.MachineCap == 0 && conf.TimeCap == 0 {
		conf.MachineCap = defaultMachineCap
		conf.SequenceCap = defaultSequenceCap
		conf.TimeCap = defaultTimeCap
	} else {
		if conf.SequenceCap+conf.MachineCap+conf.TimeCap != 63 {
			panic("invalid cap params")
		}
	}
	mid := int64(0)
	if m, err := NewMachineManager(conf.DbUrl); err != nil {
		panic("invalid Db")
	} else {
		if mid, err = m.NewMId(); err != nil {
			panic("NewMId error")
		}
	}
	sequenceSize := 1 << conf.SequenceCap
	rb := &RBuffer{
		g:           id_gen.NewGenerator(),
		size:        int64(sequenceSize),
		iDs:         make([]int64, sequenceSize),
		flags:       make([]bool, sequenceSize),
		mid:         mid,
		sequenceCap: conf.SequenceCap,
		machineCap:  conf.MachineCap,
		timeCap:     conf.TimeCap,
	}
	for i := 0; i < sequenceSize; i++ {
		rb.createID()
	}
	rb.cycleStuff()
	return rb
}

func (r *RBuffer) GetID() (res int64) {
	if r.flags[r.flagCursor] {
		res = r.iDs[r.iDsCursor]
		r.flags[r.flagCursor] = false
		r.flagCursor++
		r.iDsCursor++
		if r.flagCursor == r.size {
			r.flagCursor = 0
		}
		if r.iDsCursor == r.size {
			r.iDsCursor = 0
		}
		return
	} else {
		res = -1
	}
	return res
}

func (r *RBuffer) cycleStuff() {
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				for {
					if r.flags[r.flagTail] == false {
						r.createID()
					} else {
						break
					}
				}
			}
		}
	}()
}

func (r *RBuffer) createID() {
	//IDs的补充
	r.iDs[r.iDsTail] = r.g.GetID()<<r.timeCap + r.mid<<r.machineCap + r.iDsTail
	r.iDsTail++
	if r.iDsTail >= r.size {
		r.iDsTail = 0
	}
	//flags的复位
	r.flags[r.flagTail] = true
	r.flagTail++
	if r.flagTail >= r.size {
		r.flagTail = 0
	}
	return
}
