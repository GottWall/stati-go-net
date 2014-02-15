package stati_net

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Create authentication header
func MakeAuthHeader(timestamp int64, signature string, project string, solt_base int) string {
	// TODO: convert project to base64 to prevent spaces in project variable
	return fmt.Sprintf(AUTH_HEADER_TEMPLATE, timestamp, signature, solt_base, project)
}

// Serialize metric message
func SerializeMessage(project string, action string, name string, timestamp int64, value float64) (serialized_message []byte, err error) {

	message := Message{
		Action:    action,
		Name:      name,
		Project:   project,
		Timestamp: timestamp,
		Value:     value,
	}

	serialized_message, err = json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return serialized_message, nil
}

// Get solt from timestamp
func GetSolt(timestamp int64, solt_base int) int64 {
	return int64((timestamp / int64(solt_base)) * int64(solt_base))
}

// Make signature message
func MakeSignMsg(public_key string, solt int64) string {
	return fmt.Sprintf("%s%d", public_key, solt)
}

// Make hexdecimal signature for sign_msg
func MakeSign(private_key string, sign_msg string) string {
	mac := hmac.New(md5.New, []byte(private_key))
	mac.Write([]byte(sign_msg))
	return hex.EncodeToString(mac.Sum(nil))
}
