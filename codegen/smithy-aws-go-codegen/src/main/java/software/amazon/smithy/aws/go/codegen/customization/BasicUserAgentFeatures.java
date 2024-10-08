package software.amazon.smithy.aws.go.codegen.customization;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

import java.util.List;
import java.util.Map;
import java.util.function.BiPredicate;
import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.GoCodegenContext;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoTypes;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.protocol.traits.Rpcv2CborTrait;

/**
 * Adds user agent tracking for basic features - i.e. simple model-based ones that do not require any additional in-code
 * checks, such as a particular protocol.
 */
public class BasicUserAgentFeatures implements GoIntegration {
    private static final List<Feature> FEATURES = List.of(
           new Feature("ProtocolRPCV2CBOR", (model, service) -> service.hasTrait(Rpcv2CborTrait.class))
    );

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return FEATURES.stream().map(Feature::getPlugin).toList();
    }

    @Override
    public void writeAdditionalFiles(GoCodegenContext ctx) {
        var model = ctx.model();
        var service = ctx.settings().getService(model);
        ctx.writerDelegator().useFileWriter("api_client.go", ctx.settings().getModuleName(),
                GoWriter.ChainWritable.of(
                        FEATURES.stream()
                                .filter(it -> it.servicePredicate.test(model, service))
                                .map(Feature::getAddMiddleware)
                                .toList()
                ).compose());
    }

    private static final class Feature {
        public final Symbol featureId;
        public final BiPredicate<Model, ServiceShape> servicePredicate;

        public Feature(String id, BiPredicate<Model, ServiceShape> servicePredicate) {
            this.featureId = AwsGoDependency.AWS_MIDDLEWARE.constSymbol("UserAgentFeature" + id);
            this.servicePredicate = servicePredicate;
        }

        public RuntimeClientPlugin getPlugin() {
            return RuntimeClientPlugin.builder()
                    .servicePredicate(servicePredicate)
                    .registerMiddleware(
                            MiddlewareRegistrar.builder()
                                    .resolvedFunction(buildPackageSymbol("add" + featureId.getName()))
                                    .useClientOptions()
                                    .build()
                    )
                    .build();
        }

        public GoWriter.Writable getAddMiddleware() {
            return goTemplate("""
                    func add$featureName:L(stack $stack:P, options Options) error {
                        ua, err := getOrAddRequestUserAgent(stack)
                        if err != nil {
                            return err
                        }

                        ua.AddUserAgentFeature($featureEnum:T)
                        return nil
                    }
                    """,
                    Map.of(
                            "stack", SmithyGoTypes.Middleware.Stack,
                            "featureName", featureId.getName(),
                            "featureEnum", featureId
                    ));
        }
    }
}
