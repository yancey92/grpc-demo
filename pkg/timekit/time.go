// Copyright © 2024 Yangxinxin. All right reserved.
// Confidential and Proprietary
package timekit

import "time"

func TimeLocalFormat() string {
	// timelocal, _ := time.LoadLocation("Asia/Chongqing")
	// time.Local = timelocal
	// return time.Now().Local().Format("2006-01-02 15:04:05 CST")
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}
