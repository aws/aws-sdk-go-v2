package json

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// GetErrorInfo util looks for code, __type, and message members in the
// json body. These members are optionally available, and the function
// returns the value of member if it is available. This function is useful to
// identify the error code, msg in a REST JSON error response.
func GetErrorInfo(decoder *json.Decoder) (string, string, error) {

	var code, typeCode, msg string

	startToken, err := decoder.Token()
	if err == io.EOF {
		return "", "", nil
	}
	if err != nil {
		return "", "", err
	}

	if t, ok := startToken.(json.Delim); !ok || t.String() != "{" {
		return "", "", fmt.Errorf("expected start token to be {")
	}

	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			return "", "", err
		}

		switch t {
		case "code":
			v, err := decoder.Token()
			if err != nil {
				return "", "", err
			}
			code = v.(string)
			break
		case "message":
			v, err := decoder.Token()
			if err != nil {
				return "", "", err
			}
			msg = v.(string)
			break
		case "__type":
			v, err := decoder.Token()
			if err != nil {
				return "", "", err
			}
			typeCode = v.(string)
			break
		default:
			DiscardUnknownField(decoder)
			break
		}
	}

	endToken, err := decoder.Token()
	if err != nil {
		return "", "", err
	}

	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return "", "", fmt.Errorf("expected end token to be }")
	}

	if len(code) == 0 {
		return typeCode, msg, nil
	}
	return code, msg, nil
}

// SanitizeErrorCode sanitizes the errorCode string .
// The rule for sanitizing is if a `:` character is present, then take only the
// contents before the first : character in the value.
// If a # character is present, then take only the contents after the
// first # character in the value.
func SanitizeErrorCode(errorCode string) string {
	if strings.ContainsAny(errorCode, ":") {
		errorCode = strings.SplitN(errorCode, ":", 2)[0]
	}

	if strings.ContainsAny(errorCode, "#") {
		errorCode = strings.SplitN(errorCode, "#", 2)[1]
	}

	return errorCode
}

// DiscardUnknownField discards unknown fields from decoder body.
// This function is useful while deserializing json body with additional
// unknown information that should be discarded.
func DiscardUnknownField(decoder *json.Decoder) error {
	v, err := decoder.Token()
	if err == io.EOF {
		return nil
	}
	if err != nil{
		return err
	}

	if _, ok := v.(json.Delim); ok {
		for decoder.More() {
			err = DiscardUnknownField(decoder)
		}
		endToken, err := decoder.Token()
		if err != nil{
			return err
		}
		if _, ok := endToken.(json.Delim); !ok {
			return fmt.Errorf("invalid JSON : expected json delimiter, found %T %v",
				endToken, endToken)
		}
	}

	return err
}
