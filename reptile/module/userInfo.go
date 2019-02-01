package module

import (
	"log"
	"strings"

	"ethan.com/fruit/reptile/basis"
	"github.com/PuerkitoBio/goquery"
)

// UserInfo 用户信息
type UserInfo struct {
	name             string
	authenticate     string
	talent           string
	sex              string
	region           string
	brith            string
	sexo             string
	relationship     string
	authenticateInfo string
	introduction     string
	label            string
}

// ViewUserInfo 查看用户信息
func (userInfo *UserInfo) ViewUserInfo() string {
	var userInfoStr string
	if userInfo.name != "" {
		// fmt.Println(userInfo.name)
		userInfoStr = userInfoStr + userInfo.name + " "
	}
	if userInfo.authenticate != "" {
		// fmt.Println(userInfo.authenticate)
		userInfoStr = userInfoStr + userInfo.authenticate + " "

	}
	if userInfo.talent != "" {
		// fmt.Println(userInfo.talent)
		userInfoStr = userInfoStr + userInfo.talent + " "

	}
	if userInfo.sex != "" {
		// fmt.Println(userInfo.sex)
		userInfoStr = userInfoStr + userInfo.sex + " "

	}
	if userInfo.region != "" {
		// fmt.Println(userInfo.region)
		userInfoStr = userInfoStr + userInfo.region + " "

	}
	if userInfo.brith != "" {
		// fmt.Println(userInfo.brith)
		userInfoStr = userInfoStr + userInfo.brith + " "

	}
	if userInfo.sexo != "" {
		// fmt.Println(userInfo.sexo)
		userInfoStr = userInfoStr + userInfo.sexo + " "

	}
	if userInfo.relationship != "" {
		// fmt.Println(userInfo.relationship)
		userInfoStr = userInfoStr + userInfo.relationship + " "

	}
	if userInfo.authenticateInfo != "" {
		// fmt.Println(userInfo.authenticateInfo)
		userInfoStr = userInfoStr + userInfo.authenticateInfo + " "

	}
	if userInfo.introduction != "" {
		// fmt.Println(userInfo.introduction)
		userInfoStr = userInfoStr + userInfo.introduction + " "

	}
	if userInfo.label != "" {
		// fmt.Println(userInfo.label)
		userInfoStr = userInfoStr + userInfo.label + " "

	}
	// fmt.Println("---------------------------")
	return userInfoStr
}

// RegionCount 统计各区域人数
func (userInfo *UserInfo) RegionCount(regionMapParam map[string]int) {

	if userInfo.region != "" {
		detailRegion := strings.Split(userInfo.region, ":")[1]
		detailRegion = strings.Split(detailRegion, " ")[0]
		if _, ok := regionMapParam[detailRegion]; ok {
			regionMapParam[detailRegion]++
		}
	}
}

// SexCount 统计性别人数
func (userInfo *UserInfo) SexCount(men, women *int) {
	if userInfo.sex == "性别:女" {
		*women++
	}
	if userInfo.sex == "性别:男" {
		*men++
	}

}

// BirthCount 统计年龄段
func (userInfo *UserInfo) BirthCount(brithMapParam map[string]int) {
	if userInfo.brith != "" {
		strBrith := string([]byte(userInfo.brith)[:10])
		switch strBrith {
		case "生日:199":
			brithMapParam["90后"]++
		case "生日:200":
			brithMapParam["00后"]++
		case "生日:198":
			brithMapParam["80后"]++
		case "生日:197":
			brithMapParam["70后"]++
		case "生日:201":
			brithMapParam["10后"]++
		}
		// fmt.Println(strBrith)
	}
}
func parseArgs(str string) (userInfo *UserInfo) {
	// 字符串预处理
	okStr := pretreatment(str)

	flag := 0
	k := 1
	o := 1
	mapLabel := make(map[string]string)
	label := []string{
		"昵称:",
		"认证:",
		"达人:",
		"性别:",
		"地区:",
		"生日:",
		"性取向：",
		"感情状况：",
		"认证信息：",
		"简介:",
		"标签:",
		"更多>>",
	}
	var m int

	// str := "昵称:一条勤奋的咸鱼认证:心动超越超话主持人性别:男地区:广东 广州生日:1993-06-12认证信息：心动超越超话主持人简介:devout Christian标签:新闻热点 更多>>"
	// str := "昵称:芃芃其麦杨超越认证:知名娱乐博主 村民超话主持人性别:女地区:江苏生日:0000-11-06认证信息：知名娱乐博主 村民超话主持人简介:我行其野，芃芃其麦。"
	p := strings.Index(okStr, label[0])
	if p != -1 {
		strlen := len([]byte(okStr))

		for i := 1; i < len(label); i++ {
			if q := strings.Index(okStr, label[i]); q != -1 {
				s := string([]byte(okStr)[p:q])

				mapLabel[label[i-o]] = s
				// fmt.Println(s, "---------", label[i-o])
				p = q
				flag = 0
				k = 1
				m = i

			} else {
				flag++
			}
			o = k + flag
		}
		s := string([]byte(okStr)[p:strlen])
		mapLabel[label[m]] = s
		userInfo = &UserInfo{
			name:             mapLabel["昵称:"],
			authenticate:     mapLabel["认证:"],
			talent:           mapLabel["达人:"],
			sex:              mapLabel["性别:"],
			region:           mapLabel["地区:"],
			brith:            mapLabel["生日:"],
			sexo:             mapLabel["性取向："],
			relationship:     mapLabel["感情状况："],
			authenticateInfo: mapLabel["认证信息："],
			introduction:     mapLabel["简介:"],
			label:            mapLabel["标签:"],
		}
	}
	return
}

