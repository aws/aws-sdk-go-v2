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
    public static final GoDependency S3CONTROL_CUSTOMIZATION = aws("service/s3control/internal/customizations", "s3controlcust");
    public static final GoDependency APIGATEWAY_CUSTOMIZATION = aws(
            "service/apigateway/internal/customizations", "agcust");
    public static final GoDependency GLACIER_CUSTOMIZATION = aws(
            "service/glacier/internal/customizations", "glaciercust");
    public static final GoDependency S3_SHARED_CUSTOMIZATION = awsModuleDep(
            "service/internal/s3shared", null, Versions.INTERNAL_S3SHARED, "s3shared");
    public static final GoDependency ACCEPT_ENCODING_CUSTOMIZATION = awsModuleDep(
            "service/internal/accept-encoding", null, Versions.INTERNAL_ACCEPTENCODING, "acceptencodingcust");
    public static final GoDependency KINESIS_CUSTOMIZATION = aws(
            "service/kinesis/internal/customizations", "kinesiscust");
    public static final GoDependency MACHINE_LEARNING_CUSTOMIZATION = aws(
            "service/machinelearning/internal/customizations", "mlcust");
    public static final GoDependency ROUTE53_CUSTOMIZATION = aws(
            "service/route53/internal/customizations", "route53cust");
    public static final GoDependency PRESIGNEDURL_CUSTOMIZATION = awsModuleDep(
            "service/internal/presigned-url", null, Versions.INTERNAL_PRESIGNURL, "presignedurlcust");
    public static final GoDependency EVENTBRIDGE_CUSTOMIZATION = aws("service/eventbridge/internal/customizations", "ebcust");

    private AwsCustomGoDependency() {
        super();
    }

    private static final class Versions {
        private static final String INTERNAL_S3SHARED = "v1.2.3";
        private static final String INTERNAL_ACCEPTENCODING = "v1.0.5";
        private static final String INTERNAL_PRESIGNURL = "v1.0.7";
    }
}
