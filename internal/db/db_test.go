package db

import (
	"testing"

	"github.com/hauntedness/giom/internal/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Hello struct {
	gorm.Model
	Name string
}

func Test_connect(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	if err != nil {
		t.Errorf("connect() error = %v", err)
		return
	}
	h := &Hello{}
	err = db.AutoMigrate(h)
	if err != nil {
		t.Error(err)
		return
	}
	model := db.Model(h)
	model.Exec("delete from hellos where 1 = 1")
	if err := model.Error; err != nil {
		t.Error(err)
		return
	}
	model.Create(&Hello{Name: "amy"})
	if model.Error != nil {
		t.Error(model.Error)
		return
	}
	hellos := make([]Hello, 7)
	model.Find(h, map[string]any{"name": "amy"}).Scan(&hellos)
	if model.Error != nil {
		t.Error(model.Error)
		return
	}
	for _, h := range hellos {
		log.Infos("hello", h)
	}
	hello2 := make([]Hello, 7)
	d := model.Raw("select * from hellos where name = ?", "amy").Scan(&hello2)
	if d.Error != nil {
		t.Error(d.Error)
	}
	for _, h2 := range hello2 {
		log.Infos("hello2", h2)
	}
}
