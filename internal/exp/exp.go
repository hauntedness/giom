package exp

import (
	"encoding/json"
	"time"
	"unsafe"
)

type Book struct {
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Name      string     `json:"name,omitempty"`
	Age       *int       `json:"age,omitempty"`
}

type BookReplacer struct {
	CreatedAt *Time  `json:"created_at,omitempty"`
	Name      string `json:"name,omitempty"`
	Age       *int   `json:"age,omitempty"`
}

func (b *Book) MarshalJSON() ([]byte, error) {
	replacer := (*BookReplacer)(unsafe.Pointer(b))
	return json.Marshal(replacer)
}

func (b *Book) UnmarshalJSON(data []byte) error {
	replacer := (*BookReplacer)(unsafe.Pointer(b))
	err := json.Unmarshal(data, &replacer)
	return err
}

type Time time.Time

func (t *Time) MarshalJSON() ([]byte, error) {
	text := (*time.Time)(t).Format(time.RFC1123)
	return json.Marshal(text)
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var text string
	err := json.Unmarshal(data, &text)
	if err != nil {
		return err
	}
	t2, err := time.Parse(time.RFC1123, text)
	if err != nil {
		return err
	}
	*t = Time(t2)
	return nil
}
