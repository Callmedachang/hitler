package hitler

import (
	"time"
)

const (
	defaultSequenceCap = 20
	defaultMachineCap  = 15
	defaultTimeCap     = 28
)

type RBuffer struct {
	idsCursor  int64
	idsTail    int64
	ids        []int64
	flagCursor int64
	flagTail   int64
	flags      []bool
	size       int64
	mid        int64 //currentBuffer MachineID
	//Notice sequenceCap+machineCap+timeCap=63
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
	if m, err := newMachineManager(conf.DbUrl); err != nil {
		panic("invalid Db")
	} else {
		if mid, err = m.NewMId(); err != nil {
			panic("NewMId error")
		}
	}
	sequenceSize := 1 << conf.SequenceCap
	rb := &RBuffer{
		size:        int64(sequenceSize),
		ids:         make([]int64, sequenceSize),
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
		res = r.ids[r.idsCursor]
		r.flags[r.flagCursor] = false
		r.flagCursor++
		r.idsCursor++
		if r.flagCursor == r.size {
			r.flagCursor = 0
		}
		if r.idsCursor == r.size {
			r.idsCursor = 0
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
	r.ids[r.idsTail] = time.Now().Unix()<<r.timeCap + r.mid<<r.machineCap + r.idsTail
	r.idsTail++
	if r.idsTail >= r.size {
		r.idsTail = 0
	}
	//flags的复位
	r.flags[r.flagTail] = true
	r.flagTail++
	if r.flagTail >= r.size {
		r.flagTail = 0
	}
	return
}
