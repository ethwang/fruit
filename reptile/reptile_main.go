package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"ethan.com/fruit/reptile/module"
	"ethan.com/fruit/reptile/util"
)

/*
	数据源头URL:https://weibo.cn/comment/GDULCaK7D?uid=1977278337&rl=0&page=713
	---> commentsStr=GDULCaK7D ---> userID=1977278337
*/

// const count = 10
// const commentsID = "4296473875115737"
// const commentsCookie = "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548294434; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQzVvWiyamobldPXiJG-H3iDKIR0-TvlPiY091BPCFR24.; SUB=_2A25xJfzdDeRhGeVH6VsW9SbEzTyIHXVS6YSVrDV6PUJbktAKLXLZkW1NT3AvulTG2_x1nzGrufDPoHMAWfIb9JuX; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5K-hUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0Ry6R7PJoeAvuj; SSOLoginState=1545702541; WEIBOCN_FROM=1110006030; MLOGIN=1"
var userCookieConst = "_T_WM=9529dfe60c989460597a6b912a6c7c54; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQznJIC0e6xTseFKs46M4lrf0qZU0UomdlebuhZoLj3PM.; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5KMhUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUB=_2A25xUHRfDeRhGeVH6VsW9SbEzTyIHXVSuxwXrDV6PUJbkdAKLVCkkW1NT3AvumFsBVlmNOhRqCUsSiLdB-7BT4-9; SUHB=079zbPsHY6jhjy; SSOLoginState=1549009935"

