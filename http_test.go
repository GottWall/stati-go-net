package stati_net

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
)

var (
	host   string = "127.0.0.1"
	port   int16  = 8890
	proto  string = "http"
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

// Debug server for listen testing host ans port
func newLocalListener(host string, port int16) net.Listener {
	var serve string = fmt.Sprintf("%s:%d", host, port)
	l, err := net.Listen("tcp", serve)

	if err != nil {
		panic(fmt.Sprintf("httptest: failed to listen on %v: %v", serve, err))
	}
	return l
}

func NewTestHttpServer(handler http.Handler) *httptest.Server {
	ts := NewUnstartedServer(handler)
	ts.Start()
	return ts
}

func NewUnstartedServer(handler http.Handler) *httptest.Server {
	return &httptest.Server{
		Listener: newLocalListener(host, port),
		Config:   &http.Server{Handler: handler},
	}
}

func TestRequest(t *testing.T) {

	ts := NewTestHttpServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))

	defer ts.Close()
	if http_client.Request("incr", "test_name", 10.0, http_client.CurrentTS(), nil) != true {
		t.Fatalf("Invalid request")
	}
	ts.Close()
}

func TestInvalidRequest(t *testing.T) {
	ts := NewTestHttpServer(http.HandlerFunc(http.NotFound))
	defer ts.Close()
	if http_client.Request("incr", "test_name", 10.0, http_client.CurrentTS(), nil) != false {
		t.Fatalf("Invalid request")
	}
	ts.Close()
}
