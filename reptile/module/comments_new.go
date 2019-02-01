package module

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"ethan.com/fruit/reptile/basis"
	"ethan.com/fruit/reptile/util"

	"github.com/PuerkitoBio/goquery"
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

// HTTPGetData 多线程爬取数据入口
func HTTPGetData(j, hotj int, commentsStr, uID, userCookie string, done chan *Statistics) {
	chanStatistics := make(chan *Statistics, 1)

	if hotj != -1 {

		go func() {
			var hotcommentPageCounts, hotuserMenPageCounts, hotuserWomenPageCounts int

			hotregionMap := RetMap()

			hotageMap := InitAGE()
			// https://weibo.cn/comment/hot/FfwVQAqJ4?rl=2&page=50
			strPage := strconv.Itoa(hotj)

			strURL := "https://weibo.cn/comment/hot/" + commentsStr + "?rl=2&page=" + strPage

			req := &basis.Req{
				URL:    strURL,
				Cookie: userCookie,
			}
			resp := req.GetResp()
			defer resp.Body.Close()

			dfn := &util.DIYFileName{
				PageID: "hot_" + strPage,
				DirID:  commentsStr + "_1",
			}
			file := util.InitFile(dfn)
			defer file.Close()

			doc, err := goquery.NewDocumentFromResponse(resp)
			if err != nil {
				log.Fatal("goquery.NewDocumentFromResponse", err)
			}

			t := doc.Find("div.c")
			for i := 0; i < t.Length(); i++ {
				var userToFile, commentToFile, awesomeToFile string

				if id, _ := t.Eq(i).Attr("id"); id != "" {

					s := t.Eq(i).Find("a")

					// 赞
					cc := t.Eq(i).Find("span.cc").First()

					// uid
					userInfoStr, exists := s.First().Attr("href")
					if exists {

						// 分割uid
						stru := strings.Split(userInfoStr, "/")

						// 从uid获取相应用户信息
						switch len(stru) {
						case 3:

							// 过滤掉首条微博博主的信息
							if cc.Text() != "" {
								// fmt.Println(stru[2])

								userInfo := HTTPGetUser(stru[2], userCookie)
								if userInfo == nil {
									break
								}

								userInfo.SexCount(&hotuserMenPageCounts, &hotuserWomenPageCounts)
								userInfo.RegionCount(hotregionMap)
								userInfo.BirthCount(hotageMap)
								userToFile = userInfo.ViewUserInfo()
								commentToFile = t.Eq(i).Find("span.ctt").Text()
								awesomeToFile = cc.Text()
								err = util.WriteTo(file, userToFile+"\r\n评论："+commentToFile+"\r\n"+awesomeToFile+"\r\n")
								if err != nil {
									log.Fatal("writeTo userInfo", err)
								}
								_ = util.WriteTo(file, "------------------------\r\n")
								hotcommentPageCounts++
								// fmt.Println(s.First().Text(), "------", userInfoStr, "-----", t.Eq(i).Find("span.ctt").Text(), "------", cc.Text())
							}

						case 2:
							// 这个接口来访问uid拿不到的情况
							// fmt.Println(stru[1])
						}
					}
				}
			}
			os.Chmod(file.Name(), 0444)
			// fmt.Println("hotcommentPageCounts=====", hotcommentPageCounts)
			chanstatistics := &Statistics{
				CommentCount:    hotcommentPageCounts,
				UserMenCount:    hotuserMenPageCounts,
				UserWomenCount:  hotuserWomenPageCounts,
				UserRegionCount: hotregionMap,
				UserBrithCount:  hotageMap,
			}
			chanStatistics <- chanstatistics
		}()
	}
	// 统计一条微博下一页的评论，男女数量
	var commentPageCounts, userMenPageCounts, userWomenPageCounts int
	// 统计一条微博下一页的各地区数量
	regionMap := RetMap()

	ageMap := InitAGE()

	strPage := strconv.Itoa(j)
	strURL := "https://weibo.cn/comment/" + commentsStr + "?uid=" + uID + "&rl=1&page=" + strPage + "&rand=31374h"
	// strURL :="https://weibo.cn/comment/FfwVQAqJ4?uid=5644764907&rl=1&page=1"
	// strURL := "https://m.weibo.cn/api/comments/show?id=" + commentsID + "&page=" + strPage
	// commentsCookie := "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548294434; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQzVvWiyamobldPXiJG-H3iDKIR0-TvlPiY091BPCFR24.; SUB=_2A25xJfzdDeRhGeVH6VsW9SbEzTyIHXVS6YSVrDV6PUJbktAKLXLZkW1NT3AvulTG2_x1nzGrufDPoHMAWfIb9JuX; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5K-hUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0Ry6R7PJoeAvuj; SSOLoginState=1545702541; WEIBOCN_FROM=1110006030; MLOGIN=1"

	req := &basis.Req{
		URL:    strURL,
		Cookie: userCookie,
	}
	resp := req.GetResp()
	defer resp.Body.Close()

	dfn := &util.DIYFileName{
		PageID: strPage,
		DirID:  commentsStr + "_1",
	}
	file := util.InitFile(dfn)
	defer file.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal("goquery.NewDocumentFromResponse", err)
	}
	t := doc.Find("div.c")
	for i := 0; i < t.Length(); i++ {
		var userToFile, commentToFile, awesomeToFile string

		if id, _ := t.Eq(i).Attr("id"); id != "" {

			s := t.Eq(i).Find("a")

			kt := t.Eq(i).Find("span.kt").First().Text()

			if kt != "[热门]" {
				cc := t.Eq(i).Find("span.cc").First()
				userInfoStr, exists := s.First().Attr("href")
				if exists {

					stru := strings.Split(userInfoStr, "/")

					switch len(stru) {
					case 3:

						// 过滤掉首条微博博主的信息
						if cc.Text() != "" {
							// fmt.Println(stru[2])

							userInfo := HTTPGetUser(stru[2], userCookie)
							if userInfo == nil {
								break
							}

							userInfo.SexCount(&userMenPageCounts, &userWomenPageCounts)
							userInfo.RegionCount(regionMap)
							userInfo.BirthCount(ageMap)
							userToFile = userInfo.ViewUserInfo()
							commentToFile = t.Eq(i).Find("span.ctt").Text()
							awesomeToFile = cc.Text()
							err = util.WriteTo(file, userToFile+"\r\n评论："+commentToFile+"\r\n"+awesomeToFile+"\r\n")
							if err != nil {
								log.Fatal("writeTo userInfo", err)
							}
							_ = util.WriteTo(file, "------------------------\r\n")
							commentPageCounts++
							// fmt.Println(s.First().Text(), "------", userInfoStr, "-----", t.Eq(i).Find("span.ctt").Text(), "------", cc.Text())
						}

					case 2:
						// 这个接口来访问uid拿不到的情况
						// fmt.Println(stru[1])

					}
				}
			}
		}
	}

	os.Chmod(file.Name(), 0444)
	if hotj != -1 {
		hotStatistics := <-chanStatistics
		// fmt.Println(ageMap)

		for j := 0; j < len(REGION); j++ {
			regionMap[REGION[j]] = regionMap[REGION[j]] + hotStatistics.UserRegionCount[REGION[j]]
		}

		for j := 0; j < len(AGE); j++ {
			ageMap[AGE[j]] = ageMap[AGE[j]] + hotStatistics.UserBrithCount[AGE[j]]
		}
		commentPageCounts = commentPageCounts + hotStatistics.CommentCount
		userMenPageCounts = userMenPageCounts + hotStatistics.UserMenCount
		userWomenPageCounts = userWomenPageCounts + hotStatistics.UserWomenCount
	}
	statistics := &Statistics{
		CommentCount:    commentPageCounts,
		UserMenCount:    userMenPageCounts,
		UserWomenCount:  userWomenPageCounts,
		UserRegionCount: regionMap,
		UserBrithCount:  ageMap,
	}
	done <- statistics
}

