package hitler

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
	"time"
)

type Machine struct {
	Id          int64
	CreatedTime time.Time
}

type MachineManager struct {
	engine *xorm.Engine
}

func newMachineManager(mysqlUrl string) (*MachineManager, error) {
	engine, err := xorm.NewEngine("mysql", mysqlUrl)
	return &MachineManager{engine: engine}, err
}

func (m *MachineManager) NewMId() (int64, error) {
	newM := &Machine{CreatedTime: time.Now()}
	_, err := m.engine.Insert(newM)
	if err == nil && newM.Id > 0 {
		m.engine.Delete(newM)
	}
	return newM.Id, err
}
