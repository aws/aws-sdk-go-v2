---
title: "Amazon CloudFront Utilities"
linkTitle: "Amazon CloudFront"
date: "2020-11-12"
description: "Using the AWS SDK for Go V2 Amazon CloudFront Utilities"
---

## {{% alias service="CFlong" %}} URL Signer

The {{% alias service="CFlong" %}} URL signer simplifies the process of creating
signed URLs. A signed URL includes information, such as an expiration
date and time, that enables you to control access to your content.
Signed URLs are useful when you want to distribute content through the
internet, but want to restrict access to certain users (for example, to
users who have paid a fee).

To sign a URL, create a `URLSigner` instance with your {{% alias service="CF" %}} key pair ID and the associated private
key. Then call the
`Sign` or `SignWithPolicy` method and include the URL to sign. For more information about {{% alias service="CFlong" %}}
key pairs,
see [Creating CloudFront Key Pairs for Your Trusted Signers](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/private-content-trusted-signers.html#private-content-creating-cloudfront-key-pairs)
in the {{% alias service="CF" %}} Developer Guide.

The following example creates a signed URL that's valid for one hour
after it is created.

```go
import "github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"

// ...

signer := sign.NewURLSigner(keyID, privKey)

signedURL, err := signer.Sign(rawURL, time.Now().Add(1*time.Hour))
if err != nil {
    log.Fatalf("Failed to sign url, err: %s\n", err.Error())
    return
}
```

For more information about the signing utility, see the [sign]({{% apiref "feature/cloudfront/sign" %}}) package in the
{{% alias sdk-api %}}.

