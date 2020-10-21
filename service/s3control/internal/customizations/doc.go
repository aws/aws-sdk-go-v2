/*
Package customizations provides customizations for the Amazon S3-Control API client.

This package provides support for following S3-Control customizations

    UpdateEndpoint Middleware: resolves a custom endpoint as per s3-control config options


Dualstack support

By default dualstack support for s3-control client is disabled. By enabling `UseDualstack`
option on s3-control client, you can enable dualstack endpoint support.


UpdateEndpoint middleware handler for modifying resolved endpoint needs to be
executed after request serialization.

 Middleware layering:

 HTTP Request -> operation serializer -> Update-Endpoint customization -> next middleware

Customization option:
 UseDualstack (Disabled by Default)

*/
package customizations
