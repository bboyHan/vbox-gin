package utils

import "math/rand"

type RandomList struct {
	List     []string
	indexMap map[string]int
}

func NewRandomList(values []string) *RandomList {
	indexMap := make(map[string]int)
	for i, val := range values {
		indexMap[val] = i
	}
	return &RandomList{
		List:     values,
		indexMap: indexMap,
	}
}

func (rl *RandomList) GetRandom() string {
	if len(rl.List) == 0 {
		return ""
	}

	index := rand.Intn(len(rl.List))
	value := rl.List[index]

	// 将最后一个元素移到被取出的位置，更新索引映射
	lastIndex := len(rl.List) - 1
	lastValue := rl.List[lastIndex]
	rl.List[index] = lastValue
	rl.indexMap[lastValue] = index

	// 删除最后一个元素
	rl.List = rl.List[:lastIndex]
	delete(rl.indexMap, value)

	return value
}
