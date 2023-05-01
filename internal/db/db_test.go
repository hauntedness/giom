package db

import (
	"fmt"
	stdlog "log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hauntedness/giom/internal/log"
	"github.com/hauntedness/giom/internal/runtime"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

type Book struct {
	Name     string
	CreateAt *time.Time
	Words    int
}

type BookMap struct {
	schema []string
	Data   []map[string]any
}

func (bm *BookMap) Schema() []string {
	if len(bm.schema) != 0 {
		return bm.schema
	}
	if len(bm.Data) == 0 {
		return nil
	}
	for k := range bm.Data[0] {
		bm.schema = append(bm.schema, k)
	}
	return bm.schema
}

func TestInsertWithMap(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.New(
			stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags),
			logger.Config{
				SlowThreshold:             0,
				Colorful:                  false,
				IgnoreRecordNotFoundError: false,
				ParameterizedQueries:      false,
				LogLevel:                  logger.Info,
			},
		),
	})
	if err != nil {
		t.Errorf("connect() error = %v", err)
		return
	}
	book := &Book{}
	err = db.AutoMigrate(book)
	if err != nil {
		t.Error(err)
		return
	}
	model := db.Model(book)
	for i := 0; i < 7; i++ {
		randTime := time.Now().Add(time.Duration(rand.Uint64()))
		b := &Book{
			Name:     RandWord(),
			CreateAt: &randTime,
			Words:    rand.Intn(100),
		}
		ret := model.Create(b)
		if ret.Error != nil {
			t.Error(ret.Error)
			return
		}
	}
	bm := []map[string]any{}
	d := model.Scan(&bm)
	if d.Error != nil {
		t.Error(d.Error)
		return
	}
	if len(bm) == 0 {
		err = fmt.Errorf("%s: %s", runtime.Source(), "didn't get any data")
		t.Error(err)
		return
	}
	for _, v := range bm {
		if v == nil {
			t.Errorf("%s: %s", runtime.Source(), "empty value")
		}
	}
	sql := model.ToSQL(func(tx *gorm.DB) *gorm.DB {
		d2 := tx.Where("1 = ?", 1).Delete(book)
		if d2.Error != nil {
			t.Error(d2.Error)
			return db
		}
		return d2
	})
	log.Info(sql)
	books := []Book{}
	d3 := model.Find(&books)
	if d3.Error != nil {
		t.Error(d3.Error)
		return
	} else if len(books) > 0 {
		t.Error(fmt.Errorf("%s: %s", runtime.Source(), "expect no data"))
		return
	}
	d4 := model.Create(bm)
	if d4.Error != nil {
		t.Error(d4.Error)
		return
	}
	d5 := model.Find(&books)
	if d5.Error != nil {
		t.Error(d5.Error)
		return
	}
	if len(books) == 0 {
		t.Errorf("%s:, %s", runtime.Source(), "could not get any data")
	}
	empty := Book{}
	for i := range books {
		if books[i] == empty {
			t.Errorf("%s, %s", runtime.Source(), "empty source")
		}
	}
}

func RandWord() string {
	words := `Just click in the “Toggle Side Bar Visibility” box and remap the keybinding to Ctrl+Shift+B (since that doesn't get used by anything by default in Vim).  This menu is also really useful for just discovering all of the different keyboard shortcuts in Vscode.  I'm not really a fan of going crazy and totally customizing hotkeys in Vscode just because it makes support much harder but in this case it was easy enough and works really well with my workflow.`
	s := strings.Split(words, " ")
	random := rand.Intn(len(s) - 1)
	return s[random]
}
