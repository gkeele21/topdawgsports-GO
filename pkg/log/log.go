package log

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http/httputil"
	"os"
	"runtime"
	"time"
)

var l = log.New(os.Stderr, "", 0)

// Info defines the json structure for the log
type Info struct {
	Host    string    `json:"host"`
	Message string    `json:"message"`
	Errors  []string  `json:"errors"`
	Time    time.Time `json:"time"`
	LogType string    `json:"type"`
}

// Println creates json formatted log for use in logstash
func Println(logType, message string, args ...interface{}) {
	info := Info{
		Message: message,
	}

	for _, arg := range args {
		info.Errors = append(info.Errors, fmt.Sprintf("%v", arg))
	}

	host, _ := os.Hostname()
	info.Host = host

	info.Time = time.Now()
	info.LogType = logType

	data, err := json.Marshal(info)
	if err != nil {
		fmt.Println("couldn't json info:", err)
	}

	l.Println(string(data))
}

func LogRequestData(req echo.Context) {
	fmt.Printf("==== Request Dump %s ====\n", time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
	function, _, _, _ := runtime.Caller(1)
	fmt.Printf("Caller : %#v\n", runtime.FuncForPC(function).Name())
	includeBody := true
	if req == nil {
		includeBody = false
	}
	requestDump, err := httputil.DumpRequest(req.Request(), includeBody)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
	fmt.Println("==== End Request Dump ====")
}
