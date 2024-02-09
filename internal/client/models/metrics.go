package models

import "time"

type MetricsGRPC struct {
	ProcessBegin time.Time
	ProcessEnd   time.Time

	// Domain resolution time
	DomainResolutionBegin time.Time
	DomainResolutionEnd   time.Time

	HostName   string   //
	ResolvedIP []string // The resolved ip
	UsedIP     string   // The real IP address used

	// Connection establishment time
	ConnectedBegin time.Time
	ConnectedEnd   time.Time

	// 握手时间（其中包含有ssl的、非ssl的，取决于是否启用tls方式通信）
	HandshakeBegin time.Time
	HandshakeEnd   time.Time

	// SSL handshake time
	TLSHandshakeBegin time.Time
	TLSHandshakeEnd   time.Time

	// Request time
	RequestBegin time.Time
	RequestEnd   time.Time

	// Response time
	ResponseBegin    time.Time
	FirstResponseEnd time.Time // TODO: first response time
	ResponseEnd      time.Time

	// The length of request
	SendDataLength int64 // bytes
	// The length of Receive
	ReceiveDataLength int64 // bytes
}
