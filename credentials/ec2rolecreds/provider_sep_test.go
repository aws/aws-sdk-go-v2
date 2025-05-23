package ec2rolecreds

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type sepTestCase struct {
	Summary string
	Config  struct {
		EC2InstanceProfileName string
		EnvVars                map[string]string
	}
	Expectations []struct {
		Get      string
		Response struct {
			Status int
			Body   any
		}
	}
	Outcomes []struct {
		Result    string
		AccountID string
	}
}

type sepMockIMDS struct {
	testCase  sepTestCase
	callIndex int
}

func (m *sepMockIMDS) GetMetadata(ctx context.Context, in *imds.GetMetadataInput, opts ...func(*imds.Options)) (*imds.GetMetadataOutput, error) {
	callIndex := m.callIndex
	next := m.testCase.Expectations[callIndex]
	m.callIndex++

	expectPath := strings.TrimPrefix(next.Get, "/latest/meta-data") // the real IMDS client injects this
	if strings.HasSuffix(expectPath, "security-credentials") ||
		strings.HasSuffix(expectPath, "security-credentials-extended") {
		expectPath += "/" // we've always had the trailing / on these
	}
	if expectPath != in.Path {
		return nil, fmt.Errorf("unexpected path in call %d: expect %s, got %s", callIndex, expectPath, in.Path)
	}

	if next.Response.Status != 200 {
		return nil, mockResponseError(next.Response.Status)
	}

	switch v := next.Response.Body.(type) {
	case string:
		return &imds.GetMetadataOutput{
			Content: io.NopCloser(strings.NewReader(v)),
		}, nil
	case map[string]any:
		j, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("unhandled response marshal failure in sep test case: %v", v)
		}
		return &imds.GetMetadataOutput{
			Content: io.NopCloser(bytes.NewReader(j)),
		}, nil
	default:
		return nil, fmt.Errorf("unhandled body type in sep test case: %T", next.Response.Body)
	}
}

func mockResponseError(status int) error {
	return &smithyhttp.ResponseError{
		Response: &smithyhttp.Response{
			Response: &http.Response{
				StatusCode: status,
				Body:       http.NoBody,
			},
		},
	}
}

