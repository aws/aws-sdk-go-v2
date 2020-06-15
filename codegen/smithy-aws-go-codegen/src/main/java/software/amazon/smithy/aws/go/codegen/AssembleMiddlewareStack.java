package software.amazon.smithy.aws.go.codegen;

import java.util.List;
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.aws.traits.auth.UnsignedPayloadTrait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.traits.AuthTrait;
import software.amazon.smithy.utils.ListUtils;

public class AssembleMiddlewareStack implements GoIntegration{
    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        Symbol middlewareStepAfter = SymbolUtils.createValueSymbolBuilder(
                "After", SmithyGoDependency.SMITHY_MIDDLEWARE
        ).build();

        Symbol middlewareStepBefore = SymbolUtils.createValueSymbolBuilder(
                "Before", SmithyGoDependency.SMITHY_MIDDLEWARE
        ).build();

        Symbol RequestInvocationIDMiddleware = SymbolUtils.createValueSymbolBuilder(
                "RequestInvocationIDMiddleware", AwsGoDependency.AWS_MIDDLEWARE
        ).build();

        Symbol attemptClockSkewMiddleware = SymbolUtils.createValueSymbolBuilder(
                "AttemptClockSkewMiddleware", AwsGoDependency.AWS_MIDDLEWARE
        ).build();

        Symbol newAttemptMiddleware = SymbolUtils.createValueSymbolBuilder(
                "NewAttemptMiddleware", AwsGoDependency.AWS_RETRY_MIDDLEWARE
        ).build();

        Symbol metricHeaderMiddleware = SymbolUtils.createValueSymbolBuilder(
                "MetricsHeaderMiddleware", AwsGoDependency.AWS_RETRY_MIDDLEWARE
        ).build();

        Symbol unsignedPayloadSignerMiddleware = SymbolUtils.createValueSymbolBuilder(
                "UnsignedPayloadMiddleware", AwsGoDependency.AWS_V4SIGNER_MIDDLEWARE
        ).build();

        Symbol computePayloadSHA256Middleware = SymbolUtils.createValueSymbolBuilder(
                "ComputePayloadSHA256Middleware", AwsGoDependency.AWS_V4SIGNER_MIDDLEWARE
        ).build();

        Symbol newSignHTTPRequestMiddleware = SymbolUtils.createValueSymbolBuilder(
                "NewSignHTTPRequestMiddleware", AwsGoDependency.AWS_V4SIGNER_MIDDLEWARE
        ).build();

        return ListUtils.of(
                // Add RequestInvocationIDMiddleware to operation stack
                RuntimeClientPlugin.builder()
                    .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                        writer.write("$L.Initialize.Add($T{}, $T)", stackOperand,
                                RequestInvocationIDMiddleware, middlewareStepAfter);
                    }).build(),

                // Add serializer middleware to operation stack
                RuntimeClientPlugin.builder()
                        .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                            if (protocolGenerator == null){
                                return;
                            }
                            String serializerMiddlewareName = ProtocolGenerator.getSerializeMiddlewareName(operation.getId(),
                                    protocolGenerator.getProtocolName());
                            writer.write("$L.Serialize.Add(&$L{}, $T)",
                                    stackOperand, serializerMiddlewareName, middlewareStepAfter);
                        })
                        .build(),

                // Add deserializer middleware to operation stack
                RuntimeClientPlugin.builder()
                        .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                            if (protocolGenerator == null) {
                                return;
                            }
                            String deserializerMiddlewareName = ProtocolGenerator.getDeserializeMiddlewareName(operation.getId(),
                                    protocolGenerator.getProtocolName());
                            writer.write("$L.Deserialize.Add(&$L{}, $T)",
                                    stackOperand, deserializerMiddlewareName, middlewareStepAfter);
                        })
                        .build(),

                // Add attemptClockSkew middleware to operation stack
                RuntimeClientPlugin.builder()
                        .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                            writer.write("$L.Deserialize.Add($T{}, $T)", stackOperand,
                                    attemptClockSkewMiddleware, middlewareStepAfter);
                        })
                        .build(),

                // Add newAttempt middleware to operation stack
                RuntimeClientPlugin.builder()
                        .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                            writer.write("$L.Finalize.Add($T(options.Retryer), $T)", stackOperand,
                                    newAttemptMiddleware, middlewareStepAfter);
                        })
                        .build(),

                // Add retry middleware to operation stack
                RuntimeClientPlugin.builder()
                        .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                            writer.write("$L.Finalize.Add($T{}, $T)", stackOperand,
                                    metricHeaderMiddleware, middlewareStepAfter);
                        })
                        .build(),

                // Add unsigned payload middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            if (operation.hasTrait(UnsignedPayloadTrait.class)) {
                                return true;
                            }
                            return false;
                        })
                        .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                            writer.write("$L.Finalize.Add($T{}, $T)", stackOperand,
                                    unsignedPayloadSignerMiddleware, middlewareStepAfter);
                        })
                        .build(),

                // Add SigV4 middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            if (service.hasTrait(SigV4Trait.class) && (!operation.hasTrait(UnsignedPayloadTrait.class))
                                   && (operation.hasTrait(SigV4Trait.class) || !operation.hasTrait(AuthTrait.class))
                            ){
                                return true;
                            }
                            return false;
                        })
                        .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                            writer.write("$L.Finalize.Add($T{}, $T)", stackOperand,
                                    computePayloadSHA256Middleware, middlewareStepBefore);
                            writer.write("stack.Finalize.Add(&$T(options.Signer), middleware.After)",
                                    newSignHTTPRequestMiddleware);
                        })
                        .build()
        );
    }
}
