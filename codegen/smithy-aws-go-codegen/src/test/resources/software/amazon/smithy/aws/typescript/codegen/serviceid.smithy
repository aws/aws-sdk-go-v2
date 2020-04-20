namespace smithy.example

@aws.api#service(sdkId: "Not Same")
@protocols([{name: "aws.rest-json-1.1"}])
service OriginalName {
    version: "2019-10-15",
}