const sepTestCaseJSON = `[
  {
    "summary": "Test IMDS credentials provider with env vars { AWS_EC2_METADATA_DISABLED=true } returns no credentials",
    "config": {
      "ec2InstanceProfileName": null,
      "envVars": {
        "AWS_EC2_METADATA_DISABLED": "true"
      }
    },
    "expectations": [],
    "outcomes": [
      {
        "result": "no credentials"
      }
    ]
  },
  {
    "summary": "Test IMDS credentials provider returns valid credentials with account ID",
    "config": {
      "ec2InstanceProfileName": null
    },
    "expectations": [
      {
        "get": "/latest/meta-data/iam/security-credentials-extended",
        "response": {
          "status": 200,
          "body": "my-profile-0001"
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0001",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-12T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-12T21:53:17.832308Z",
            "UnexpectedElement1": {
              "Name": "ignore-me-1"
            },
            "AccountId": "123456789101"
          }
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0001",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-12T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-12T21:53:17.832308Z",
            "UnexpectedElement1": {
              "Name": "ignore-me-1"
            },
            "AccountId": "123456789101"
          }
        }
      }
    ],
    "outcomes": [
      {
        "result": "credentials",
        "accountId": "123456789101"
      },
      {
        "result": "credentials",
        "accountId": "123456789101"
      }
    ]
  },
  {
    "summary": "Test IMDS credentials provider with a given profile name returns valid credentials with account ID",
    "config": {
      "ec2InstanceProfileName": "my-profile-0002"
    },
    "expectations": [
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0002",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-13T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-13T21:53:17.832308Z",
            "UnexpectedElement2": {
              "Name": "ignore-me-2"
            },
            "AccountId": "234567891011"
          }
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0002",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-13T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-13T21:53:17.832308Z",
            "UnexpectedElement2": {
              "Name": "ignore-me-2"
            },
            "AccountId": "234567891011"
          }
        }
      }
    ],
    "outcomes": [
      {
        "result": "credentials",
        "accountId": "234567891011"
      },
      {
        "result": "credentials",
        "accountId": "234567891011"
      }
    ]
  },
  {
    "summary": "Test IMDS credentials provider when profile is unstable returns valid credentials with account ID",
    "config": {
      "ec2InstanceProfileName": null
    },
    "expectations": [
      {
        "get": "/latest/meta-data/iam/security-credentials-extended",
        "response": {
          "status": 200,
          "body": "my-profile-0003"
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0003",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-14T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-14T21:53:17.832308Z",
            "UnexpectedElement3": {
              "Name": "ignore-me-3"
            },
            "AccountId": "345678910112"
          }
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0003",
        "response": {
          "status": 404
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials-extended",
        "response": {
          "status": 200,
          "body": "my-profile-0003-b"
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0003-b",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-14T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-14T21:53:17.832308Z",
            "UnexpectedElement3": {
              "Name": "ignore-me-3"
            },
            "AccountId": "314253647589"
          }
        }
      }
    ],
    "outcomes": [
      {
        "result": "credentials",
        "accountId": "345678910112"
      },
      {
        "result": "credentials",
        "accountId": "314253647589"
      }
    ]
  },
  {
    "summary": "Test IMDS credentials provider with a given profile name when profile is invalid throws an error",
    "config": {
      "ec2InstanceProfileName": "my-profile-0004"
    },
    "expectations": [
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0004",
        "response": {
          "status": 404
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials/my-profile-0004",
        "response": {
          "status": 404
        }
      }
    ],
    "outcomes": [
      {
        "result": "invalid profile"
      }
    ]
  },
  {
    "summary": "Test IMDS credentials provider when account ID is unavailable returns valid credentials",
    "config": {
      "ec2InstanceProfileName": null
    },
    "expectations": [
      {
        "get": "/latest/meta-data/iam/security-credentials-extended",
        "response": {
          "status": 200,
          "body": "my-profile-0005"
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0005",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-16T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-16T21:53:17.832308Z",
            "UnexpectedElement5": {
              "Name": "ignore-me-5"
            }
          }
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0005",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-16T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-16T21:53:17.832308Z",
            "UnexpectedElement5": {
              "Name": "ignore-me-5"
            }
          }
        }
      }
    ],
    "outcomes": [
      {
        "result": "credentials"
      },
      {
        "result": "credentials"
      }
    ]
  },
  {
    "summary": "Test IMDS credentials provider with a given profile name when account ID is unavailable returns valid credentials",
    "config": {
      "ec2InstanceProfileName": "my-profile-0006"
    },
    "expectations": [
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0006",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-17T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-17T21:53:17.832308Z",
            "UnexpectedElement6": {
              "Name": "ignore-me-6"
            }
          }
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0006",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-17T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-17T21:53:17.832308Z",
            "UnexpectedElement6": {
              "Name": "ignore-me-6"
            }
          }
        }
      }
    ],
    "outcomes": [
      {
        "result": "credentials"
      },
      {
        "result": "credentials"
      }
    ]
  },
  {
    "summary": "Test IMDS credentials provider when account ID is unavailable when profile is unstable returns valid credentials",
    "config": {
      "ec2InstanceProfileName": null
    },
    "expectations": [
      {
        "get": "/latest/meta-data/iam/security-credentials-extended",
        "response": {
          "status": 200,
          "body": "my-profile-0007"
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0007",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-18T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-18T21:53:17.832308Z",
            "UnexpectedElement7": {
              "Name": "ignore-me-7"
            }
          }
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0007",
        "response": {
          "status": 404
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials-extended",
        "response": {
          "status": 200,
          "body": "my-profile-0007-b"
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0007-b",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-18T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-18T21:53:17.832308Z",
            "UnexpectedElement7": {
              "Name": "ignore-me-7"
            }
          }
        }
      }
    ],
    "outcomes": [
      {
        "result": "credentials"
      },
      {
        "result": "credentials"
      }
    ]
  },
  {
    "summary": "Test IMDS credentials provider with a given profile name when account ID is unavailable when profile is invalid throws an error",
    "config": {
      "ec2InstanceProfileName": "my-profile-0008"
    },
    "expectations": [
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0008",
        "response": {
          "status": 404
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials/my-profile-0008",
        "response": {
          "status": 404
        }
      }
    ],
    "outcomes": [
      {
        "result": "invalid profile"
      }
    ]
  },
  {
    "summary": "Test IMDS credentials provider against legacy API returns valid credentials",
    "config": {
      "ec2InstanceProfileName": null
    },
    "expectations": [
      {
        "get": "/latest/meta-data/iam/security-credentials-extended",
        "response": {
          "status": 404
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials",
        "response": {
          "status": 200,
          "body": "my-profile-0009"
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials/my-profile-0009",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-20T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-20T21:53:17.832308Z"
          }
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials/my-profile-0009",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-20T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-20T21:53:17.832308Z"
          }
        }
      }
    ],
    "outcomes": [
      {
        "result": "credentials"
      },
      {
        "result": "credentials"
      }
    ]
  },
  {
    "summary": "Test IMDS credentials provider with a given profile name against legacy API returns valid credentials",
    "config": {
      "ec2InstanceProfileName": "my-profile-0010"
    },
    "expectations": [
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0010",
        "response": {
          "status": 404
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials/my-profile-0010",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-21T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-21T21:53:17.832308Z"
          }
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials/my-profile-0010",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-21T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-21T21:53:17.832308Z"
          }
        }
      }
    ],
    "outcomes": [
      {
        "result": "credentials"
      },
      {
        "result": "credentials"
      }
    ]
  },
  {
    "summary": "Test IMDS credentials provider against legacy API when profile is unstable returns valid credentials",
    "config": {
      "ec2InstanceProfileName": null
    },
    "expectations": [
      {
        "get": "/latest/meta-data/iam/security-credentials-extended",
        "response": {
          "status": 404
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials",
        "response": {
          "status": 200,
          "body": "my-profile-0011"
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials/my-profile-0011",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-22T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-22T21:53:17.832308Z"
          }
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials/my-profile-0011",
        "response": {
          "status": 404
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials",
        "response": {
          "status": 200,
          "body": "my-profile-0011-b"
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials/my-profile-0011-b",
        "response": {
          "status": 200,
          "body": {
            "Code": "Success",
            "LastUpdated": "2025-03-22T20:53:17.832308Z",
            "Type": "AWS-HMAC",
            "AccessKeyId": "ASIAIOSFODNN7EXAMPLE",
            "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
            "Token": "AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT+FvwqnKw...(truncated)",
            "Expiration": "2025-03-22T21:53:17.832308Z"
          }
        }
      }
    ],
    "outcomes": [
      {
        "result": "credentials"
      },
      {
        "result": "credentials"
      }
    ]
  },
  {
    "summary": "Test IMDS credentials provider with a given profile name against legacy API when profile is invalid throws an error",
    "config": {
      "ec2InstanceProfileName": "my-profile-0012"
    },
    "expectations": [
      {
        "get": "/latest/meta-data/iam/security-credentials-extended/my-profile-0012",
        "response": {
          "status": 404
        }
      },
      {
        "get": "/latest/meta-data/iam/security-credentials/my-profile-0012",
        "response": {
          "status": 404
        }
      }
    ],
    "outcomes": [
      {
        "result": "invalid profile"
      }
    ]
  }
]`

