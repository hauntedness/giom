package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hauntedness/giom/internal/log"
	"github.com/hauntedness/httputil"
	"github.com/sourcegraph/conc/iter"
)

//go:embed testdata/*
var htmlList embed.FS

func main() {
	headers := httputil.H{
		"accept":       "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Content-Type": "text/html",
	}
	errs := make([]error, 446)
	iter.ForEachIdx(errs, func(id int, box *error) {
		url := fmt.Sprintf("https://www.ddxs.com/xiaoyaoxiaoguixu/%d.html", id)
		data, err := httputil.Get(url, headers)
		if err != nil {
			log.Errors(err, "url", url)
			*box = err
			return
		}
		path := filepath.Join("testdata", strconv.Itoa(id)+".html")
		err = os.WriteFile(path, data, os.ModePerm)
		if err != nil {
			log.Errors(err, "path", path)
			*box = err
			return
		}
	})
	for index, err := range errs {
		if err != nil {
			log.Errors(err, "index", index)
		}
	}
}
