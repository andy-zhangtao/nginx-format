package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var tab = 4
const spaceLine = " "
const blankLine = "[+n+]"

type ngxString struct {
	data string
	tab  int
}

// ordrFormat 格式化输入的文件
// 顺序读取每一行，如果遇到:
// 1. 以'{'结尾。 下行tab+1
// 2. 以'}'结尾。 当前行需要回缩tab值,下行保持不变
// 3. 以';'结尾。下行tab值保持不变
// 4. 以非上面两类结尾，下行tab值保持不变
func ordrFormat(f *os.File) (ngs []ngxString, err error) {
	defer f.Close()

	var tabIdx = 0

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		switch theSuffixChar(text) {
		case ';':
			ngs = append(ngs, ngxString{
				data: text,
				tab:  tab * tabIdx,
			})
		case '{':
			ngs = append(ngs, ngxString{
				data: text,
				tab:  tab * tabIdx,
			})
			tabIdx++
		case '}':
			tabIdx--
			ngs = append(ngs, ngxString{
				data: text,
				tab:  tab * tabIdx,
			})
		default:
			ngs = append(ngs, ngxString{
				data: text,
				tab:  tab * tabIdx,
			})
		}
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}

func theSuffixChar(str string) byte {
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return 0
	}
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] != ' ' {
			return str[i]
		}
	}

	return 0
}

func output(ngxs []ngxString) string {
	str := ""
	for _, n := range ngxs {
		space := ""
		for i := 0; i < n.tab; i++ {
			space += spaceLine
		}

		if n.data == blankLine {
			str = fmt.Sprintf("%s%s\n", str, space)
		} else {
			str = fmt.Sprintf("%s%s%s\n", str, space, n.data)
		}

	}
	return str
}