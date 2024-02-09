// Copyright © 2024 Yangxinxin. All right reserved.
// Confidential and Proprietary
package strkit_test

import (
	"testing"

	"demo.test/grpc-demo/pkg/strkit"
)

func Test_ValidateXss(t *testing.T) {
	t.Helper()

	if strkit.ValidateXss(`"script`) {
		t.Fail()
	}

	if strkit.ValidateXss(`<script>`) {
		t.Fail()
	}

	if strkit.ValidateXss(`{script}`) {
		t.Fail()
	}

	if strkit.ValidateXss(`abc?`) {
		t.Fail()
	}

	if !strkit.ValidateXss(`script`) {
		t.Fail()
	}
}
