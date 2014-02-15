package stati_net

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

var (
	host   string = "127.0.0.1"
	port   int16  = 8890
	proto  string = "https"
	prefix string = "/custom_prefix"
)

var http_client *HTTPClient = HTTPClientInit(
	project, private_key, public_key, host, port, proto, prefix)

func TestHttpClient(t *testing.T) {

	if http_client.SoltBase != DEFAULT_SOLT_BASE {
		t.Fatalf("Wrong solt base")
	}

	http_client.SetSoltBase(1200)

	if http_client.SoltBase != 1200 {
		t.Fatalf("Invalid solt base")
	}

	http_client.SetSoltBase(DEFAULT_SOLT_BASE)

	h := http_client.GetProjectHash()
	if h != project {
		t.Fatalf("Invalid project hash: %s != %s", h, project)
	}

	var check_res string = fmt.Sprintf("%s://%s:%d/%s/api/v1/%s/%s",
		proto, host, port, "custom_prefix", h, "incr")

	if check_res != http_client.GetUrl("incr") {
		t.Fatalf("Invalid url")
	}
}

func TestUserAgent(t *testing.T) {
	var user_agent string = fmt.Sprintf(USER_AGENT_TEMPLATE, VERSION.ToString(), runtime.GOOS, runtime.GOARCH, "go", runtime.Version(), http_client.PublicKey)
	if http_client.GetUserAgent() != user_agent {
		t.Fatalf("Invalid user agent")
	}
}

func TestRequest(t *testing.T) {
	if http_client.Request("incr", "test_name", 10.0, time.Now().UTC().Unix(), nil) != true {
		t.Fatalf("Invalid request")
	}
}
