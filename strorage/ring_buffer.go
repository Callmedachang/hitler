package strorage

import (
	"fmt"
	"hitler/id_gen"
	"hitler/utils"
	"strconv"
	"time"
)

const minCap = 1024

type RBuffer struct {
	g          *id_gen.Generator
	iDsCursor  int64
	iDsTail    int64
	iDs        []int64
	flagCursor int64
	flagTail   int64
	flags      []bool
	size       int64
	binCap     int
	mid        string
}
type RBufferConfig struct {
	Size       int64
	MachineCap int
	DbUrl      string
}

func NewRBuffer(conf *RBufferConfig) *RBuffer {
	if conf.Size < 1 {
		panic("RBuffer size must more than 1")
	}
	conf.Size = utils.NextPowOf2(conf.Size)
	if conf.Size < minCap {
		conf.Size = minCap
	}
	mid := int64(0)
	if m, err := NewMachineManager(conf.DbUrl); err != nil {
		panic("invalid Db")
	} else {
		if mid, err = m.NewMId(); err != nil {
			panic("NewMId error")
		}
	}
	machineId := utils.ConvertToBin(mid)
	if len([]rune(machineId)) < conf.MachineCap {
		for i := len(machineId); i < conf.MachineCap; i++ {
			machineId = "0" + machineId
		}
	}
	rb := &RBuffer{
		g:          id_gen.NewGenerator(),
		size:       conf.Size,
		binCap:     len(utils.ConvertToBin(conf.Size)),
		iDs:        make([]int64, conf.Size),
		flags:      make([]bool, conf.Size),
		mid:        machineId,
	}
	for i := 0; i < int(conf.Size); i++ {
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
	sequence := utils.ConvertToBin(r.iDsTail)
	if len([]rune(sequence)) < r.binCap {
		for i := len(sequence); i < r.binCap; i++ {
			sequence = "0" + sequence
		}
	}
	newIDStr := fmt.Sprintf("%s%s%s", r.g.GetID(), r.mid, sequence)

	r.iDs[r.iDsTail], _ = strconv.ParseInt(newIDStr, 2, 64)
	r.iDsTail++
	if r.iDsTail >= r.size {
		r.iDsTail = 0
	}
	r.flags[r.flagTail] = true
	r.flagTail++
	if r.flagTail >= r.size {
		r.flagTail = 0
	}
	return
}
