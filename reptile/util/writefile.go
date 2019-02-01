package util

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// DIYFileName 自定义文件夹和文件
type DIYFileName struct {
	PageID string
	DirID  string
}

// InitFile 初始化文件夹及文件
func InitFile(dfn *DIYFileName) *os.File {

	file, err := os.OpenFile(
		"./test_"+dfn.DirID+"/test"+dfn.PageID+".txt",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal("os.OpenFile", err)
	}
	return file
}

// WriteTo 写入文件
func WriteTo(file *os.File, val string) (err error) {
	byteSlice := []byte(val)
	_, err = file.Write(byteSlice)
	return
}

// Decimal float保留两位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
