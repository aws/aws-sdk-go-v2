/*
Package aws provides the core SDK's utilities and shared types. Use this package's
utilities to simplify setting and reading API operations parameters.

Value and Pointer Conversion Utilities

This package includes a helper conversion utility for each scalar type the SDK's
API use. These utilities make getting a pointer of the scalar, and dereferencing
a pointer easier.

Each conversion utility comes in two forms. Value to Pointer and Pointer to Value.
The Pointer to value will safely dereference the pointer and return its value.
If the pointer was nil, the scalar's zero value will be returned.

The value to pointer functions will be named after the scalar type. So get a
*string from a string value use the "String" function. This makes it easy to
to get pointer of a literal string value, because getting the address of a
literal requires assigning the value to a variable first.

   var strPtr *string

   // Without the SDK's conversion functions
   str := "my string"
   strPtr = &str

   // With the SDK's conversion functions
   strPtr = aws.String("my string")

   // Convert *string to string value
   str = aws.StringValue(strPtr)

In addition to scalars the aws package also includes conversion utilities for
map and slice for commonly types used in API parameters. The map and slice
conversion functions use similar naming pattern as the scalar conversion
functions.

   var strPtrs []*string
   var strs []string = []string{"Go", "Gophers", "Go"}

   // Convert []string to []*string
   strPtrs = aws.StringSlice(strs)

   // Convert []*string to []string
   strs = aws.StringValueSlice(strPtrs)

SDK Default HTTP Client

The SDK will use the http.DefaultClient if a HTTP client is not provided to
the SDK's Session, or service client constructor. This means that if the
http.DefaultClient is modified by other components of your application the
modifications will be picked up by the SDK as well.

In some cases this might be intended, but it is a better practice to create
a custom HTTP Client to share explicitly through your application. You can
configure the SDK to use the custom HTTP Client by setting the HTTPClient
value of the SDK's Config type when creating a Session or service client.

SDK Credentials

The CredentialsLoader is the primary method of getting access to and managing
credentials Values. Using dependency injection retrieval of the credential
values is handled by a object which satisfies the CredentialsProvider interface.

By default the CredentialsLoader.Get() will cache the successful result of a
CredentialsProvider's Retrieve() until CredentialsProvider.IsExpired() returns true. At which
point CredentialsLoader will call CredentialsProvider's Retrieve() to get new credential Credentials.

The CredentialsProvider is responsible for determining when credentials Credentials have expired.
It is also important to note that CredentialsLoader will always call Retrieve the
first time CredentialsLoader.Get() is called.

Example of using the environment variable credentials.

    creds := aws.NewEnvCredentials()

    // Retrieve the credentials value
    credValue, err := creds.Get()
    if err != nil {
        // handle error
    }

Example of forcing credentials to expire and be refreshed on the next Get().
This may be helpful to proactively expire credentials and refresh them sooner
than they would naturally expire on their own.

    creds := aws.NewCredentials(&ec2rolecreds.EC2RoleProvider{})
    creds.Expire()
    credsValue, err := creds.Get()
    // New credentials will be retrieved instead of from cache.


Custom CredentialsProvider

Each CredentialsProvider built into this package also provides a helper method to generate
a CredentialsLoader pointer setup with the CredentialsProvider. To use a custom CredentialsProvider just
create a type which satisfies the CredentialsProvider interface and pass it to the
NewCredentials method.

    type MyProvider struct{}
    func (m *MyProvider) Retrieve() (Credentials, error) {...}
    func (m *MyProvider) IsExpired() bool {...}

    creds := aws.NewCredentials(&MyProvider{})
    credValue, err := creds.Get()
*/
package aws
