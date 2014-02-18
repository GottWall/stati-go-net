package stati_net

import (
	"net"
	"strings"
)

const (

	// Auth header separator
	DEFAULT_UDP_AUTH_DELIMITER = "--chunk-auth--"

	// Metric data part separator
	DEFAULT_UDP_CHUNK_DELIMITER = "--chunk--"

	// Setup default udp client port
	DEFAULT_UDP_PORT int16 = 8897
)

type UDPClient struct {
	*Client        // Base client structure
	AuthDelimiter  string
	ChunkDelimiter string
}

// Construct UDPClient instance

func UDPClientInit(project string, private_key, public_key string, host string, port int16, auth_delimiter string, chunk_delimiter string) *UDPClient {

	client := &UDPClient{
		Client: ClientInit(project, private_key, public_key, host, port),

		AuthDelimiter:  auth_delimiter,
		ChunkDelimiter: chunk_delimiter,
	}
	return client
}

func (c *UDPClient) MakePacket(message *Message) string {
	serialized_message, err := c.SerializeMessage(message)

	if err != nil {
		return ""
	}

	return strings.Join([]string{c.GetAuthHeader(), c.AuthDelimiter,
		string(*serialized_message), c.ChunkDelimiter}, "")
}

func (c *UDPClient) Incr(name string, value float64, timestamp int64, filters map[string]interface{}) bool {
	return c.Request(c.MakePacket(c.MakeMessage("incr", name, value, timestamp, filters)))
}

func (c *UDPClient) Decr(name string, value float64, timestamp int64, filters map[string]interface{}) bool {
	return c.Request(c.MakePacket(c.MakeMessage("decr", name, value, timestamp, filters)))
}

func (c *UDPClient) GetUDPConnection() (net.PacketConn, error) {
	return net.ListenPacket("udp", "127.0.0.1:0")
}

// Make new resolve or get already resolved connection
func (c *UDPClient) GetAddr() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", c.Addr)
}

func (c *UDPClient) Request(packet string) bool {

	if packet == "" {
		return false
	}

	conn, err := c.GetUDPConnection()
	defer conn.Close()

	ra, err := c.GetAddr()

	_, err = conn.(*net.UDPConn).WriteToUDP([]byte(packet), ra)

	if err != nil {
		return false
	}

	return true
}
