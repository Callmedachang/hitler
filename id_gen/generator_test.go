package id_gen

import (
	"log"
	"testing"
	"time"
)

func TestNewGenerator(t *testing.T) {
	s:=NewGenerator()
	for i:=0;i<10;i++{
		time.Sleep(time.Second)
		log.Println(s.GetID())
	}
}