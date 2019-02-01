package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"ethan.com/reptile/module"
	"ethan.com/reptile/util"
)

/*
	数据源头URL:https://m.weibo.cn/u/5644764907?uid=5644764907&luicode=10000011&lfid=100103type%3D1%26q%3D%E6%9D%A8%E8%B6%85%E8%B6%8A
	---> commentsID ---> userID
*/

// const count = 10
// const commentsID = "4296473875115737"
// const commentsCookie = "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548294434; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQzVvWiyamobldPXiJG-H3iDKIR0-TvlPiY091BPCFR24.; SUB=_2A25xJfzdDeRhGeVH6VsW9SbEzTyIHXVS6YSVrDV6PUJbktAKLXLZkW1NT3AvulTG2_x1nzGrufDPoHMAWfIb9JuX; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5K-hUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0Ry6R7PJoeAvuj; SSOLoginState=1545702541; WEIBOCN_FROM=1110006030; MLOGIN=1"
// const userCookie = "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548294434; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQzVvWiyamobldPXiJG-H3iDKIR0-TvlPiY091BPCFR24.; SUB=_2A25xJfzdDeRhGeVH6VsW9SbEzTyIHXVS6YSVrDV6PUJbktAKLXLZkW1NT3AvulTG2_x1nzGrufDPoHMAWfIb9JuX; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5K-hUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0Ry6R7PJoeAvuj; SSOLoginState=1545702541"