// HTTPGetUser 获取用户信息
func HTTPGetUser(uid string, userCookie string) (userInfo *UserInfo) {
	// uid := strconv.FormatFloat(id, 'f', -1, 64)
	strURL := "https://weibo.cn/" + uid + "/info"

	// userCookie := "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548294434; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQzVvWiyamobldPXiJG-H3iDKIR0-TvlPiY091BPCFR24.; SUB=_2A25xJfzdDeRhGeVH6VsW9SbEzTyIHXVS6YSVrDV6PUJbktAKLXLZkW1NT3AvulTG2_x1nzGrufDPoHMAWfIb9JuX; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5K-hUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0Ry6R7PJoeAvuj; SSOLoginState=1545702541"

	reqArg := &basis.Req{
		URL:    strURL,
		Cookie: userCookie,
	}
	resp := reqArg.GetResp()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal("goquery.NewDocumentFromResponse", err)
	}

	// fmt.Println(doc.Html())
	// fmt.Println(doc.Find("title").Text())
	t := doc.Find("div")
	for i := 0; i < t.Length(); i++ {
		if t.Eq(i).Text() == "基本信息" {
			// fmt.Println(t.Eq(i + 1).Text())
			userInfo = parseArgs(t.Eq(i + 1).Text())
			break
		}
	}
	return
}

// 字符串预处理 过滤掉<性别:>之后的<认证:>
func pretreatment(str string) (okStr string) {
	okStr = str
	if ethan := strings.Index(str, "性别:"); ethan != -1 {
		strlen := len([]byte(str))

		s := string([]byte(str)[ethan:strlen])

		if x := strings.Index(s, "认证:"); x != -1 {
			okStr = string([]byte(str)[:ethan]) + strings.Split(s, "认证:")[0] + strings.Split(s, "认证:")[1]
			// fmt.Println(str)
		}
	}
	return
}

// func HTTPGetUser(id float64, userCookie string) (userInfo *UserInfo) {
// 	uid := strconv.FormatFloat(id, 'f', -1, 64)
// 	strURL := "https://weibo.cn/" + uid + "/info"

// 	// userCookie := "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548294434; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQzVvWiyamobldPXiJG-H3iDKIR0-TvlPiY091BPCFR24.; SUB=_2A25xJfzdDeRhGeVH6VsW9SbEzTyIHXVS6YSVrDV6PUJbktAKLXLZkW1NT3AvulTG2_x1nzGrufDPoHMAWfIb9JuX; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5K-hUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0Ry6R7PJoeAvuj; SSOLoginState=1545702541"

// 	reqArg := &basis.Req{
// 		URL:    strURL,
// 		Cookie: userCookie,
// 	}
// 	resp := reqArg.GetResp()

// 	doc, err := goquery.NewDocumentFromResponse(resp)
// 	if err != nil {
// 		log.Fatal("goquery.NewDocumentFromResponse", err)
// 	}

// 	// fmt.Println(doc.Html())
// 	// fmt.Println(doc.Find("title").Text())
// 	t := doc.Find("div")
// 	for i := 0; i < t.Length(); i++ {
// 		if t.Eq(i).Text() == "基本信息" {
// 			// fmt.Println(t.Eq(i + 1).Text())
// 			userInfo = parseArgs(t.Eq(i + 1).Text())
// 			break
// 		}
// 	}
// 	return
// }
