// Package customizations provides customizations for the Amazon Route53 API client.
//
// This package provides support for following customizations
//  Process Response Middleware: used for custom error deserializing
//
//
// Process Response Middleware
//
// Route53 operation "ChangeResourceRecordSets" can have an error response returned in
// a slightly different format. This customization is only applicable to
// ChangeResourceRecordSets operation of Route53.
//
// Here's a sample error response:
//     <?xml version="1.0" encoding="UTF-8"?>
//         <InvalidChangeBatch xmlns="https://route53.amazonaws.com/doc/2013-04-01/">
//         <Messages>
//             <Message>Tried to create resource record set duplicate.example.com. type A, but it already exists</Message>
//         </Messages>
//     </InvalidChangeBatch>
//
//
// The processResponse middleware customizations enables SDK to check for an error
// response starting with `InvalidChangeBatch` tag prior to deserialization.
//
// As this check in error response needs to be performed earlier than response
// deserialization. Since the behavior of Deserialization is in
// reverse order to the other stack steps its easier to consider that "after" means
// "before".
//
//  Middleware layering:
//
// 	HTTP Response -> process response error -> deserialize
//
//
// In case the returned error response has `InvalidChangeBatch` format, the error is
// deserialized and returned. The operation deserializer does not attempt to deserialize
// as an error is returned by the process response error middleware.
//
package customizations
