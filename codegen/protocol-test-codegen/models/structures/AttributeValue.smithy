$version: "2"

namespace com.amazonaws.sdk.benchmark

@documentation("""
    The famous recursive structure from Amazon DynamoDB.
""")
union AttributeValue {
    S: String
    N: String
    B: Blob
    SS: StringSet
    NS: NumberSet
    BS: BinarySet
    M: AttributeValueMap
    L: AttributeValueList
    NULL: Boolean
    BOOL: Boolean
}

list StringSet {
    member: String
}

list NumberSet {
    member: String
}

list BinarySet {
    member: Blob
}

map AttributeValueMap {
    key: String
    value: AttributeValue
}

list AttributeValueList {
    member: AttributeValue
}