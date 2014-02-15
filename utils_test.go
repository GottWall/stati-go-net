package stati_net

import (
	"fmt"
	"testing"
)

var (
	timestamp int64  = 1392454739
	signature string = "signature"
)

func TestMakeAuthHeader(t *testing.T) {

	res := MakeAuthHeader(timestamp, signature, project, solt_base)
	check_res := fmt.Sprintf(AUTH_HEADER_TEMPLATE, timestamp, signature, solt_base, project)

	if res != check_res {
		t.Fatalf("Auth header failed: %s != %s", res, check_res)
	}
}

func TestSerializeMessage(t *testing.T) {
	var check_res string
	serialized_message, _ := SerializeMessage(project, "incr", "test_name", timestamp, 10.1)

	check_res = fmt.Sprintf("{\"a\":\"%s\",\"n\":\"%s\",\"p\":\"%s\",\"ts\":%d,\"v\":%.1f,\"f\":null}",
		"incr", "test_name", project, timestamp, 10.1)

	if string(serialized_message) != check_res {
		t.Fatalf("Invalid serialization: %s != %s", serialized_message, check_res)
	}
}

func TestGetSolt(t *testing.T) {

	if GetSolt(timestamp, solt_base) != 1392454000 {
		t.Fatalf("Invalid solt: %s != %s", GetSolt(timestamp, solt_base), 1392454000)
	}

}

func TestMakeSignMsg(t *testing.T) {

	if MakeSignMsg("public_key", 1392454000) != fmt.Sprintf("%s%d", "public_key", 1392454000) {
		t.Fatalf("Invalid make signature")
	}
}

func TestMakeSign(t *testing.T) {

	var sign_msg string = MakeSignMsg(public_key, GetSolt(timestamp, solt_base))
	var sign string = MakeSign(private_key, sign_msg)
	if sign != "5134790630290ed958e5413d526d1915" {
		t.Fatalf("Invalid signature: %s != %s", sign, "5134790630290ed958e5413d526d1915")
	}
}
