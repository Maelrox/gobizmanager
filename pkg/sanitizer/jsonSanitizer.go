package sanitizer

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

type captureWriter struct {
	http.ResponseWriter
	status int
	buf    bytes.Buffer
}

func (w *captureWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

func (w *captureWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

// Sanitize output (for HTML, adjust policy as needed)
func SanitizeOutput(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := &captureWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(cw, r)

		p := bluemonday.UGCPolicy()
		clean := p.SanitizeBytes(cw.buf.Bytes())

		w.WriteHeader(cw.status)
		w.Write(clean)
	})
}

func SanitizeInput(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == "application/json" {
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
				p := bluemonday.StrictPolicy()
				for k, v := range body {
					if str, ok := v.(string); ok {
						body[k] = strings.TrimSpace(p.Sanitize(str))
					}
				}
				b, _ := json.Marshal(body)
				r.Body = io.NopCloser(bytes.NewReader(b))
			}
		}
		next.ServeHTTP(w, r)
	})
}
