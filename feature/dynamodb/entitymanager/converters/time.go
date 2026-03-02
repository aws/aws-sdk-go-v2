package converters

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var _ AttributeConverter[time.Time] = (*TimeConverter)(nil)

// taken from go src/time/format.go
var knowTimeFormats = []string{
	time.Layout,      //= "01/02 03:04:05PM '06 -0700" // The reference time, in numerical order.
	time.ANSIC,       //= "Mon Jan _2 15:04:05 2006"
	time.UnixDate,    //= "Mon Jan _2 15:04:05 MST 2006"
	time.RubyDate,    //= "Mon Jan 02 15:04:05 -0700 2006"
	time.RFC822,      //= "02 Jan 06 15:04 MST"
	time.RFC822Z,     //= "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	time.RFC850,      //= "Monday, 02-Jan-06 15:04:05 MST"
	time.RFC1123,     //= "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z,    //= "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	time.RFC3339,     //= "2006-01-02T15:04:05Z07:00"
	time.RFC3339Nano, //= "2006-01-02T15:04:05.999999999Z07:00"
	time.Kitchen,     //= "3:04PM"
	time.Stamp,       //= "Jan _2 15:04:05"
	time.StampMilli,  //= "Jan _2 15:04:05.000"
	time.StampMicro,  //= "Jan _2 15:04:05.000000"
	time.StampNano,   //= "Jan _2 15:04:05.000000000"
	time.DateTime,    //= "2006-01-02 15:04:05"
	time.DateOnly,    //= "2006-01-02"
	time.TimeOnly,    //= "15:04:05"
}

// defaultTimeFormat is the layout used when encoding time values as strings
// if no explicit "format" option is supplied.
var defaultTimeFormat = time.RFC3339Nano

// TimeConverter implements AttributeConverter for time.Time values.
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
// Zero value handling: Decoding to a zero time returns the zero value (not an error). There is
// no nil sentinel for non-pointer time.Time values.
//
// Fractional seconds when encoding as number are trimmed of trailing zeros.
// The converter is stateless and safe for concurrent use.
type TimeConverter struct{}

// FromAttributeValue converts a DynamoDB AttributeValue into a time.Time using either
// string or numeric representations.
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
// Returns zero time if parsed value is zero. Zero time is treated as absence but
// still returned (not an error). Consumers may check t.IsZero().
func (tc TimeConverter) FromAttributeValue(v types.AttributeValue, opts []string) (time.Time, error) {
	t := time.Time{}

	switch av := v.(type) {
	case *types.AttributeValueMemberS:
		// when calling with v = (*types.AttributeValueMemberS)(nil) then v == nil is false
		// e.g. tc.FromAttributeValue((*types.AttributeValueMemberS)(nil), nil) -> panics
		if av == nil {
			return time.Time{}, ErrNilValue
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
			return time.Time{}, fmt.Errorf("unable to process time %s with format(s): %v", av.Value, formats)
		}
	case *types.AttributeValueMemberN:
		// when calling with v = (*types.AttributeValueMemberN)(nil) then v == nil is false
		// e.g. tc.FromAttributeValue((*types.AttributeValueMemberN)(nil), nil) -> panics
		if av == nil {
			return time.Time{}, ErrNilValue
		}

		parts := strings.Split(av.Value, ".")
		// format is "000" or "000.000", anything else is an issue
		if len(parts) > 2 {
			return time.Time{}, fmt.Errorf("unsupported format for number inside of types.AttributeValueMemberN: %v", av.Value)
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
				return time.Time{}, fmt.Errorf("error parsing int: %v", parts[i])
			}
		}

		t = time.Unix(ps[0], ps[1])
	default:
		return time.Time{}, unsupportedType(
			v,
			(*types.AttributeValueMemberS)(nil),
			(*types.AttributeValueMemberN)(nil),
		)
	}

	if t.IsZero() {
		return time.Time{}, nil
	}

	if tz := getOpt(opts, "TZ"); tz != "" {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			return time.Time{}, fmt.Errorf(`error loading timezone "%s" data: %v`, tz, err)
		}
		t = t.In(loc)
	}

	return t, nil
}

// ToAttributeValue converts a time.Time into a DynamoDB AttributeValue.
//
// Options:
//
//	as:     "number" -> encodes seconds[.nanos] in AttributeValueMemberN, trimming trailing zeros in nanos.
//	        "string" or empty -> encodes formatted string using layout from "format" opt or defaultTimeFormat.
//	format: Go layout string used when encoding as string (ignored for number encoding). Falls back to defaultTimeFormat.
//
// Errors:
//   - Unknown "as" value -> error listing expected values
//
// Returned AttributeValue will be *types.AttributeValueMemberS or *types.AttributeValueMemberN depending on representation.
func (tc TimeConverter) ToAttributeValue(v time.Time, opts []string) (types.AttributeValue, error) {
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
		return nil, fmt.Errorf(`unknown value for time format: expected "", "string" or "number", got "%v"`, as)
	}
}
