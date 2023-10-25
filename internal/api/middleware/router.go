package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"runtime/debug"
	"time"

	"test-zero-agency/internal/entity"

	chiM "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type IMiddleware interface {
	PanicRecovery(next http.Handler) http.Handler
	RequestLogger(next func(w http.ResponseWriter, r *http.Request) ([]byte, int, *entity.LogicError)) func(http.ResponseWriter, *http.Request)
	ContentTypeJSON(next http.Handler) http.Handler
}

type middleware struct {
	logger *zap.Logger
}

func NewMiddleware(logger *zap.Logger) IMiddleware {
	middleware := new(middleware)
	middleware.logger = logger
	return middleware
}

func (m *middleware) ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) PanicRecovery(next http.Handler) http.Handler {
	timeStart := time.Now()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				resp := []byte("{\"error\": \"InternalServerError\"}")
				respCode := 500

				// Ответ клиенту
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(respCode)
				w.Write(resp)

				// Записываем в logger
				requestDump, _ := httputil.DumpRequest(r, true)

				m.logger.DPanic("Panic Recovery",
					zap.String("LeadTime", fmt.Sprintf("%.3f", time.Duration(time.Now().UnixNano()-timeStart.UnixNano()).Seconds())),
					zap.String("RequestMethod", r.Method),
					zap.Any("LogicError", err),
					zap.String("URL", r.URL.RequestURI()),
					zap.Int32("ResponseCode", int32(respCode)),
					zap.String("ResponseBody", string(resp)),
					zap.String("RemoteAddr", r.RemoteAddr),
					zap.String("UserAgent", r.UserAgent()),
					zap.String("RequestDump", string(requestDump)),
					zap.String("StackTrace", string(debug.Stack())),
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) RequestLogger(next func(w http.ResponseWriter, r *http.Request) ([]byte, int, *entity.LogicError)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		timeStart := time.Now()
		requestDump, _ := httputil.DumpRequest(r, true)

		response, code, logicError := next(w, r)

		fields := []zapcore.Field{
			zap.String("RequestId", chiM.GetReqID(r.Context())),
			zap.String("LeadTime", fmt.Sprintf("%.3f", time.Duration(time.Now().UnixNano()-timeStart.UnixNano()).Seconds())),
			zap.String("RequestMethod", r.Method),
			zap.String("URL", r.URL.RequestURI()),
			zap.Int("ResponseCode", code),
			zap.String("ResponseBody", string(response)),
			zap.String("RemoteAddr", r.RemoteAddr),
			zap.String("UserAgent", r.UserAgent()),
			zap.String("RequestDump", string(requestDump)),
		}

		if logicError == nil {

			m.logger.Info("Request Logger",
				fields...,
			)
		} else {
			fields = append(fields, []zapcore.Field{
				zap.String("LogicError", logicError.Error()),
				zap.String("StackTrace", logicError.StackTrace)}...,
			)

			m.logger.Error("Request Logger",
				fields...,
			)
		}
	}
}
