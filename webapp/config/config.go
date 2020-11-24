//Package config :general config of the server/logging
package config

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

type key int

const (
	//PORT :
	PORT = ":5221"
	//INDEX :
	INDEX = "templates/index.html"
	//NCOOKIE : Name of cookie
	NCOOKIE = "vbscookie"
	//KEY : TLS KEY FILE
	KEY = "key.pem"
	//CERT : TLS CERT FILE
	CERT = "cert.pem"
	//LOG : Logging file
	LOG              = "log.txt"
	requestIDKey key = 0
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../")
	//KeyPath : TLS KEY PATH
	KeyPath = Root + "/" + KEY
	//CertPath : TLS CERT PATH
	CertPath = Root + "/" + CERT
)

//Logging :
type Logging struct {
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

// CreateLogging : create custom logging obj for controllers.
//models logger can also be created seperately
func CreateLogging() *Logging {
	logFile, err := os.OpenFile(LOG, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	logging := Logging{}

	trace := log.New(mw, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	info := log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warning := log.New(mw, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	error := log.New(io.MultiWriter(mw), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	logging.Trace = trace
	logging.Info = info
	logging.Warning = warning
	logging.Error = error
	return &logging
}

//Infologging :closure for http handler info logging
func (logger *Logging) Infologging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Trace.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

//Tracing : closure for http handler to generate request id and maintain context
func Tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
