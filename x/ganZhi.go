package x

import (
	"fmt"
	"github.com/starainrt/astro/basic"
	"github.com/starainrt/astro/calendar"
	"sort"
	"time"
)

//干支信息
type GanZhi struct {
	YGZ string `json:"ygz"`
	MGZ string `json:"mgz"`
	DGZ string `json:"dgz"`
	HGZ string `json:"hgz"`
}

func NewGanZhi(year, month, day, hour int) {
	cust := time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.Local) //精确到时
	lcb := fixLiChun(year, cust)
	fmt.Println(lcb)

}

var (
	//十天干 甲1 ...癸10
	Gan = [11]string{"err", "甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	//十二地支 子1 ...亥12
	Zhi = [13]string{"err", "子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
)

//##############################################s
//计算年干支
//##############################################

//年干支
func GetYGZ(year, month, day, hour int) string {
	cust := time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.Local) //精确到时
	lcb := fixLiChun(year, cust)
	g, z := YearGZ(year, lcb)
	return g + z
}

//传入阳历年 立春布尔值 返回年干 年支 年干支
//年干支
func YearGZ(year int, lcb bool) (string, string) {
	var aliasGan, aliasZhi string
	switch lcb {
	case false: //日期在立春之前
		//干
		g := 1 + (year+6)%10
		if g -= 1; g < 1 {
			g += 10
		}
		aliasGan = Gan[g]
		//支
		z := 1 + (year+8)%12
		if z -= 1; z < 1 {
			z += 12
		}
		aliasZhi = Zhi[z]
	case true: //日期在立春之后
		yearg := 1 + (year+6)%10
		yearz := 1 + (year+8)%12
		aliasGan = Gan[yearg]
		aliasZhi = Zhi[yearz]
	}

	return aliasGan, aliasZhi

}

//立春修正
func fixLiChun(year int, cust time.Time) bool {
	lct, _ := getJie12T(year)
	lct = time.Date(lct.Year(), lct.Month(), lct.Day(), lct.Hour(), 0, 0, 0, time.Local)
	var b bool
	if cust.Equal(lct) || cust.After(lct) {
		b = true //当前时间在立春之后
	} else {
		b = false //当前时间在立春之前
	}
	return b
}

//传入阳历年数字 返回本年立春阳历时间戳 12节时间戳数组(上一年冬至到本年冬至)
//获取本年立春时间戳
func getJie12T(year int) (time.Time, []time.Time) {

	year -= 1 //k:1-->上一年冬至时间 k:25-->本年冬至时间 k:4--本年立春
	jq := basic.GetOneYearJQ(year)
	var keys []int
	for k := range jq {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	//k:1上一年冬至...k4:立春... k:25本年冬至
	/*
		"冬至", "小寒", "大寒", "立春", "雨水", "惊蛰",
		"春分", "清明", "谷雨", "立夏", "小满", "芒种",
		"夏至", "小暑", "大暑", "立秋", "处暑", "白露",
		"秋分", "寒露", "霜降", "立冬", "小雪", "大雪", "冬至",
	*/
	var jieArr []time.Time //12节
	var lct time.Time
	for _, v := range keys {
		//fmt.Printf("k:%v -本地时区%v\n", v, calendar.JDE2Date(jq[v]))
		if v%2 == 0 {
			jieArr = append(jieArr, calendar.JDE2Date(jq[v]))
		}
		if v == 4 {
			lct = calendar.JDE2Date(jq[v])
		}
	}

	//12节
	// 小寒  立春  惊蛰  清明  立夏  芒种  小暑  立秋  白露  寒露  立冬  大雪
	//排序后对应的k值 2 4 6 8 10 12 14 16 18 20 22 24
	return lct, jieArr
}

//##############################################s
//计算月干支
//##############################################

//传入阳历时间 返回月干支
//以12节气定月干支
func MonthGZ(year, month, day, hour int) string {
	cust := time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.Local)
	arrT := getJieArr(year)
	b, index := findJie(cust, arrT)
	lcb := fixLiChun(year, cust)

	yg, _ := YearGZ(year, lcb)
	gzArr := mgzArr(yg)

	if (b == false && index == 0) && lcb == false { //在本年立春之前
		index -= 1
		if index < 0 {
			index += 12
		}
		index -= 1
	} else if (b == false && index == 0) && lcb == true {
		index -= 1
		if index < 0 {
			index += 12
		}
	} else if b == true {
		index -= 1
		if index < 0 {
			index += 12
		}
	}
	//fmt.Printf("年干:%s index:%d 月干支:%s\n", yg, index, gzArr[index])
	return gzArr[index]
}

//正月立春节 二月惊蛰节 三月清明节 四月立夏节 五月忙钟节 六月小暑节
//七月立秋节 八月白露节 九月寒露节 十月立东节 冬月大雪节 腊月小寒节
//0:上一年小寒 1今年立春...11大雪 12本年小寒 13下年立春
//上一年小寒到下一年节气的时间戳数组
func getJieArr(year int) []time.Time {
	_, j12arr := getJie12T(year)
	_, j4arr := getJie12T(year + 1)

	var arrT []time.Time
	arrT = append(arrT, j12arr...)
	arrT = append(arrT, j4arr...)
	return arrT
}

//true节气之后 false节气之前 节气计算精确到日
func findJie(cust time.Time, jarrT []time.Time) (bool, int) {
	cust = time.Date(cust.Year(), cust.Month(), cust.Day(), 0, 0, 0, 0, time.Local)
	var b bool
	var index int
	for i := 0; i < len(jarrT)-1; i++ {
		j0 := jarrT[i]
		j1 := jarrT[i+1]
		j0 = time.Date(j0.Year(), j0.Month(), j0.Day(), 0, 0, 0, 0, time.Local) //精确到日
		j1 = time.Date(j1.Year(), j1.Month(), j1.Day(), 0, 0, 0, 0, time.Local)
		if (cust.Equal(j0) || cust.After(j0)) && cust.Before(j1) {
			index = i
			b = true //节气之后
			break
		}
	}
	return b, index
}

// 五虎盾元 甲己之年丙作初，乙庚之歲戊為頭，丙辛歲首從庚起，丁壬壬位順流行，若問戊癸何方法，甲寅之上好推求
// 传入年干 返回本年月干支数组
// 月干支数组
func mgzArr(yg string) []string {
	gan := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhi := []string{"寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥", "子", "丑"} //正月建寅

	var arr []string //月干支数组
	switch yg {
	case gan[0], gan[5]: //甲己 丙寅
		end := gan[:2]
		head := gan[2:]
		arr = whd(gan, zhi, head, end)
	case gan[1], gan[6]: //乙庚 戊寅
		end := gan[:4]
		head := gan[4:]
		arr = whd(gan, zhi, head, end)
	case gan[2], gan[7]: //丙辛 庚寅
		end := gan[:6]
		head := gan[6:]
		arr = whd(gan, zhi, head, end)
	case gan[3], gan[8]: //丁壬 壬寅
		end := gan[:8]
		head := gan[8:]
		arr = whd(gan, zhi, head, end)
	case gan[4], gan[9]: //戊癸 甲寅
		end := gan
		head := gan
		arr = whd(gan, zhi, head, end)
	}
	return arr
}

//五虎遁干支数组
func whd(gan, zhi, head, end []string) []string {
	var arr []string
	gan = append(head, end...)
	gan = append(gan, gan...)
	for i := 0; i < len(gan); i++ {
		for j := i; j < len(zhi); j++ {
			s := gan[j] + zhi[j]
			arr = append(arr, s)
			if j == i {
				break
			}

		}
	}
	return arr
}

//##############################################s
//计算日干支
//##############################################
