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
public class AwsGoDependency {
    public static final GoDependency AWS_REST_JSON_PROTOCOL = aws("aws/protocol/restjson");
    public static final GoDependency AWS_QUERY_PROTOCOL = aws("aws/protocol/query");
    public static final GoDependency AWS_EC2QUERY_PROTOCOL = aws("aws/protocol/ec2query");
    public static final GoDependency AWS_CORE = aws("aws");
    public static final GoDependency AWS_MIDDLEWARE = aws("aws/middleware", "awsmiddleware");
    public static final GoDependency AWS_RETRY = aws("aws/retry");
    public static final GoDependency AWS_SIGNER_V4 = aws("aws/signer/v4");
    public static final GoDependency AWS_ENDPOINTS = aws("internal/endpoints");
    public static final GoDependency AWS_XML = aws("aws/protocol/xml", "awsxml");
    public static final GoDependency AWS_HTTP_TRANSPORT = aws("aws/transport/http", "awshttp");
    public static final GoDependency AWSTESTING_UNIT = aws("internal/awstesting/unit");

    public static final GoDependency S3_SHARED_CONFIG = aws("service/internal/s3shared/config", "s3sharedconfig");

    public static final GoDependency REGEXP = SmithyGoDependency.stdlib("regexp");

    public static final String AWS_SOURCE_PATH = "github.com/aws/aws-sdk-go-v2";

    protected AwsGoDependency() {
    }

    protected static GoDependency aws(String relativePath) {
        return aws(relativePath, null);
    }

    protected static GoDependency aws(String relativePath, String alias) {
        return module(AWS_SOURCE_PATH, relativePath, Versions.AWS_SDK, alias);
    }

    /**
     * awsModuleDep returns a GoDependency relative to the version of AWS_SDK core.
     *
     * @param moduleImportPath the module path within aws sdk to be added as go mod dependency.
     * @param relativePath the relative path which will be used as import path relative to aws sdk path.
     * @param version the version of the aws module dependency to be imported
     * @param alias the go import alias.
     * @return GoDependency
     */
    protected static GoDependency awsModuleDep(
         String moduleImportPath,
         String relativePath,
         String version,
         String alias
    ) {
        moduleImportPath = AWS_SOURCE_PATH+ "/" + moduleImportPath;
        return module(moduleImportPath, relativePath, version, alias);
    }

    protected static GoDependency module(
            String moduleImportPath,
            String relativePath,
            String version,
            String alias
    ) {
        String importPath = moduleImportPath;
        if (relativePath != null) {
            importPath = importPath + "/" + relativePath;
        }
        return GoDependency.moduleDependency(moduleImportPath, importPath, version, alias);
    }

    private static final class Versions {
        private static final String AWS_SDK = "v1.0.1-0.20210122214637-6cf9ad2f8e2f";
    }
}
