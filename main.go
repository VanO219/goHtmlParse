package main

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	myLog     *infoLogger
	z         *html.Tokenizer
	tgname    string
	tgtext    string
	tableRows []Information
	depth     = 0
	data      = []string{}
)

func main() {
	f, _ := os.Open(`source.html`)
	defer f.Close()
	myLog = newInfoLogger("./output.txt")
	defer myLog.closeFile()

	z = html.NewTokenizer(f)

loop:
	for {
		tt := z.Next()
		tgtext = standardizeSpaces(strings.TrimSpace(string(z.Text())))

		switch tt {
		case html.ErrorToken:
			fmt.Println(z.Err())
			return
		case html.TextToken:
			if depth > 0 {
				//if tgtext != "" {
				//	myLog.write(fmt.Sprintf("TEXT, %s, %s", tgname, tgtext))
				//}
				continue
			}
		case html.StartTagToken, html.EndTagToken:
			if tt == html.StartTagToken {
				tgname = tagname(z)
				depth++
				if tgname == `tbody` {
					tableBodyParse()
					break loop
				}
			} else {
				//myLog.write(fmt.Sprintf("Close, %s, ", tagname(z)))
				depth--
				continue
			}
		case html.SelfClosingTagToken:
			continue
		default:
			continue
		}
	}
	fmt.Println(len(tableRows))
	for i, j := range tableRows {
		fmt.Println(i, j)
	}
}

func tableBodyParse() {
loop:
	for {
		tt := z.Next()
		tgtext = standardizeSpaces(strings.TrimSpace(string(z.Text())))

		switch tt {
		case html.ErrorToken:
			fmt.Println(z.Err())
			return
		case html.TextToken:
			if depth > 0 {
				if tgtext != "" {
					//myLog.write(fmt.Sprintf("TEXT, %s, %s", tgname, tgtext))
				}
				continue
			}
		case html.StartTagToken, html.EndTagToken:
			if tt == html.StartTagToken {
				tgname = tagname(z)
				depth++
				if tgname == `tr` {
					rawTableParse()
				}
			} else {
				depth--
				tgname = tagname(z)
				if tgname == `tbody` {
					break loop
				}
				continue
			}
		case html.SelfClosingTagToken:
			continue
		default:
			continue
		}
	}
}

func rawTableParse() {
loop:
	for {
		tt := z.Next()
		tgtext = standardizeSpaces(strings.TrimSpace(string(z.Text())))

		switch tt {
		case html.ErrorToken:
			fmt.Println(z.Err())
			return
		case html.TextToken:
			if depth > 0 {
				if tgtext != "" {
					//myLog.write(fmt.Sprintf("TEXT, %s, %s", tgname, tgtext))
					data = append(data, tgtext)
				}
				continue
			}
		case html.StartTagToken, html.EndTagToken:
			if tt == html.StartTagToken {
				tgname = tagname(z)
				depth++

			} else {
				depth--
				tgname = tagname(z)
				if tgname == `td` {
					cellParse()
				} else if tgname == `tr` {
					//myLog.write(`-------------------------------`)
					break loop
				}
				continue
			}
		case html.SelfClosingTagToken:
			continue
		default:
			continue
		}
	}
	if len(data) < 14 && len(data) > 0 {
		parseNotFull()
		data = []string{}
	}else if len(data) > 0 && len(data) == 14 {
		parseFull()
		data = []string{}
	}
}

func cellParse() {
loop:
	for {
		tt := z.Next()
		tgtext = standardizeSpaces(strings.TrimSpace(string(z.Text())))

		switch tt {
		case html.ErrorToken:
			fmt.Println(z.Err())
			return
		case html.TextToken:
			if depth > 0 {
				if tgtext != "" {
					//myLog.write(fmt.Sprintf("TEXT, %s, %s", tgname, tgtext))
					data = append(data, tgtext)
				}
				continue
			}
		case html.StartTagToken, html.EndTagToken:
			if tt == html.StartTagToken {
				tgname = tagname(z)
				depth++

			} else {
				depth--
				tgname = tagname(z)
				if tgname == `td` {
					break loop
				}
			}
		case html.SelfClosingTagToken:
			continue
		default:
			continue
		}
	}

}

