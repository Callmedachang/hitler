package id_gen

import (
	"hitler/utils"
	"time"
)

type Generator struct {
}

func (g *Generator) GetID() string {
	h := utils.ConvertToBin(time.Now().Unix())
	start:=len(h)-28
	res := string([]rune(h)[start:])
	return res
}

func NewGenerator() *Generator {
	return &Generator{}
}
