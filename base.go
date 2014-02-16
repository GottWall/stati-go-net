package stati_net

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Version struct {
	Major int16
	Minor int16
	Patch int16
}

// Convert version to string representation
func (v *Version) ToString() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

var VERSION *Version = &Version{0, 1, 0}

const (
	DEFAULT_SOLT_BASE int = 1000

	// Auth header in format: GottWallS2 {timestamp} {sign} {solt_base} {project}
	AUTH_HEADER_TEMPLATE string = "GottWallS2 %d %s %d %s"
)

var (
	RequestError = errors.New("Request error")
)

type Client struct {
	PrivateKey string
	PublicKey  string
	Project    string // project name
	SoltBase   int    // GottWall auth header solt base
}

// Metric message structure
type Message struct {
	Action    string                 `json:"a"`
	Name      string                 `json:"n"`
	Project   string                 `json:"p"`
	Timestamp int64                  `json:"ts"`
	Value     float64                `json:"v"`
	Filters   map[string]interface{} `json:"f"`
}

type client interface {
	Incr() bool
	Decr() bool
}

func (c *Client) GetSign(timestamp int64) string {
	return MakeSign(c.PrivateKey,
		MakeSignMsg(c.PublicKey, GetSolt(timestamp, c.SoltBase)))
}

func (c *Client) SetSoltBase(base int) int {
	c.SoltBase = base
	return c.SoltBase
}

// Get auth header
func (c *Client) GetAuthHeader() string {

	var timestamp int64 = time.Now().UTC().Unix()

	return MakeAuthHeader(timestamp, c.GetSign(timestamp), c.GetProjectHash(), c.SoltBase)
}

func (c *Client) GetProjectHash() string {
	return c.Project
}

func (c *Client) SerializeMessage(message *Message) (*[]byte, error) {
	serialized_message, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return &serialized_message, nil
}

// Get current time as seconds int UTC
func (c *Client) CurrentTS() int64 {
	return time.Now().UTC().Unix()
}
