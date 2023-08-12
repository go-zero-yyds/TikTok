/**
 * @Author: FxShadow
 * @Description:
 * @Date: 2023/08/12 19:13
 */

package utils

import (
	"log"
	"strings"
)

// MatchUID 校验用户名参数是否合规（暂时废弃）
func MatchUID(uid int64) (ok bool, err error) {
	if strings.Count(string(uid), "") != 64 {
		log.Println(strings.Count(string(uid), ""))
		return false, err
	}

	return true, err
}
