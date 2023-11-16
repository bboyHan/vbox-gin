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

func J3Contains(target string) bool {
	set := map[string]bool{
		"2000": true,
		"2001": true,
		"2002": true,
		"2003": true,
		"2004": true,
		"2005": true,
		"2006": true,
		"2007": true,
		"2008": true,
		"2009": true,
	}
	_, found := set[target]
	return found
}
