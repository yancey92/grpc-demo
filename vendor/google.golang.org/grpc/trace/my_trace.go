// Copyright © 2024 yangxinxin_mail@163.com. All right reserved.
package trace

// unique type to prevent assignment.
type HandshakeHandlerContextKey struct{}

type HandshakeHandler struct {
	HandshakeStart    func()
	HandshakeDone     func(error)
	TLSHandshakeStart func(serverName string)
	TLSHandshakeDone  func(serverName string, err error)
}

// unique type to prevent assignment.
type RequestHandlerContextKey struct{}

type RequestHandler struct {
	RequestStart  func()
	RequestDone   func(error)
	ResponseStart func()
	ResponseDone  func(error)
}
