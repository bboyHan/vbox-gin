package utils

import "strconv"

func ConvertStringSliceToUintSlice(strSlice []string) ([]uint, error) {
	uintSlice := make([]uint, 0, len(strSlice))
	for _, str := range strSlice {
		num, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, err
		}
		uintSlice = append(uintSlice, uint(num))
	}
	return uintSlice, nil
}
