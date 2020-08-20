// +build disabled

package external

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/aws/stscreds"
)

func ExampleWithMFATokenFunc() {
	cfg, err := external.LoadDefaultAWSConfig(
		// Set the provider function for the MFA token.
		external.WithMFATokenFunc(stscreds.StdinTokenProvider),

		// Optionally, specify the shared configuration profile to load.
		external.WithSharedConfigProfile("exampleProfile"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config, %v", err)
		os.Exit(1)
	}

	// If assume role credentials with MFA enabled are specified in the shared
	// 	configuration the MFA token provider function will be called to retrieve
	// the MFA token for the assume role API call.
	fmt.Println(cfg.Credentials.Retrieve(context.Background()))
}
