package main

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/hauntedness/giom/internal/log"
	"github.com/hauntedness/giom/internal/runtime"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func Test_main(t *testing.T) {
	entryList, err := htmlList.ReadDir("testdata")
	if err != nil {
		t.Error(err)
		return
	}
	var sb strings.Builder
	slices.SortFunc(entryList, func(left, right fs.DirEntry) bool {
		if left.IsDir() {
			return false
		}
		l, _ := strconv.Atoi(strings.Split(left.Name(), ".")[0])
		r, _ := strconv.Atoi(strings.Split(right.Name(), ".")[0])
		return l < r
	})
	for i := range entryList {
		entry := entryList[i]
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			t.Error(err)
			return
		}
		filename := strings.Join([]string{"testdata", info.Name()}, "/")
		err = WriteChapter(&sb, filename)
		if err != nil {
			t.Error(err)
			return
		}
	}
	target, err := os.Create("testdata/novel/逍遥小贵婿.txt")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = target.Write([]byte(sb.String()))
	if err != nil {
		t.Error(err)
		return
	}
}

func WriteChapter(sb *strings.Builder, fileName string) error {
	file, err := htmlList.Open(fileName)
	if err != nil {
		return fmt.Errorf("%s: %w", runtime.Source(), err)
	}
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return fmt.Errorf("%s: %w", runtime.Source(), err)
	}
	doc.Find("h1").
		Each(
			func(_ int, s *goquery.Selection) {
				s.Text()
				sb.WriteString(s.Text())
				sb.WriteRune('\n')
				sb.WriteRune('\n')
			},
		)

	doc.Find("#contents > p").
		Each(
			func(_ int, s *goquery.Selection) {
				for _, n := range s.Nodes {
					if n.Type == html.ElementNode && n.DataAtom == atom.P {
						if n.FirstChild.Type == html.TextNode && n.FirstChild.DataAtom == atom.Atom(0) {
							sb.WriteString(n.FirstChild.Data)
							sb.WriteRune('\n')
						}
					}
				}
			},
		)

	return nil
}

func TestTransform(t *testing.T) {
	urls := make([]string, 10)
	lo.Async0(func() {
	})
	ch := lo.SliceToChannel(3, urls)
	for e := range ch {
		go func(elem string) {
			log.Info(elem)
		}(e)
	}
	time.Sleep(time.Second * 10)
}
