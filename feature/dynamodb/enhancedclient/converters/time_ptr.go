package converters

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var _ AttributeConverter[*time.Time] = (*TimePtrConverter)(nil)

// TimePtrConverter implements AttributeConverter for *time.Time values.
//
// Supported DynamoDB representations:
//   - String (AttributeValueMemberS) using a layout provided via the "format" option or any of knowTimeFormats.
//   - Number (AttributeValueMemberN) encoding seconds and optional fractional nanoseconds as "<sec>[.<nanos>]".
//
// Options (case-sensitive keys):
//
//	format: Custom time.Parse / time.Format layout string (Go reference layout). Overrides fallback list when decoding.
//	TZ:     IANA timezone name (e.g., "UTC", "America/New_York") applied after successful parse when decoding.
//	as:     When encoding, chooses representation: "string" (default) or "number".
//
// Nil handling:
//   - Encoding a nil *time.Time returns ErrNilValue.
//   - Decoding a zero time (t.IsZero()) returns (nil, nil) signaling absence.
//
// Fractional seconds when encoding as number are trimmed of trailing zeros.
// The converter is stateless and safe for concurrent use.
type TimePtrConverter struct{}

// FromAttributeValue converts a DynamoDB AttributeValue into a *time.Time using
// either string or numeric representations.
//
// Decoding logic:
//
//	String: Attempts parse with provided "format" opt if present, else tries knowTimeFormats in order.
//	        On parse failure across all formats returns an error listing attempted formats.
//	Number: Splits on '.', first part seconds, second (optional) part nanoseconds (truncated to 9 digits).
//	        Rejects more than one '.' (i.e., len(parts) > 2).
//
// Options:
//
//	format: Single layout used instead of fallback list (string input only).
//	TZ:     IANA timezone applied post-parse; errors if unknown.
//
// Error cases:
//   - Nil underlying AttributeValueMemberS/N -> ErrNilValue
//   - Unsupported AttributeValue type -> unsupportedType error
//   - Invalid numeric format or parse errors -> descriptive errors
//   - Unknown timezone -> error
//
// Returns (nil, nil) if parsed time is the zero value.
func (tc TimePtrConverter) FromAttributeValue(v types.AttributeValue, opts []string) (*time.Time, error) {
	t := time.Time{}

	switch av := v.(type) {
	case *types.AttributeValueMemberS:
		// when calling with v = (*types.AttributeValueMemberS)(nil) then v == nil is false
		// e.g. tc.FromAttributeValue((*types.AttributeValueMemberS)(nil), nil) -> panics
		if av == nil {
			return nil, ErrNilValue
		}

		format := getOpt(opts, "format")
		var formats []string
		if format != "" {
			formats = []string{format}
		} else {
			formats = knowTimeFormats
		}

		var err error
		for _, f := range formats {
			t, err = time.Parse(f, av.Value)
			if err == nil {
				break
			}
		}

		// err will be populated only if all time.Parse() attempts returned an error
		if err != nil {
			return nil, fmt.Errorf("unable to process time %s with format(s): %v", av.Value, formats)
		}
	case *types.AttributeValueMemberN:
		// when calling with v = (*types.AttributeValueMemberN)(nil) then v == nil is false
		// e.g. tc.FromAttributeValue((*types.AttributeValueMemberN)(nil), nil) -> panics
		if av == nil {
			return nil, ErrNilValue
		}

		parts := strings.Split(av.Value, ".")
		// format is "000" or "000.000", anything else is an issue
		if len(parts) > 2 {
			return nil, fmt.Errorf("unsupported format for number inside of types.AttributeValueMemberN: %v", av.Value)
		}

		var err error
		ps := make([]int64, 2)

		for i := range parts {
			// microseconds can be at most 9 chars long, otherwise they overflow into the seconds part
			if i == 1 && len(parts[i]) > 9 {
				parts[i] = parts[i][0:9]
			}
			ps[i], err = strconv.ParseInt(parts[i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing int: %v", parts[i])
			}
		}

		t = time.Unix(ps[0], ps[1])
	default:
		return nil, unsupportedType(
			v,
			(*types.AttributeValueMemberS)(nil),
			(*types.AttributeValueMemberN)(nil),
		)
	}

	if t.IsZero() {
		return nil, nil
	}

	if tz := getOpt(opts, "TZ"); tz != "" {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			return nil, fmt.Errorf(`error loading timezone "%s" data: %v`, tz, err)
		}
		t = t.In(loc)
	}

	return &t, nil
}

// ToAttributeValue converts a *time.Time into a DynamoDB AttributeValue.
//
// Options:
//
//	as:     "number" -> encodes seconds[.nanos] in AttributeValueMemberN, trimming trailing zeros in nanos.
//	        "string" or empty -> encodes formatted string using layout from "format" opt or defaultTimeFormat.
//	format: Go layout string used when encoding as string (ignored for number encoding). Falls back to defaultTimeFormat.
//
// Errors:
//   - Nil input -> ErrNilValue
//   - Unknown "as" value -> error listing expected values
//
// Returned AttributeValue will be *types.AttributeValueMemberS or *types.AttributeValueMemberN depending on representation.
func (tc TimePtrConverter) ToAttributeValue(v *time.Time, opts []string) (types.AttributeValue, error) {
	if v == nil {
		return nil, ErrNilValue
	}

	as := getOpt(opts, "as")
	format := getOpt(opts, "format")
	if format == "" {
		format = defaultTimeFormat
	}

	switch as {
	case "number":
		parts := []string{
			fmt.Sprintf("%v", v.Unix()),
		}
		if v.Nanosecond() != 0 {
			parts = append(parts, strings.TrimRight(fmt.Sprintf("%v", v.Nanosecond()), "0"))
		}
		return &types.AttributeValueMemberN{
			Value: strings.Join(parts, "."),
		}, nil
	case "string", "":
		return &types.AttributeValueMemberS{
			Value: v.Format(format),
		}, nil
	default:
		return nil, fmt.Errorf(`unknown time format: expected "", "string" or "number", got "%v"`, as)
	}
}
