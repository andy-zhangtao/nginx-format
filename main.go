package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var f *os.File
var err error

func main() {
	f, err = os.OpenFile("test/nginx-test.conf", os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = format(f)

	if err != nil {
		log.Fatal(err)
	}
}

func format(f *os.File) error {
	var ngx []string

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		ngx = append(ngx, text)
	}

	ngs, err := ordrFormat(ngx)
	if err != nil {
		return err
	}

	fmt.Printf("%s", output(ngs))
	return nil
}

//func format(f *os.File) error {
//	//ngx, err := getDelimLine(f)
//	ngx, err := getCustDelimLine(f)
//	if err != nil {
//		return err
//	}
//
//	var ns []ngxStr
//	level := 1
//	for _, n := range ngx {
//		ns = append(ns, parseStr(n, level)...)
//	}
//
//	fmt.Printf("%s", output(ns))
//	return nil
//}
