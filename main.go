package main

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"io"
	"log"
	"os"
	"path/filepath"
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
loop:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			fmt.Println(z.Err())
			break loop
		case html.TextToken:
			if depth > 0 {
				// emitBytes should copy the []byte it receives,
				// if it doesn't process it immediately.
				myLog.write(fmt.Sprintf("TEXT %s", string(z.Text())))
			}
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
			if tt == html.StartTagToken {
				if string(tn) == `thead` {
					parseColumns(f)
				}
				//myLog.write(fmt.Sprintf(`OPEN %s %d`, string(tn), depth))
				depth++
			} else {
				depth--
				myLog.write(fmt.Sprintf(`CLOSE %s %d`, string(tn), depth))
			}
		case html.SelfClosingTagToken:
			tn, _ := z.TagName()
			myLog.write(fmt.Sprintf(`EMPTY %s %d`, string(tn), depth))
		default:
			tn, _ := z.TagName()
			myLog.write(fmt.Sprintf(`SKIP %s %d %s`, string(tn), depth, z.Token().Type))
		}
	}
}

func parseColumns(reader io.Reader) {
	z := html.NewTokenizer(reader)
	tt := z.Next()
	if tt != html.StartTagToken {
		return
	}
	tn, _ := z.TagName()
	if string(tn) != `thead` {
		return
	}
	tt = z.Next()
	if tt != html.StartTagToken {
		return
	}
	tn, _ = z.TagName()
	if string(tn) != `tr` {
		return
	}
	tt = z.Next()
	if tt != html.StartTagToken {
		return
	}
	//tn, _ = z.TagName()
	//if string(tn) != `th` {
	//	return
	//}
	
	//cols:=[]string{}
	//loop:
	for {
		tt = z.Next()
		switch tt {
		case html.StartTagToken:
			if tt == html.TextToken {
				myLog.write(fmt.Sprintf("TEXT %s", string(z.Text())))
			}
			continue
		case html.EndTagToken:
			if tn, _ = z.TagName(); string(tn) != `tr` {
				return
			}
			continue
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
