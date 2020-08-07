package hitler

import (
	"log"
	"testing"
)

func TestMachineManager_NewMId(t *testing.T) {
	//root:root@/orm_test?charset=utf8
	m, err := newMachineManager("root:Dachang1234!@(127.0.0.1:3306)/id_gen?charset=utf8mb4")
	log.Println(err)
	newId, err := m.NewMId()
	log.Println(newId)
	log.Println(err)

}
