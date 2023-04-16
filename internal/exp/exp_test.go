package exp

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/hauntedness/giom/internal/log"
)

func TestMarshal(t *testing.T) {
	time_, err := time.Parse(time.DateTime, "2023-04-07 22:04:05")
	if err != nil {
		t.Error(err)
		return
	}
	book := Book{
		CreatedAt: &time_,
		Name:      "amy",
		Age:       NewPointer(18),
	}
	data, err := json.Marshal(&book)
	if err != nil {
		t.Error(err)
		return
	}
	log.Info(string(data))
	book2 := Book{}
	err = json.Unmarshal(data, &book2)
	if err != nil {
		t.Error(err)
		return
	}
	if book2.CreatedAt == nil {
		t.Errorf("actual: %v, want: %v", book2.CreatedAt, nil)
	}
}

func TestMeta(t *testing.T) {
	src := reflect.TypeOf(Book{})
	dst := reflect.TypeOf(BookReplacer{})
	if src.NumField() != dst.NumField() {
		t.Errorf("actual: %v, want: %v", "equals", "not equals")
		return
	}
	for i := 0; i < src.NumField(); i++ {
		srcField := src.Field(i)
		dstField := dst.Field(i)
		if !identical(srcField, dstField) {
			continue
		}
		log.Error("not identical", "src field", srcField)
		log.Error("not identical", "dst field", dstField)
		t.FailNow()
	}
}

func identical(src reflect.StructField, dst reflect.StructField) bool {
	if src.Name != dst.Name {
		return false
	}
	if src.Tag != dst.Tag {
		return false
	}
	if src.Type != dst.Type && !src.Type.ConvertibleTo(dst.Type) {
		return false
	}
	return true
}

func NewPointer[T any](value T) *T {
	return &value
}
