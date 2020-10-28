package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var tab = 4

const blankLine = "[+n+]"
const spaceLine = " "
const delim = ';'
const colon = "\""
const leftCurlyBrackets = '{'
const rightCurlyBrackets = '}'

type ngxStr struct {
	data    string
	level   int
	convert bool
}

func output(ngxs []ngxStr) string {
	str := ""
	for _, n := range ngxs {
		space := ""
		for i := 0; i < (n.level-1)*tab; i++ {
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

func parseStr(str string, level int) []ngxStr {
	//	 判断是单语句还是多行语句
	strLines, single := isSingle(str)
	if single {
		return []ngxStr{
			{
				data:  strLines[0],
				level: level,
			},
		}
	}

	return parseMultipleLine(strLines, level)
}

func isComment(s string) bool {
	s = strings.TrimSpace(s)
	if s[0] == '#' {
		return true
	}
	return false
}

func parseMultipleLine(str []string, level int) []ngxStr {
	var ns []ngxStr

	for _, s := range str {

		s = strings.TrimSpace(s)

		if s == blankLine {
			ns = append(ns, ngxStr{
				data:    s,
				level:   level,
				convert: false,
			})
			continue
		}

		if s[0] == '}' {
			level--
		}

		ns = append(ns, ngxStr{
			data:    s,
			level:   level,
			convert: false,
		})

		if s[0] == '{' || s[len(s)-1] == '{' {
			level++
		}

	}

	return ns
}

func isSingle(s string) ([]string, bool) {

	str := strings.Split(s, "\n")
	if len(str) == 0 {
		return []string{blankLine}, true
	}

	if len(str) == 1 {
		return []string{strings.TrimSpace(str[0])}, true
	}

	var ms []string
	for _, _s := range str {
		if strings.TrimSpace(_s) == "" {
			ms = append(ms, blankLine)
		} else {
			ms = append(ms, strings.TrimSpace(_s))
		}

	}

	return ms, false
}

// getDelimLine 使用';'进行分割，不适用"... ;" ;这种场景
func getDelimLine(f *os.File) (str []string, err error) {

	r := bufio.NewReader(f)

	for {
		_str, err := r.ReadString(delim)
		if err != nil {
			if err == io.EOF {
				//fmt.Printf("%#v\n", str)
				return str, nil
			}

			return nil, err
		}

		str = append(str, _str)
	}

	return
}

// getCustDelimLine 自定义规则的字符串分割函数. 用来支持分割" ... ; ";这种场景
// 当遇到 xxxxx " ....;" ; 语句时，返回[]string{"xxxxx \" ....;\""}
//
func getCustDelimLine(f *os.File) (str []string, err error) {
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var temp []string //保存换行字符串

	for scanner.Scan() {

		text := strings.TrimSpace(scanner.Text())
		switch strings.Count(text, colon) {
		case 0:

			switch theSuffixChar(text) {
			case delim:
				if len(temp) == 0 {
					//	不包含" 并且以; 结尾，同时没有缓存数据
					str = append(str, text)
				} else {
					//	不包含" 并且以; 结尾，但有缓存数据，继续缓存
					temp = append(temp, text)
				}
			case leftCurlyBrackets:
				temp = append(temp, text)
			case rightCurlyBrackets:
				temp = append(temp, text)
				str = append(str, strings.Join(temp, "\n"))
				temp = temp[:0]
			default:
				// 不存在可以引起歧义的"，那么可以直接返回此行数据
				str = append(str, strings.Join(temp, "\n"))
				str = append(str, text)
				temp = temp[:0]
			}

		case 1:
			if theSuffixChar(text) != delim {
				//	只有一个" ,如果不是以;结尾，那么就需要暂存起来
				temp = append(temp, text)
			} else {
				if len(temp) > 0 {
					//	只有一个" ,但以;结尾，同时已经存在暂存数据,那么就需要作为一条语句释放出来
					temp = append(temp, text)
					str = append(str, strings.Join(temp, "\n"))
					temp = temp[:0]
				} else {
					//	只有一个" ,但以;结尾. 当前没有暂存数据，那么就需要暂存起来
					temp = append(temp, text)
				}
			}
		case 2:
			//	包含 " ",并且以;结尾
			if theSuffixChar(text) == delim {
				// 如果有缓存数据，则说明和缓存数据同属于一个代码块
				if len(temp) == 0 {
					str = append(str, text)
				} else {
					temp = append(temp, text)
				}

			} else {
				// 不是以;结尾，所以暂存一下
				temp = append(temp, text)
			}
		default:
			//同一行存在多于2个"的情况
			if strings.Count(text, colon)%2 == 0 {
				//	包含相互配对的" ",并且以;结尾
				if theSuffixChar(text) == delim {
					if len(temp) == 0 {
						str = append(str, text)
					} else {
						temp = append(temp, text)
						//str = append(str, strings.Join(temp, "\n"))
						//temp = temp[:0]
					}
				} else {
					temp = append(temp, text)
				}
			} else {
				//	因为是"数量为单数，所以需要暂存起来. 但需要考虑是否会和上一行匹配成一行
				if theSuffixChar(text) == delim {
					temp = append(temp, text)
					str = append(str, strings.Join(temp, "\n"))
					temp = temp[:0]
				} else {
					temp = append(temp, text)
				}
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
