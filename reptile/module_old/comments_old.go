package module

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"

	"ethan.com/reptile/basis"
	"ethan.com/reptile/util"
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

// HTTPGetData 获取数据入口
func HTTPGetData(j int, commentsID, commentsCookie, userCookie string, done chan *Statistics) {

	// 统计一条微博下一页的评论，男女数量
	var commentPageCounts, userMenPageCounts, userWomenPageCounts int
	// 统计一条微博下一页的各地区数量
	regionMap := RetMap()

	ageMap := InitAGE()

	strPage := strconv.Itoa(j)

	strURL := "https://m.weibo.cn/api/comments/show?id=" + commentsID + "&page=" + strPage
	// commentsCookie := "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548294434; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQzVvWiyamobldPXiJG-H3iDKIR0-TvlPiY091BPCFR24.; SUB=_2A25xJfzdDeRhGeVH6VsW9SbEzTyIHXVS6YSVrDV6PUJbktAKLXLZkW1NT3AvulTG2_x1nzGrufDPoHMAWfIb9JuX; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5K-hUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0Ry6R7PJoeAvuj; SSOLoginState=1545702541; WEIBOCN_FROM=1110006030; MLOGIN=1"

	req := &basis.Req{
		URL:    strURL,
		Cookie: commentsCookie,
	}
	resp := req.GetResp()
	defer resp.Body.Close()

	dfn := &util.DIYFileName{
		PageID: strPage,
		DirID:  commentsID,
	}
	file := util.InitFile(dfn)
	defer file.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("ioutil.ReadAll fatal:", err)
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		log.Fatalln("json.Unmarshal fatal:", err)
	}

	for keyA, valA := range dat {

		if keyA == "data" {
			for keyB, valB := range valA.(map[string]interface{}) {
				if keyB == "data" {
					interfaceArray := valB.([]interface{})
					for i := 0; i < len(interfaceArray); i++ {

						var str, src string
						var flag bool
						for keyC, valC := range interfaceArray[i].(map[string]interface{}) {
							if keyC == "user" {
								for keyD, valD := range valC.(map[string]interface{}) {
									if keyD == "id" {
										userInfo := HTTPGetUser(valD.(string), userCookie)

										userInfo.SexCount(&userMenPageCounts, &userWomenPageCounts)
										userInfo.RegionCount(regionMap)
										userInfo.BirthCount(ageMap)

										str = userInfo.ViewUserInfo()
									}
								}
							}

							if keyC == "text" {
								re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
								src = re.ReplaceAllString(valC.(string), "")
							}

							if str != "" && src != "" {
								if !flag {
									err = util.WriteTo(file, str+"\r\n")
									if err != nil {
										log.Fatal("writeTo userInfo", err)
									}

									err = util.WriteTo(file, "评论: "+src+"\r\n")
									if err != nil {
										log.Fatal("writeTo C", err)
									}

									flag = true
								}
							}
						}

						err = util.WriteTo(file, "----------------------\r\n")
						if err != nil {
							log.Fatal("writeTo n", err)
						}
						commentPageCounts++
					}
				}
				if keyB == "hot_data" {
					interfaceArray := valB.([]interface{})
					for i := 0; i < len(interfaceArray); i++ {

						var str, src string
						var flag bool
						for keyC, valC := range interfaceArray[i].(map[string]interface{}) {
							if keyC == "user" {
								for keyD, valD := range valC.(map[string]interface{}) {
									if keyD == "id" {
										userInfo := HTTPGetUser(valD.(string), userCookie)

										userInfo.SexCount(&userMenPageCounts, &userWomenPageCounts)
										userInfo.RegionCount(regionMap)
										userInfo.BirthCount(ageMap)

										str = userInfo.ViewUserInfo()
									}

								}
							}

							if keyC == "text" {
								re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
								src = re.ReplaceAllString(valC.(string), "")

							}

							if str != "" && src != "" {
								if !flag {
									err = util.WriteTo(file, str+"\r\n")
									if err != nil {
										log.Fatal("writeTo userInfo", err)
									}

									err = util.WriteTo(file, "评论: "+src+"\r\n")
									if err != nil {
										log.Fatal("writeTo C", err)
									}

									flag = true
								}
							}
						}

						err = util.WriteTo(file, "------------------------\r\n")
						if err != nil {
							log.Fatal("writeTo hot n", err)
						}
						commentPageCounts++
					}
				}
			}
		}
	}
	os.Chmod(file.Name(), 0444)

	// fmt.Println(ageMap)
	statistics := &Statistics{
		CommentCount:    commentPageCounts,
		UserMenCount:    userMenPageCounts,
		UserWomenCount:  userWomenPageCounts,
		UserRegionCount: regionMap,
		UserBrithCount:  ageMap,
	}
	done <- statistics
}
