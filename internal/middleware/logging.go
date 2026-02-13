package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	rec.body.Write(b)
	return rec.ResponseWriter.Write(b)
}

// LoggingMiddleware — та самая функция, которую ищет main.go
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Читаем тело запроса, если оно есть
		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Создаем обертку для ответа
		rec := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           bytes.NewBufferString(""),
		}

		next.ServeHTTP(rec, r)

		// Вывод подробного лога в консоль
		log.Printf("\n--- [API LOG] ---\nPath: %s\nMethod: %s\nReq Body: %s\nStatus: %d\nRes Body: %s\nDuration: %v\n-----------------",
			r.URL.Path,
			r.Method,
			string(requestBody),
			rec.statusCode,
			rec.body.String(), // ТУТ БУДЕТ ОШИБКА БАЗЫ
			time.Since(start),
		)
	})
}