func main() {
	var commentCount, userMenCount, userWomenCount, totalRegion, totalAge int

	scount := flag.Int("scount", 1, "开始爬取普通评论页")
	count := flag.Int("count", 10, "爬取普通评论页数")
	hotcount := flag.Int("hotcount", 1, "爬取热点评论页数")
	uID := flag.String("uID", "5644764907", "微博用户ID")
	// commentsID := flag.String("commentsID", "4296473875115737", "微博评论ID")
	commentsStr := flag.String("commentsStr", "FfwVQAqJ4", "微博评论Str")
	// commentsCookie := flag.String("commentsCookie", "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548294434; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQzVvWiyamobldPXiJG-H3iDKIR0-TvlPiY091BPCFR24.; SUB=_2A25xJfzdDeRhGeVH6VsW9SbEzTyIHXVS6YSVrDV6PUJbktAKLXLZkW1NT3AvulTG2_x1nzGrufDPoHMAWfIb9JuX; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5K-hUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0Ry6R7PJoeAvuj; SSOLoginState=1545702541; WEIBOCN_FROM=1110006030; MLOGIN=1", "微博评论Cookie")
	userCookie := flag.String("userCookie", userCookieConst, "微博用户Cookie")
	sign := flag.String("sign", "0", "微博评论爬取标志,默认为单线程")

	flag.Parse()

	if *hotcount > *count {
		log.Fatal("Failed, invaild: hotcount > count")
		return
	}

	regionMap := module.RetMap()
	ageMap := module.InitAGE()

	var icount int

	if (*sign) == "1" {

		err := os.RemoveAll("./test_" + *commentsStr + "_1")
		if err != nil {
			log.Fatal("os.Remove", err)
		}

		os.Mkdir("./test_"+*commentsStr+"_1/", 0777)

		reqDone := make(chan *module.Statistics, *count)
		for icount = 1; icount <= *hotcount; icount++ {
			go module.HTTPGetData(icount, icount, *commentsStr, *uID, *userCookie, reqDone)
		}
		for ; icount <= *count; icount++ {
			go module.HTTPGetData(icount, -1, *commentsStr, *uID, *userCookie, reqDone)
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
	} else if (*sign) == "0" {

		os.Mkdir("./test_"+*commentsStr+"_0/", 0777)
		if *scount <= *hotcount {
			for icount = *scount; icount <= *hotcount; icount++ {
				st, err := module.HTTPGetDataSingle(icount, icount, *commentsStr, *uID, *userCookie)
				if err != nil {
					break
				}

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

			for ; icount <= *count; icount++ {
				st, err := module.HTTPGetDataSingle(icount, -1, *commentsStr, *uID, *userCookie)
				if err != nil {
					break
				}
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
		} else {
			for icount = *scount; icount <= *count; icount++ {
				st, err := module.HTTPGetDataSingle(icount, -1, *commentsStr, *uID, *userCookie)
				if err != nil {
					break
				}
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
		}
	} else {
		log.Fatal("Failed, invaild sign")
		return
	}

	for j := 0; j < len(module.REGION); j++ {
		totalRegion = regionMap[module.REGION[j]] + totalRegion
	}
	for j := 0; j < len(module.AGE); j++ {
		totalAge = ageMap[module.AGE[j]] + totalAge
	}

	file, err := os.OpenFile(
		"./test_"+*commentsStr+"_"+*sign+"/count_"+strconv.Itoa(icount-1)+".txt",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	defer file.Close()
	if err != nil {
		log.Fatal("os.OpenFile count.txt", err)
	}

	util.WriteTo(file, "----评论总数："+strconv.Itoa(commentCount)+"----\r\n\r\n")
	util.WriteTo(file, "----男性人数："+strconv.Itoa(userMenCount)+", 女性人数："+strconv.Itoa(userWomenCount)+", 总人数："+strconv.Itoa(userWomenCount+userMenCount)+"\r\n")
	util.WriteTo(file, "男性占比："+strconv.FormatFloat(util.Decimal(float64(userMenCount)/float64(userMenCount+userWomenCount)), 'f', -1, 64)+"\r\n")
	util.WriteTo(file, "女性占比："+strconv.FormatFloat(util.Decimal(float64(userWomenCount)/float64(userMenCount+userWomenCount)), 'f', -1, 64)+"\r\n\r\n\r\n")

	fmt.Println("----commentCount", commentCount, "-----\a")
	fmt.Println("-----userMenCount", userMenCount, "-----totalSex---", userMenCount+userWomenCount, "----\a")
	fmt.Println("-----userWomenCount", userWomenCount, "-----totalSex---", userMenCount+userWomenCount, "----\a")
	fmt.Printf("男性占比：%.2f \n", float64(userMenCount)/float64(userMenCount+userWomenCount))
	fmt.Printf("女性占比：%.2f \n", float64(userWomenCount)/float64(userMenCount+userWomenCount))

	fmt.Println("----userRegionCount", regionMap, "--------totalRegion-----", totalRegion, "-----\a")
	for j := 0; j < len(module.REGION); j++ {
		fmt.Printf("省份：%s,占比：%.2f \n", module.REGION[j], float64(regionMap[module.REGION[j]])/float64(totalRegion))
		util.WriteTo(file, "省份："+module.REGION[j]+", 个数："+strconv.Itoa(regionMap[module.REGION[j]])+", 总数："+strconv.Itoa(totalRegion)+", 占比："+strconv.FormatFloat(util.Decimal(float64(regionMap[module.REGION[j]])/float64(totalRegion)), 'f', -1, 64)+"\r\n\r\n")
	}
	util.WriteTo(file, "\r\n")
	fmt.Println("----userBrithCount", ageMap, "------totalAge------", totalAge, "------\a")
	for j := 0; j < len(module.AGE); j++ {
		fmt.Printf("人群：%s,占比：%.2f \n", module.AGE[j], float64(ageMap[module.AGE[j]])/float64(totalAge))
		util.WriteTo(file, "人群："+module.AGE[j]+", 人数："+strconv.Itoa(ageMap[module.AGE[j]])+", 总数："+strconv.Itoa(totalAge)+", 占比："+strconv.FormatFloat(util.Decimal(float64(ageMap[module.AGE[j]])/float64(totalAge)), 'f', -1, 64)+"\r\n\r\n")

	}
	os.Chmod(file.Name(), 0444)

	// 结束提示音
	for i := 0; i < 5; i++ {
		fmt.Println("\a")
		time.Sleep(time.Second)
	}
}
