package stati_net

import (
	"strconv"
	"strings"
	"testing"
)

var (
	solt_base   int    = DEFAULT_SOLT_BASE
	project     string = "test_project"
	private_key string = "private_key"
	public_key  string = "public_key"
)

func TestClient(t *testing.T) {

	var client *Client = &Client{
		PrivateKey: private_key,
		PublicKey:  public_key,
		SoltBase:   solt_base,
		Project:    project,
	}

	var header string = client.GetAuthHeader()
	var header_pieces []string = strings.Split(header, " ")

	timestamp, _ := strconv.Atoi(header_pieces[1])
	solt_base, _ := strconv.Atoi(header_pieces[3])

	if MakeSign(private_key, MakeSignMsg(public_key, GetSolt(int64(timestamp), solt_base))) != header_pieces[2] {

		t.Fatalf("Invalid auth header generator")
	}
}
