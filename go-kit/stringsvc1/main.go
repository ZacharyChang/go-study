package main

import (
	// "log"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	svc := stringService{}

	uppercase := makeUppercaseEndpoint(svc)
	// decorator middleware
	uppercase = loggingMiddleware(log.With(logger, "method", "uppercase"))(uppercase)
	uppercaseHandler := httptransport.NewServer(
		uppercase,
		decodeUppercaseRequest,
		encodeResponse,
	)

	count := makeCountEndpoint(svc)
	count = loggingMiddleware(log.With(logger, "method", "count"))(count)
	countHandler := httptransport.NewServer(
		count,
		decodeCountRequest,
		encodeResponse,
	)
	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.ListenAndServe(":8080", nil)
}
