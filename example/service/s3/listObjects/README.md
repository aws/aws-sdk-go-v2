# Amazon S3 ListObjectsV2 Example

This is an example using the AWS SDK for Go to list objects in a S3 bucket.

### Usage

The example uses the bucket name provided, and lists all object keys in a bucket.
Optionally taking a prefix to filter object with that prefix, and separator.

```
go run listObjects.go -bucket <bucket-name> [-prefix <string>] [-delimiter <string>] [-max-keys <int>]

  -bucket name
        The name of the S3 bucket to list objects from.
  -delimiter object key delimiter
        The optional object key delimiter used by S3 List objects to group object keys.
  -max-keys keys per page
        The maximum number of keys per page to retrieve at once.
  -prefix object prefix
        The optional object prefix of the S3 Object keys to list.
```

### Output:

```
Objects:
Object: myKey
Object: mykey.txt
Object: resources/0001/item-01
Object: resources/0001/item-02
Object: resources/0001/item-03
Object: resources/0002/item-01
Object: resources/0002/item-02
Object: resources/0002/item-03
Object: resources/0002/item-04
Object: resources/0002/item-05
```
