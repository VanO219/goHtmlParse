package main

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	myLog  *infoLogger
	z      *html.Tokenizer
	tgname string
	tgtext string
	res    []Information
	depth  = 0
	data   = []string{}
)

func main() {
	query()
}



type Information struct {
	ProductName     string
	Category        string
	NumberOfReviews int64
	SKU             int64
	Seller          string
	Brand           string
	QuantityInStock int64
	Price           float64
	Discount        int64
	OldPrice        float64
	ACP             float64
	LP              float64
	AmountOfSales   int64
	Revenue         float64
}



type infoLogger struct {
	file *os.File
}

func newInfoLogger(path string) *infoLogger {
	var err error
	if path == "" {
		err = errors.New("newInfoLogger() не задан путь к файлу")
		log.Fatalln(err)
	}
	abPath, err := filepath.Abs(path)
	if err != nil {
		err = errors.Wrap(err, "newInfoLogger() не корректный путь к файлу")
		log.Fatalln(err)
	}
	file, err := os.OpenFile(abPath, os.O_RDWR|os.O_CREATE|os.O_SYNC|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		err = errors.Wrap(err, `newInfoLogger() не удалось открыть файл`)
		log.Fatalln(err)
	}
	return &infoLogger{file: file}
}

func (i *infoLogger) write(input string) {
	_, err := i.file.WriteString(fmt.Sprintf("%s \n", input))
	if err != nil {
		err = errors.Wrap(err, "write() ошибка записи в файл")
		log.Fatalln(err)
	}
}

func (i *infoLogger) closeFile() {
	if i.file == nil {
		return
	}
	if err := i.file.Close(); err != nil {
		log.Fatalln(err)
	}
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
