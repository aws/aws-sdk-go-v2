package v4

import "testing"

func TestAllowedQueryHoisting(t *testing.T) {
	cases := map[string]struct {
		Header      string
		ExpectHoist bool
	}{
		"object-lock": {
			Header:      "X-Amz-Object-Lock-Mode",
			ExpectHoist: false,
		},
		"s3 metadata": {
			Header:      "X-Amz-Meta-SomeName",
			ExpectHoist: false,
		},
		"another header": {
			Header:      "X-Amz-SomeOtherHeader",
			ExpectHoist: true,
		},
		"non X-AMZ header": {
			Header:      "X-SomeOtherHeader",
			ExpectHoist: false,
		},
		"checksum SHA256": {
			Header:      "X-Amz-Checksum-Sha256",
			ExpectHoist: false,
		},
		"checksum CRC32": {
			Header:      "X-Amz-Checksum-Crc32",
			ExpectHoist: false,
		},
		"checksum CRC32C": {
			Header:      "X-Amz-Checksum-Crc32c",
			ExpectHoist: false,
		},
		"checksum CRC64NVME": {
			Header:      "X-Amz-Checksum-Crc64nvme",
			ExpectHoist: false,
		},
		"checksum SHA1": {
			Header:      "X-Amz-Checksum-Sha1",
			ExpectHoist: false,
		},
		"checksum SHA512": {
			Header:      "X-Amz-Checksum-Sha512",
			ExpectHoist: false,
		},
		"checksum future algorithm": {
			Header:      "X-Amz-Checksum-Somefuturealgorithm",
			ExpectHoist: false,
		},
		// x-amz-checksum-mode is a read-side request flag on GetObject/HeadObject,
		// not a checksum value header. It must be hoisted (not signed) so presigned
		// GET URLs work with plain HTTP clients that never send the header.
		// Regression test for https://github.com/aws/aws-sdk-go-v2/issues/3528.
		"checksum mode": {
			Header:      "X-Amz-Checksum-Mode",
			ExpectHoist: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if e, a := c.ExpectHoist, AllowedQueryHoisting.IsValid(c.Header); e != a {
				t.Errorf("expect hoist %v, was %v", e, a)
			}
		})
	}
}

func TestIgnoredHeaders(t *testing.T) {
	cases := map[string]struct {
		Header        string
		ExpectIgnored bool
	}{
		"expect": {
			Header:        "Expect",
			ExpectIgnored: true,
		},
		"authorization": {
			Header:        "Authorization",
			ExpectIgnored: true,
		},
		"X-AMZ header": {
			Header:        "X-Amz-Content-Sha256",
			ExpectIgnored: false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if e, a := c.ExpectIgnored, IgnoredHeaders.IsValid(c.Header); e == a {
				t.Errorf("expect ignored %v, was %v", e, a)
			}
		})
	}
}
