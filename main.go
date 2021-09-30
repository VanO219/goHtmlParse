package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

var (
	myLog *infoLogger
)

func main() {
	f, _ := os.Open(`source.html`)
	defer f.Close()
	myLog = newInfoLogger("./output.txt")
	defer myLog.closeFile()

	z := html.NewTokenizer(f)

	depth := 0
	//loop:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			fmt.Println(z.Err())
			return
		//case html.TextToken:
		//	if depth > 0 {
		//		// emitBytes should copy the []byte it receives,
		//		// if it doesn't process it immediately.
		//		myLog.write(fmt.Sprintf("TEXT %s", string(z.Text())))
		//	}
		case html.StartTagToken, html.EndTagToken:
			//tn, _ := z.TagName()
			if tt == html.StartTagToken {
				//myLog.write(tagname(z))
				//myLog.write(fmt.Sprintf(`OPEN %s %d`, string(tn), depth))
				//log.Println(fmt.Sprintf(`OPEN %s %d`, string(tn), depth))
				fmt.Println(z.Text())
				if tagname(z) == `table` {
					parseColumns(z)
				}
				depth++
			} else {
				depth--
				//myLog.write(fmt.Sprintf(`CLOSE %s %d`, string(tn), depth))
				continue
			}
		case html.SelfClosingTagToken:
			//tn, _ := z.TagName()
			//myLog.write(fmt.Sprintf(`EMPTY %s %d`, string(tn), depth))
			continue
		default:
			//tn, _ := z.TagName()
			//myLog.write(fmt.Sprintf(`SKIP %s %d %s`, string(tn), depth, z.Token().Type))
			continue
		}
	}
}
func tagname(t *html.Tokenizer) (out string) {
	bs, _ := t.TagName()
	return string(bs)
}
func parseColumns(z *html.Tokenizer) {
	//z := html.NewTokenizer(reader)
	tt := z.Next()
	if tt != html.StartTagToken {
		return
	}
	tn, _ := z.TagName()
	if string(tn) != `thead` {
		return
	}
	headers := parseHeaders(z)
	myLog.write(fmt.Sprintf(`%v`, headers))
	tt = z.Next()
	if tt != html.EndTagToken {
		tn, _ := z.TagName()
		if string(tn) != `thead` {
			return
		}
	}
}

func parseHeaders(z *html.Tokenizer) (out []string) {
	//z := html.NewTokenizer(in)
	tt := z.Next()
	if tt != html.StartTagToken {
		return
	}
	tn, _ := z.TagName()
	if string(tn) != `tr` {
		return
	}
	tt = z.Next()
	if tt != html.StartTagToken {
		return
	}
loop:
	for {
		tt = z.Next()
		myLog.write(fmt.Sprintf("swich %s", tt.String()))
		switch tt {
		case html.StartTagToken:
			bs, _ := z.TagName()
			if string(bs) != `th` {
				continue loop
			}
			tt = z.Next()
			if tt != html.TextToken {
				myLog.write(string(z.Text()))
				// TODO отработать вложенные елементыы
				continue loop
			}
			out = append(out, string(z.Text()))
		case html.EndTagToken:
			if tn, _ = z.TagName(); string(tn) != `tr` {
				return
			}
		}
	}

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
