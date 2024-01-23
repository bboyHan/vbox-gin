package product

import (
	"fmt"
	"strings"
)

//var rawURL = "https://security.seasungame.com/security_extend_server/helper/balance/queryBalance?
//gameCode=jx3&account=18210889498&accountType=&zoneCode=z05&SN=98710648156&remark=&sign=36A360706FD189A2BF867D70F656C7BE"

// 校验传入卡密合法性

func ParseJWCardRecord(ext string) (string, error) {

	if ext == "" {
		return "", fmt.Errorf("卡密不合法")
	}
	if !strings.Contains(ext, "_") {
		return "", fmt.Errorf("卡密不合法")
	}

	split := strings.Split(ext, "_")
	if len(split) != 2 {
		return "", fmt.Errorf("卡密不合法")
	} else {
		return ext, nil
	}

}
