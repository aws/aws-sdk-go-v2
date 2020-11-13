---
title: "Enforcing TLS Version 1.2 in AWS SDK for Go V2"
linkTitle: "Enforcing TLS 1.2"
date: "2020-11-09"
---

To add increased security when communicating with AWS services, you should configure your client to use TLS 1.2 or
later.

## How do I set my TLS version?

You can set the TLS version to 1.2 using the following code.

1. Create a custom HTTP transport to require a minimum version of TLS 1.2
   ```go
   tr := &http.Transport{
       TLSClientConfig: &tls.Config{
           MinVersion: tls.VersionTLS12,
       },
   }
   ```
1. Configure the transport.
   ```go
   // In Go versions earlier than 1.13
   err := http2.ConfigureTransport(tr)
   if err != nil {
      fmt.Println("Got an error configuring HTTP transport")
      fmt.Println(err)
      return
   }

    // In Go versions later than 1.13
    tr.ForceAttemptHTTP2 = true
   ```
1. Create an HTTP client with the configured transport, and use that to create a session. REGION is the AWS Region, such
   as `us-west-2`.
   ```go
   client := http.Client{Transport: tr}

   sess := session.Must(session.NewSession(&aws.Config{
   Region:     &REGION,
   HTTPClient: &client,
   }))
   ```

1. Use the following function to confirm your TLS version.
   ```go
   func GetTLSVersion(tr *http.Transport) string {
        switch tr.TLSClientConfig.MinVersion {
        case tls.VersionTLS10:
            return "TLS 1.0"
        case tls.VersionTLS11:
            return "TLS 1.1"
        case tls.VersionTLS12:
            return "TLS 1.2"
        case tls.VersionTLS13:
            return "TLS 1.3"
        }

        return "Unknown"
    }
   ```

5. Confirm your TLS version by calling `GetTLSVersion`.
   ```go
   if tr, ok := s3Client.Config.HTTPClient.Transport.(*http.Transport); ok {
      log.Printf("Client uses %v", GetTLSVersion(tr))
   }
   ```
       
