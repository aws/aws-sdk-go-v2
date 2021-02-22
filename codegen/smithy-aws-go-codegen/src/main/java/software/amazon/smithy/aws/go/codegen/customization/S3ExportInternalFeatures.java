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

import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;

/**
 * Exports internal functionality from the s3shared package.
 */
public class S3ExportInternalFeatures implements GoIntegration {
    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator
    ) {
        ServiceShape service = settings.getService(model);
        if (!requiresCustomization(model, service)) {
            return;
        }

        goDelegator.useShapeWriter(service, writer -> {
            writeResponseErrorInterface(writer);
            writeGetHostIDWrapper(writer);
        });
    }

    private void writeGetHostIDWrapper(GoWriter writer) {
        Symbol metadata = SymbolUtils.createPointableSymbolBuilder("Metadata",
                SmithyGoDependency.SMITHY_MIDDLEWARE).build();
        Symbol getHostID = SymbolUtils.createValueSymbolBuilder("GetHostIDMetadata",
                AwsCustomGoDependency.S3_SHARED_CUSTOMIZATION).build();
        writer.writeDocs("GetHostIDMetadata retrieves the host id from middleware metadata "
                + "returns host id as string along with a boolean indicating presence of "
                + "hostId on middleware metadata.");
        writer.openBlock("func GetHostIDMetadata(metadata $T) (string, bool) {", "}", metadata, () -> {
            writer.write("return $T(metadata)", getHostID);
        });
    }

    private void writeResponseErrorInterface(GoWriter writer) {
        writer.writeDocs("ResponseError provides the HTTP centric error type wrapping the underlying error "
                + "with the HTTP response value and the deserialized RequestID.");
        writer.openBlock("type ResponseError interface {", "}", () -> {
            writer.write("error").write("");
            writer.write("ServiceHostID() string");
            writer.write("ServiceRequestID() string");
        }).write("");
        writer.write("var _ ResponseError = ($P)(nil)", SymbolUtils.createPointableSymbolBuilder("ResponseError",
                AwsCustomGoDependency.S3_SHARED_CUSTOMIZATION).build());
    }

    // returns true if service is either s3 or s3 control
    private static boolean requiresCustomization(Model model, ServiceShape service) {
        return S3ModelUtils.isServiceS3(model, service) || S3ModelUtils.isServiceS3Control(model, service);
    }
}
