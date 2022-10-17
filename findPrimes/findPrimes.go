package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
Требуется разработать консольную утилиту, которая принимает несколько диапазонов чисел,
имя файла для вывода, таймаут ограничивающий исполнение команды по времени.
Утилита находит в заданных дипазонах все простые числа и выводит их в файл.
Например для диапазона 11:20 простые числа это 11, 13, 17, 19
*/

// arrayRange custom type implemting interface flag.Value
// need for multiple flags (--range)
type arrayRange []string

// this method is stub
func (a *arrayRange) String() string {
	return ""
}

// this work method for append multiple flags
func (a *arrayRange) Set(value string) error {
	*a = append(*a, value)
	return nil
}

// global variables
var (
	fileName   string
	timeout    string
	arrayFlags arrayRange
)

// Parse flags
func init() {
	flag.StringVar(&fileName, "file", "testfile.txt", "use this flag for naming exist file")
	flag.StringVar(&timeout, "timeout", "10s", "setup timeout")
	flag.Var(&arrayFlags, "range", "use multiple flags")
	flag.Parse()
}

func main() {
	duration, err := time.ParseDuration(timeout) // timeout duration that breaking gourutines
	if err != nil {
		log.Fatalln("Cannot convert string to time duration")
	}

	chanInt := make(chan []int64, len(arrayFlags)) // interaction channel

	var wg sync.WaitGroup // wait group for gourutines that calls from main routine
	wg.Add(2)

	contextTimeout, cancelFunc := context.WithTimeout(context.Background(), duration) // context with timeout
	defer cancelFunc()

	go func(waitGroup *sync.WaitGroup) {
		var wg sync.WaitGroup // wait group for iteration routines
		wg.Add(len(arrayFlags))
		for _, val := range arrayFlags { // routines processing flags value and send in channel
			go func(wg *sync.WaitGroup, val string) {
				startInt, endInt := convertStringToInt(val)
				chanInt <- findPrimes(startInt, endInt)
				wg.Done()
			}(&wg, val)
		}
		wg.Wait()
		close(chanInt)
		waitGroup.Done()
	}(&wg)

	go func() {
		defer wg.Done()
		var result string
	loop1:
		for {
			select {
			case <-contextTimeout.Done(): // if received ctx from ctx channel break loop
				break loop1
			case val, ok := <-chanInt:
				if !ok {
					break loop1
				} else if ok {
					result += concatData(val) // concatenation strings in result
				}
			}
		}
		fileSave(result) // save result in file
	}()

	wg.Wait()
}

func convertStringToInt(val string) (int64, int64) {
	str := strings.Split(val, ":")
	startInt, err := strconv.Atoi(str[0])
	if err != nil {
		log.Fatalln("failed to convert string to int")
	}
	endInt, err := strconv.Atoi(str[1])
	if err != nil {
		log.Fatalln("failed to convert string to int")
	}
	return int64(startInt), int64(endInt)
}

func findPrimes(start, end int64) []int64 {
	var result []int64
	for i := start; i <= end; i++ {
		if big.NewInt(i).ProbablyPrime(0) {
			result = append(result, i)
		}
	}
	return result
}

func concatData(primesArray []int64) string {
	result := "["
	for i, val := range primesArray {
		if i == len(primesArray)-1 {
			result += strconv.Itoa(int(val))
		} else {
			result += strconv.Itoa(int(val)) + " "
		}
	}
	result += "]"

	return result
}

func fileSave(result string) error {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("Cannot create file")
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	fmt.Fprint(w, result)

	return nil
}
