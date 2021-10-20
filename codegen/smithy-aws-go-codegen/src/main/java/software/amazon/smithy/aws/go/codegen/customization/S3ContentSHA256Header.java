package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import java.util.Optional;
import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.aws.traits.auth.UnsignedPayloadTrait;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.StreamingTrait;
import software.amazon.smithy.utils.ListUtils;


public class S3ContentSHA256Header implements GoIntegration {

    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first.
     *
     * @return Returns the sort order, defaults to -40.
     */
    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        // Only add the middleware if UnsignedPayloadTrait is not specified, as this middleware
                        // will have already been added.
                        .operationPredicate((model, service, operation) -> {
                            if (!(S3ModelUtils.isServiceS3(model, service)
                                    || S3ModelUtils.isServiceS3Control(model, service))) {
                                return false;
                            }
                            return !operation.hasTrait(UnsignedPayloadTrait.class);
                        })
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddContentSHA256HeaderMiddleware",
                                        AwsGoDependency.AWS_SIGNER_V4
                                ).build())
                                .build())
                        .build(),
                RuntimeClientPlugin.builder()
                        // If a S3 operation has a streaming payload but is not event stream payload,
                        // client swaps signing middleware to use dynamic payload signing middleware.
                        // This enables client to use unsigned payload when TLS is enabled, and switch
                        // to signed payload for security when TLS is disabled.
                        .operationPredicate(((model, service, operation) -> {
                            if (!(S3ModelUtils.isServiceS3(model, service))) {
                                return false;
                            }

                            Optional<ShapeId> input = operation.getInput();
                            if (!input.isPresent()) {
                                return false;
                            }

                            StructureShape inputShape = model.expectShape(input.get(), StructureShape.class);
                            for (MemberShape memberShape : inputShape.getAllMembers().values()) {
                                Shape targetShape = model.expectShape(memberShape.getTarget());
                                if (targetShape.hasTrait(StreamingTrait.class)
                                        && !StreamingTrait.isEventStream(model, memberShape)) {
                                    return true;
                                }
                            }
                            return false;
                        }))
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "UseDynamicPayloadSigningMiddleware", AwsGoDependency.AWS_SIGNER_V4
                                ).build())
                                .build())
                        .build()
        );
    }
}
