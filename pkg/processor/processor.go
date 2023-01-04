package processor

import (
	"bufio"
	"fmt"
	"github.com/Bambelbl/go-counter/pkg/HadlerFile"
	"github.com/Bambelbl/go-counter/pkg/HadlerUrl"
	"log"
	"net/url"
	"os"
	"regexp"
	"sync"
	"sync/atomic"
)

// Тут по тз не совсем поняла, что имеется в виду под строкой "Go", т.к. если отделять пробелами, то значение
// получается больше 9, а если отделять \n, то просто 0
const stringForPattern = " Go "

type Source interface {
	Handler(stringPattern *regexp.Regexp) (count uint64, err error)
}

type Processor struct {
	ch            chan bool       // Буферизированный канал для отслеживания числа запущенных горутин
	wg            *sync.WaitGroup // Синхронизация горутин
	count         uint64          // Суммарное искомое число вхождений
	stringPattern *regexp.Regexp  // Искомая строка
}

func NewProcessor(k uint) (*Processor, error) {
	pattern, err := regexp.Compile(stringForPattern)
	if err != nil {
		return nil, err
	}
	return &Processor{
		ch:            make(chan bool, k),
		wg:            &sync.WaitGroup{},
		count:         0,
		stringPattern: pattern,
	}, nil
}

// ProcessOneSrc Подсчет вхождений в одном источнике (url/файл)
func (p *Processor) ProcessOneSrc(src string) {
	defer func() {
		<-p.ch
		p.wg.Done()
	}()
	var srcType Source
	_, isUrl := url.ParseRequestURI(src)
	if isUrl == nil {
		srcType = HadlerUrl.UrlSource{Url: src}
	} else {
		srcType = HadlerFile.FileSource{Filename: src}
	}
	count, err := CountOccurrence(srcType, p.stringPattern)
	if err != nil {
		log.Printf("error in counting: %s in src: %s", err.Error(), src)
	}
	atomic.AddUint64(&p.count, count)
	_, err = fmt.Printf("Count for %s: %d\n", src, count)
	if err != nil {
		log.Printf("error in printf: %s in src: %s", err.Error(), src)
	}
}

func CountOccurrence(src Source, stringPattern *regexp.Regexp) (uint64, error) {
	return src.Handler(stringPattern)
}

func (p *Processor) Process() error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		p.ch <- true
		p.wg.Add(1)
		go p.ProcessOneSrc(scanner.Text())
	}
	p.wg.Wait()

	if scanner.Err() != nil {
		return fmt.Errorf("error in scanner: %s", scanner.Err().Error())
	}
	_, err := fmt.Println("Total: ", p.count)
	if err != nil {
		return fmt.Errorf("error in print of result: %s", err.Error())
	}
	return nil
}
