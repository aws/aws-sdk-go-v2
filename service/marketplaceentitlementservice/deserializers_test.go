package marketplaceentitlementservice

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/marketplaceentitlementservice/types"
	"github.com/stretchr/testify/require"
)

func Test_awsAwsjson11_deserializeOpDocumentGetEntitlementsOutput(t *testing.T) {
	const responseBody = `{
	  "Entitlements": [
		{
		  "CustomerIdentifier": "***",
		  "Dimension": "numNodes",
		  "ExpirationDate": 1.684888527303E9,
		  "ProductCode": "***",
		  "Value": {
			"EntitlementValueType": "int",
			"IntegerValue": 2
		  }
		}
	  ]
	}`

	decoder := json.NewDecoder(strings.NewReader(responseBody))
	decoder.UseNumber()

	var shape interface{}
	require.NoError(t, decoder.Decode(&shape))

	var output *GetEntitlementsOutput
	require.NoError(t, awsAwsjson11_deserializeOpDocumentGetEntitlementsOutput(&output, shape))

	require.Equal(t, &types.EntitlementValueMemberIntegerValue{Value: 2}, output.Entitlements[0].Value)
}
