$version: "2"

namespace com.amazonaws.sdk.benchmark

use aws.protocols#awsJson1_0
use smithy.protocols#rpcv2Cbor
use smithy.test#httpRequestTests

@documentation("""
    From Amazon DynamoDB.
    Serialization of recursive structures.
""")
@http(method: "POST", uri: "/PutItem", code: 200)
@httpRequestTests([
    // section: json
    {
        id: "awsJson1_0_PutItemRequest_Baseline"
        protocol: awsJson1_0
        documentation: """
        This test gives baseline of serializing a minimal
        amount of data for a data-plane write.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: {
                    S: "test-id"
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_PutItemRequest_ShallowMap_S"
        protocol: awsJson1_0
        documentation: """
        Serializing a map (small) with many keys but minimal nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                attr1: { S: "value1" }
                attr2: { S: "value2" }
                attr3: { S: "value3" }
                attr4: { S: "value4" }
                attr5: { S: "value5" }
                attr6: { S: "value6" }
                attr7: { S: "value7" }
                attr8: { S: "value8" }
                attr9: { S: "value9" }
                attr10: { S: "value10" }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_PutItemRequest_ShallowMap_M"
        protocol: awsJson1_0
        documentation: """
        Serializing a map (medium) with many keys but minimal nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                attr1: { S: "value1" }
                attr2: { S: "value2" }
                attr3: { S: "value3" }
                attr4: { S: "value4" }
                attr5: { S: "value5" }
                attr6: { S: "value6" }
                attr7: { S: "value7" }
                attr8: { S: "value8" }
                attr9: { S: "value9" }
                attr10: { S: "value10" }
                attr11: { S: "value11" }
                attr12: { S: "value12" }
                attr13: { S: "value13" }
                attr14: { S: "value14" }
                attr15: { S: "value15" }
                attr16: { S: "value16" }
                attr17: { S: "value17" }
                attr18: { S: "value18" }
                attr19: { S: "value19" }
                attr20: { S: "value20" }
                attr21: { S: "value21" }
                attr22: { S: "value22" }
                attr23: { S: "value23" }
                attr24: { S: "value24" }
                attr25: { S: "value25" }
                attr26: { S: "value26" }
                attr27: { S: "value27" }
                attr28: { S: "value28" }
                attr29: { S: "value29" }
                attr30: { S: "value30" }
                attr31: { S: "value31" }
                attr32: { S: "value32" }
                attr33: { S: "value33" }
                attr34: { S: "value34" }
                attr35: { S: "value35" }
                attr36: { S: "value36" }
                attr37: { S: "value37" }
                attr38: { S: "value38" }
                attr39: { S: "value39" }
                attr40: { S: "value40" }
                attr41: { S: "value41" }
                attr42: { S: "value42" }
                attr43: { S: "value43" }
                attr44: { S: "value44" }
                attr45: { S: "value45" }
                attr46: { S: "value46" }
                attr47: { S: "value47" }
                attr48: { S: "value48" }
                attr49: { S: "value49" }
                attr50: { S: "value50" }
                attr51: { S: "value51" }
                attr52: { S: "value52" }
                attr53: { S: "value53" }
                attr54: { S: "value54" }
                attr55: { S: "value55" }
                attr56: { S: "value56" }
                attr57: { S: "value57" }
                attr58: { S: "value58" }
                attr59: { S: "value59" }
                attr60: { S: "value60" }
                attr61: { S: "value61" }
                attr62: { S: "value62" }
                attr63: { S: "value63" }
                attr64: { S: "value64" }
                attr65: { S: "value65" }
                attr66: { S: "value66" }
                attr67: { S: "value67" }
                attr68: { S: "value68" }
                attr69: { S: "value69" }
                attr70: { S: "value70" }
                attr71: { S: "value71" }
                attr72: { S: "value72" }
                attr73: { S: "value73" }
                attr74: { S: "value74" }
                attr75: { S: "value75" }
                attr76: { S: "value76" }
                attr77: { S: "value77" }
                attr78: { S: "value78" }
                attr79: { S: "value79" }
                attr80: { S: "value80" }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_PutItemRequest_ShallowMap_L"
        protocol: awsJson1_0
        documentation: """
        Serializing a map (large) with many keys but minimal nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                attr1: { S: "value1" }
                attr2: { S: "value2" }
                attr3: { S: "value3" }
                attr4: { S: "value4" }
                attr5: { S: "value5" }
                attr6: { S: "value6" }
                attr7: { S: "value7" }
                attr8: { S: "value8" }
                attr9: { S: "value9" }
                attr10: { S: "value10" }
                attr11: { S: "value11" }
                attr12: { S: "value12" }
                attr13: { S: "value13" }
                attr14: { S: "value14" }
                attr15: { S: "value15" }
                attr16: { S: "value16" }
                attr17: { S: "value17" }
                attr18: { S: "value18" }
                attr19: { S: "value19" }
                attr20: { S: "value20" }
                attr21: { S: "value21" }
                attr22: { S: "value22" }
                attr23: { S: "value23" }
                attr24: { S: "value24" }
                attr25: { S: "value25" }
                attr26: { S: "value26" }
                attr27: { S: "value27" }
                attr28: { S: "value28" }
                attr29: { S: "value29" }
                attr30: { S: "value30" }
                attr31: { S: "value31" }
                attr32: { S: "value32" }
                attr33: { S: "value33" }
                attr34: { S: "value34" }
                attr35: { S: "value35" }
                attr36: { S: "value36" }
                attr37: { S: "value37" }
                attr38: { S: "value38" }
                attr39: { S: "value39" }
                attr40: { S: "value40" }
                attr41: { S: "value41" }
                attr42: { S: "value42" }
                attr43: { S: "value43" }
                attr44: { S: "value44" }
                attr45: { S: "value45" }
                attr46: { S: "value46" }
                attr47: { S: "value47" }
                attr48: { S: "value48" }
                attr49: { S: "value49" }
                attr50: { S: "value50" }
                attr51: { S: "value51" }
                attr52: { S: "value52" }
                attr53: { S: "value53" }
                attr54: { S: "value54" }
                attr55: { S: "value55" }
                attr56: { S: "value56" }
                attr57: { S: "value57" }
                attr58: { S: "value58" }
                attr59: { S: "value59" }
                attr60: { S: "value60" }
                attr61: { S: "value61" }
                attr62: { S: "value62" }
                attr63: { S: "value63" }
                attr64: { S: "value64" }
                attr65: { S: "value65" }
                attr66: { S: "value66" }
                attr67: { S: "value67" }
                attr68: { S: "value68" }
                attr69: { S: "value69" }
                attr70: { S: "value70" }
                attr71: { S: "value71" }
                attr72: { S: "value72" }
                attr73: { S: "value73" }
                attr74: { S: "value74" }
                attr75: { S: "value75" }
                attr76: { S: "value76" }
                attr77: { S: "value77" }
                attr78: { S: "value78" }
                attr79: { S: "value79" }
                attr80: { S: "value80" }
                attr81: { S: "value81" }
                attr82: { S: "value82" }
                attr83: { S: "value83" }
                attr84: { S: "value84" }
                attr85: { S: "value85" }
                attr86: { S: "value86" }
                attr87: { S: "value87" }
                attr88: { S: "value88" }
                attr89: { S: "value89" }
                attr90: { S: "value90" }
                attr91: { S: "value91" }
                attr92: { S: "value92" }
                attr93: { S: "value93" }
                attr94: { S: "value94" }
                attr95: { S: "value95" }
                attr96: { S: "value96" }
                attr97: { S: "value97" }
                attr98: { S: "value98" }
                attr99: { S: "value99" }
                attr100: { S: "value100" }
                attr101: { S: "value101" }
                attr102: { S: "value102" }
                attr103: { S: "value103" }
                attr104: { S: "value104" }
                attr105: { S: "value105" }
                attr106: { S: "value106" }
                attr107: { S: "value107" }
                attr108: { S: "value108" }
                attr109: { S: "value109" }
                attr110: { S: "value110" }
                attr111: { S: "value111" }
                attr112: { S: "value112" }
                attr113: { S: "value113" }
                attr114: { S: "value114" }
                attr115: { S: "value115" }
                attr116: { S: "value116" }
                attr117: { S: "value117" }
                attr118: { S: "value118" }
                attr119: { S: "value119" }
                attr120: { S: "value120" }
                attr121: { S: "value121" }
                attr122: { S: "value122" }
                attr123: { S: "value123" }
                attr124: { S: "value124" }
                attr125: { S: "value125" }
                attr126: { S: "value126" }
                attr127: { S: "value127" }
                attr128: { S: "value128" }
                attr129: { S: "value129" }
                attr130: { S: "value130" }
                attr131: { S: "value131" }
                attr132: { S: "value132" }
                attr133: { S: "value133" }
                attr134: { S: "value134" }
                attr135: { S: "value135" }
                attr136: { S: "value136" }
                attr137: { S: "value137" }
                attr138: { S: "value138" }
                attr139: { S: "value139" }
                attr140: { S: "value140" }
                attr141: { S: "value141" }
                attr142: { S: "value142" }
                attr143: { S: "value143" }
                attr144: { S: "value144" }
                attr145: { S: "value145" }
                attr146: { S: "value146" }
                attr147: { S: "value147" }
                attr148: { S: "value148" }
                attr149: { S: "value149" }
                attr150: { S: "value150" }
                attr151: { S: "value151" }
                attr152: { S: "value152" }
                attr153: { S: "value153" }
                attr154: { S: "value154" }
                attr155: { S: "value155" }
                attr156: { S: "value156" }
                attr157: { S: "value157" }
                attr158: { S: "value158" }
                attr159: { S: "value159" }
                attr160: { S: "value160" }
                attr161: { S: "value161" }
                attr162: { S: "value162" }
                attr163: { S: "value163" }
                attr164: { S: "value164" }
                attr165: { S: "value165" }
                attr166: { S: "value166" }
                attr167: { S: "value167" }
                attr168: { S: "value168" }
                attr169: { S: "value169" }
                attr170: { S: "value170" }
                attr171: { S: "value171" }
                attr172: { S: "value172" }
                attr173: { S: "value173" }
                attr174: { S: "value174" }
                attr175: { S: "value175" }
                attr176: { S: "value176" }
                attr177: { S: "value177" }
                attr178: { S: "value178" }
                attr179: { S: "value179" }
                attr180: { S: "value180" }
                attr181: { S: "value181" }
                attr182: { S: "value182" }
                attr183: { S: "value183" }
                attr184: { S: "value184" }
                attr185: { S: "value185" }
                attr186: { S: "value186" }
                attr187: { S: "value187" }
                attr188: { S: "value188" }
                attr189: { S: "value189" }
                attr190: { S: "value190" }
                attr191: { S: "value191" }
                attr192: { S: "value192" }
                attr193: { S: "value193" }
                attr194: { S: "value194" }
                attr195: { S: "value195" }
                attr196: { S: "value196" }
                attr197: { S: "value197" }
                attr198: { S: "value198" }
                attr199: { S: "value199" }
                attr200: { S: "value200" }
                attr201: { S: "value201" }
                attr202: { S: "value202" }
                attr203: { S: "value203" }
                attr204: { S: "value204" }
                attr205: { S: "value205" }
                attr206: { S: "value206" }
                attr207: { S: "value207" }
                attr208: { S: "value208" }
                attr209: { S: "value209" }
                attr210: { S: "value210" }
                attr211: { S: "value211" }
                attr212: { S: "value212" }
                attr213: { S: "value213" }
                attr214: { S: "value214" }
                attr215: { S: "value215" }
                attr216: { S: "value216" }
                attr217: { S: "value217" }
                attr218: { S: "value218" }
                attr219: { S: "value219" }
                attr220: { S: "value220" }
                attr221: { S: "value221" }
                attr222: { S: "value222" }
                attr223: { S: "value223" }
                attr224: { S: "value224" }
                attr225: { S: "value225" }
                attr226: { S: "value226" }
                attr227: { S: "value227" }
                attr228: { S: "value228" }
                attr229: { S: "value229" }
                attr230: { S: "value230" }
                attr231: { S: "value231" }
                attr232: { S: "value232" }
                attr233: { S: "value233" }
                attr234: { S: "value234" }
                attr235: { S: "value235" }
                attr236: { S: "value236" }
                attr237: { S: "value237" }
                attr238: { S: "value238" }
                attr239: { S: "value239" }
                attr240: { S: "value240" }
                attr241: { S: "value241" }
                attr242: { S: "value242" }
                attr243: { S: "value243" }
                attr244: { S: "value244" }
                attr245: { S: "value245" }
                attr246: { S: "value246" }
                attr247: { S: "value247" }
                attr248: { S: "value248" }
                attr249: { S: "value249" }
                attr250: { S: "value250" }
                attr251: { S: "value251" }
                attr252: { S: "value252" }
                attr253: { S: "value253" }
                attr254: { S: "value254" }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_PutItemRequest_Nested_M"
        protocol: awsJson1_0
        documentation: """
        A narrow item with moderate nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                nested: {
                    M: {
                        level1: {
                            M: {
                                level2: {
                                    M: {
                                        level3: {
                                            M: {
                                                level4: {
                                                    M: {
                                                        deepValue: { S: "deep-nested-value" }
                                                    }
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_PutItemRequest_Nested_L"
        protocol: awsJson1_0
        documentation: """
        A narrow item with deep nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                nested: {
                    M: {
                        level1: {
                            L: [{
                                M: {
                                    level3: {
                                        L: [{
                                            M: {
                                                level5: {
                                                    L: [{
                                                        M: {
                                                            level7: {
                                                                L: [{
                                                                    M: {
                                                                        level9: {
                                                                            L: [{
                                                                                M: {
                                                                                    level11: {
                                                                                        L: [{
                                                                                            M: {
                                                                                                level13: {
                                                                                                    L: [{
                                                                                                        M: {
                                                                                                            level15: {
                                                                                                                L: [{
                                                                                                                    M: {
                                                                                                                        level17: {
                                                                                                                            L: [{
                                                                                                                                M: {
                                                                                                                                    level19: {
                                                                                                                                        L: [{
                                                                                                                                            M: {
                                                                                                                                                level21: {
                                                                                                                                                    L: [{
                                                                                                                                                        M: {
                                                                                                                                                            level23: {
                                                                                                                                                                L: [{
                                                                                                                                                                    M: {
                                                                                                                                                                        level25: {
                                                                                                                                                                            L: [{
                                                                                                                                                                                M: {
                                                                                                                                                                                    level27: {
                                                                                                                                                                                        L: [{
                                                                                                                                                                                            M: {
                                                                                                                                                                                                deepValue: { S: "smithy parser limit reached" }
                                                                                                                                                                                            }
                                                                                                                                                                                        }]
                                                                                                                                                                                    }
                                                                                                                                                                                }
                                                                                                                                                                            }]
                                                                                                                                                                        }
                                                                                                                                                                    }
                                                                                                                                                                }]
                                                                                                                                                            }
                                                                                                                                                        }
                                                                                                                                                    }]
                                                                                                                                                }
                                                                                                                                            }
                                                                                                                                        }]
                                                                                                                                    }
                                                                                                                                }
                                                                                                                            }]
                                                                                                                        }
                                                                                                                    }
                                                                                                                }]
                                                                                                            }
                                                                                                        }
                                                                                                    }]
                                                                                                }
                                                                                            }
                                                                                        }]
                                                                                    }
                                                                                }
                                                                            }]
                                                                        }
                                                                    }
                                                                }]
                                                            }
                                                        }
                                                    }]
                                                }
                                            }
                                        }]
                                    }
                                }
                            }]
                        }
                    }
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_PutItemRequest_MixedItem_S"
        protocol: awsJson1_0
        documentation: """
        An item (small) that uses mixed AttributeValue types, including nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                stringSet: { SS: ["item1", "item2", "item3"] }
                numberSet: { NS: ["1", "2", "3"] }
                list: {
                    L: [
                        { S: "listItem1" }
                        { N: "42" }
                        { BOOL: true }
                    ]
                }
                mixedMap: {
                    M: {
                        stringAttr: { S: "value" }
                        numberAttr: { N: "123" }
                        boolAttr: { BOOL: false }
                    }
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_PutItemRequest_MixedItem_M"
        protocol: awsJson1_0
        documentation: """
        An item (medium) that uses mixed AttributeValue types, including nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                stringSet: { SS: ["item1", "item2", "item3", "item4", "item5", "item6", "item7", "item8", "item9", "item10", "item11", "item12", "item13", "item14", "item15", "item16", "item17", "item18", "item19", "item20", "item21", "item22", "item23", "item24", "item25", "item26", "item27", "item28", "item29", "item30"] }
                numberSet: { NS: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30"] }
                list: {
                    L: [
                        { S: "listItem1" }
                        { N: "42" }
                        { BOOL: true }
                        { S: "listItem4" }
                        { N: "100" }
                        { BOOL: false }
                        { S: "listItem7" }
                        { N: "200" }
                        { BOOL: true }
                        { S: "listItem10" }
                        { N: "300" }
                        { BOOL: false }
                        { S: "listItem13" }
                        { N: "400" }
                        { BOOL: true }
                        { S: "listItem16" }
                        { N: "500" }
                        { BOOL: false }
                        { S: "listItem19" }
                        { N: "600" }
                        { BOOL: true }
                        { S: "listItem22" }
                        { N: "700" }
                        { BOOL: false }
                        { S: "listItem25" }
                        { N: "800" }
                        { BOOL: true }
                        { S: "listItem28" }
                        { N: "900" }
                        { BOOL: false }
                    ]
                }
                mixedMap: {
                    M: {
                        stringAttr1: { S: "value1" }
                        numberAttr1: { N: "123" }
                        boolAttr1: { BOOL: false }
                        stringAttr2: { S: "value2" }
                        numberAttr2: { N: "456" }
                        boolAttr2: { BOOL: true }
                        stringAttr3: { S: "value3" }
                        numberAttr3: { N: "789" }
                        boolAttr3: { BOOL: false }
                        stringAttr4: { S: "value4" }
                        numberAttr4: { N: "101" }
                        boolAttr4: { BOOL: true }
                        stringAttr5: { S: "value5" }
                        numberAttr5: { N: "202" }
                        boolAttr5: { BOOL: false }
                        stringAttr6: { S: "value6" }
                        numberAttr6: { N: "303" }
                        boolAttr6: { BOOL: true }
                        stringAttr7: { S: "value7" }
                        numberAttr7: { N: "404" }
                        boolAttr7: { BOOL: false }
                        stringAttr8: { S: "value8" }
                        numberAttr8: { N: "505" }
                        boolAttr8: { BOOL: true }
                        stringAttr9: { S: "value9" }
                        numberAttr9: { N: "606" }
                        boolAttr9: { BOOL: false }
                        stringAttr10: { S: "value10" }
                        numberAttr10: { N: "707" }
                        boolAttr10: { BOOL: true }
                    }
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_PutItemRequest_MixedItem_L"
        protocol: awsJson1_0
        documentation: """
        An item (large) that uses mixed AttributeValue types, including nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                stringSet: { SS: ["item1", "item2", "item3", "item4", "item5", "item6", "item7", "item8", "item9", "item10", "item11", "item12", "item13", "item14", "item15", "item16", "item17", "item18", "item19", "item20", "item21", "item22", "item23", "item24", "item25", "item26", "item27", "item28", "item29", "item30"] }
                numberSet: { NS: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30"] }
                list: {
                    L: [
                        { S: "listItem1" }
                        { N: "42" }
                        { BOOL: true }
                        { S: "listItem4" }
                        { N: "100" }
                        { BOOL: false }
                        { S: "listItem7" }
                        { N: "200" }
                        { BOOL: true }
                        { S: "listItem10" }
                        { N: "300" }
                        { BOOL: false }
                        { S: "listItem13" }
                        { N: "400" }
                        { BOOL: true }
                        { S: "listItem16" }
                        { N: "500" }
                        { BOOL: false }
                        { S: "listItem19" }
                        { N: "600" }
                        { BOOL: true }
                        { S: "listItem22" }
                        { N: "700" }
                        { BOOL: false }
                        { S: "listItem25" }
                        { N: "800" }
                        { BOOL: true }
                        { S: "listItem28" }
                        { N: "900" }
                        { BOOL: false }
                        {
                            M: {
                                id: { S: "test-id" }
                                stringSet: { SS: ["item1", "item2", "item3", "item4", "item5", "item6", "item7", "item8", "item9", "item10", "item11", "item12", "item13", "item14", "item15", "item16", "item17", "item18", "item19", "item20", "item21", "item22", "item23", "item24", "item25", "item26", "item27", "item28", "item29", "item30"] }
                                numberSet: { NS: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30"] }
                                list: {
                                    L: [
                                        { S: "listItem1" }
                                        { N: "42" }
                                        { BOOL: true }
                                        { S: "listItem4" }
                                        { N: "100" }
                                        { BOOL: false }
                                        { S: "listItem7" }
                                        { N: "200" }
                                        { BOOL: true }
                                        { S: "listItem10" }
                                        { N: "300" }
                                        { BOOL: false }
                                        { S: "listItem13" }
                                        { N: "400" }
                                        { BOOL: true }
                                        { S: "listItem16" }
                                        { N: "500" }
                                        { BOOL: false }
                                        { S: "listItem19" }
                                        { N: "600" }
                                        { BOOL: true }
                                        { S: "listItem22" }
                                        { N: "700" }
                                        { BOOL: false }
                                        { S: "listItem25" }
                                        { N: "800" }
                                        { BOOL: true }
                                        { S: "listItem28" }
                                        { N: "900" }
                                        { BOOL: false }
                                    ]
                                }
                                mixedMap: {
                                    M: {
                                        stringAttr1: { S: "value1" }
                                        numberAttr1: { N: "123" }
                                        boolAttr1: { BOOL: false }
                                        stringAttr2: { S: "value2" }
                                        numberAttr2: { N: "456" }
                                        boolAttr2: { BOOL: true }
                                        stringAttr3: { S: "value3" }
                                        numberAttr3: { N: "789" }
                                        boolAttr3: { BOOL: false }
                                        stringAttr4: { S: "value4" }
                                        numberAttr4: { N: "101" }
                                        boolAttr4: { BOOL: true }
                                        stringAttr5: { S: "value5" }
                                        numberAttr5: { N: "202" }
                                        boolAttr5: { BOOL: false }
                                        stringAttr6: { S: "value6" }
                                        numberAttr6: { N: "303" }
                                        boolAttr6: { BOOL: true }
                                        stringAttr7: { S: "value7" }
                                        numberAttr7: { N: "404" }
                                        boolAttr7: { BOOL: false }
                                        stringAttr8: { S: "value8" }
                                        numberAttr8: { N: "505" }
                                        boolAttr8: { BOOL: true }
                                        stringAttr9: { S: "value9" }
                                        numberAttr9: { N: "606" }
                                        boolAttr9: { BOOL: false }
                                        stringAttr10: { S: "value10" }
                                        numberAttr10: { N: "707" }
                                        boolAttr10: { BOOL: true }
                                    }
                                }
                            }
                        }
                    ]
                }
                mixedMap: {
                    M: {
                        stringAttr1: { S: "value1" }
                        numberAttr1: { N: "123" }
                        boolAttr1: { BOOL: false }
                        stringAttr2: { S: "value2" }
                        numberAttr2: { N: "456" }
                        boolAttr2: { BOOL: true }
                        stringAttr3: { S: "value3" }
                        numberAttr3: { N: "789" }
                        boolAttr3: { BOOL: false }
                        stringAttr4: { S: "value4" }
                        numberAttr4: { N: "101" }
                        boolAttr4: { BOOL: true }
                        stringAttr5: { S: "value5" }
                        numberAttr5: { N: "202" }
                        boolAttr5: { BOOL: false }
                        stringAttr6: { S: "value6" }
                        numberAttr6: { N: "303" }
                        boolAttr6: { BOOL: true }
                        stringAttr7: { S: "value7" }
                        numberAttr7: { N: "404" }
                        boolAttr7: { BOOL: false }
                        stringAttr8: { S: "value8" }
                        numberAttr8: { N: "505" }
                        boolAttr8: { BOOL: true }
                        stringAttr9: { S: "value9" }
                        numberAttr9: { N: "606" }
                        boolAttr9: { BOOL: false }
                        stringAttr10: { S: "value10" }
                        numberAttr10: { N: "707" }
                        boolAttr10: { BOOL: true }
                        mediumMixedItem: {
                            M: {
                                id: { S: "test-id" }
                                stringSet: { SS: ["item1", "item2", "item3", "item4", "item5", "item6", "item7", "item8", "item9", "item10", "item11", "item12", "item13", "item14", "item15", "item16", "item17", "item18", "item19", "item20", "item21", "item22", "item23", "item24", "item25", "item26", "item27", "item28", "item29", "item30"] }
                                numberSet: { NS: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30"] }
                                list: {
                                    L: [
                                        { S: "listItem1" }
                                        { N: "42" }
                                        { BOOL: true }
                                        { S: "listItem4" }
                                        { N: "100" }
                                        { BOOL: false }
                                        { S: "listItem7" }
                                        { N: "200" }
                                        { BOOL: true }
                                        { S: "listItem10" }
                                        { N: "300" }
                                        { BOOL: false }
                                        { S: "listItem13" }
                                        { N: "400" }
                                        { BOOL: true }
                                        { S: "listItem16" }
                                        { N: "500" }
                                        { BOOL: false }
                                        { S: "listItem19" }
                                        { N: "600" }
                                        { BOOL: true }
                                        { S: "listItem22" }
                                        { N: "700" }
                                        { BOOL: false }
                                        { S: "listItem25" }
                                        { N: "800" }
                                        { BOOL: true }
                                        { S: "listItem28" }
                                        { N: "900" }
                                        { BOOL: false }
                                    ]
                                }
                                mixedMap: {
                                    M: {
                                        stringAttr1: { S: "value1" }
                                        numberAttr1: { N: "123" }
                                        boolAttr1: { BOOL: false }
                                        stringAttr2: { S: "value2" }
                                        numberAttr2: { N: "456" }
                                        boolAttr2: { BOOL: true }
                                        stringAttr3: { S: "value3" }
                                        numberAttr3: { N: "789" }
                                        boolAttr3: { BOOL: false }
                                        stringAttr4: { S: "value4" }
                                        numberAttr4: { N: "101" }
                                        boolAttr4: { BOOL: true }
                                        stringAttr5: { S: "value5" }
                                        numberAttr5: { N: "202" }
                                        boolAttr5: { BOOL: false }
                                        stringAttr6: { S: "value6" }
                                        numberAttr6: { N: "303" }
                                        boolAttr6: { BOOL: true }
                                        stringAttr7: { S: "value7" }
                                        numberAttr7: { N: "404" }
                                        boolAttr7: { BOOL: false }
                                        stringAttr8: { S: "value8" }
                                        numberAttr8: { N: "505" }
                                        boolAttr8: { BOOL: true }
                                        stringAttr9: { S: "value9" }
                                        numberAttr9: { N: "606" }
                                        boolAttr9: { BOOL: false }
                                        stringAttr10: { S: "value10" }
                                        numberAttr10: { N: "707" }
                                        boolAttr10: { BOOL: true }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_PutItemRequest_BinaryData_S"
        protocol: awsJson1_0
        documentation: """
        An item (small) with binary data.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: {
                    S: "test-id"
                }
                binary: {
                    B: "data data data data data data data data data data data data data data data data data data"
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_PutItemRequest_BinaryData_M"
        protocol: awsJson1_0
        documentation: """
        An item (medium) with binary data.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: {
                    S: "test-id"
                }
                binary1: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary2: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary3: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary4: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary5: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary6: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary7: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary8: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary9: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary10: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "awsJson1_0_PutItemRequest_BinaryData_L"
        protocol: awsJson1_0
        documentation: """
        An item (large) with binary data.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: {
                    S: "test-id"
                }
                map: {
                    M: {
                        binarySet: {
                            BS: [
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                            ]
                        }
                        binary: {
                            B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                        }
                        map: {
                            M: {
                                binarySet: {
                                    BS: [
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                    ]
                                }
                                binary: {
                                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                }
                                map: {
                                    M: {
                                        binarySet: {
                                            BS: [
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                            ]
                                        }
                                        binary: {
                                            B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    // section: cbor
    {
        id: "rpcv2Cbor_PutItemRequest_Baseline"
        protocol: rpcv2Cbor
        documentation: """
        This test gives baseline of serializing a minimal
        amount of data for a data-plane write.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: {
                    S: "test-id"
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_PutItemRequest_ShallowMap_S"
        protocol: rpcv2Cbor
        documentation: """
        Serializing a map (small) with many keys but minimal nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                attr1: { S: "value1" }
                attr2: { S: "value2" }
                attr3: { S: "value3" }
                attr4: { S: "value4" }
                attr5: { S: "value5" }
                attr6: { S: "value6" }
                attr7: { S: "value7" }
                attr8: { S: "value8" }
                attr9: { S: "value9" }
                attr10: { S: "value10" }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_PutItemRequest_ShallowMap_M"
        protocol: rpcv2Cbor
        documentation: """
        Serializing a map (medium) with many keys but minimal nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                attr1: { S: "value1" }
                attr2: { S: "value2" }
                attr3: { S: "value3" }
                attr4: { S: "value4" }
                attr5: { S: "value5" }
                attr6: { S: "value6" }
                attr7: { S: "value7" }
                attr8: { S: "value8" }
                attr9: { S: "value9" }
                attr10: { S: "value10" }
                attr11: { S: "value11" }
                attr12: { S: "value12" }
                attr13: { S: "value13" }
                attr14: { S: "value14" }
                attr15: { S: "value15" }
                attr16: { S: "value16" }
                attr17: { S: "value17" }
                attr18: { S: "value18" }
                attr19: { S: "value19" }
                attr20: { S: "value20" }
                attr21: { S: "value21" }
                attr22: { S: "value22" }
                attr23: { S: "value23" }
                attr24: { S: "value24" }
                attr25: { S: "value25" }
                attr26: { S: "value26" }
                attr27: { S: "value27" }
                attr28: { S: "value28" }
                attr29: { S: "value29" }
                attr30: { S: "value30" }
                attr31: { S: "value31" }
                attr32: { S: "value32" }
                attr33: { S: "value33" }
                attr34: { S: "value34" }
                attr35: { S: "value35" }
                attr36: { S: "value36" }
                attr37: { S: "value37" }
                attr38: { S: "value38" }
                attr39: { S: "value39" }
                attr40: { S: "value40" }
                attr41: { S: "value41" }
                attr42: { S: "value42" }
                attr43: { S: "value43" }
                attr44: { S: "value44" }
                attr45: { S: "value45" }
                attr46: { S: "value46" }
                attr47: { S: "value47" }
                attr48: { S: "value48" }
                attr49: { S: "value49" }
                attr50: { S: "value50" }
                attr51: { S: "value51" }
                attr52: { S: "value52" }
                attr53: { S: "value53" }
                attr54: { S: "value54" }
                attr55: { S: "value55" }
                attr56: { S: "value56" }
                attr57: { S: "value57" }
                attr58: { S: "value58" }
                attr59: { S: "value59" }
                attr60: { S: "value60" }
                attr61: { S: "value61" }
                attr62: { S: "value62" }
                attr63: { S: "value63" }
                attr64: { S: "value64" }
                attr65: { S: "value65" }
                attr66: { S: "value66" }
                attr67: { S: "value67" }
                attr68: { S: "value68" }
                attr69: { S: "value69" }
                attr70: { S: "value70" }
                attr71: { S: "value71" }
                attr72: { S: "value72" }
                attr73: { S: "value73" }
                attr74: { S: "value74" }
                attr75: { S: "value75" }
                attr76: { S: "value76" }
                attr77: { S: "value77" }
                attr78: { S: "value78" }
                attr79: { S: "value79" }
                attr80: { S: "value80" }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_PutItemRequest_ShallowMap_L"
        protocol: rpcv2Cbor
        documentation: """
        Serializing a map (large) with many keys but minimal nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                attr1: { S: "value1" }
                attr2: { S: "value2" }
                attr3: { S: "value3" }
                attr4: { S: "value4" }
                attr5: { S: "value5" }
                attr6: { S: "value6" }
                attr7: { S: "value7" }
                attr8: { S: "value8" }
                attr9: { S: "value9" }
                attr10: { S: "value10" }
                attr11: { S: "value11" }
                attr12: { S: "value12" }
                attr13: { S: "value13" }
                attr14: { S: "value14" }
                attr15: { S: "value15" }
                attr16: { S: "value16" }
                attr17: { S: "value17" }
                attr18: { S: "value18" }
                attr19: { S: "value19" }
                attr20: { S: "value20" }
                attr21: { S: "value21" }
                attr22: { S: "value22" }
                attr23: { S: "value23" }
                attr24: { S: "value24" }
                attr25: { S: "value25" }
                attr26: { S: "value26" }
                attr27: { S: "value27" }
                attr28: { S: "value28" }
                attr29: { S: "value29" }
                attr30: { S: "value30" }
                attr31: { S: "value31" }
                attr32: { S: "value32" }
                attr33: { S: "value33" }
                attr34: { S: "value34" }
                attr35: { S: "value35" }
                attr36: { S: "value36" }
                attr37: { S: "value37" }
                attr38: { S: "value38" }
                attr39: { S: "value39" }
                attr40: { S: "value40" }
                attr41: { S: "value41" }
                attr42: { S: "value42" }
                attr43: { S: "value43" }
                attr44: { S: "value44" }
                attr45: { S: "value45" }
                attr46: { S: "value46" }
                attr47: { S: "value47" }
                attr48: { S: "value48" }
                attr49: { S: "value49" }
                attr50: { S: "value50" }
                attr51: { S: "value51" }
                attr52: { S: "value52" }
                attr53: { S: "value53" }
                attr54: { S: "value54" }
                attr55: { S: "value55" }
                attr56: { S: "value56" }
                attr57: { S: "value57" }
                attr58: { S: "value58" }
                attr59: { S: "value59" }
                attr60: { S: "value60" }
                attr61: { S: "value61" }
                attr62: { S: "value62" }
                attr63: { S: "value63" }
                attr64: { S: "value64" }
                attr65: { S: "value65" }
                attr66: { S: "value66" }
                attr67: { S: "value67" }
                attr68: { S: "value68" }
                attr69: { S: "value69" }
                attr70: { S: "value70" }
                attr71: { S: "value71" }
                attr72: { S: "value72" }
                attr73: { S: "value73" }
                attr74: { S: "value74" }
                attr75: { S: "value75" }
                attr76: { S: "value76" }
                attr77: { S: "value77" }
                attr78: { S: "value78" }
                attr79: { S: "value79" }
                attr80: { S: "value80" }
                attr81: { S: "value81" }
                attr82: { S: "value82" }
                attr83: { S: "value83" }
                attr84: { S: "value84" }
                attr85: { S: "value85" }
                attr86: { S: "value86" }
                attr87: { S: "value87" }
                attr88: { S: "value88" }
                attr89: { S: "value89" }
                attr90: { S: "value90" }
                attr91: { S: "value91" }
                attr92: { S: "value92" }
                attr93: { S: "value93" }
                attr94: { S: "value94" }
                attr95: { S: "value95" }
                attr96: { S: "value96" }
                attr97: { S: "value97" }
                attr98: { S: "value98" }
                attr99: { S: "value99" }
                attr100: { S: "value100" }
                attr101: { S: "value101" }
                attr102: { S: "value102" }
                attr103: { S: "value103" }
                attr104: { S: "value104" }
                attr105: { S: "value105" }
                attr106: { S: "value106" }
                attr107: { S: "value107" }
                attr108: { S: "value108" }
                attr109: { S: "value109" }
                attr110: { S: "value110" }
                attr111: { S: "value111" }
                attr112: { S: "value112" }
                attr113: { S: "value113" }
                attr114: { S: "value114" }
                attr115: { S: "value115" }
                attr116: { S: "value116" }
                attr117: { S: "value117" }
                attr118: { S: "value118" }
                attr119: { S: "value119" }
                attr120: { S: "value120" }
                attr121: { S: "value121" }
                attr122: { S: "value122" }
                attr123: { S: "value123" }
                attr124: { S: "value124" }
                attr125: { S: "value125" }
                attr126: { S: "value126" }
                attr127: { S: "value127" }
                attr128: { S: "value128" }
                attr129: { S: "value129" }
                attr130: { S: "value130" }
                attr131: { S: "value131" }
                attr132: { S: "value132" }
                attr133: { S: "value133" }
                attr134: { S: "value134" }
                attr135: { S: "value135" }
                attr136: { S: "value136" }
                attr137: { S: "value137" }
                attr138: { S: "value138" }
                attr139: { S: "value139" }
                attr140: { S: "value140" }
                attr141: { S: "value141" }
                attr142: { S: "value142" }
                attr143: { S: "value143" }
                attr144: { S: "value144" }
                attr145: { S: "value145" }
                attr146: { S: "value146" }
                attr147: { S: "value147" }
                attr148: { S: "value148" }
                attr149: { S: "value149" }
                attr150: { S: "value150" }
                attr151: { S: "value151" }
                attr152: { S: "value152" }
                attr153: { S: "value153" }
                attr154: { S: "value154" }
                attr155: { S: "value155" }
                attr156: { S: "value156" }
                attr157: { S: "value157" }
                attr158: { S: "value158" }
                attr159: { S: "value159" }
                attr160: { S: "value160" }
                attr161: { S: "value161" }
                attr162: { S: "value162" }
                attr163: { S: "value163" }
                attr164: { S: "value164" }
                attr165: { S: "value165" }
                attr166: { S: "value166" }
                attr167: { S: "value167" }
                attr168: { S: "value168" }
                attr169: { S: "value169" }
                attr170: { S: "value170" }
                attr171: { S: "value171" }
                attr172: { S: "value172" }
                attr173: { S: "value173" }
                attr174: { S: "value174" }
                attr175: { S: "value175" }
                attr176: { S: "value176" }
                attr177: { S: "value177" }
                attr178: { S: "value178" }
                attr179: { S: "value179" }
                attr180: { S: "value180" }
                attr181: { S: "value181" }
                attr182: { S: "value182" }
                attr183: { S: "value183" }
                attr184: { S: "value184" }
                attr185: { S: "value185" }
                attr186: { S: "value186" }
                attr187: { S: "value187" }
                attr188: { S: "value188" }
                attr189: { S: "value189" }
                attr190: { S: "value190" }
                attr191: { S: "value191" }
                attr192: { S: "value192" }
                attr193: { S: "value193" }
                attr194: { S: "value194" }
                attr195: { S: "value195" }
                attr196: { S: "value196" }
                attr197: { S: "value197" }
                attr198: { S: "value198" }
                attr199: { S: "value199" }
                attr200: { S: "value200" }
                attr201: { S: "value201" }
                attr202: { S: "value202" }
                attr203: { S: "value203" }
                attr204: { S: "value204" }
                attr205: { S: "value205" }
                attr206: { S: "value206" }
                attr207: { S: "value207" }
                attr208: { S: "value208" }
                attr209: { S: "value209" }
                attr210: { S: "value210" }
                attr211: { S: "value211" }
                attr212: { S: "value212" }
                attr213: { S: "value213" }
                attr214: { S: "value214" }
                attr215: { S: "value215" }
                attr216: { S: "value216" }
                attr217: { S: "value217" }
                attr218: { S: "value218" }
                attr219: { S: "value219" }
                attr220: { S: "value220" }
                attr221: { S: "value221" }
                attr222: { S: "value222" }
                attr223: { S: "value223" }
                attr224: { S: "value224" }
                attr225: { S: "value225" }
                attr226: { S: "value226" }
                attr227: { S: "value227" }
                attr228: { S: "value228" }
                attr229: { S: "value229" }
                attr230: { S: "value230" }
                attr231: { S: "value231" }
                attr232: { S: "value232" }
                attr233: { S: "value233" }
                attr234: { S: "value234" }
                attr235: { S: "value235" }
                attr236: { S: "value236" }
                attr237: { S: "value237" }
                attr238: { S: "value238" }
                attr239: { S: "value239" }
                attr240: { S: "value240" }
                attr241: { S: "value241" }
                attr242: { S: "value242" }
                attr243: { S: "value243" }
                attr244: { S: "value244" }
                attr245: { S: "value245" }
                attr246: { S: "value246" }
                attr247: { S: "value247" }
                attr248: { S: "value248" }
                attr249: { S: "value249" }
                attr250: { S: "value250" }
                attr251: { S: "value251" }
                attr252: { S: "value252" }
                attr253: { S: "value253" }
                attr254: { S: "value254" }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_PutItemRequest_Nested_M"
        protocol: rpcv2Cbor
        documentation: """
        A narrow item with moderate nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                nested: {
                    M: {
                        level1: {
                            M: {
                                level2: {
                                    M: {
                                        level3: {
                                            M: {
                                                level4: {
                                                    M: {
                                                        deepValue: { S: "deep-nested-value" }
                                                    }
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_PutItemRequest_Nested_L"
        protocol: rpcv2Cbor
        documentation: """
        A narrow item with deep nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                nested: {
                    M: {
                        level1: {
                            L: [{
                                M: {
                                    level3: {
                                        L: [{
                                            M: {
                                                level5: {
                                                    L: [{
                                                        M: {
                                                            level7: {
                                                                L: [{
                                                                    M: {
                                                                        level9: {
                                                                            L: [{
                                                                                M: {
                                                                                    level11: {
                                                                                        L: [{
                                                                                            M: {
                                                                                                level13: {
                                                                                                    L: [{
                                                                                                        M: {
                                                                                                            level15: {
                                                                                                                L: [{
                                                                                                                    M: {
                                                                                                                        level17: {
                                                                                                                            L: [{
                                                                                                                                M: {
                                                                                                                                    level19: {
                                                                                                                                        L: [{
                                                                                                                                            M: {
                                                                                                                                                level21: {
                                                                                                                                                    L: [{
                                                                                                                                                        M: {
                                                                                                                                                            level23: {
                                                                                                                                                                L: [{
                                                                                                                                                                    M: {
                                                                                                                                                                        level25: {
                                                                                                                                                                            L: [{
                                                                                                                                                                                M: {
                                                                                                                                                                                    level27: {
                                                                                                                                                                                        L: [{
                                                                                                                                                                                            M: {
                                                                                                                                                                                                deepValue: { S: "smithy parser limit reached" }
                                                                                                                                                                                            }
                                                                                                                                                                                        }]
                                                                                                                                                                                    }
                                                                                                                                                                                }
                                                                                                                                                                            }]
                                                                                                                                                                        }
                                                                                                                                                                    }
                                                                                                                                                                }]
                                                                                                                                                            }
                                                                                                                                                        }
                                                                                                                                                    }]
                                                                                                                                                }
                                                                                                                                            }
                                                                                                                                        }]
                                                                                                                                    }
                                                                                                                                }
                                                                                                                            }]
                                                                                                                        }
                                                                                                                    }
                                                                                                                }]
                                                                                                            }
                                                                                                        }
                                                                                                    }]
                                                                                                }
                                                                                            }
                                                                                        }]
                                                                                    }
                                                                                }
                                                                            }]
                                                                        }
                                                                    }
                                                                }]
                                                            }
                                                        }
                                                    }]
                                                }
                                            }
                                        }]
                                    }
                                }
                            }]
                        }
                    }
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_PutItemRequest_MixedItem_S"
        protocol: rpcv2Cbor
        documentation: """
        An item (small) that uses mixed AttributeValue types, including nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                stringSet: { SS: ["item1", "item2", "item3"] }
                numberSet: { NS: ["1", "2", "3"] }
                list: {
                    L: [
                        { S: "listItem1" }
                        { N: "42" }
                        { BOOL: true }
                    ]
                }
                mixedMap: {
                    M: {
                        stringAttr: { S: "value" }
                        numberAttr: { N: "123" }
                        boolAttr: { BOOL: false }
                    }
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_PutItemRequest_MixedItem_M"
        protocol: rpcv2Cbor
        documentation: """
        An item (medium) that uses mixed AttributeValue types, including nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                stringSet: { SS: ["item1", "item2", "item3", "item4", "item5", "item6", "item7", "item8", "item9", "item10", "item11", "item12", "item13", "item14", "item15", "item16", "item17", "item18", "item19", "item20", "item21", "item22", "item23", "item24", "item25", "item26", "item27", "item28", "item29", "item30"] }
                numberSet: { NS: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30"] }
                list: {
                    L: [
                        { S: "listItem1" }
                        { N: "42" }
                        { BOOL: true }
                        { S: "listItem4" }
                        { N: "100" }
                        { BOOL: false }
                        { S: "listItem7" }
                        { N: "200" }
                        { BOOL: true }
                        { S: "listItem10" }
                        { N: "300" }
                        { BOOL: false }
                        { S: "listItem13" }
                        { N: "400" }
                        { BOOL: true }
                        { S: "listItem16" }
                        { N: "500" }
                        { BOOL: false }
                        { S: "listItem19" }
                        { N: "600" }
                        { BOOL: true }
                        { S: "listItem22" }
                        { N: "700" }
                        { BOOL: false }
                        { S: "listItem25" }
                        { N: "800" }
                        { BOOL: true }
                        { S: "listItem28" }
                        { N: "900" }
                        { BOOL: false }
                    ]
                }
                mixedMap: {
                    M: {
                        stringAttr1: { S: "value1" }
                        numberAttr1: { N: "123" }
                        boolAttr1: { BOOL: false }
                        stringAttr2: { S: "value2" }
                        numberAttr2: { N: "456" }
                        boolAttr2: { BOOL: true }
                        stringAttr3: { S: "value3" }
                        numberAttr3: { N: "789" }
                        boolAttr3: { BOOL: false }
                        stringAttr4: { S: "value4" }
                        numberAttr4: { N: "101" }
                        boolAttr4: { BOOL: true }
                        stringAttr5: { S: "value5" }
                        numberAttr5: { N: "202" }
                        boolAttr5: { BOOL: false }
                        stringAttr6: { S: "value6" }
                        numberAttr6: { N: "303" }
                        boolAttr6: { BOOL: true }
                        stringAttr7: { S: "value7" }
                        numberAttr7: { N: "404" }
                        boolAttr7: { BOOL: false }
                        stringAttr8: { S: "value8" }
                        numberAttr8: { N: "505" }
                        boolAttr8: { BOOL: true }
                        stringAttr9: { S: "value9" }
                        numberAttr9: { N: "606" }
                        boolAttr9: { BOOL: false }
                        stringAttr10: { S: "value10" }
                        numberAttr10: { N: "707" }
                        boolAttr10: { BOOL: true }
                    }
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_PutItemRequest_MixedItem_L"
        protocol: rpcv2Cbor
        documentation: """
        An item (large) that uses mixed AttributeValue types, including nesting.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: { S: "test-id" }
                stringSet: { SS: ["item1", "item2", "item3", "item4", "item5", "item6", "item7", "item8", "item9", "item10", "item11", "item12", "item13", "item14", "item15", "item16", "item17", "item18", "item19", "item20", "item21", "item22", "item23", "item24", "item25", "item26", "item27", "item28", "item29", "item30"] }
                numberSet: { NS: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30"] }
                list: {
                    L: [
                        { S: "listItem1" }
                        { N: "42" }
                        { BOOL: true }
                        { S: "listItem4" }
                        { N: "100" }
                        { BOOL: false }
                        { S: "listItem7" }
                        { N: "200" }
                        { BOOL: true }
                        { S: "listItem10" }
                        { N: "300" }
                        { BOOL: false }
                        { S: "listItem13" }
                        { N: "400" }
                        { BOOL: true }
                        { S: "listItem16" }
                        { N: "500" }
                        { BOOL: false }
                        { S: "listItem19" }
                        { N: "600" }
                        { BOOL: true }
                        { S: "listItem22" }
                        { N: "700" }
                        { BOOL: false }
                        { S: "listItem25" }
                        { N: "800" }
                        { BOOL: true }
                        { S: "listItem28" }
                        { N: "900" }
                        { BOOL: false }
                        {
                            M: {
                                id: { S: "test-id" }
                                stringSet: { SS: ["item1", "item2", "item3", "item4", "item5", "item6", "item7", "item8", "item9", "item10", "item11", "item12", "item13", "item14", "item15", "item16", "item17", "item18", "item19", "item20", "item21", "item22", "item23", "item24", "item25", "item26", "item27", "item28", "item29", "item30"] }
                                numberSet: { NS: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30"] }
                                list: {
                                    L: [
                                        { S: "listItem1" }
                                        { N: "42" }
                                        { BOOL: true }
                                        { S: "listItem4" }
                                        { N: "100" }
                                        { BOOL: false }
                                        { S: "listItem7" }
                                        { N: "200" }
                                        { BOOL: true }
                                        { S: "listItem10" }
                                        { N: "300" }
                                        { BOOL: false }
                                        { S: "listItem13" }
                                        { N: "400" }
                                        { BOOL: true }
                                        { S: "listItem16" }
                                        { N: "500" }
                                        { BOOL: false }
                                        { S: "listItem19" }
                                        { N: "600" }
                                        { BOOL: true }
                                        { S: "listItem22" }
                                        { N: "700" }
                                        { BOOL: false }
                                        { S: "listItem25" }
                                        { N: "800" }
                                        { BOOL: true }
                                        { S: "listItem28" }
                                        { N: "900" }
                                        { BOOL: false }
                                    ]
                                }
                                mixedMap: {
                                    M: {
                                        stringAttr1: { S: "value1" }
                                        numberAttr1: { N: "123" }
                                        boolAttr1: { BOOL: false }
                                        stringAttr2: { S: "value2" }
                                        numberAttr2: { N: "456" }
                                        boolAttr2: { BOOL: true }
                                        stringAttr3: { S: "value3" }
                                        numberAttr3: { N: "789" }
                                        boolAttr3: { BOOL: false }
                                        stringAttr4: { S: "value4" }
                                        numberAttr4: { N: "101" }
                                        boolAttr4: { BOOL: true }
                                        stringAttr5: { S: "value5" }
                                        numberAttr5: { N: "202" }
                                        boolAttr5: { BOOL: false }
                                        stringAttr6: { S: "value6" }
                                        numberAttr6: { N: "303" }
                                        boolAttr6: { BOOL: true }
                                        stringAttr7: { S: "value7" }
                                        numberAttr7: { N: "404" }
                                        boolAttr7: { BOOL: false }
                                        stringAttr8: { S: "value8" }
                                        numberAttr8: { N: "505" }
                                        boolAttr8: { BOOL: true }
                                        stringAttr9: { S: "value9" }
                                        numberAttr9: { N: "606" }
                                        boolAttr9: { BOOL: false }
                                        stringAttr10: { S: "value10" }
                                        numberAttr10: { N: "707" }
                                        boolAttr10: { BOOL: true }
                                    }
                                }
                            }
                        }
                    ]
                }
                mixedMap: {
                    M: {
                        stringAttr1: { S: "value1" }
                        numberAttr1: { N: "123" }
                        boolAttr1: { BOOL: false }
                        stringAttr2: { S: "value2" }
                        numberAttr2: { N: "456" }
                        boolAttr2: { BOOL: true }
                        stringAttr3: { S: "value3" }
                        numberAttr3: { N: "789" }
                        boolAttr3: { BOOL: false }
                        stringAttr4: { S: "value4" }
                        numberAttr4: { N: "101" }
                        boolAttr4: { BOOL: true }
                        stringAttr5: { S: "value5" }
                        numberAttr5: { N: "202" }
                        boolAttr5: { BOOL: false }
                        stringAttr6: { S: "value6" }
                        numberAttr6: { N: "303" }
                        boolAttr6: { BOOL: true }
                        stringAttr7: { S: "value7" }
                        numberAttr7: { N: "404" }
                        boolAttr7: { BOOL: false }
                        stringAttr8: { S: "value8" }
                        numberAttr8: { N: "505" }
                        boolAttr8: { BOOL: true }
                        stringAttr9: { S: "value9" }
                        numberAttr9: { N: "606" }
                        boolAttr9: { BOOL: false }
                        stringAttr10: { S: "value10" }
                        numberAttr10: { N: "707" }
                        boolAttr10: { BOOL: true }
                        mediumMixedItem: {
                            M: {
                                id: { S: "test-id" }
                                stringSet: { SS: ["item1", "item2", "item3", "item4", "item5", "item6", "item7", "item8", "item9", "item10", "item11", "item12", "item13", "item14", "item15", "item16", "item17", "item18", "item19", "item20", "item21", "item22", "item23", "item24", "item25", "item26", "item27", "item28", "item29", "item30"] }
                                numberSet: { NS: ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30"] }
                                list: {
                                    L: [
                                        { S: "listItem1" }
                                        { N: "42" }
                                        { BOOL: true }
                                        { S: "listItem4" }
                                        { N: "100" }
                                        { BOOL: false }
                                        { S: "listItem7" }
                                        { N: "200" }
                                        { BOOL: true }
                                        { S: "listItem10" }
                                        { N: "300" }
                                        { BOOL: false }
                                        { S: "listItem13" }
                                        { N: "400" }
                                        { BOOL: true }
                                        { S: "listItem16" }
                                        { N: "500" }
                                        { BOOL: false }
                                        { S: "listItem19" }
                                        { N: "600" }
                                        { BOOL: true }
                                        { S: "listItem22" }
                                        { N: "700" }
                                        { BOOL: false }
                                        { S: "listItem25" }
                                        { N: "800" }
                                        { BOOL: true }
                                        { S: "listItem28" }
                                        { N: "900" }
                                        { BOOL: false }
                                    ]
                                }
                                mixedMap: {
                                    M: {
                                        stringAttr1: { S: "value1" }
                                        numberAttr1: { N: "123" }
                                        boolAttr1: { BOOL: false }
                                        stringAttr2: { S: "value2" }
                                        numberAttr2: { N: "456" }
                                        boolAttr2: { BOOL: true }
                                        stringAttr3: { S: "value3" }
                                        numberAttr3: { N: "789" }
                                        boolAttr3: { BOOL: false }
                                        stringAttr4: { S: "value4" }
                                        numberAttr4: { N: "101" }
                                        boolAttr4: { BOOL: true }
                                        stringAttr5: { S: "value5" }
                                        numberAttr5: { N: "202" }
                                        boolAttr5: { BOOL: false }
                                        stringAttr6: { S: "value6" }
                                        numberAttr6: { N: "303" }
                                        boolAttr6: { BOOL: true }
                                        stringAttr7: { S: "value7" }
                                        numberAttr7: { N: "404" }
                                        boolAttr7: { BOOL: false }
                                        stringAttr8: { S: "value8" }
                                        numberAttr8: { N: "505" }
                                        boolAttr8: { BOOL: true }
                                        stringAttr9: { S: "value9" }
                                        numberAttr9: { N: "606" }
                                        boolAttr9: { BOOL: false }
                                        stringAttr10: { S: "value10" }
                                        numberAttr10: { N: "707" }
                                        boolAttr10: { BOOL: true }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_PutItemRequest_BinaryData_S"
        protocol: rpcv2Cbor
        documentation: """
        An item (small) with binary data.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: {
                    S: "test-id"
                }
                binary: {
                    B: "data data data data data data data data data data data data data data data data data data"
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_PutItemRequest_BinaryData_M"
        protocol: rpcv2Cbor
        documentation: """
        An item (medium) with binary data.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: {
                    S: "test-id"
                }
                binary1: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary2: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary3: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary4: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary5: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary6: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary7: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary8: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary9: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
                binary10: {
                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                }
            }
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "rpcv2Cbor_PutItemRequest_BinaryData_L"
        protocol: rpcv2Cbor
        documentation: """
        An item (large) with binary data.
        """
        method: "POST"
        uri: "/"
        params: {
            TableName: "test-table"
            Item: {
                id: {
                    S: "test-id"
                }
                map: {
                    M: {
                        binarySet: {
                            BS: [
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                            ]
                        }
                        binary: {
                            B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                        }
                        map: {
                            M: {
                                binarySet: {
                                    BS: [
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                    ]
                                }
                                binary: {
                                    B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                }
                                map: {
                                    M: {
                                        binarySet: {
                                            BS: [
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                                "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                            ]
                                        }
                                        binary: {
                                            B: "data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data data"
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        tags: ["serde-benchmark"]
    }
])
operation PutItem {
    input: PutItemInput
    output: PutItemOutput
}

structure PutItemInput {
    @required
    TableName: String
    @required
    Item: AttributeValueMap

    Expected: ExpectedAttributeMap
    ReturnValues: String
    ReturnConsumedCapacity: String
    ReturnItemCollectionMetrics: String
    ConditionalOperator: String
    ConditionExpression: String
    ExpressionAttributeNames: ExpressionAttributeNameMap
    ExpressionAttributeValues: ExpressionAttributeValueMap
    ReturnValuesOnConditionCheckFailure: String
}

structure PutItemOutput {
    Attributes: AttributeValueMap
    ConsumedCapacity: ConsumedCapacity
    ItemCollectionMetrics: ItemCollectionMetrics
}

map ExpectedAttributeMap {
    key: String
    value: ExpectedAttributeValue
}

structure ExpectedAttributeValue {
    Value: AttributeValue
    Exists: Boolean
    ComparisonOperator: String
    AttributeValueList: AttributeValueList
}

map ExpressionAttributeValueMap {
    key: String
    value: AttributeValue
}

structure ItemCollectionMetrics {
    ItemCollectionKey: AttributeValueMap
    SizeEstimateRangeGB: SizeEstimateRange
}

list SizeEstimateRange {
    member: Double
}
