package utils

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"time"
)

func SearchIp2Region(ip string) (region string, err error) {
	var dbPath = "utils/ip2region.xdb"
	searcher, err := xdb.NewWithFileOnly(dbPath)
	if err != nil {
		fmt.Printf("failed to create searcher: %s\n", err.Error())
		return
	}

	defer searcher.Close()

	var tStart = time.Now()
	region, err = searcher.SearchByStr(ip)
	if err != nil {
		fmt.Printf("failed to SearchIP(%s): %s\n", ip, err)
		return
	}

	fmt.Printf("{region: %s, took: %s}\n", region, time.Since(tStart))

	// 备注：并发使用，每个 goroutine 需要创建一个独立的 searcher 对象。
	return region, err
}
