package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"ethan.com/reptile/util"
)

var reg = regexp.MustCompile(`^test[1-9]\d*.txt$`)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("命令行参数要等于2个")
		return
	}

	os.Mkdir("./comments_"+os.Args[1], 0777)

	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal("文件夹错误")
		return
	}

	var testFlag int
	var contents string
	testMap := make(map[int]string)

	for _, f := range files {
		if reg.MatchString(f.Name()) {
			// fmt.Println(f.Name())
			testFlag++
			testMap[testFlag] = f.Name()
		}

	}
	fmt.Println("txt数：", testFlag)
	// fc := math.Ceil(float64(testFlag) / 100.0)
	ff := math.Floor(float64(testFlag) / 100.0)
	// fmt.Println("int(ff)", int(ff))

	done := make(chan bool, int(ff))

	if testFlag < 100 {

		file, err := os.OpenFile(
			"./comments_"+os.Args[1]+"/test_1-"+strconv.Itoa(testFlag)+".txt",
			os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
			0666,
		)
		defer file.Close()
		if err != nil {
			log.Fatal("os.OpenFile count.txt", err)
		}

		for i := 1; i <= testFlag; i++ {
			contents = commentplus(os.Args[1]+"/"+testMap[i]) + contents
		}
		util.WriteTo(file, contents)

	} else if testFlag >= 100 {
		for i := 0; i < int(ff); i++ {
			go commentsToFile(i, testMap, os.Args[1], done)
		}

		if int(ff)*100 < testFlag {
			file, err := os.OpenFile(
				"./comments_"+os.Args[1]+"/test_"+strconv.Itoa(int(ff)*100)+"-"+strconv.Itoa(testFlag)+".txt",
				os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
				0666,
			)
			defer file.Close()
			if err != nil {
				log.Fatal("os.OpenFile count.txt", err)
			}

			for i := int(ff) * 100; i <= testFlag; i++ {
				contents = commentplus(os.Args[1]+"/"+testMap[i]) + contents
			}
			util.WriteTo(file, contents)
		}
	}

	for i := 0; i < int(ff); i++ {
		<-done
	}
}

func commentsToFile(i int, testMap map[int]string, dir string, done chan bool) {
	var contents string

	file, err := os.OpenFile(
		"./comments_"+os.Args[1]+"/test_"+strconv.Itoa(i*100+1)+"-"+strconv.Itoa((i+1)*100)+".txt",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	defer file.Close()
	if err != nil {
		log.Fatal("os.OpenFile count.txt", err)
	}

	for j := i*100 + 1; j <= (i+1)*100; j++ {
		contents = commentplus(dir+"/"+testMap[j]) + contents
	}
	util.WriteTo(file, contents)

	done <- true
}
func commentplus(commentURL string) string {
	var contents string
	file, err := os.Open(commentURL)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		str := scanner.Text()
		strss := strings.Split(str, "：")

		if strss[0] == "评论" {
			s := string([]byte(strss[1])[:6])
			if s == "回复" {
				if len(strings.Split(strss[1], ":")) < 2 {
					contents = contents + strss[1]
				} else {
					content := strings.Split(strss[1], ":")[1]
					contents = contents + content
				}
			} else {
				contents = contents + strss[1]
			}
		}
	}
	return contents
}