func parseFull() {
	inf := Information{}
	for i, j := range data {
		switch i {
		case 0:
			inf.ProductName = j
			i++
		case 1:
			inf.Category = j
			i++
		case 2:
			if j == `null` {
				inf.NumberOfReviews = 0
			} else {
				tx := strings.Replace(j, "(", "", -1)
				tx = strings.Replace(tx, ")", "", -1)
				nm, err := strconv.Atoi(tx)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 2: "))
				}
				inf.NumberOfReviews = int64(nm)
			}
			i++
		case 3:
			if j == `null` {
				inf.SKU = 0
			} else {
				nm, err := strconv.Atoi(j)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 3: "))
				}
				inf.SKU = int64(nm)
			}
			i++
		case 4:
			inf.Seller = j
			i++
		case 5:
			inf.Brand = j
			i++
		case 6:
			if j == `null` {
				inf.QuantityInStock = 0
			} else {
				nm, err := strconv.Atoi(j)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 6: "))
				}
				inf.QuantityInStock = int64(nm)
			}
			i++
		case 7:
			if j == `null` {
				inf.Price = 0
			} else {
				tx := strings.Replace(j, " руб.", "", -1)
				tx = strings.Replace(tx, ",", "", -1)
				nm, err := strconv.ParseFloat(tx, 64)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 7: "))
				}
				inf.Price = nm
			}
			i++
		case 8:
			if j == `null` {
				inf.Discount = 0
			} else {
				tx := strings.Replace(j, "(", "", -1)
				tx = strings.Replace(tx, ")", "", -1)
				tx = strings.Replace(tx, "%", "", -1)
				nm, err := strconv.Atoi(tx)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 8: "))
				}
				inf.Discount = int64(nm)
			}
			i++
		case 9:
			if j == `null` {
				inf.OldPrice = 0
			} else {
				tx := strings.Replace(j, " руб.", "", -1)
				tx = strings.Replace(tx, ",", "", -1)
				nm, err := strconv.ParseFloat(tx, 64)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 9: "))
				}
				inf.OldPrice = nm
			}
			i++
		case 10:
			if j == `null` {
				inf.ACP = 0
			} else {
				tx := strings.Replace(j, " руб.", "", -1)
				tx = strings.Replace(tx, ",", "", -1)
				nm, err := strconv.ParseFloat(tx, 64)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 10: "))
				}
				inf.ACP = nm
			}
			i++
		case 11:
			if j == `null` {
				inf.LP = 0
			} else {
				tx := strings.Replace(j, " руб.", "", -1)
				tx = strings.Replace(tx, ",", "", -1)
				nm, err := strconv.ParseFloat(tx, 64)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 11: "))
				}
				inf.LP = nm
			}
			i++
		case 12:
			if j == `null` {
				inf.AmountOfSales = 0
			} else {
				nm, err := strconv.Atoi(j)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 12: "))
				}
				inf.AmountOfSales = int64(nm)
			}
			i++
		case 13:
			if j == `null` {
				inf.Revenue = 0
				i = 0
				tableRows = append(tableRows, inf)
				inf = Information{}
			} else {
				tx := strings.Replace(j, " руб.", "", -1)
				tx = strings.Replace(tx, ",", "", -1)
				nm, err := strconv.ParseFloat(tx, 64)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 13: "))
				}
				inf.Revenue = nm
				i = 0
				tableRows = append(tableRows, inf)
				inf = Information{}
			}
		}
	}
}

