package main

import (
	"fmt"
	"log"
	"strings"

	"ethan.com/reptile/basis"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	// uid := strconv.FormatFloat(id, 'f', -1, 64)
	strURL := "https://weibo.cn/comment/FfwVQAqJ4?uid=5644764907&rl=1&page=2"
	// userCookie := "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548294434; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQzVvWiyamobldPXiJG-H3iDKIR0-TvlPiY091BPCFR24.; SUB=_2A25xJfzdDeRhGeVH6VsW9SbEzTyIHXVS6YSVrDV6PUJbktAKLXLZkW1NT3AvulTG2_x1nzGrufDPoHMAWfIb9JuX; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5K-hUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0Ry6R7PJoeAvuj; SSOLoginState=1545702541"
	userCookie := "_T_WM=4f1ae964d419e3772734d3cbeeb50e31; ALF=1548294434; SCF=AlOt0mmvqNuhM7kCNwt17u7M_hlwdiaU7Bz8zFbDy0dQzVvWiyamobldPXiJG-H3iDKIR0-TvlPiY091BPCFR24.; SUB=_2A25xJfzdDeRhGeVH6VsW9SbEzTyIHXVS6YSVrDV6PUJbktAKLXLZkW1NT3AvulTG2_x1nzGrufDPoHMAWfIb9JuX; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF6y4fgufihqD0lMV9jOgOi5JpX5K-hUgL.Foe4eo.NSKnRSo52dJLoI7_8UsyydLvgCHWiKBtt; SUHB=0Ry6R7PJoeAvuj; SSOLoginState=1545702541"
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
	// doc.Find("div.c").Each(func(i int, s *goquery.Selection) {
	// 	str := s.Find("a").Text()
	// 	if str != "举报" {
	// 		userInfoStr, exists := s.Find("a").Attr("href")
	// 		// if !exists {
	// 		// 	log.Fatal("no exists")

	// 		// }
	// 		if exists {
	// 			fmt.Println(str, "------", userInfoStr)
	// 		}
	// 	}

	// 	// s.Find("span.ctt").Text()
	// 	// fmt.Println(s.Find("span.ctt").Text())

	// })
	t := doc.Find("div.c")
	for i := 0; i < t.Length(); i++ {
		if id, _ := t.Eq(i).Attr("id"); id != "" {
			// if t.Eq(i).Find("a").Text() != "举报" {

			s := t.Eq(i).Find("a")
			// for j := 0; j < s.Length(); j++ {

			// }
			kt := t.Eq(i).Find("span.kt").First().Text()
			if kt == "[热门]" {

			} else {

				cc := t.Eq(i).Find("span.cc").First()
				userInfoStr, exists := s.First().Attr("href")
				if exists {
					if cc.Text() != "" {
						fmt.Println(s.First().Text(), "------", userInfoStr, "-----", t.Eq(i).Find("span.ctt").Text(), "------", cc.Text())
					}
				}
				stru := strings.Split(userInfoStr, "/")
				// if len(stru) == 3 {
				// 	fmt.Println(stru[2])
				// }
				switch len(stru) {
				case 3:
					// 过滤掉首条微博博主的信息
					if cc.Text() != "" {
						fmt.Println(stru[2])
					}
				case 2:
					fmt.Println(stru[1])
				}
			}
		}
	}

	return
}
