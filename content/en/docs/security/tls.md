---
title: "Enforcing TLS Version 1.2 in AWS SDK for Go V2"
linkTitle: "Enforcing TLS 1.2"
date: "2020-11-09"
---

The {{% alias sdk-go %}} configures the default HTTP client used by the SDK service clients to require a minimum version
of TLS 1.2 or greater. When using the [http.BuildableClient]({{< apiref "aws/transport/http#BuildableClient" >}}) to
customize the SDK HTTP client, the minimum TLS value is configured as TLS 1.2.

{{% pageinfo color="warning" %}}
If your application constructs an HTTP client using a method other than the provided `BuildableClient`, you must
configure your client to set the minimum TLS version to 1.2.
{{% /pageinfo %}}
