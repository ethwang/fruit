package module

import "fmt"

// Statistics 统计结构体
type Statistics struct {
	CommentCount    int
	UserMenCount    int
	UserWomenCount  int
	UserRegionCount map[string]int
	UserBrithCount  map[string]int
}

// ViewSTT 查看单页统计信息
func (stt *Statistics) ViewSTT() {
	fmt.Println("==commentpageCount", stt.CommentCount, "==")
	fmt.Println("==userMenPageCount", stt.UserMenCount, "==")
	fmt.Println("==userWomenPageCount", stt.UserWomenCount, "==")
	fmt.Println("==userRegionPageCount", stt.UserRegionCount, "==")
	fmt.Println("==userAgePageCount", stt.UserBrithCount, "==")
}
