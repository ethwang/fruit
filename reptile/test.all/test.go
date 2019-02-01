package main

import (
	"fmt"
	"strings"
)

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

func main() {
	// str := "昵称:浅薄67904性别:男地区:其他"
	// str := "昵称:一条勤奋的咸鱼认证:心动超越超话主持人性别:男地区:广东 广州生日:1993-06-12认证信息：心动超越超话主持人简介:devout Christian标签:新闻热点 更多>>"
	// str := "昵称:芃芃其麦杨超越认证:知名娱乐博主 村民超话主持人性别:女地区:江苏生日:0000-11-06认证信息：知名娱乐博主 村民超话主持人简介:我行其野，芃芃其麦。"
	// str := "昵称:时光深度丨性别:女地区:河南 驻马店"
	// str := "昵称:拉翠性别:女地区:台湾 台北市生日:2000-08-23"
	// str := "昵称:天秤座的懒癌患者性别:女地区:湖南生日:2017-12-25感情状况：单身简介:可以交友哦QwQ"
	// str := "昵称:忧伤独存性别:男地区:江西生日:1999-11-03简介:一个比较随便的人"
	// str := "274467799"
	str := "昵称:IU-YUNER8性别:女地区:广东 广州生日:2005-03-21简介:微博认证:"
	userInfo := parseArgs(str)
	fmt.Println(userInfo)
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
		fmt.Println("strlen", strlen)
		fmt.Println("p", p)
		s := string([]byte(okStr)[p:strlen])
		mapLabel[label[m]] = s
		fmt.Println("label[m]", label[m])
		fmt.Println("s", s)
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

// 预处理
func pretreatment(str string) (okStr string) {
	okStr = str
	if ethan := strings.Index(str, "性别:"); ethan != -1 {
		strlen := len([]byte(str))

		s := string([]byte(str)[ethan:strlen])

		if x := strings.Index(s, "认证:"); x != -1 {
			okStr = string([]byte(str)[:ethan]) + strings.Split(s, "认证:")[0] + strings.Split(s, "认证:")[1]
			fmt.Println(str)
		}
	}
	return
}
