package global

func TxContains(target string) bool {
	set := map[string]bool{
		"1000": true,
		"1001": true,
		"1002": true,
		"1003": true,
		"1004": true,
		"1005": true,
		"1006": true,
		"1007": true,
		"1008": true,
		"1009": true,
	}
	_, found := set[target]
	return found
}
