package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var f *os.File
var err error

func main() {

	r := gin.Default()
	r.Use(cors.New(cors.Config{AllowAllOrigins: true}))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/transfer", func(c *gin.Context) {
		message := c.PostForm("nginx")

		ngs, err := format(message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Faile",
				"error":  err.Error(),
			})
			return
		}

		c.String(http.StatusOK, ngs)
	})
	r.Run()

}

func format(str string) (string, error) {
	var ngx []string

	//defer f.Close()
	//
	//scanner := bufio.NewScanner(f)
	//
	//for scanner.Scan() {
	//	text := strings.TrimSpace(scanner.Text())
	//	ngx = append(ngx, text)
	//}

	strList := strings.Split(str, "\n")
	for _, s := range strList {
		text := strings.TrimSpace(s)
		ngx = append(ngx, text)
	}

	ngs, err := ordrFormat(ngx)
	if err != nil {
		return "", err
	}

	//fmt.Printf("%s", output(ngs))
	return output(ngs), nil
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