var skipSEPTestCases = map[string]string{
	"Test IMDS credentials provider with env vars { AWS_EC2_METADATA_DISABLED=true } returns no credentials": "environment variables are not considered in these unit tests",
}

func TestProvider_SEPTestCases(t *testing.T) {
	var testCases []sepTestCase
	if err := json.Unmarshal([]byte(sepTestCaseJSON), &testCases); err != nil {
		t.Fatal(err)
	}

	for _, tt := range testCases {
		t.Run(tt.Summary, func(t *testing.T) {
			if reason, ok := skipSEPTestCases[tt.Summary]; ok {
				t.Skip(reason)
			}

			mockIMDS := &sepMockIMDS{testCase: tt}
			provider := New(func(o *Options) {
				o.ProfileName = tt.Config.EC2InstanceProfileName
				o.Client = mockIMDS
			})

			for _, expect := range tt.Outcomes {
				creds, err := provider.Retrieve(context.Background())
				switch expect.Result {
				case "credentials":
					if creds.AccessKeyID == "" {
						t.Errorf("expected credentials, got none: %v", err)
					}
					if expect.AccountID != creds.AccountID {
						t.Errorf("expected account id %q, got %q", expect.AccountID, creds.AccountID)
					}
				case "no credentials":
					fallthrough
				case "invalid profile":
					if err == nil {
						t.Error("expected error, got none")
					}
				}
			}
		})
	}
}
