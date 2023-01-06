package ssocreds

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-cmp/cmp"
)

var tokenCmpOptions = cmp.Options{
	cmp.AllowUnexported(token{}, tokenKnownFields{}, rfc3339{}),
}

func TestStandardSSOCacheTokenFilepath(t *testing.T) {
	origHomeDur := osUserHomeDur
	defer func() {
		osUserHomeDur = origHomeDur
	}()

	cases := map[string]struct {
		key            string
		osUserHomeDir  func() string
		expectFilename string
		expectErr      string
	}{
		"success": {
			key: "https://example.awsapps.com/start",
			osUserHomeDir: func() string {
				return os.TempDir()
			},
			expectFilename: filepath.Join(os.TempDir(), ".aws", "sso", "cache",
				"e8be5486177c5b5392bd9aa76563515b29358e6e.json"),
		},
		"failure": {
			key: "https://example.awsapps.com/start",
			osUserHomeDir: func() string {
				return ""
			},
			expectErr: "some error",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			osUserHomeDur = c.osUserHomeDir

			actual, err := StandardCachedTokenFilepath(c.key)
			if c.expectErr != "" {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if e, a := c.expectFilename, actual; e != a {
				t.Errorf("expect %v filename, got %v", e, a)
			}
		})
	}
}

func TestLoadCachedToken(t *testing.T) {
	cases := map[string]struct {
		filename    string
		expectToken token
		expectErr   string
	}{
		"file not found": {
			filename:  filepath.Join("testdata", "does_not_exist.json"),
			expectErr: "failed to read cached SSO token file",
		},
		"invalid json": {
			filename:  filepath.Join("testdata", "invalid_json.json"),
			expectErr: "failed to parse cached SSO token file",
		},
		"missing accessToken": {
			filename:  filepath.Join("testdata", "missing_accessToken.json"),
			expectErr: "must contain accessToken and expiresAt fields",
		},
		"missing expiresAt": {
			filename:  filepath.Join("testdata", "missing_expiresAt.json"),
			expectErr: "must contain accessToken and expiresAt fields",
		},
		"standard token": {
			filename: filepath.Join("testdata", "valid_token.json"),
			expectToken: token{
				tokenKnownFields: tokenKnownFields{
					AccessToken:  "dGhpcyBpcyBub3QgYSByZWFsIHZhbHVl",
					ExpiresAt:    (*rfc3339)(aws.Time(time.Date(2044, 4, 4, 7, 0, 1, 0, time.UTC))),
					ClientID:     "client id",
					ClientSecret: "client secret",
					RefreshToken: "refresh token",
				},
				UnknownFields: map[string]interface{}{
					"unknownField":          "some value",
					"registrationExpiresAt": "2044-04-04T07:00:01Z",
					"region":                "region",
					"startURL":              "start URL",
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actualToken, err := loadCachedToken(c.filename)
			if c.expectErr != "" {
				if err == nil {
					t.Fatalf("expect %v error, got none", c.expectErr)
				}
				if e, a := c.expectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect %v error, got %v", e, a)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if diff := cmp.Diff(c.expectToken, actualToken, tokenCmpOptions...); diff != "" {
				t.Errorf("expect tokens match\n%s", diff)
			}
		})
	}
}

func TestStoreCachedToken(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "aws-sdk-go-v2-"+t.Name())
	if err != nil {
		t.Fatalf("failed to create temporary test directory, %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("failed to cleanup temporary test directory, %v", err)
		}
	}()

	cases := map[string]struct {
		token    token
		filename string
		fileMode os.FileMode
	}{
		"standard token": {
			filename: filepath.Join(tempDir, "token_file.json"),
			fileMode: 0600,
			token: token{
				tokenKnownFields: tokenKnownFields{
					AccessToken:  "dGhpcyBpcyBub3QgYSByZWFsIHZhbHVl",
					ExpiresAt:    (*rfc3339)(aws.Time(time.Date(2044, 4, 4, 7, 0, 1, 0, time.UTC))),
					ClientID:     "client id",
					ClientSecret: "client secret",
					RefreshToken: "refresh token",
				},
				UnknownFields: map[string]interface{}{
					"unknownField":          "some value",
					"registrationExpiresAt": "2044-04-04T07:00:01Z",
					"region":                "region",
					"startURL":              "start URL",
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := storeCachedToken(c.filename, c.token, c.fileMode)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			actual, err := loadCachedToken(c.filename)
			if err != nil {
				t.Fatalf("failed to load stored token, %v", err)
			}

			if diff := cmp.Diff(c.token, actual, tokenCmpOptions...); diff != "" {
				t.Errorf("expect tokens match\n%s", diff)
			}
		})
	}
}
