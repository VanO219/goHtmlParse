package main

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func goTokens() {
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

	//диапазон цен
	priceRangee := strings.Replace(priceRange(), ".", ",", -1)
	//выручка топ 10
	top10Revenuee := strings.Replace(top10Revenue(), ".", ",", -1)
	//среднее количество комментариев
	average10Commentss := strings.Replace(average10Comments(), ".", ",", -1)
	//упущенная выручка топ 30
	lost30Revenuee := strings.Replace(lost30Revenue(), ".", ",", -1)
	//оборот 50 поставщика
	vendor50 := strings.Replace(fmt.Sprintf("%.2f", res[50].Revenue), ".", ",", -1)

	myLog.write(fmt.Sprintf("диапазон цен: \t%s ", priceRangee))
	myLog.write(fmt.Sprintf("выручка топ 10: \t%s", top10Revenuee))
	myLog.write(fmt.Sprintf("среднее количество комментариев: \t%s", average10Commentss))
	myLog.write(fmt.Sprintf("упущенная выручка топ 30: \t%s", lost30Revenuee))
	myLog.write(fmt.Sprintf("оборот 50 поставщика: \t%s", vendor50))
}

func priceRange() string {
	minPrice := res[0].Price
	maxPrice := res[0].Price
	for i := 1; i < 50; i++ {
		if res[i].Price < minPrice {
			minPrice = res[i].Price
		}
		if res[i].Price > maxPrice {
			maxPrice = res[i].Price
		}
	}
	return fmt.Sprintf(`%.2f - %.2f`, minPrice, maxPrice)
}

func top10Revenue() string {
	r := 0.0
	for i := 0; i < 10; i++ {
		r += res[i].Revenue
	}
	return fmt.Sprintf(`%.2f`, r)
}

func average10Comments() string {
	var r float64
	for i := 0; i < 10; i++ {
		r += float64(res[i].NumberOfReviews)
	}
	v := r / 10
	re := math.Round(v)
	return fmt.Sprintf(`%.2f`, re)
}

func lost30Revenue() string {
	r := 0.0
	for i := 0; i < 30; i++ {
		r += res[i].LP
	}
	return fmt.Sprintf(`%.2f`, r)
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
	} else if len(data) > 0 && len(data) == 14 {
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
	//myLog.write(fmt.Sprintf("parseFull len: %d, \n %s", len(data), data))
	inf := Information{}
	for i, j := range data {
		switch i {
		case 0:
			inf.ProductName = j
		case 1:
			inf.Category = j
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
		case 4:
			inf.Seller = j
		case 5:
			inf.Brand = j
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
		case 8:
			if j == `null` {
				inf.Discount = 0
			} else {
				tx := strings.Replace(j, "(", "", -1)
				tx = strings.Replace(tx, ")", "", -1)
				tx = strings.Replace(tx, "%", "", -1)
				tx = strings.Replace(tx, "-", "", -1)
				nm, err := strconv.Atoi(tx)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "case 8: "))
				}
				inf.Discount = int64(nm)
			}
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
		case 13:
			if j == `null` {
				inf.Revenue = 0
				i = 0
				res = append(res, inf)
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
				res = append(res, inf)
				inf = Information{}
			}
		}
	}
}

func parseNotFull() {
	//myLog.write(fmt.Sprintf("parseNotFull len: %d, \n %s", len(data), data))
	inf := Information{}
	for i, j := range data {
		switch i {
		case 0:
			inf.ProductName = j
		case 1:
			inf.Category = j
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
		case 4:
			inf.Seller = j
		case 5:
			inf.Brand = j
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
		case 11:
			if j == `null` {
				inf.Revenue = 0
				i = 0
				res = append(res, inf)
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
			}
		}
	}
	inf.Discount = 0
	inf.OldPrice = 0
	res = append(res, inf)
}

func tagname(t *html.Tokenizer) (out string) {
	bs, _ := t.TagName()
	return string(bs)
}