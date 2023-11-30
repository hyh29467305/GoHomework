package main

import "fmt"

func main() {
	TestVals := []int{1, 2, 3, 4, 5, 6, 7, 8}
	TmpResult := SliceDelete(7, TestVals)
	fmt.Printf("第1次删除%v \n", TmpResult)
}
