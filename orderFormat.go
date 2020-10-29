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
// 4. 以非上面两类结尾，如果上面一行是';'则下行tab值保持不变，同时标记下行需要tab+1. 否则tab值保持不变
// 5. 如果首字符是'#', 那么tab保存为当前值，不再进行调整
func ordrFormat(f *os.File) (ngs []ngxString, err error) {
	defer f.Close()

	var tabIdx = 0
	var makeUpTab = 0

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if thePrefixChar(text) == '#' {
			ngs = append(ngs, ngxString{
				data: text,
				tab:  tab * tabIdx,
			})
			continue
		}
		switch theSuffixChar(text) {
		case ';':
			if makeUpTab > 0 {
				ngs = append(ngs, ngxString{
					data: text,
					tab:  tab * (tabIdx + makeUpTab),
				})
				makeUpTab = 0
			} else {
				ngs = append(ngs, ngxString{
					data: text,
					tab:  tab * (tabIdx),
				})
			}

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
		case 0:
			ngs = append(ngs, ngxString{
				data: "",
				tab:  0,
			})
		default:
			if makeUpTab == 0 {
				ngs = append(ngs, ngxString{
					data: text,
					tab:  tab * tabIdx,
				})
				makeUpTab++
			} else {
				ngs = append(ngs, ngxString{
					data: text,
					tab:  tab * (tabIdx + makeUpTab),
				})
			}

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

func thePrefixChar(str string) byte {
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return 0
	}

	return str[0]
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
