package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/harboe/gogeo/geo"
)

func Logger(f geo.Fetcher) geo.Fetcher {
	return geo.FetcherFunc(func(req *http.Request) (*http.Response, error) {
		start := time.Now()
		resp, err := f.Do(req)
		end := time.Now()
		status, size := 0, float64(0)

		if resp != nil {
			status = resp.StatusCode
			size = float64(resp.ContentLength) / 1024
		}

		entry := fmt.Sprintf("%24s | %3d | %12v | %8.2fKb | %s",
			strings.Replace(strings.ToUpper(req.URL.Host), "WWW.", "", 1),
			status,
			end.Sub(start),
			size,
			req.URL,
		)

		if err != nil {
			msg := err.Error()

			if i := strings.Index(msg, ": "); i > 0 {
				msg = msg[i+2:]
			}

			entry += " " + msg
		}

		log.Println(entry)
		return resp, err
	})
}
