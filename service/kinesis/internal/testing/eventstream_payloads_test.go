package testing

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/smithy-go/encoding/cbor"
)

// Protocol ID names the eventstream fixtures know how to encode payloads for.
// These mirror the Name of the shape ID each protocol implementation returns
// from its ID() method.
const (
	protocolCBOR   = "rpcv2Cbor"
	protocolJSON10 = "awsJson1_0"
	protocolJSON11 = "awsJson1_1"
)

// clientWireFormat reports the protocol ID name the client's resolved protocol
// uses to (de)serialize eventstream payloads. The fixtures encode payloads to
// match, so they follow whatever protocol the client is configured with.
func clientWireFormat(svc *kinesis.Client) string {
	return svc.Options().Protocol.ID().Name
}

// objectPayload encodes a string-keyed object as either a JSON object or a
// CBOR map, matching the given protocol. Values are limited to what the
// eventstream fixtures need (strings today); extend as needed.
func objectPayload(protocol string, fields map[string]string) []byte {
	switch protocol {
	case protocolCBOR:
		m := cbor.Map{}
		for k, v := range fields {
			m[k] = cbor.String(v)
		}
		return cbor.Encode(m)
	case protocolJSON10, protocolJSON11:
		// Marshaling a map[string]string yields a deterministic (key-sorted)
		// JSON object, which is all the fixtures require.
		b, err := json.Marshal(fields)
		if err != nil {
			panic(fmt.Sprintf("kinesis eventstream tests: marshal json payload: %v", err))
		}
		return b
	default:
		panic(fmt.Sprintf("kinesis eventstream tests: unsupported protocol %q", protocol))
	}
}

// emptyPayload encodes an empty object, used for the initial-response message
// and empty unknown events.
func emptyPayload(protocol string) []byte {
	switch protocol {
	case protocolCBOR:
		return cbor.Encode(cbor.Map{})
	case protocolJSON10, protocolJSON11:
		return []byte(`{}`)
	default:
		panic(fmt.Sprintf("kinesis eventstream tests: unsupported protocol %q", protocol))
	}
}
