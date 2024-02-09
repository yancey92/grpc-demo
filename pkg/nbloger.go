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
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"demo.test/grpc-demo/pkg/strkit"
)

const (
	KPluginLogLevelNone = iota
	KPluginLogLevelError
	KPluginLogLevelWarn
	KPluginLogLevelInfo
	KPluginLogLevelDebug
	KPluginLogLevelAll
)

const (
	KMaxNValues = 31
	KCmdNValue  = 201
	KCmdIndex   = 202
	KCmdNStart  = 203
	KCmdError   = 1002
)

const (
	KCmdIndexLog  = 0
	KCmdIndexIp   = 1
	KCmdIndexHost = 2
	KCmdIndexDns  = 3
)

const (
	KIndexTotal = iota
	KIndexDNS
	KIndexConn
	KIndexSSL
	KIndexRequest
	KIndexFirst  // 响应的第一次读到数据的时长
	KIndexRemain // 响应第一次之后，剩下的时长
	KIndexSpeed  //
	KIndexDownloadBytes
	KIndexCipherId
	KIndexTlsv
	KIndexRequestLength
	KIndexResponseLength
)

const (
	KErrorBase      = 0
	KErrorResolve   = KErrorBase + 12007
	KErrorConnect   = KErrorBase + 12029
	KErrorHandshake = KErrorBase + 12157
	KErrorTimeout   = KErrorBase + 12002
	KErrorReqText   = KErrorBase + 12272
	KErrorMatchMD5  = KErrorBase + 12275
	KErrorHeader    = KErrorBase + 12290
)

var LogPushedUri = ""

type TPluginLoger struct {
	prefix string
	name   string
	level  int
	lock   sync.Mutex
	file   *os.File
	values []int64
}

func (d *TPluginLoger) OpenLoger() {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.values = make([]int64, KMaxNValues)

	d.level = KPluginLogLevelError
	if d.name == "" {
		return
	}

	var err error = nil
	d.file, err = os.Open(d.name)
	if err != nil {
		d.show(fmt.Sprintf("%s|open log: %s\n", d.getTime(), err.Error()))
		return
	}
}
func (d *TPluginLoger) SetLevel(level int) {
	d.level = level
}

func (d *TPluginLoger) SetPrefix(prefix string) {
	d.prefix = prefix
}

func (d *TPluginLoger) loglevel(msg string, level int) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if level != KPluginLogLevelAll && d.level < level {
		return
	}

	slevel := ""
	switch level {
	case KPluginLogLevelError:
		slevel = "E"
	case KPluginLogLevelWarn:
		slevel = "W"
	case KPluginLogLevelInfo:
		slevel = "L"
	case KPluginLogLevelDebug:
		slevel = "D"
	case KPluginLogLevelAll:
		slevel = "I"
	default:
		slevel = "?"
	}

	msg = fmt.Sprintf("%s|%s|%s|%s\n", d.getTime(), d.prefix, slevel, msg)

	d.write(msg)
}

func (d *TPluginLoger) sendMsg(msg string) {
	go func() {
		if strkit.StrIsBlank(LogPushedUri) {
			return
		}
		//d.show("--------" + msg + "\n")
		client := http.Client{Timeout: 1 * time.Second}
		reader := bytes.NewReader([]byte(msg))
		//resp, err :=
		client.Post(LogPushedUri, "text/plain", reader)
		/*
			if err != nil {
				d.show(fmt.Sprintf("%s|send log: %s\n", d.getTime(), err.Error()))
				return
			}
			body, _ := io.ReadAll(resp.Body)
			d.show(fmt.Sprintf("%s|send log: %d, %s\n", d.getTime(), resp.StatusCode, body))
			// */
	}()
	d.loglevel(msg, KPluginLogLevelDebug)
}

func (d *TPluginLoger) SetLog(msg string) {
	d.sendMsg(fmt.Sprintf("%d|%d|%s", KCmdIndex, KCmdIndexLog, msg))
}

func (d *TPluginLoger) SetIp(msg string) {
	d.sendMsg(fmt.Sprintf("%d|%d|%s", KCmdIndex, KCmdIndexIp, msg))
}

func (d *TPluginLoger) SetHost(msg string) {
	d.sendMsg(fmt.Sprintf("%d|%d|%s", KCmdIndex, KCmdIndexHost, msg))
}

func (d *TPluginLoger) SetDns(msg string) {
	d.sendMsg(fmt.Sprintf("%d|%d|%s", KCmdIndex, KCmdIndexDns, msg))
}

func (d *TPluginLoger) SetErrorCode(code int) {
	d.sendMsg(fmt.Sprintf("%d|%d|%d", KCmdError, KCmdIndexLog, code))
}

func (d *TPluginLoger) SetNStart(index int) {
	d.sendMsg(fmt.Sprintf("%d|%d|%s", KCmdNStart, index, ""))
}

func (d *TPluginLoger) SetNValue(index, value int64) {
	d.values[index] = value
	d.sendMsg(fmt.Sprintf("%d|%d|%d", KCmdNValue, index, value))
}

func (d *TPluginLoger) ClearNValue() {
	d.values = make([]int64, KMaxNValues)
}

func (d *TPluginLoger) GetNValue(index int64) int64 {
	return d.values[index]
}

func (d *TPluginLoger) Debug(msg string) {
	d.loglevel(msg, KPluginLogLevelDebug)
}

func (d *TPluginLoger) Debugf(format string, args ...any) {
	d.Debug(fmt.Sprintf(format, args...))
}

func (d *TPluginLoger) Info(msg string) {
	d.loglevel(msg, KPluginLogLevelInfo)
}

func (d *TPluginLoger) Infof(format string, args ...any) {
	d.Info(fmt.Sprintf(format, args...))
}

func (d *TPluginLoger) Warn(msg string) {
	d.loglevel(msg, KPluginLogLevelWarn)
}

func (d *TPluginLoger) Warnf(format string, args ...any) {
	d.Warn(fmt.Sprintf(format, args...))
}

func (d *TPluginLoger) Error(msg string) {
	d.loglevel(msg, KPluginLogLevelError)
}

func (d *TPluginLoger) Errorf(format string, args ...any) {
	d.Error(fmt.Sprintf(format, args...))
}

func (d *TPluginLoger) write(msg string) {

	if d.file == nil {
		d.show(msg)
		return
	}

	d.file.WriteString(msg)
}

func (d *TPluginLoger) CloseLoger() {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.file != nil {
		d.file.Close()
		d.file = nil
	}
}

func (d *TPluginLoger) getTime() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

func (d *TPluginLoger) show(msg string) {
	fmt.Printf("%s", msg)
}
