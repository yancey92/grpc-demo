// Copyright © 2024 Yangxinxin. All right reserved.
// Confidential and Proprietary
package strkit_test

import (
	"testing"

	"demo.test/grpc-demo/pkg/strkit"
)

func TestStrNotBlank(t *testing.T) {
	if !strkit.StrNotBlank("a") { //应该是true
		t.Log("TestStrNotBlank-01 fail")
		t.Fail()
	}

	if !strkit.StrNotBlank("a", "15") { //应该是true
		t.Log("TestStrNotBlank-02 fail")
		t.Fail()
	}

	if strkit.StrNotBlank("a", "", "     ") { //应该是false
		t.Log("TestStrNotBlank-03 fail")
		t.Fail()
	}

	if strkit.StrNotBlank("") { //应该是false
		t.Log("TestStrNotBlank-04 fail")
		t.Fail()
	}

	if strkit.StrNotBlank() { //应该是false
		t.Log("TestStrNotBlank-05 fail")
		t.Fail()
	}

	if !strkit.StrNotBlank("   ") { //应该是true
		t.Log("TestStrNotBlank-06 fail")
		t.Fail()
	}
}

func TestStrIsBlank(t *testing.T) {
	if strkit.StrIsBlank("a") { //应该是false
		t.Log("TestStrIsBlank-01 fail")
		t.Fail()
	}

	if strkit.StrIsBlank("a", "15") { //应该是false
		t.Log("TestStrIsBlank-02 fail")
		t.Fail()
	}

	if strkit.StrIsBlank("a", "", "     ") { //应该是false
		t.Log("TestStrIsBlank-03 fail")
		t.Fail()
	}

	if strkit.StrIsBlank("", "   ") { //应该是false
		t.Log("TestStrIsBlank-04 fail")
		t.Fail()
	}

	if strkit.StrIsBlank() { //应该是false
		t.Log("TestStrIsBlank-05 fail")
		t.Fail()
	}

	if !strkit.StrIsBlank("", "") { //应该是true

		t.Log("TestStrIsBlank-06 fail")
		t.Fail()
	}
}

func TestStrJoin(t *testing.T) {
	if !(strkit.StrJoin("hello ", "world", "", " ", "is go write") == "hello world is go write") {
		t.Log("TestStrJoin-01 fail")
		t.Fail()
	}
}

func TestStringBuilder_ToString(t *testing.T) {
	sb := strkit.StringBuilder{}
	content := sb.Append("hello").Append(" world is ").Append("go write").ToString()
	if !(content == strkit.StrJoin("hello ", "world", "", " ", "is go write")) {
		t.Log("TestStrJoin-01 fail")
		t.Fail()
	}
}

func TestGetStrLen(t *testing.T) {
	enLen := strkit.GetStrLen("hello ")
	if enLen != 6 {
		t.Fail()
	}

	cnLen := strkit.GetStrLen("你好,Go")
	if cnLen != 5 {
		t.Fail()
	}
}
