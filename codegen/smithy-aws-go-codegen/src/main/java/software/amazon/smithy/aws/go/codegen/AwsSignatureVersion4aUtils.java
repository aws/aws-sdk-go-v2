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

import software.amazon.smithy.aws.go.codegen.customization.AwsCustomGoDependency;
import software.amazon.smithy.go.codegen.GoDependency;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;

/**
 * Generates Client Configuration, Middleware, and Config Resolvers for AWS Signature Version 4a support.
 */
public final class AwsSignatureVersion4aUtils {
    public static final String RESOLVE_CREDENTIAL_PROVIDER = "resolveCredentialProvider";
    public static final String REGISTER_MIDDLEWARE_FUNCTION = "swapWithCustomHTTPSignerMiddleware";
    public static final String V4A_SIGNER_INTERFACE_NAME = "httpSignerV4a";
    public static final String SIGNER_OPTION_FIELD_NAME = V4A_SIGNER_INTERFACE_NAME;
    public static final String NEW_SIGNER_FUNC_NAME = "newDefaultV4aSigner";
    public static final String SIGNER_RESOLVER = "resolveHTTPSignerV4a";

    public static void writeCredentialProviderResolver(GoWriter writer) {
        writer.pushState();
        writer.putContext("resolverName", RESOLVE_CREDENTIAL_PROVIDER);
        writer.putContext("fieldName", AddAwsConfigFields.CREDENTIALS_CONFIG_NAME);
        writer.putContext("credType", SymbolUtils.createPointableSymbolBuilder("CredentialsProvider",
                AwsCustomGoDependency.INTERNAL_SIGV4A).build());
        writer.putContext("anonType", SymbolUtils.createPointableSymbolBuilder("AnonymousCredentials",
                AwsCustomGoDependency.AWS_CORE).build());
        writer.putContext("adapType", SymbolUtils.createPointableSymbolBuilder("SymmetricCredentialAdaptor",
                AwsCustomGoDependency.INTERNAL_SIGV4A).build());
        writer.write("""
                     func $resolverName:L(o *Options) {
                         if o.$fieldName:L == nil {
                             return
                         }
                         
                         if _, ok := o.$fieldName:L.($credType:T); ok {
                             return
                         }
                         
                         switch o.$fieldName:L.(type) {
                         case $anonType:T, $anonType:P:
                             return
                         }
                         
                         o.$fieldName:L = &$adapType:T{SymmetricProvider: o.$fieldName:L}
                     }
                     """);
        writer.popState();
    }

    public static void writerSignerInterface(GoWriter writer) {
        writer.pushState();
        writer.putContext("ifaceName", V4A_SIGNER_INTERFACE_NAME);
        writer.putContext("contextType", SymbolUtils.createValueSymbolBuilder("Context",
                SmithyGoDependency.CONTEXT).build());
        writer.putContext("credType", SymbolUtils.createValueSymbolBuilder("Credentials",
                AwsGoDependency.INTERNAL_SIGV4A).build());
        writer.putContext("reqType", SymbolUtils.createPointableSymbolBuilder("Request",
                SmithyGoDependency.NET_HTTP).build());
        writer.putContext("timeType", SymbolUtils.createPointableSymbolBuilder("Time",
                SmithyGoDependency.TIME).build());
        writer.putContext("optionsType", SymbolUtils.createPointableSymbolBuilder("SignerOptions",
                AwsGoDependency.INTERNAL_SIGV4A).build());
        writer.write("""
                     type $ifaceName:L interface {
                         SignHTTP(ctx $contextType:T, credentials $credType:T, r $reqType:P, payloadHash,
                             service string, regionSet []string, signingTime $timeType:T,
                             optFns ...func($optionsType:P)) error
                     }
                     """);
        writer.popState();
    }

