/*
*
network-plugin-sm
Tingyun.com
ibetter
2024.01.02
//
*/
package pkg

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	KHttpHeaderAcceptEncoding  = "Accept-Encoding"
	KHttpHeaderContentEncoding = "Content-Encoding"
	KHttpHeaderContentType     = "Content-Type"
	KHttpEncodingIdentity      = "identity"
	KHttpEncodingGzipDeflate   = "gzip, deflate"
)

const (
	Base64DataHeader = "data:text/plain;base64,"
)

func IsIpv6(hostname string) bool {
	ip := net.ParseIP(hostname)
	return ip != nil && ip.To16() != nil && ip.To4() == nil
}

func IsIpv4(hostname string) bool {
	ip := net.ParseIP(hostname)
	return ip != nil && ip.To4() != nil
}

func IsIp(hostname string) bool {
	return IsIpv4(hostname) || IsIpv6(hostname)
}

func CeilMicroToMilli(value int64) int64 {
	if value%1000 > 0 {
		value += 1000
	}
	return value / 1000
}

func CeilFloat(value float64) int64 {
	return CeilMicroToMilli(int64(value * 1000))
}

func CeilFloatInt(value float64) int {
	return int(CeilFloat(value))
}

func SinceMicroSecond(start, end time.Time) int64 {
	return end.Sub(start).Microseconds()
}

func SinceMiliSecond(start, end time.Time) float64 {
	var mic = SinceMicroSecond(start, end)

	return float64(mic) / 1000
}

func CeilSinceMiliSecond(start, end time.Time) int {
	var mic = int64(SinceMiliSecond(start, end) * 1000)

	if mic%1000 > 0 {
		mic += 1000
	}
	return int(mic / 1000)
}

func Hostname2IpPort(hostname string, port *string) string {
	hostname = strings.Trim(hostname, "[]")
	if IsIpv6(hostname) {
		return hostname
	}

	var index = strings.LastIndex(hostname, ":")
	if index < 0 {
		return hostname
	}

	if port != nil {
		*port = hostname[index+1:]
	}
	return strings.Trim(hostname[:index], "[]")
}

func ParseInt64(value string, def int64) int64 {
	number, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err == nil {
		return number
	}
	return def
}

func ParseInt(value string, def int) int {
	return int(ParseInt64(value, int64(def)))
}

func ParseBool(value string, def bool) bool {
	ok, err := strconv.ParseBool(value)
	if err != nil {
		return def
	}

	return ok
}

func UgzipBuffer(data *[]byte) []byte {
	var result []byte = nil
	defer func() {
		recover()
	}()

	buf := new(bytes.Buffer)
	buf.Write(*data)

	zr, err := gzip.NewReader(buf)
	if err != nil {
		return result
	}
	defer zr.Close()

	result, _ = io.ReadAll(zr)

	return result
}

func FromBase64(value string) []byte {
	value = strings.TrimPrefix(value, Base64DataHeader)
	result, _ := base64.StdEncoding.DecodeString(value)
	return result
}

func GetBufMD5(value []byte) string {
	hash := md5.Sum(value)
	return strings.ToLower(hex.EncodeToString(hash[:]))
}

func SplitHttp2(header string) (string, string) {
	index := strings.Index(header, ":")
	if index == 0 {
		index = strings.Index(header[1:], ":") + 1
		if index == 1 {
			return "", ""
		}
	}

	if index == -1 {
		return "", ""
	}

	return strings.TrimSpace(header[:index]), strings.TrimSpace(header[index+1:])
}
