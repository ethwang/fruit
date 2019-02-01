package basis

import (
	"log"
	"net/http"
)

// Req 请求参数
type Req struct {
	URL    string
	Cookie string
}

// GetResp :req --> Resp
func (req *Req) GetResp() (resp *http.Response) {
	// resp, err := http.Get("https://m.weibo.cn/api/comments/show?id=4248590911655823&page=1")
	// strURL := url

	// fmt.Println(url)

	// jar, _ := cookiejar.New(nil)
	// client := &http.Client{Jar: jar}
	client := &http.Client{}

	reqest, err := http.NewRequest("GET", req.URL, nil)
	if err != nil {
		panic(err)
	}

	// 增加header选项
	reqest.Header.Add("Cookie", req.Cookie)
	// reqest.Header.Add("Cookie", "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548212292; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQf4MQo7Yt3VsW2PbOiWWNvz2YPiUpXDD0ODmkKgpSJhU.; SUB=_2A25xJDsVDeRhGeVH6VsW9SbEzTyIHXVS50VdrDV6PUJbktAKLWzkkW1NT3AvugY-8qzK-sQXfz0UazHw_PkluYH0; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5KzhUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0laV7SlW8Hji0O; SSOLoginState=1545620293; WEIBOCN_FROM=1110006030; MLOGIN=1")
	//处理返回结果
	// resp, _ := client.Get(strURL)

	resp, err = client.Do(reqest)
	if err != nil {
		log.Fatalln("http.Get fatal:", err)
	}
	return
}
