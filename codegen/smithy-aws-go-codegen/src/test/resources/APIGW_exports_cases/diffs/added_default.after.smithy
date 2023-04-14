$version: "2.0"

namespace com.amazonaws.greengrass

use aws.api#service

@service(sdkId: "Greengrass")
service Greengrass {}

boolean __boolean

integer __integer

@default(false)
boolean NewDefaultedBoolean

@default(0)
integer NewDefaultedInteger