func parseNotFull()  {
	inf := Information{}
	for i, j := range data {
		switch i {
		case 0:
			inf.ProductName = j
			i++
		case 1:
			inf.Category = j
			i++
		case 2:
			if j == `null` {
				inf.NumberOfReviews = 0
			} else {
				tx := strings.Replace(j, "(", "", -1)
				tx = strings.Replace(tx, ")", "", -1)
				nm, err := strconv.Atoi(tx)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 2: "))
				}
				inf.NumberOfReviews = int64(nm)
			}
			i++
		case 3:
			if j == `null` {
				inf.SKU = 0
			} else {
				nm, err := strconv.Atoi(j)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 3: "))
				}
				inf.SKU = int64(nm)
			}
			i++
		case 4:
			inf.Seller = j
			i++
		case 5:
			inf.Brand = j
			i++
		case 6:
			if j == `null` {
				inf.QuantityInStock = 0
			} else {
				nm, err := strconv.Atoi(j)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 6: "))
				}
				inf.QuantityInStock = int64(nm)
			}
			i++
		case 7:
			if j == `null` {
				inf.Price = 0
			} else {
				tx := strings.Replace(j, " руб.", "", -1)
				tx = strings.Replace(tx, ",", "", -1)
				nm, err := strconv.ParseFloat(tx, 64)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 7: "))
				}
				inf.Price = nm
			}
			i++
		case 8:
			if j == `null` {
				inf.ACP = 0
			} else {
				tx := strings.Replace(j, " руб.", "", -1)
				tx = strings.Replace(tx, ",", "", -1)
				nm, err := strconv.ParseFloat(tx, 64)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 10: "))
				}
				inf.ACP = nm
			}
			i++
		case 9:
			if j == `null` {
				inf.LP = 0
			} else {
				tx := strings.Replace(j, " руб.", "", -1)
				tx = strings.Replace(tx, ",", "", -1)
				nm, err := strconv.ParseFloat(tx, 64)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 11: "))
				}
				inf.LP = nm
			}
			i++
		case 10:
			if j == `null` {
				inf.AmountOfSales = 0
			} else {
				nm, err := strconv.Atoi(j)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 12: "))
				}
				inf.AmountOfSales = int64(nm)
			}
			i++
		case 11:
			if j == `null` {
				inf.Revenue = 0
				i = 0
				tableRows = append(tableRows, inf)
				inf = Information{}
			} else {
				tx := strings.Replace(j, " руб.", "", -1)
				tx = strings.Replace(tx, ",", "", -1)
				nm, err := strconv.ParseFloat(tx, 64)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 13: "))
				}
				inf.Revenue = nm
				i = 0

				inf = Information{}
			}
		}
	}
	inf.Discount = 0
	inf.OldPrice = 0
	tableRows = append(tableRows, inf)
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

//func parseColumns(z *html.Tokenizer) {
//	//z := html.NewTokenizer(reader)
//	tt := z.Next()
//	if tt != html.StartTagToken {
//		return
//	}
//	tn, _ := z.TagName()
//	if string(tn) != `thead` {
//		return
//	}
//	headers := parseHeaders(z)
//	myLog.write(fmt.Sprintf(`%v`, headers))
//	tt = z.Next()
//	if tt != html.EndTagToken {
//		tn, _ := z.TagName()
//		if string(tn) != `thead` {
//			return
//		}
//	}
//}

//func parseHeaders(z *html.Tokenizer) (out []string) {
//	//z := html.NewTokenizer(in)
//	tt := z.Next()
//	if tt != html.StartTagToken {
//		return
//	}
//	tn, _ := z.TagName()
//	if string(tn) != `tr` {
//		return
//	}
//	tt = z.Next()
//	if tt != html.StartTagToken {
//		return
//	}
//loop:
//	for {
//		tt = z.Next()
//		myLog.write(fmt.Sprintf("swich %s", tt.String()))
//		switch tt {
//		case html.StartTagToken:
//			bs, _ := z.TagName()
//			if string(bs) != `th` {
//				continue loop
//			}
//			tt = z.Next()
//			if tt != html.TextToken {
//				myLog.write(string(z.Text()))
//				// TODO отработать вложенные елементыы
//				continue loop
//			}
//			out = append(out, string(z.Text()))
//		case html.EndTagToken:
//			if tn, _ = z.TagName(); string(tn) != `tr` {
//				return
//			}
//		}
//	}
//
//}

func tagname(t *html.Tokenizer) (out string) {
	bs, _ := t.TagName()
	return string(bs)
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