func main() {
	var commentCount, userMenCount, userWomenCount, totalRegion, totalAge int

	count := flag.Int("count", 10, "爬取页数")
	commentsID := flag.String("commentsID", "4296473875115737", "微博评论ID")
	commentsCookie := flag.String("commentsCookie", "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548294434; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQzVvWiyamobldPXiJG-H3iDKIR0-TvlPiY091BPCFR24.; SUB=_2A25xJfzdDeRhGeVH6VsW9SbEzTyIHXVS6YSVrDV6PUJbktAKLXLZkW1NT3AvulTG2_x1nzGrufDPoHMAWfIb9JuX; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5K-hUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0Ry6R7PJoeAvuj; SSOLoginState=1545702541; WEIBOCN_FROM=1110006030; MLOGIN=1", "微博评论Cookie")
	userCookie := flag.String("userCookie", "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548294434; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQzVvWiyamobldPXiJG-H3iDKIR0-TvlPiY091BPCFR24.; SUB=_2A25xJfzdDeRhGeVH6VsW9SbEzTyIHXVS6YSVrDV6PUJbktAKLXLZkW1NT3AvulTG2_x1nzGrufDPoHMAWfIb9JuX; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5K-hUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0Ry6R7PJoeAvuj; SSOLoginState=1545702541", "微博用户Cookie")
	flag.Parse()

	err := os.RemoveAll("./test_" + *commentsID + "")
	if err != nil {
		log.Fatal("os.Remove", err)
	}
	regionMap := module.RetMap()
	ageMap := module.InitAGE()

	reqDone := make(chan *module.Statistics, *count)

	os.Mkdir("./test_"+*commentsID+"/", 0777)

	for i := 1; i <= *count; i++ {
		go module.HTTPGetData(i, *commentsID, *commentsCookie, *userCookie, reqDone)
	}
	for i := 1; i <= *count; i++ {
		st := <-reqDone
		commentCount = st.CommentCount + commentCount
		userMenCount = st.UserMenCount + userMenCount
		userWomenCount = st.UserWomenCount + userWomenCount
		for g := 0; g < len(module.REGION); g++ {
			if st.UserRegionCount[module.REGION[g]] != 0 {
				regionMap[module.REGION[g]] = st.UserRegionCount[module.REGION[g]] + regionMap[module.REGION[g]]
			}
		}

		for a := 0; a < len(module.AGE); a++ {
			if st.UserBrithCount[module.AGE[a]] != 0 {
				ageMap[module.AGE[a]] = st.UserBrithCount[module.AGE[a]] + ageMap[module.AGE[a]]
			}
		}
	}

	for j := 0; j < len(module.REGION); j++ {
		totalRegion = regionMap[module.REGION[j]] + totalRegion
	}
	for j := 0; j < len(module.AGE); j++ {
		totalAge = ageMap[module.AGE[j]] + totalAge
	}

	file, err := os.OpenFile(
		"./test_"+*commentsID+"/count.txt",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	defer file.Close()
	if err != nil {
		log.Fatal("os.OpenFile count.txt", err)
	}

	util.WriteTo(file, "----评论总数："+strconv.Itoa(commentCount)+"----\r\n\r\n")
	util.WriteTo(file, "----男性人数："+strconv.Itoa(userMenCount)+" --女性总数："+strconv.Itoa(userWomenCount)+" --男女总数："+strconv.Itoa(userWomenCount+userMenCount)+"----\r\n")
	util.WriteTo(file, "男性占比："+strconv.FormatFloat(util.Decimal(float64(userMenCount)/float64(userMenCount+userWomenCount)), 'f', -1, 64)+"\r\n")
	util.WriteTo(file, "女性占比："+strconv.FormatFloat(util.Decimal(float64(userWomenCount)/float64(userMenCount+userWomenCount)), 'f', -1, 64)+"\r\n\r\n\r\n")

	fmt.Println("----commentCount", commentCount, "-----")
	fmt.Println("-----userMenCount", userMenCount, "-----totalSex---", userMenCount+userWomenCount, "----")
	fmt.Println("-----userWomenCount", userWomenCount, "-----totalSex---", userMenCount+userWomenCount, "----")
	fmt.Printf("男性占比：%.2f \n", float64(userMenCount)/float64(userMenCount+userWomenCount))
	fmt.Printf("女性占比：%.2f \n", float64(userWomenCount)/float64(userMenCount+userWomenCount))

	fmt.Println("----userRegionCount", regionMap, "--------totalRegion-----", totalRegion, "-----")
	for j := 0; j < len(module.REGION); j++ {
		fmt.Printf("省份：%s,占比：%.2f \n", module.REGION[j], float64(regionMap[module.REGION[j]])/float64(totalRegion))
		util.WriteTo(file, "省份："+module.REGION[j]+" ,个数："+strconv.Itoa(regionMap[module.REGION[j]])+" ,总数："+strconv.Itoa(totalRegion)+" ,占比："+strconv.FormatFloat(util.Decimal(float64(regionMap[module.REGION[j]])/float64(totalRegion)), 'f', -1, 64)+"\r\n\r\n")
	}
	util.WriteTo(file, "\r\n")
	fmt.Println("----userBrithCount", ageMap, "------totalAge------", totalAge, "------")
	for j := 0; j < len(module.AGE); j++ {
		fmt.Printf("人群：%s,占比：%.2f \n", module.AGE[j], float64(ageMap[module.AGE[j]])/float64(totalAge))
		util.WriteTo(file, "人群："+module.AGE[j]+" ,人数："+strconv.Itoa(ageMap[module.AGE[j]])+" ,总数："+strconv.Itoa(totalAge)+" ,占比："+strconv.FormatFloat(util.Decimal(float64(ageMap[module.AGE[j]])/float64(totalAge)), 'f', -1, 64)+"\r\n\r\n")

	}

	os.Chmod(file.Name(), 0444)
}

// func Decimal(value float64) float64 {
// 	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
// 	return value
// }

// func httpGetCommentss(j int) {
// 	strPage := strconv.Itoa(j)

// 	strURL := "https://m.weibo.cn/api/comments/show?id=4296473875115737&page=" + strPage

// 	commentsCookie := "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548294434; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQzVvWiyamobldPXiJG-H3iDKIR0-TvlPiY091BPCFR24.; SUB=_2A25xJfzdDeRhGeVH6VsW9SbEzTyIHXVS6YSVrDV6PUJbktAKLXLZkW1NT3AvulTG2_x1nzGrufDPoHMAWfIb9JuX; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5K-hUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0Ry6R7PJoeAvuj; SSOLoginState=1545702541; WEIBOCN_FROM=1110006030; MLOGIN=1"

// 	resp := GetResp(strURL, commentsCookie)
// 	defer resp.Body.Close()

// 	newFile, err := os.Create("test" + strPage + ".txt")
// 	if err != nil {
// 		log.Fatal("os.Create", err)
// 	}
// 	defer newFile.Close()

// 	file, err := os.OpenFile(
// 		"test"+strPage+".txt",
// 		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
// 		0666,
// 	)
// 	if err != nil {
// 		log.Fatal("os.OpenFile", err)
// 	}
// 	defer file.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalln("ioutil.ReadAll fatal:", err)
// 	}
// 	var dat map[string]interface{}
// 	if err := json.Unmarshal(body, &dat); err != nil {
// 		log.Fatalln("json.Unmarshal fatal:", err)
// 	}

// 	for keyA, valA := range dat {

// 		if keyA == "data" {
// 			for keyB, valB := range valA.(map[string]interface{}) {
// 				if keyB == "data" {
// 					interfaceArray := valB.([]interface{})
// 					for i := 0; i < len(interfaceArray); i++ {
// 						for keyC, valC := range interfaceArray[i].(map[string]interface{}) {
// 							if keyC == "user" {
// 								for keyD, valD := range valC.(map[string]interface{}) {
// 									if keyD == "id" {
// 										_ = httpGetUser(valD.(float64))
// 										// fmt.Println(userInfo.region)
// 									}

// 									if keyD == "screen_name" {
// 										err = writeTo(file, keyD+": "+valD.(string)+" ")
// 										if err != nil {
// 											log.Fatal("writeTo D", err)
// 										}
// 									}
// 								}
// 							}

// 							if keyC == "text" {
// 								re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
// 								src := re.ReplaceAllString(valC.(string), "")

// 								err = writeTo(file, keyC+": "+src+" ")
// 								if err != nil {
// 									log.Fatal("writeTo C", err)
// 								}
// 							}
// 						}

// 						err = writeTo(file, "\n")
// 						if err != nil {
// 							log.Fatal("writeTo n", err)
// 						}
// 					}
// 				}
// 				if keyB == "hot_data" {
// 					interfaceArray := valB.([]interface{})
// 					for i := 0; i < len(interfaceArray); i++ {
// 						for keyC, valC := range interfaceArray[i].(map[string]interface{}) {
// 							if keyC == "user" {
// 								for keyD, valD := range valC.(map[string]interface{}) {
// 									if keyD == "id" {
// 										_ = httpGetUser(valD.(float64))
// 										// fmt.Println(userInfo.region)
// 									}

// 									if keyD == "screen_name" {
// 										err = writeTo(file, keyD+": "+valD.(string)+" ")
// 										if err != nil {
// 											log.Fatal("writeTo hot D", err)
// 										}

// 									}
// 								}
// 							}

// 							if keyC == "text" {
// 								re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
// 								src := re.ReplaceAllString(valC.(string), "")

// 								err = writeTo(file, keyC+": "+src+" ")
// 								if err != nil {
// 									log.Fatal("writeTo hot C", err)
// 								}
// 							}
// 						}

// 						err = writeTo(file, "\n")
// 						if err != nil {
// 							log.Fatal("writeTo hot n", err)
// 						}

// 					}
// 				}
// 			}
// 		}
// 	}
// 	// done <- true
// }

// func httpPost() {
//     resp, err := http.Post("http://www.01happy.com/demo/accept.php",
//         "application/x-www-form-urlencoded",
//         strings.NewReader("name=cjb"))
//     if err != nil {
//         fmt.Println(err)
//     }

//     defer resp.Body.Close()
//     body, err := ioutil.ReadAll(resp.Body)
//     if err != nil {
//         // handle error
//     }

//     fmt.Println(string(body))
// }
