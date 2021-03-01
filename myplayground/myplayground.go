// Package myplayground ...
package myplayground

import (
	"fmt"
	"go-api-sooon/app"
)

type play struct {
}

// DumpAnything implements IDump interface
// interface practice
func (p *play) DumpAnything(i interface{}) {
	fmt.Println("-------------------------")
	fmt.Printf("%#v\n", i)
	fmt.Println("-------------------------")
}

// Play ...
var Play *play

func init() {
	Play = &play{}
}

// TwoSum 這題主要目的在於練習HashMap/Dictionary的應用。
// 只要每一次從map中確認當下target-num是否在map中，
// 在的話就表示找到了，可以將結果取得，
// 沒找到的話，就將一組(num, index)放進map中，
// 依此流程，最壞的狀況整個array遍歷後，就可以得到答案。
// key, value對調放進map, 目標值-陣列value的差值 存在於 map的key 就可以取得 index
func TwoSum(nums []int, target int) []int {

	revereMap := make(map[int]int)

	for i := 0; i < len(nums); i++ {
		if j, ok := revereMap[target-nums[i]]; ok {
			return []int{j, i}
		}
		revereMap[nums[i]] = i
	}
	// mymap := make(map[int]int)
	// for i := 0; i < len(nums); i++ {
	// 	j, ok := mymap[target-nums[i]] // target-nums[i] 差值即是答案取得j即是原始num的index
	// 	if ok {
	// 		result := []int{j, i} // 答案回傳index
	// 		return result
	// 	}
	// 	mymap[nums[i]] = i
	// }
	result := []int{-1, -1}
	return result
}

// ArrReverse 陣列反轉 不套用函式
// 技巧 j := len(in) - i - 1 來反轉陣列的元素
func (p *play) ArrReverse(in []int) {

	for i := 0; i < len(in)/2; i++ {
		j := len(in) - i - 1        // j與i為頭尾對應index
		in[i], in[j] = in[j], in[i] // 對調
	}
	app.SFunc.Dump(p, in)
}

// Fibonacci1 費氏數列
// F_{0}=0}F_{0}=0
// F_{1}=1}F_{1}=1
// F_{n}=F_{n-1}+F_{n-2}}F_{n}=F_{{n-1}}+F_{{n-2}}（n≧2）
// 用文字來說，就是費氏數列由0和1開始，之後的費波那契數就是由之前的兩數相加而得出
// 0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233…
func (p *play) Fibonacci1() func(int) int {
	return func(x int) int {
		if x < 2 {
			return x
		}
		return p.Fibonacci1()(x-2) + p.Fibonacci1()(x-1) // 前兩個數字相加
	}
}

// Reverse ...
func (p *play) Reverse(x int) int {
	var MaxInt int32 = 2147483647

	if x > int(MaxInt) || x < -int(MaxInt) {
		return 0
	}

	tmp := 0
	run := x
	if x < 0 {
		run = -run
	}
	for run > 0 {
		tmp *= 10
		digit := run % 10
		tmp += digit
		run /= 10
	}

	if tmp > int(MaxInt) {
		return 0
	}

	if x < 0 {
		return -tmp
	}
	return tmp
}
