package main

import (
	"fmt"
	"log"
	"os"
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
	//ngx, err := getDelimLine(f)
	ngx, err := getCustDelimLine(f)
	if err != nil {
		return err
	}

	var ns []ngxStr
	level := 1
	for _, n := range ngx {
		ns = append(ns, parseStr(n, level)...)
	}

	fmt.Printf("%s", output(ns))
	return nil
}
