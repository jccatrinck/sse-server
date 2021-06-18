package send

import (
	"fmt"
	"net/http"

	"sse-server/handlers/sse"
)

func Handler(rw http.ResponseWriter, req *http.Request) {
	msg := req.URL.Query()["msg"][0]

	sse.Send(msg)

	fmt.Fprintf(rw, "%s", "sent msg")
}
