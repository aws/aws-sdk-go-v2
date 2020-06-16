package software.amazon.smithy.aws.go.codegen;

import java.util.List;
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.aws.traits.auth.UnsignedPayloadTrait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.traits.AuthTrait;
import software.amazon.smithy.utils.ListUtils;

public class AssembleMiddlewareStack implements GoIntegration{

    private static final String INITIALIZE_MIDDLEWARE = "Initialize";
    private static final String SERIALIZE_MIDDLEWARE = "Serialize";
    private static final String BUILD_MIDDLEWARE = "Build";
    private static final String FINALIZE_MIDDLEWARE = "Finalize";
    private static final String DESERIALIZE_MIDDLEWARE = "Deserialize";

    /**
     * Generates code to add middleware at the end in operation stack step.
     * @param writer writer used to write Go code.
     * @param stackstep stack step where the middleware is to be added.
     * @param middlewareSymbol middleware symbol corresponding to middleware to be added.
     * @param content Gowriter content used for generation.
     * @param stackOperand stack operand to which middleware is to be added.
     */
    private void writeAddMiddlewareAfter(
        GoWriter writer,
        String stackstep,
        Symbol middlewareSymbol,
        String content,
        String stackOperand
    ) {
        String st = String.format("%s.%s.Add(%s, middleware.After)", stackOperand, stackstep, content);
        writer.write(st, middlewareSymbol);
        writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);
    }

    /**
     * Generates code to add middleware at the end in operation stack step.
     * @param writer writer used to write Go code.
     * @param stackstep stack step where the middleware is to be added.
     * @param middlewareName middleware name corresponding to middleware to be added.
     * @param content Gowriter content used for generation.
     * @param stackOperand stack operand to which middleware is to be added.
     */
    private void writeAddMiddlewareAfter(
            GoWriter writer,
            String stackstep,
            String middlewareName,
            String content,
            String stackOperand
    ) {
        String st = String.format("%s.%s.Add(%s, middleware.After)", stackOperand, stackstep, content);
        writer.write(st, middlewareName);
        writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);
    }

    /**
     * Generates code to add middleware at the beginning in operation stack step.
     * @param writer writer used to write Go code.
     * @param stackstep stack step where the middleware is to be added.
     * @param middlewareSymbol middleware symbol corresponding to middleware to be added.
     * @param content Gowriter content used for generation.
     * @param stackOperand stack operand to which middleware is to be added.
     */
    private void writeAddMiddlewareBefore(
            GoWriter writer,
            String stackstep,
            Symbol middlewareSymbol,
            String content,
            String stackOperand
    ) {
        String st = String.format("%s.%s.Add(%s, middleware.Before)", stackOperand, stackstep, content);
        writer.write(st, middlewareSymbol);
        writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                // Add RequestInvocationIDMiddleware to operation stack
                RuntimeClientPlugin.builder()
                    .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                        // RequestInvocationIDMiddleware
                        Symbol RequestInvocationIDMiddleware = SymbolUtils.createValueSymbolBuilder(
                                "RequestInvocationIDMiddleware", AwsGoDependency.AWS_MIDDLEWARE
                        ).build();

                        writeAddMiddlewareAfter(writer, BUILD_MIDDLEWARE,
                                RequestInvocationIDMiddleware, "$T{}", stackOperand);
                    }).build(),

                // Add serializer middleware to operation stack
                RuntimeClientPlugin.builder()
                        .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                            if (protocolGenerator == null){
                                return;
                            }
                            String serializerMiddlewareName = ProtocolGenerator.getSerializeMiddlewareName(operation.getId(),
                                    protocolGenerator.getProtocolName());
                            writeAddMiddlewareAfter(writer, SERIALIZE_MIDDLEWARE,
                                    serializerMiddlewareName, "&$L{}", stackOperand);
                        }).build(),

                // Add deserializer middleware to operation stack
                RuntimeClientPlugin.builder()
                        .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                            if (protocolGenerator == null) {
                                return;
                            }
                            String deserializerMiddlewareName = ProtocolGenerator.getDeserializeMiddlewareName(operation.getId(),
                                    protocolGenerator.getProtocolName());
                            writeAddMiddlewareAfter(writer, DESERIALIZE_MIDDLEWARE,
                                    deserializerMiddlewareName, "&$L{}", stackOperand);
                        }).build(),

                // Add attemptClockSkew middleware to operation stack
                RuntimeClientPlugin.builder()
                        .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                            // attemptClockSkewMiddleware
                            Symbol attemptClockSkewMiddleware = SymbolUtils.createValueSymbolBuilder(
                                    "AttemptClockSkewMiddleware", AwsGoDependency.AWS_MIDDLEWARE
                            ).build();

                            writeAddMiddlewareAfter(writer, DESERIALIZE_MIDDLEWARE,
                                    attemptClockSkewMiddleware, "$T{}", stackOperand);
                        }).build(),

                // Add newAttempt middleware to operation stack
                RuntimeClientPlugin.builder()
                        .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                            // newAttemptMiddleware
                            Symbol newAttemptMiddleware = SymbolUtils.createValueSymbolBuilder(
                                    "NewAttemptMiddleware", AwsGoDependency.AWS_RETRY_MIDDLEWARE
                            ).build();

                            writeAddMiddlewareAfter(writer, FINALIZE_MIDDLEWARE,
                                    newAttemptMiddleware, "$T(options.Retryer, smithyhttp.RequestCloner)", stackOperand);
                        }).build(),

                // Add retry middleware to operation stack
                RuntimeClientPlugin.builder()
                        .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                            // metricHeaderMiddleware
                            Symbol metricHeaderMiddleware = SymbolUtils.createValueSymbolBuilder(
                                    "MetricsHeaderMiddleware", AwsGoDependency.AWS_RETRY_MIDDLEWARE
                            ).build();

                            writeAddMiddlewareAfter(writer, FINALIZE_MIDDLEWARE,
                                    metricHeaderMiddleware, "$T{}", stackOperand);
                        }).build(),

                // Add unsigned payload middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            if (operation.hasTrait(UnsignedPayloadTrait.class)) {
                                return true;
                            }
                            return false;
                        })
                        .buildMiddlewareStack((writer, service, operation, protocolGenerator, stackOperand) -> {
                            // unsignedPayloadSignerMiddleware
                            Symbol unsignedPayloadSignerMiddleware = SymbolUtils.createValueSymbolBuilder(
                                    "UnsignedPayloadMiddleware", AwsGoDependency.AWS_V4SIGNER_MIDDLEWARE
                            ).build();

                            writeAddMiddlewareAfter(writer, FINALIZE_MIDDLEWARE,
                                    unsignedPayloadSignerMiddleware, "$T{}", stackOperand);
                        }).build(),

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
                            // computePayloadSHA256Middleware
                            Symbol computePayloadSHA256Middleware = SymbolUtils.createValueSymbolBuilder(
                                    "ComputePayloadSHA256Middleware", AwsGoDependency.AWS_V4SIGNER_MIDDLEWARE
                            ).build();

                            writeAddMiddlewareBefore(writer, FINALIZE_MIDDLEWARE,
                                    computePayloadSHA256Middleware, "$T{}", stackOperand);

                            // newSignHttpRequestMiddleware
                            Symbol newSignHTTPRequestMiddleware = SymbolUtils.createValueSymbolBuilder(
                                    "NewSignHTTPRequestMiddleware", AwsGoDependency.AWS_V4SIGNER_MIDDLEWARE
                            ).build();

                            writeAddMiddlewareAfter(writer, FINALIZE_MIDDLEWARE,
                                    newSignHTTPRequestMiddleware, "&$T(options.Signer)", stackOperand);
                        }).build()
        );
    }
}
