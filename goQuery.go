package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"os"
)

func query() {
	f, _ := os.Open(`source.html`)
	defer f.Close()
	myLog = newInfoLogger("./output.txt")
	defer myLog.closeFile()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatalln(err)
	}

	doc.Find(`table id="model-table" class="table table-striped"`).Find("tbody").Find(`tr`).Each(func(i int, s *goquery.Selection) {
		fmt.Println(i)
		s.Find(`td`).Each(func(_ int, s *goquery.Selection) {
			fmt.Println(s.Html())
		})
	})

}
