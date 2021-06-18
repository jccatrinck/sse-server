package handlers

import (
	"sse-server/handlers/send"
	"sse-server/handlers/sse"
)

var SSE = sse.Handler
var Send = send.Handler
