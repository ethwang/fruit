package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"ethan.com/reptile/module"
	"ethan.com/reptile/util"
)

var reg = regexp.MustCompile(`^count_[1-9]\d*.txt$`)

// Count 终极统计
type Count struct {
	Total           int
	Men             int
	Women           int
	Regions         map[string]int
	Generation      map[string]int
	GenerationCount int
}

func main() {
	if len(os.Args) < 4 {
		log.Fatal("命令行参数要大于等于4个")
		return
	}
	countURLs := []string{}

	if os.Args[2] == "0" {
		for i := 3; i < len(os.Args); i++ {
			countURLs = append(countURLs, os.Args[i])
		}
	} else if os.Args[2] == "1" {
		files, err := ioutil.ReadDir(os.Args[3])
		if err != nil {
			log.Fatal("文件夹错误")
			return
		}

		for _, f := range files {
			fmt.Println(f.Name())
			if reg.MatchString(f.Name()) {
				countURLs = append(countURLs, os.Args[3]+"/"+f.Name())
			}
		}
	} else {
		log.Fatal("参数错误")
		return
	}

	totalCount := &Count{}
	totalCount.Regions = make(map[string]int)
	totalCount.Generation = make(map[string]int)

	// create and open file
	file, err := os.OpenFile(
		"./count_"+os.Args[1]+".txt",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	defer file.Close()
	if err != nil {
		log.Fatal("os.OpenFile count.txt", err)
	}

	s := len(countURLs)
	done := make(chan *Count, s)
	for i := 0; i < s; i++ {
		go ccount(countURLs[i], done)
	}

	for j := 0; j < s; j++ {
		cd := <-done
		totalCount.Total = cd.Total + totalCount.Total
		totalCount.Men = cd.Men + totalCount.Men
		totalCount.Women = cd.Women + totalCount.Women
		totalCount.GenerationCount = cd.GenerationCount + totalCount.GenerationCount
		for l := 0; l < len(module.REGION); l++ {
			totalCount.Regions[module.REGION[l]] = cd.Regions[module.REGION[l]] + totalCount.Regions[module.REGION[l]]
		}
		for m := 0; m < len(module.AGE); m++ {
			totalCount.Generation[module.AGE[m]] = cd.Generation[module.AGE[m]] + totalCount.Generation[module.AGE[m]]
		}
	}
	util.WriteTo(file, "----评论总数："+strconv.Itoa(totalCount.Total)+"----\r\n\r\n")
	util.WriteTo(file, "----男性人数："+strconv.Itoa(totalCount.Men)+", 女性人数："+strconv.Itoa(totalCount.Women)+", 总人数："+strconv.Itoa(totalCount.Total)+"\r\n")
	util.WriteTo(file, "男性占比："+strconv.FormatFloat(util.Decimal(float64(totalCount.Men)/float64(totalCount.Total)), 'f', -1, 64)+"\r\n")
	util.WriteTo(file, "女性占比："+strconv.FormatFloat(util.Decimal(float64(totalCount.Women)/float64(totalCount.Total)), 'f', -1, 64)+"\r\n\r\n\r\n")

	for j := 0; j < len(module.REGION); j++ {
		util.WriteTo(file, "省份："+module.REGION[j]+", 个数："+strconv.Itoa(totalCount.Regions[module.REGION[j]])+", 总数："+strconv.Itoa(totalCount.Total)+", 占比："+strconv.FormatFloat(util.Decimal(float64(totalCount.Regions[module.REGION[j]])/float64(totalCount.Total)), 'f', -1, 64)+"\r\n\r\n")
	}
	util.WriteTo(file, "\r\n")
	for j := 0; j < len(module.AGE); j++ {
		util.WriteTo(file, "人群："+module.AGE[j]+", 人数："+strconv.Itoa(totalCount.Generation[module.AGE[j]])+", 总数："+strconv.Itoa(totalCount.GenerationCount)+", 占比："+strconv.FormatFloat(util.Decimal(float64(totalCount.Generation[module.AGE[j]])/float64(totalCount.GenerationCount)), 'f', -1, 64)+"\r\n\r\n")

	}
	os.Chmod(file.Name(), 0444)
	fmt.Println(totalCount)
}
func ccount(countURL string, countDone chan *Count) {
	file, err := os.Open(countURL)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	counts := &Count{}
	counts.Regions = make(map[string]int)
	counts.Generation = make(map[string]int)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	//是否有下一行
	for scanner.Scan() {
		// fmt.Println("read string:", scanner.Text())
		str := scanner.Text()
		strss := strings.Split(str, "：")
		switch strss[0] {
		case "----男性人数":
			persons := strings.Split(str, ",")
			if len(persons) != 3 {
				log.Fatal("Get count persons failed!")
				return
			}
			men := strings.Split(persons[0], "：")[1]
			if men != "" {

				menInt, err := strconv.Atoi(men)
				if err != nil {
					log.Fatal("men to menInt fail.", err)
					return
				}
				counts.Men = menInt
			}
			women := strings.Split(persons[1], "：")[1]
			if women != "" {
				womenInt, err := strconv.Atoi(women)
				if err != nil {
					log.Fatal("women to womenInt fail.", err)
					return
				}
				counts.Women = womenInt

			}
			total := strings.Split(persons[2], "：")[1]
			if total != "" {
				totalInt, err := strconv.Atoi(total)
				if err != nil {
					log.Fatal("total to totalInt fail.", err)
					return
				}
				counts.Total = totalInt
			}

		case "省份":
			regions := strings.Split(str, ",")
			if len(regions) != 4 {
				log.Fatal("Get count regions failed.")
				return
			}
			rname := strings.Split(regions[0], "：")[1]
			rcount := strings.Split(regions[1], "：")[1]
			if rname != "" && rcount != "" {
				rcountInt, err := strconv.Atoi(rcount)
				if err != nil {
					log.Fatal("rcount to rcountInt fail.", err)
					return
				}
				counts.Regions[rname] = rcountInt
			}
		case "人群":
			generation := strings.Split(str, ",")
			if len(generation) != 4 {
				log.Fatal("Get count regions failed.")
				return
			}
			gccount := strings.Split(generation[2], "：")[1]
			if gccount != "" {
				gccountInt, err := strconv.Atoi(gccount)
				if err != nil {
					log.Fatal("gccount to gccountInt fail.", err)
					return

				}
				counts.GenerationCount = gccountInt
			}
			gname := strings.Split(generation[0], "：")[1]
			gcount := strings.Split(generation[1], "：")[1]
			if gname != "" && gcount != "" {
				gcountInt, err := strconv.Atoi(gcount)
				if err != nil {
					log.Fatal("gcount to gcountInt fail.", err)
					return
				}
				counts.Generation[gname] = gcountInt
			}
		}
	}
	countDone <- counts
}
