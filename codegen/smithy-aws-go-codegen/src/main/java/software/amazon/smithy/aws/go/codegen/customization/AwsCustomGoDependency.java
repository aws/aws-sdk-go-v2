/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *  http://aws.amazon.com/apache2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.go.codegen.GoDependency;

/**
 * A class of constants for dependencies used by this package.
 */
public final class AwsCustomGoDependency extends AwsGoDependency {
    public static final GoDependency DYNAMODB_CUSTOMIZATION = aws(
            "service/dynamodb/internal/customizations", "ddbcust");
    public static final GoDependency S3_CUSTOMIZATION = aws("service/s3/internal/customizations", "s3cust");
    public static final GoDependency APIGATEWAY_CUSTOMIZATION = aws(
            "service/apigateway/internal/customizations", "agcust");
    public static final GoDependency GLACIER_CUSTOMIZATION = aws(
            "service/glacier/internal/customizations", "glaciercust");
    public static final GoDependency S3_SHARED_CUSTOMIZATION = awsModuleDep(
            "service/internal/s3shared", null, "s3shared");

    private AwsCustomGoDependency() {
        super();
    }
}
