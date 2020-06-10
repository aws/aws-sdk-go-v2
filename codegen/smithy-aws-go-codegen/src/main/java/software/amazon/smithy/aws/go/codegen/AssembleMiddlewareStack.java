package software.amazon.smithy.aws.go.codegen;

import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.aws.traits.auth.UnsignedPayloadTrait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDependency;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;

public class AssembleMiddlewareStack implements GoIntegration {
    @Override
    public void assembleMiddlewareStack(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoWriter writer,
            OperationShape operationShape
    ){
        // build middleware
        Symbol requestInvocationIDMiddleware = SymbolUtils.createValueSymbolBuilder(
             "RequestInvocationIDMiddleware", GoDependency.AWS_MIDDLEWARE).build();
        writer.write("stack.Build.Add($T{}, middleware.After)", requestInvocationIDMiddleware);

        // deserialize middleware
        Symbol attemptClockSkewMiddleware = SymbolUtils.createValueSymbolBuilder(
                "AttemptClockSkewMiddleware", GoDependency.AWS_MIDDLEWARE).build();
        writer.write("stack.Deserialize.Add($T{}, middleware.After)", attemptClockSkewMiddleware);

        // retry middleware
        Symbol newAttemptMiddleware = SymbolUtils.createValueSymbolBuilder(
                "NewAttemptMiddleware", GoDependency.AWS_RETRY_MIDDLEWARE).build();
        writer.write("stack.Finalize.Add($T(options.Retryer), middleware.After)", newAttemptMiddleware);

        // retry metric middleware
        Symbol metricsHeaderMiddleware = SymbolUtils.createValueSymbolBuilder(
                "MetricsHeaderMiddleware", GoDependency.AWS_RETRY_MIDDLEWARE).build();
        writer.write("stack.Finalize.Add($T{}, middleware.After)", metricsHeaderMiddleware);

        // signer middleware
        if (operationShape.hasTrait(UnsignedPayloadTrait.class)) {
            // unsigned payload middleware
            Symbol unsignedPayloadSignerMiddleware = SymbolUtils.createValueSymbolBuilder(
                    "UnsignedPayloadMiddleware", GoDependency.AWS_V4SIGNER_MIDDLEWARE).build();
            writer.write("stack.Finalize.Add(&$T{}, middleware.After)", unsignedPayloadSignerMiddleware);
        } else if (operationShape.hasTrait(SigV4Trait.class)) {
            // sigV4 signer middleware
            Symbol computePayloadSHA256Middleware = SymbolUtils.createValueSymbolBuilder(
                    "ComputePayloadSHA256Middleware", GoDependency.AWS_V4SIGNER_MIDDLEWARE).build();
            writer.write("stack.Finalize.Add(&$T{}, middleware.Before)", computePayloadSHA256Middleware);

            Symbol newSignHTTPRequestMiddleware = SymbolUtils.createValueSymbolBuilder(
                    "NewSignHTTPRequestMiddleware", GoDependency.AWS_V4SIGNER_MIDDLEWARE).build();
            writer.write("stack.Finalize.Add(&$T(options.Signer), middleware.After)",
                    newSignHTTPRequestMiddleware);
        } else {
            // v2 signer middleware
            writer.write("// TODO: Which middleware to add in case it's not a sigV4 supported service?");
        }
    }
}
