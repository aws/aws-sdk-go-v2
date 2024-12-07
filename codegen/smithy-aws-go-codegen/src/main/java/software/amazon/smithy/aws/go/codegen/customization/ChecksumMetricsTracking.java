package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.aws.traits.HttpChecksumTrait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;

import java.util.List;
import java.util.Map;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

public class ChecksumMetricsTracking implements GoIntegration {
    private static final MiddlewareRegistrar RequestMIDDLEWARE = MiddlewareRegistrar.builder()
            .resolvedFunction(buildPackageSymbol("addRequestChecksumMetricsTracking"))
            .useClientOptions()
            .build();
    private static final MiddlewareRegistrar ResponseMIDDLEWARE = MiddlewareRegistrar.builder()
            .resolvedFunction(buildPackageSymbol("addResponseChecksumMetricsTracking"))
            .useClientOptions()
            .build();

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return List.of(
                RuntimeClientPlugin.builder()
                        .operationPredicate((m, s, o) -> {
                            if (!hasInputChecksumTrait(m, s, o)) {
                                return false;
                            }
                        return true;
                        })
                        .registerMiddleware(RequestMIDDLEWARE)
                        .build(),
                RuntimeClientPlugin.builder()
                        .operationPredicate((m, s, o) -> {
                            if (!hasOutputChecksumTrait(m, s, o)) {
                                return false;
                            }
                            return true;
                        })
                        .registerMiddleware(ResponseMIDDLEWARE)
                        .build()
        );
    }

    private static boolean hasChecksumTrait(Model model, ServiceShape service, OperationShape operation) {
                            return operation.hasTrait(HttpChecksumTrait.class);
    }

    private static boolean hasInputChecksumTrait(Model model, ServiceShape service, OperationShape operation) {
        if (!hasChecksumTrait(model, service, operation)) {
            return false;
        }
        HttpChecksumTrait trait = operation.expectTrait(HttpChecksumTrait.class);
        return trait.isRequestChecksumRequired() || trait.getRequestAlgorithmMember().isPresent();
    }

    private static boolean hasOutputChecksumTrait(Model model, ServiceShape service, OperationShape operation) {
        if (!hasChecksumTrait(model, service, operation)) {
            return false;
        }
        HttpChecksumTrait trait = operation.expectTrait(HttpChecksumTrait.class);
        return trait.getRequestValidationModeMember().isPresent() && !trait.getResponseAlgorithms().isEmpty();
    }

    @Override
    public void writeAdditionalFiles(GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator) {
        ServiceShape service = settings.getService(model);
        boolean supportsComputeInputChecksumsWorkflow = false;
        boolean supportsChecksumValidationWorkflow = false;

        for (OperationShape operation : TopDownIndex.of(model).getContainedOperations(service)) {
            if (!hasChecksumTrait(model, service, operation)) {
                continue;
            }

            if (hasInputChecksumTrait(model, service, operation)) {
                supportsComputeInputChecksumsWorkflow = true;
            }

            if (hasOutputChecksumTrait(model, service, operation)) {
                supportsChecksumValidationWorkflow = true;
            }
        }

        if (supportsComputeInputChecksumsWorkflow) {
            goDelegator.useFileWriter("api_client.go", settings.getModuleName(), goTemplate("""
                func addRequestChecksumMetricsTracking(stack $stack:P, options Options) error {
                    ua, err := getOrAddRequestUserAgent(stack)
                    if err != nil {
                        return err
                    }

                    return stack.Build.Insert(&$requestMetricsTracking:P{
                        RequestChecksumCalculation: options.RequestChecksumCalculation,
                        UserAgent: ua,
                    }, "UserAgent", $before:T)
                }""",
                    Map.of(
                            "stack", SmithyGoDependency.SMITHY_MIDDLEWARE.struct("Stack"),
                            "requestMetricsTracking", AwsGoDependency.SERVICE_INTERNAL_CHECKSUM.valueSymbol("RequestChecksumMetricsTracking"),
                            "before", SmithyGoDependency.SMITHY_MIDDLEWARE
                                    .valueSymbol("Before")
                    )));
        }

        if (supportsChecksumValidationWorkflow) {
            goDelegator.useFileWriter("api_client.go", settings.getModuleName(), goTemplate("""
                func addResponseChecksumMetricsTracking(stack $stack:P, options Options) error {
                    ua, err := getOrAddRequestUserAgent(stack)
                    if err != nil {
                        return err
                    }

                    return stack.Build.Insert(&$responseMetricsTracking:P{
                        ResponseChecksumValidation: options.ResponseChecksumValidation,
                        UserAgent: ua,
                    }, "UserAgent", $before:T)
                }""",
                    Map.of(
                            "stack", SmithyGoDependency.SMITHY_MIDDLEWARE.struct("Stack"),
                            "responseMetricsTracking", AwsGoDependency.SERVICE_INTERNAL_CHECKSUM.valueSymbol("ResponseChecksumMetricsTracking"),
                            "before", SmithyGoDependency.SMITHY_MIDDLEWARE
                                    .valueSymbol("Before")
                    )));
        }
    }
}
