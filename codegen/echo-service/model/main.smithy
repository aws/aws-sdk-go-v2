$version: "2.0"

namespace smithy.sdk.example

use aws.api#service
use aws.auth#sigv4
use smithy.protocols#rpcv2Cbor
 
@title("Echo Service")
@service(sdkId: "Echo")
@sigv4(name: "echoservice")
@rpcv2Cbor
// comment cbor and uncomment this for enabling awsJson
// @aws.protocols#awsJson1_0
service EchoService {
    version: "2020-07-02"
    operations: [EchoOperation]
}

operation EchoOperation {
    input := with [ComplexStructureMixin] {}
    output := with [ComplexStructureMixin] {}
}

structure ComplexStructure with [ComplexStructureMixin] {}

@mixin()
structure ComplexStructureMixin {
    booleanMember: Boolean
    stringMember: String
    integerMember: Integer
    longMember: Long
    floatMember: Float
    doubleMember: Double
    timestampMember: Timestamp
    blobMember: Blob
    listOfStringsMember: ListOfStrings
    listOfComplexObjectMember: ListOfComplexStructure
    mapOfStringToStringMember: MapOfStringToString
    complexStructMember: ComplexStructure
}

map MapOfStringToString {
    key: String
    value: String
}

list ListOfStrings {
    member: String
}

list ListOfComplexStructure {
    member: ComplexStructure
}
