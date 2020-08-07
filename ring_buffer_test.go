package hitler

import (
	"log"
	"testing"
	"time"
)

func TestRBuffer_GetID(t *testing.T) {
	rb := NewRBuffer(&RBufferConfig{DbUrl: "root:Dachang1234!@(127.0.0.1:3306)/id_gen?charset=utf8mb4"})
	for i := 0; i < 20000; i++ {
		time.Sleep(time.Millisecond)
		log.Println(rb.GetID())
	}
}
