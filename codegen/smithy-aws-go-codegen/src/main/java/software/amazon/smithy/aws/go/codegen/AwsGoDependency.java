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

package software.amazon.smithy.aws.go.codegen;

import software.amazon.smithy.go.codegen.GoDependency;
import software.amazon.smithy.go.codegen.SmithyGoDependency;

/**
 * A class of constants for dependencies used by this package.
 */
public final class AwsGoDependency {
    public static final GoDependency AWS_REST_JSON_PROTOCOL = aws("aws/protocol/restjson");
    public static final GoDependency AWS_QUERY_PROTOCOL = aws("aws/protocol/query");
    public static final GoDependency AWS_CORE = aws("aws");
    public static final GoDependency AWS_MIDDLEWARE = aws("aws/middleware", "awsmiddleware");
    public static final GoDependency AWS_RETRY = aws("aws/retry");
    public static final GoDependency AWS_SIGNER_V4 = aws("aws/signer/v4");
    public static final GoDependency AWS_ENDPOINTS = aws("aws/endpoints/v2", "endpoints");

    public static final GoDependency REGEXP = SmithyGoDependency.stdlib("regexp");

    public static final String AWS_SOURCE_PATH = "github.com/aws/aws-sdk-go-v2";

    private AwsGoDependency() {
    }

    private static GoDependency aws(String relativePath) {
        return aws(relativePath, null);
    }

    private static GoDependency aws(String relativePath, String alias) {
        String importPath = AWS_SOURCE_PATH;
        if (relativePath != null) {
            importPath = importPath + "/" + relativePath;
        }
        return GoDependency.moduleDependency(AWS_SOURCE_PATH, importPath, Versions.AWS_SDK, alias);
    }

    private static final class Versions {
    private static final String AWS_SDK = "v0.0.0-20200720171838-25cc9a8769cf";
    }
}