// HTTPGetDataSingle 单线程爬取数据入口
func HTTPGetDataSingle(j, hotj int, commentsStr, uID, userCookie string) (statistics *Statistics, err error) {
	chanStatistics := make(chan *Statistics, 1)

	if hotj != -1 {

		go func() {
			var hotcommentPageCounts, hotuserMenPageCounts, hotuserWomenPageCounts int

			hotregionMap := RetMap()

			hotageMap := InitAGE()
			// https://weibo.cn/comment/hot/FfwVQAqJ4?rl=2&page=50
			strPage := strconv.Itoa(hotj)

			strURL := "https://weibo.cn/comment/hot/" + commentsStr + "?rl=2&page=" + strPage

			req := &basis.Req{
				URL:    strURL,
				Cookie: userCookie,
			}
			resp := req.GetResp()
			defer resp.Body.Close()

			dfn := &util.DIYFileName{
				PageID: "hot_" + strPage,
				DirID:  commentsStr + "_0",
			}
			file := util.InitFile(dfn)
			defer file.Close()

			doc, err := goquery.NewDocumentFromResponse(resp)
			if err != nil {
				log.Fatal("goquery.NewDocumentFromResponse", err)
			}

			t := doc.Find("div.c")
			for i := 0; i < t.Length(); i++ {
				var userToFile, commentToFile, awesomeToFile string

				if id, _ := t.Eq(i).Attr("id"); id != "" {

					s := t.Eq(i).Find("a")

					// 赞
					cc := t.Eq(i).Find("span.cc").First()

					// uid
					userInfoStr, exists := s.First().Attr("href")
					if exists {

						// 分割uid
						stru := strings.Split(userInfoStr, "/")

						// 从uid获取相应用户信息
						switch len(stru) {
						case 3:

							// 过滤掉首条微博博主的信息
							if cc.Text() != "" {
								// fmt.Println(stru[2])

								userInfo := HTTPGetUser(stru[2], userCookie)
								if userInfo == nil {
									err = errors.New("userInfo is nil")
									return
								}

								userInfo.SexCount(&hotuserMenPageCounts, &hotuserWomenPageCounts)
								userInfo.RegionCount(hotregionMap)
								userInfo.BirthCount(hotageMap)
								userToFile = userInfo.ViewUserInfo()
								commentToFile = t.Eq(i).Find("span.ctt").Text()
								awesomeToFile = cc.Text()
								err = util.WriteTo(file, userToFile+"\r\n评论："+commentToFile+"\r\n"+awesomeToFile+"\r\n")
								if err != nil {
									log.Fatal("writeTo userInfo", err)
								}
								_ = util.WriteTo(file, "------------------------\r\n")
								hotcommentPageCounts++
								// fmt.Println(s.First().Text(), "------", userInfoStr, "-----", t.Eq(i).Find("span.ctt").Text(), "------", cc.Text())
							}

						case 2:
							// 这个接口来访问uid拿不到的情况
							// fmt.Println(stru[1])
						}
					}
				}
			}
			// os.Chmod(file.Name(), 0444)
			// fmt.Println("hotcommentPageCounts=====", hotcommentPageCounts)
			chanstatistics := &Statistics{
				CommentCount:    hotcommentPageCounts,
				UserMenCount:    hotuserMenPageCounts,
				UserWomenCount:  hotuserWomenPageCounts,
				UserRegionCount: hotregionMap,
				UserBrithCount:  hotageMap,
			}
			chanStatistics <- chanstatistics
		}()
	}
	// 统计一条微博下一页的评论，男女数量
	var commentPageCounts, userMenPageCounts, userWomenPageCounts int
	// 统计一条微博下一页的各地区数量
	regionMap := RetMap()

	ageMap := InitAGE()

	strPage := strconv.Itoa(j)
	strURL := "https://weibo.cn/comment/" + commentsStr + "?uid=" + uID + "&rl=1&page=" + strPage + "&rand=31374h"
	// strURL :="https://weibo.cn/comment/FfwVQAqJ4?uid=5644764907&rl=1&page=1"
	// strURL := "https://m.weibo.cn/api/comments/show?id=" + commentsID + "&page=" + strPage
	// commentsCookie := "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548294434; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQzVvWiyamobldPXiJG-H3iDKIR0-TvlPiY091BPCFR24.; SUB=_2A25xJfzdDeRhGeVH6VsW9SbEzTyIHXVS6YSVrDV6PUJbktAKLXLZkW1NT3AvulTG2_x1nzGrufDPoHMAWfIb9JuX; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5K-hUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0Ry6R7PJoeAvuj; SSOLoginState=1545702541; WEIBOCN_FROM=1110006030; MLOGIN=1"

	req := &basis.Req{
		URL:    strURL,
		Cookie: userCookie,
	}
	resp := req.GetResp()
	defer resp.Body.Close()

	dfn := &util.DIYFileName{
		PageID: strPage,
		DirID:  commentsStr + "_0",
	}
	file := util.InitFile(dfn)
	defer file.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal("goquery.NewDocumentFromResponse", err)
	}
	t := doc.Find("div.c")
	for i := 0; i < t.Length(); i++ {
		var userToFile, commentToFile, awesomeToFile string

		if id, _ := t.Eq(i).Attr("id"); id != "" {

			s := t.Eq(i).Find("a")

			kt := t.Eq(i).Find("span.kt").First().Text()

			if kt != "[热门]" {
				cc := t.Eq(i).Find("span.cc").First()
				userInfoStr, exists := s.First().Attr("href")
				if exists {

					stru := strings.Split(userInfoStr, "/")

					switch len(stru) {
					case 3:

						// 过滤掉首条微博博主的信息
						if cc.Text() != "" {
							// fmt.Println(stru[2])

							userInfo := HTTPGetUser(stru[2], userCookie)
							if userInfo == nil {
								err = errors.New("userInfo is nil")
								return
							}

							userInfo.SexCount(&userMenPageCounts, &userWomenPageCounts)
							userInfo.RegionCount(regionMap)
							userInfo.BirthCount(ageMap)
							userToFile = userInfo.ViewUserInfo()
							commentToFile = t.Eq(i).Find("span.ctt").Text()
							awesomeToFile = cc.Text()
							err = util.WriteTo(file, userToFile+"\r\n评论："+commentToFile+"\r\n"+awesomeToFile+"\r\n")
							if err != nil {
								log.Fatal("writeTo userInfo", err)
							}
							_ = util.WriteTo(file, "------------------------\r\n")
							commentPageCounts++
							// fmt.Println(s.First().Text(), "------", userInfoStr, "-----", t.Eq(i).Find("span.ctt").Text(), "------", cc.Text())
						}

					case 2:
						// 这个接口来访问uid拿不到的情况
						// fmt.Println(stru[1])

					}
				}
			}
		}
	}

	// os.Chmod(file.Name(), 0444)
	if hotj != -1 {
		hotStatistics := <-chanStatistics
		// fmt.Println(ageMap)

		for j := 0; j < len(REGION); j++ {
			regionMap[REGION[j]] = regionMap[REGION[j]] + hotStatistics.UserRegionCount[REGION[j]]
		}

		for j := 0; j < len(AGE); j++ {
			ageMap[AGE[j]] = ageMap[AGE[j]] + hotStatistics.UserBrithCount[AGE[j]]
		}
		commentPageCounts = commentPageCounts + hotStatistics.CommentCount
		userMenPageCounts = userMenPageCounts + hotStatistics.UserMenCount
		userWomenPageCounts = userWomenPageCounts + hotStatistics.UserWomenCount
	}
	statistics = &Statistics{
		CommentCount:    commentPageCounts,
		UserMenCount:    userMenPageCounts,
		UserWomenCount:  userWomenPageCounts,
		UserRegionCount: regionMap,
		UserBrithCount:  ageMap,
	}
	return
}