    public static void writerConfigFieldResolver(GoWriter writer, ServiceShape serviceShape) {
        writer.pushState();
        writer.putContext("resolverName", SIGNER_RESOLVER);
        writer.putContext("optionName", SIGNER_OPTION_FIELD_NAME);
        writer.putContext("newSigner", NEW_SIGNER_FUNC_NAME);
        writer.write("""
                     func $resolverName:L(o *Options) {
                         if o.$optionName:L != nil {
                             return
                         }
                         o.$optionName:L = $newSigner:L(*o)
                     }
                     """);
        writer.popState();
    }

    public static void writeNewV4ASignerFunc(GoWriter writer, ServiceShape serviceShape) {
        writeNewV4ASignerFunc(writer, serviceShape, false);
    }

    public static void writeNewV4ASignerFunc(
            GoWriter writer,
            ServiceShape serviceShape,
            boolean disableURIPathEscaping
    ) {
        writer.pushState();
        writer.putContext("funcName", NEW_SIGNER_FUNC_NAME);
        writer.putContext("signerType", SymbolUtils.createPointableSymbolBuilder("Signer",
                AwsCustomGoDependency.INTERNAL_SIGV4A).build());
        writer.putContext("newSigner", SymbolUtils.createValueSymbolBuilder("NewSigner",
                AwsCustomGoDependency.INTERNAL_SIGV4A).build());
        writer.putContext("signerOptions", SymbolUtils.createPointableSymbolBuilder("SignerOptions",
                AwsCustomGoDependency.INTERNAL_SIGV4A).build());
        writer.putContext("loggerField", AddAwsConfigFields.LOGGER_CONFIG_NAME);
        writer.putContext("modeField", AddAwsConfigFields.LOG_MODE_CONFIG_NAME);
        writer.putContext("disableEscape", disableURIPathEscaping);
        writer.write("""
                     func $funcName:L(o Options) $signerType:P {
                         return $newSigner:T(func(so $signerOptions:P){
                             so.Logger = o.$loggerField:L
                             so.LogSigning = o.$modeField:L.IsSigning()
                             so.DisableURIPathEscaping = $disableEscape:L
                         })
                     }
                     """);
        writer.popState();
    }

    public static void writeMiddlewareRegister(
            Model model,
            GoWriter writer,
            ServiceShape serviceShape,
            GoDependency signerMiddleware
    ) {
        writer.pushState();
        writer.putContext("funcName", REGISTER_MIDDLEWARE_FUNCTION);
        writer.putContext("stackType", SymbolUtils.createPointableSymbolBuilder("Stack",
                SmithyGoDependency.SMITHY_MIDDLEWARE).build());
        writer.putContext("newMiddleware", SymbolUtils.createValueSymbolBuilder(
                "NewSignHTTPRequestMiddleware", signerMiddleware).build());
        writer.putContext("middleOptions", SymbolUtils.createValueSymbolBuilder(
                "SignHTTPRequestMiddlewareOptions", signerMiddleware).build());
        writer.putContext("registerMiddleware", SymbolUtils.createValueSymbolBuilder(
                "RegisterSigningMiddleware", signerMiddleware).build());
        writer.putContext("credFileName", AddAwsConfigFields.CREDENTIALS_CONFIG_NAME);
        writer.putContext("v4Signer", AwsSignatureVersion4.SIGNER_CONFIG_FIELD_NAME);
        writer.putContext("v4aSigner", SIGNER_OPTION_FIELD_NAME);
        writer.putContext("logMode", AddAwsConfigFields.LOG_MODE_CONFIG_NAME);
        writer.write("""
                     func $funcName:L(stack $stackType:P, o Options) error {
                         mw := $newMiddleware:T($middleOptions:T{
                             CredentialsProvider: o.$credFileName:L,
                             V4Signer: o.$v4Signer:L,
                             V4aSigner: o.$v4aSigner:L,
                             LogSigning: o.$logMode:L.IsSigning(),
                         })
                         
                         return $registerMiddleware:T(stack, mw)
                     }
                     """);
        writer.popState();
    }
}
