/*
 * Copyright 2023 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

import software.amazon.smithy.aws.go.codegen.customization.AwsCustomGoDependency;
import software.amazon.smithy.codegen.core.Symbol;

/**
 * Collection of Symbol constants for types in the aws-sdk-go-v2 runtime.
 */
public final class SdkGoTypes {
    private SdkGoTypes() { }

    public static final class Aws {
        public static final Symbol String = AwsGoDependency.AWS_CORE.valueSymbol("String");
        public static final Symbol Bool = AwsGoDependency.AWS_CORE.valueSymbol("Bool");

        public static final Symbol FIPSEndpointStateEnabled = AwsGoDependency.AWS_CORE.valueSymbol("FIPSEndpointStateEnabled");
        public static final Symbol DualStackEndpointStateEnabled = AwsGoDependency.AWS_CORE.valueSymbol("DualStackEndpointStateEnabled");

        public static final Symbol IsCredentialsProvider = AwsGoDependency.AWS_CORE.valueSymbol("IsCredentialsProvider");
        public static final Symbol AnonymousCredentials = AwsGoDependency.AWS_CORE.pointableSymbol("AnonymousCredentials");
        public static final Symbol AccountIDEndpointMode = AwsGoDependency.AWS_CORE.valueSymbol("AccountIDEndpointMode");
        public static final Symbol AccountIDEndpointModeUnset = AwsGoDependency.AWS_CORE.valueSymbol("AccountIDEndpointModeUnset");
        public static final Symbol AccountIDEndpointModePreferred = AwsGoDependency.AWS_CORE.valueSymbol("AccountIDEndpointModePreferred");
        public static final Symbol AccountIDEndpointModeRequired = AwsGoDependency.AWS_CORE.valueSymbol("AccountIDEndpointModeRequired");
        public static final Symbol AccountIDEndpointModeDisabled = AwsGoDependency.AWS_CORE.valueSymbol("AccountIDEndpointModeDisabled");

        public static final Symbol RequestChecksumCalculation = AwsGoDependency.AWS_CORE.valueSymbol("RequestChecksumCalculation");
        public static final Symbol ResponseChecksumValidation = AwsGoDependency.AWS_CORE.valueSymbol("ResponseChecksumValidation");

        public static final class Middleware {
            public static final Symbol GetRequiresLegacyEndpoints = AwsGoDependency.AWS_MIDDLEWARE.valueSymbol("GetRequiresLegacyEndpoints");
            public static final Symbol GetSigningName = AwsGoDependency.AWS_MIDDLEWARE.valueSymbol("GetSigningName");
            public static final Symbol GetSigningRegion = AwsGoDependency.AWS_MIDDLEWARE.valueSymbol("GetSigningRegion");
            public static final Symbol SetSigningName = AwsGoDependency.AWS_MIDDLEWARE.valueSymbol("SetSigningName");
            public static final Symbol SetSigningRegion = AwsGoDependency.AWS_MIDDLEWARE.valueSymbol("SetSigningRegion");
        }


        public static final class Retry {
            public static final Symbol Attempt = AwsGoDependency.AWS_RETRY.pointableSymbol("Attempt");
            public static final Symbol MetricsHeader = AwsGoDependency.AWS_RETRY.pointableSymbol("MetricsHeader");
        }
    }

    public static final class Internal {
        public static final class Auth {
            public static final Symbol HTTPAuthScheme = AwsGoDependency.INTERNAL_AUTH.pointableSymbol("HTTPAuthScheme");
            public static final Symbol NewHTTPAuthScheme = AwsGoDependency.INTERNAL_AUTH.valueSymbol("NewHTTPAuthScheme");

            public static final class Smithy {
                public static final Symbol CredentialsAdapter = AwsGoDependency.INTERNAL_AUTH_SMITHY.pointableSymbol("CredentialsAdapter");
                public static final Symbol CredentialsProviderAdapter = AwsGoDependency.INTERNAL_AUTH_SMITHY.pointableSymbol("CredentialsProviderAdapter");
                public static final Symbol V4SignerAdapter = AwsGoDependency.INTERNAL_AUTH_SMITHY.pointableSymbol("V4SignerAdapter");
                public static final Symbol BearerTokenAdapter = AwsGoDependency.INTERNAL_AUTH_SMITHY.pointableSymbol("BearerTokenAdapter");
                public static final Symbol BearerTokenProviderAdapter = AwsGoDependency.INTERNAL_AUTH_SMITHY.pointableSymbol("BearerTokenProviderAdapter");
                public static final Symbol BearerTokenSignerAdapter = AwsGoDependency.INTERNAL_AUTH_SMITHY.pointableSymbol("BearerTokenSignerAdapter");
            }
        }

        public static final class Context {
            public static final Symbol SetS3Backend = AwsGoDependency.INTERNAL_CONTEXT.valueSymbol("SetS3Backend");
        }

        public static final class Endpoints {
            public static final Symbol MapFIPSRegion = AwsGoDependency.INTERNAL_ENDPOINTS.valueSymbol("MapFIPSRegion");
        }

        public static final class V4A {
            public static final Symbol SymmetricCredentialAdaptor = AwsGoDependency.INTERNAL_SIGV4A.pointableSymbol("SymmetricCredentialAdaptor");

            public static final Symbol CredentialsAdapter = AwsGoDependency.INTERNAL_SIGV4A.pointableSymbol("CredentialsAdapter");
            public static final Symbol CredentialsProviderAdapter = AwsGoDependency.INTERNAL_SIGV4A.pointableSymbol("CredentialsProviderAdapter");
            public static final Symbol SignerAdapter = AwsGoDependency.INTERNAL_SIGV4A.pointableSymbol("SignerAdapter");
        }
    }

    public static final class ServiceCustomizations {
        public static final class S3 {
            public static final Symbol SetSignerVersion = AwsCustomGoDependency.S3_CUSTOMIZATION.valueSymbol("SetSignerVersion");
            public static final Symbol ExpressIdentityResolver = AwsCustomGoDependency.S3_CUSTOMIZATION.valueSymbol("ExpressIdentityResolver");
            public static final Symbol ExpressSigner = AwsCustomGoDependency.S3_CUSTOMIZATION.valueSymbol("ExpressSigner");
            public static final Symbol GetPropertiesBackend = AwsCustomGoDependency.S3_CUSTOMIZATION.valueSymbol("GetPropertiesBackend");
            public static final Symbol AddExpressDefaultChecksumMiddleware = AwsCustomGoDependency.S3_CUSTOMIZATION.valueSymbol("AddExpressDefaultChecksumMiddleware");
        }

        public static final class S3Control {
            public static final Symbol AddDisableHostPrefixMiddleware = AwsCustomGoDependency.S3CONTROL_CUSTOMIZATION.valueSymbol("AddDisableHostPrefixMiddleware");
        }

        public static final class EventBridge {
            public static final Symbol SetSignerVersion = AwsCustomGoDependency.EVENTBRIDGE_CUSTOMIZATION.valueSymbol("SetSignerVersion");
        }
    }

    public static final class ServiceInternal {
        public static final class AcceptEncoding {
            public static final Symbol DisableGzip = AwsCustomGoDependency.ACCEPT_ENCODING_CUSTOMIZATION.pointableSymbol("DisableGzip");
        }
    }
}
