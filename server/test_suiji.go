package main

//import (
//	"fmt"
//	"math/rand"
//)
//
//type RandomList struct {
//	list     []string
//	indexMap map[string]int
//}
//
//func NewRandomList(values []string) *RandomList {
//	indexMap := make(map[string]int)
//	for i, val := range values {
//		indexMap[val] = i
//	}
//	return &RandomList{
//		list:     values,
//		indexMap: indexMap,
//	}
//}
//
//func (rl *RandomList) GetRandom() string {
//	if len(rl.list) == 0 {
//		return ""
//	}
//
//	index := rand.Intn(len(rl.list))
//	value := rl.list[index]
//
//	// 将最后一个元素移到被取出的位置，更新索引映射
//	lastIndex := len(rl.list) - 1
//	lastValue := rl.list[lastIndex]
//	rl.list[index] = lastValue
//	rl.indexMap[lastValue] = index
//
//	// 删除最后一个元素
//	rl.list = rl.list[:lastIndex]
//	delete(rl.indexMap, value)
//
//	return value
//}
//
//func main() {
//	resList := []string{"A", "B", "C", "D", "E"}
//
//	randomList := NewRandomList(resList)
//
//	for len(randomList.list) > 0 {
//		value := randomList.GetRandom()
//		fmt.Println("取出的值:", value)
//	}
//	value := randomList.GetRandom()
//	fmt.Println("取出的值:", value)
//}
