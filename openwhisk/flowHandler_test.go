package openwhisk

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func doFlow(ts *httptest.Server, value string, workflow string) {
	if value == "" {
		value = `{"name":"Meteion"}`
	}
	resp, status, err := doPost(ts.URL+"/flow",
		`{ "value": `+value+`, "workflow": `+workflow+`}`)
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
	doFlow(ts, `{"name":"Meteion"}`, `{"enabled":true}`)
	time.Sleep(2 * time.Second)
	doFlow(ts, `{"name":"World"}`, `{"enabled":true}`)
	stopTestServer(ts, cur, log)
}
