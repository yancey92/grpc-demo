// Copyright © 2024 Yangxinxin. All right reserved.
// Confidential and Proprietary
package strkit

import (
	"bytes"
	"strconv"
	"unicode/utf8"
)

// 判断多个字符串是否不为空, 有一个为空则返回false，全部不为空才返回true
// StrNotBlank("a", "15") return true
func StrNotBlank(strs ...string) bool {
	if len(strs) == 0 {
		return false
	}
	for _, v := range strs {
		if v == "" {
			return false
		}
	}
	return true
}

// 判断多个字符串是否为空，全为空则为true；如果有一个不为空则返回false
// StrIsBlank("a", "") return false
func StrIsBlank(strs ...string) bool {
	if len(strs) == 0 {
		return false
	}
	for _, v := range strs {
		if v != "" {
			return false
		}
	}
	return true
}

// Int64 string to int64
func StrToInt64(str string) (int64, error) {
	v, err := strconv.ParseInt(str, 10, 64)
	return int64(v), err
}

// Int string to int
func StrToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

// 多个字符串拼接
// StrJoin("hello ", "world", "", " ", "is go write")
// return hello world is go write
func StrJoin(strs ...string) string {
	if len(strs) == 0 {
		return ""
	}
	var strBuffer bytes.Buffer
	for _, v := range strs {
		if StrNotBlank(v) {
			strBuffer.WriteString(v)
		}
	}
	return strBuffer.String()
}

// 获取字符串长度
// GetStrLen("hello ")
// return 6
func GetStrLen(str string) int {
	return utf8.RuneCountInString(str)
}

// 字符串构建对象
// 实例构造方法: StringBuilder{}
type StringBuilder struct {
	buf bytes.Buffer
}

// 添加字符串到字符串构建实例里面,空字符串将会被忽略
// 实例构造方法: sb.Append("hello").Append(" world")
func (sb *StringBuilder) Append(str string) *StringBuilder {
	if StrNotBlank(str) {
		sb.buf.WriteString(str)
	}
	return sb
}

// 输出字符串构建实例里面的所有字符串,空字符串将会被忽略
// 实例构造方法: sb.Append("hello").Append(" world").ToString()
func (sb *StringBuilder) ToString() string {
	return sb.buf.String()
}

// 获取字符串长度
// start开始下标，end结束下标（包含）
// SubStr("20170620120101", 0,6), return 201706
func SubStr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return ""
	}

	if end < 0 || end > length {
		return ""
	}
	return string(rs[start:end])
}
