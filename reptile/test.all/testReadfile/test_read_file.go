package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// AGE 年龄段
var AGE = []string{
	"70后",
	"80后",
	"90后",
	"00后",
	"10后",
}

// InitAGE 初始化age
func InitAGE() (ageMap map[string]int) {
	ageMap = make(map[string]int)
	for i := 0; i < len(AGE); i++ {
		ageMap[AGE[i]] = 0
	}
	return
}

// REGION 省份
var REGION = []string{
	"安徽",
	"北京",
	"重庆",
	"广东",
	"广西",
	"贵州",
	"海南",
	"河北",
	"黑龙江",
	"河南",
	"湖北",
	"湖南",
	"内蒙古",
	"江苏",
	"江西",
	"吉林",
	"辽宁",
	"宁夏",
	"青海",
	"山西",
	"山东",
	"上海",
	"四川",
	"天津",
	"西藏",
	"新疆",
	"云南",
	"浙江",
	"陕西",
	"台湾",
	"香港",
	"澳门",
	"海外",
	"其他",
	"福建",
	"甘肃",
}

// RetMap 初始化map
func RetMap() (regionMap map[string]int) {
	regionMap = make(map[string]int)
	for i := 0; i < len(REGION); i++ {
		regionMap[REGION[i]] = 0
	}
	return
}

type Count struct {
	Total           int
	Men             int
	Women           int
	Regions         map[string]int
	Generation      map[string]int
	GenerationCount int
}

func main() {
	// file, err := os.Open("../../test_GyHeQ23UA_0/count_813.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// counts := &Count{}
	// counts.Regions = make(map[string]int)
	// counts.Generation = make(map[string]int)

	// scanner := bufio.NewScanner(file)
	// scanner.Split(bufio.ScanLines)

	// //是否有下一行
	// for scanner.Scan() {
	// 	// fmt.Println("read string:", scanner.Text())
	// 	str := scanner.Text()
	// 	strss := strings.Split(str, "：")
	// 	switch strss[0] {
	// 	case "----男性人数":
	// 		persons := strings.Split(str, ",")
	// 		if len(persons) != 3 {
	// 			log.Fatal("Get count persons failed!")
	// 			return
	// 		}
	// 		men := strings.Split(persons[0], "：")[1]
	// 		if men != "" {

	// 			menInt, err := strconv.Atoi(men)
	// 			if err != nil {
	// 				log.Fatal("men to menInt fail.", err)
	// 				return
	// 			}
	// 			counts.Men = menInt
	// 		}
	// 		women := strings.Split(persons[1], "：")[1]
	// 		if women != "" {
	// 			womenInt, err := strconv.Atoi(women)
	// 			if err != nil {
	// 				log.Fatal("women to womenInt fail.", err)
	// 				return
	// 			}
	// 			counts.Women = womenInt

	// 		}
	// 		total := strings.Split(persons[2], "：")[1]
	// 		if total != "" {
	// 			totalInt, err := strconv.Atoi(total)
	// 			if err != nil {
	// 				log.Fatal("total to totalInt fail.", err)
	// 				return
	// 			}
	// 			counts.Total = totalInt
	// 		}

	// 	case "省份":
	// 		regions := strings.Split(str, ",")
	// 		if len(regions) != 4 {
	// 			log.Fatal("Get count regions failed.")
	// 			return
	// 		}
	// 		rname := strings.Split(regions[0], "：")[1]
	// 		rcount := strings.Split(regions[1], "：")[1]
	// 		if rname != "" && rcount != "" {
	// 			rcountInt, err := strconv.Atoi(rcount)
	// 			if err != nil {
	// 				log.Fatal("rcount to rcountInt fail.", err)
	// 				return
	// 			}
	// 			counts.Regions[rname] = rcountInt
	// 		}
	// 	case "人群":
	// 		generation := strings.Split(str, ",")
	// 		if len(generation) != 4 {
	// 			log.Fatal("Get count regions failed.")
	// 			return
	// 		}
	// 		gccount := strings.Split(generation[2], "：")[1]
	// 		if gccount != "" {
	// 			gccountInt, err := strconv.Atoi(gccount)
	// 			if err != nil {
	// 				log.Fatal("gccount to gccountInt fail.", err)
	// 				return

	// 			}
	// 			counts.GenerationCount = gccountInt
	// 		}
	// 		gname := strings.Split(generation[0], "：")[1]
	// 		gcount := strings.Split(generation[1], "：")[1]
	// 		if gname != "" && gcount != "" {
	// 			gcountInt, err := strconv.Atoi(gcount)
	// 			if err != nil {
	// 				log.Fatal("gcount to gcountInt fail.", err)
	// 				return
	// 			}
	// 			counts.Generation[gname] = gcountInt
	// 		}
	// 	}
	// }

	// fmt.Println(counts)
	totalCount := &Count{}
	totalCount.Regions = make(map[string]int)
	totalCount.Generation = make(map[string]int)
	countURLs := []string{
		"../../test_GyHeQ23UA_0/count_813.txt",
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
		for l := 0; l < len(REGION); l++ {
			totalCount.Regions[REGION[l]] = cd.Regions[REGION[l]] + totalCount.Regions[REGION[l]]
		}
		for m := 0; m < len(AGE); m++ {
			totalCount.Generation[AGE[m]] = cd.Generation[AGE[m]] + totalCount.Generation[AGE[m]]
		}
	}

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
