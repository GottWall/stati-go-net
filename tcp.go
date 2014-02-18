package stati_net

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

const (

	// Stream auth header delimiter
	DEFAULT_TCP_AUTH_DELIMITER string = "--stream-auth--"

	// Chunks delimiter
	DEFAULT_TCP_CHUNK_DELIMITER string = "--chunk--"

	// Setup default TCP backend port
	DEFAULT_TCP_PORT int16 = 8897
)

type TCPClient struct {
	*Client
	AuthDelimiter  string
	ChunkDelimiter string
	Authenticated  bool
	CurrentConn    *net.Conn
}

func TCPClientInit(project string, private_key string, public_key string, host string,
	port int16, auth_delimiter string, chunk_delimiter string) *TCPClient {

	client := &TCPClient{
		Client: ClientInit(project, private_key, public_key, host, port),

		AuthDelimiter:  auth_delimiter,
		ChunkDelimiter: chunk_delimiter,
	}
	return client
}

func (c *TCPClient) Incr(name string, value float64, timestamp int64, filters map[string]interface{}) bool {
	return c.Request(c.MakePacket(c.MakeMessage("incr", name, value, timestamp, filters)))
}

func (c *TCPClient) Decr(name string, value float64, timestamp int64, filters map[string]interface{}) bool {
	return c.Request(c.MakePacket(c.MakeMessage("incr", name, value, timestamp, filters)))
}

// Marshal Message to json object
func (c *TCPClient) MakePacket(message *Message) string {
	serialized_message, err := c.SerializeMessage(message)
	if err != nil {
		return ""
	}

	return strings.Join([]string{string(*serialized_message), c.ChunkDelimiter}, "")
}

// Send metric package to connection
func (c *TCPClient) Request(packet string) bool {
	conn, err := c.GetConnection()
	if err != nil {
		return false
	}

	if !c.isAuthenticated() {
		return false
	}

	fmt.Fprintf(*conn, packet)
	return true
}

// Check that connection is authenticated
func (c *TCPClient) isAuthenticated() bool {
	if c.Authenticated == false {
		return c.Authenticate()
	}
	return true
}

// Check that already established connection not closed
func (c *TCPClient) ok() bool {
	if c.CurrentConn == nil {
		return false
	}

	var buf []byte

	if _, err := bufio.NewReader(*c.CurrentConn).Read(buf); err == io.EOF {
		return false
	}
	return true
}

// Get point to current established TCP connection or create new
func (c *TCPClient) GetConnection() (*net.Conn, error) {

	if c.CurrentConn != nil && c.ok() {
		return c.CurrentConn, nil
	}

	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		return nil, err
	}

	c.CurrentConn = &conn

	// Reset authentication for new connection
	c.Authenticated = false
	return c.CurrentConn, nil
}

// Authenticate connection
func (c *TCPClient) Authenticate() bool {

	if c.Authenticated == false {

		conn, err := c.GetConnection()

		fmt.Fprintf(*conn, c.GetAuthHeaderChunk())
		status, err := bufio.NewReader(*conn).Peek(2)

		if err != nil {
			c.Authenticated = false
			return true
		}

		if string(status) == "OK" {
			c.Authenticated = true
			return true
		} else {
			c.Authenticated = false
			return false
		}
	}
	return true
}

func (c *TCPClient) GetAuthHeaderChunk() string {
	return c.GetAuthHeader() + c.AuthDelimiter
}
