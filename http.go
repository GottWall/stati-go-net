package stati_net

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

const (
	DEFAULT_PROTO string = "http"

	// User agent in format: GottWall-Stati/0.1.0(linux, go/version, public_key)
	USER_AGENT_TEMPLATE string = "GottWall-Stati/%s(%s; %s; %s/%s; %s)"

	// Setup default port for http backend
	DEFAULT_HTTP_PORT int16 = 8890
)

type HTTPClient struct {
	*Client              // Base client
	Prefix        string // GottWall api prefix
	Proto         string // http/https protocol
	RequestClient *http.Client
	UserAgent     string // Request user agent
}

// Construct HTTPClient instance
func HTTPClientInit(project string, private_key string, public_key string, host string, port int16, proto string, prefix string) *HTTPClient {

	prefix = strings.TrimSuffix(prefix, "/")

	if prefix != "" && !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	client := &HTTPClient{
		// Construct base client intance
		Client: ClientInit(project, private_key, public_key, host, port),

		Prefix:        prefix,
		Proto:         proto,
		RequestClient: &http.Client{},
	}
	client.SetUpClient()
	client.SetUserAgent(client.GetDefaultUserAgent())
	return client
}

func (c *HTTPClient) SetUpClient() {
	c.RequestClient = &http.Client{Transport: &http.Transport{DisableCompression: true}}
}

// Construct API url
func (c *HTTPClient) GetUrl(action string) string {
	return fmt.Sprintf("%s://%s:%d%s/api/v1/%s/%s",
		c.Proto, c.Host, c.Port, c.Prefix, c.GetProjectHash(), action)
}

func (c *HTTPClient) SetUserAgent(user_agent string) {
	c.UserAgent = user_agent
}

func (c *HTTPClient) GetUserAgent() string {
	return c.UserAgent
}

func (c *HTTPClient) GetDefaultUserAgent() string {
	return fmt.Sprintf(USER_AGENT_TEMPLATE, VERSION.ToString(), runtime.GOOS, runtime.GOARCH, "go", runtime.Version(), c.PublicKey)
}

func (c *HTTPClient) Incr(name string, value float64, timestamp int64, filters map[string]interface{}) bool {
	return c.Request("incr", name, value, timestamp, filters)
}

func (c *HTTPClient) Decr(name string, value float64, timestamp int64, filters map[string]interface{}) bool {
	return c.Request("decr", name, value, timestamp, filters)
}

func (c *HTTPClient) Request(action string, name string, value float64, timestamp int64, filters map[string]interface{}) bool {

	serialized_message, err := c.SerializeMessage(&Message{
		Action:    action,
		Name:      name,
		Project:   c.Project,
		Timestamp: timestamp,
		Value:     value,
	})

	req, err := http.NewRequest("POST", c.GetUrl(action), bytes.NewReader(*serialized_message))
	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("User-Agent", c.GetUserAgent())
	req.Header.Add("X-GottWall-Auth", c.GetAuthHeader())

	resp, err := c.RequestClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return false
	}

	return true
}
