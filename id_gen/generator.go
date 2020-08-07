package id_gen

import (
	"time"
)

type Generator struct {
}

func (g *Generator) GetID() int64 {
	return time.Now().Unix()
}

func NewGenerator() *Generator {
	return &Generator{}
}
