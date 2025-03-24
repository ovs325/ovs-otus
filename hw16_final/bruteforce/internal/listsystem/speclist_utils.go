package listsystem

import (
	"fmt"
	"net"

	tp "bruteforce/internal/types"
)

// Проверка принадлежности IP к подсети.
func SearchEngineIP(subnetMap map[string]tp.IsBlack, ipStr string) (isFound bool, isBlack tp.IsBlack, err error) {
	for subnetStr, typeList := range subnetMap {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return false, false, fmt.Errorf("некорректный IP-адрес: %s", ipStr)
		}
		_, subnet, err := net.ParseCIDR(subnetStr)
		if err != nil {
			continue
		}
		if subnet.Contains(ip) {
			return true, typeList, nil
		}
	}
	return false, false, nil
}

// проверка правильности наименования подсети.
func SubnetNameVerifier(subnetStr string) (ok bool) {
	if _, _, err := net.ParseCIDR(subnetStr); err != nil {
		return false
	}
	return true
}

// Базовая проверка принадлежности id к какому-либо списку.
func SearchEngineCommon(specListMap map[string]tp.IsBlack, id string) (isFound bool, isBlack tp.IsBlack, err error) {
	isBlack, isFound = specListMap[id]
	return isFound, isBlack, nil
}

// Базовая проверка правильности наименования id.
func CommonNameVerifier(_ string) (ok bool) {
	return true
}
