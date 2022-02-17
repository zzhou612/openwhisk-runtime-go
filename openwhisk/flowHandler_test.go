package openwhisk

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func doFlow(ts *httptest.Server, message string) {
	if message == "" {
		message = `{"name":"Meteion"}`
	}
	resp, status, err := doPost(ts.URL+"/flow", `{ "value": `+message+`}`)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%d %s", status, resp)
	}
	if !strings.HasSuffix(resp, "\n") {
		fmt.Println()
	}
}

func TestFlow(t *testing.T) {
	ts, cur, log := startTestServer("")
	buf, _ := Zip("_test/pysample")
	doInit(ts, initBytes(buf, ""))
	doFlow(ts, `{"name":"Meteion"}`)
	time.Sleep(2 * time.Second)
	doFlow(ts, `{"name":"World"}`)
	stopTestServer(ts, cur, log)
}
