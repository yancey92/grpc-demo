// Copyright © 2024 Yangxinxin. All right reserved.
// Confidential and Proprietary
package timekit

import (
	"fmt"
	"testing"
)

func TestTimeLocalFormat(t *testing.T) {
	timeStr := TimeLocalFormat()
	fmt.Println(timeStr)
}
