package software.amazon.smithy.aws.go.codegen;

import java.util.List;
import java.util.logging.Logger;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.aws.traits.clientendpointdiscovery.ClientEndpointDiscoveryTrait;
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
import software.amazon.smithy.utils.ListUtils;

/**
 * Registers additional client specific configuration fields
 * TODO: This needs to refactored so that we aren't defining "pseudo-config fields"
 */
public class ResolveClientConfigFromSources implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(AddAwsConfigFields.class.getName());

    private static final String CONFIG_SOURCE_CONFIG_NAME = "ConfigSources";

    private static final String USE_ARN_REGION_OPTION = "UseARNRegion";
    private static final String USE_ARN_REGION_CONFIG_RESOLVER = "resolveUseARNRegion";
    private static final String RESOLVE_USE_ARN_REGION = "ResolveUseARNRegion";

    private static final String DISABLE_MRAP_OPTION = "DisableMultiRegionAccessPoints";
    private static final String DISABLE_MRAP_CONFIG_RESOLVER = "resolveDisableMultiRegionAccessPoints";
    private static final String RESOLVE_DISABLE_MRAP = "ResolveDisableMultiRegionAccessPoints";

    // EndpointDiscovery options
    private static final String ENDPOINT_DISCOVERY_OPTION = "EndpointDiscovery";
    private static final Symbol ENDPOINT_DISCOVERY_OPTION_TYPE = SymbolUtils.createValueSymbolBuilder(
            "EndpointDiscoveryOptions").build();

    // Enable EndpointDiscovery
    private static final String ENABLE_ENDPOINT_DISCOVERY_OPTION = "EnableEndpointDiscovery";
    private static final String ENABLE_ENDPOINT_DISCOVERY_CONFIG_RESOLVER = "resolveEnableEndpointDiscoveryFromConfigSources";
    private static final String RESOLVE_ENABLE_ENDPOINT_DISCOVERY = "ResolveEnableEndpointDiscovery";

    // UseDualStack
    private static final String DUAL_STACK_ENDPOINT_CONFIG_RESOLVER = "resolveUseDualStackEndpoint";
    private static final String RESOLVE_USE_DUAL_STACK_ENDPOINT = "ResolveUseDualStackEndpoint";
    private static final String USE_FIPS_ENDPOINT_CONFIG_RESOLVER = "resolveUseFIPSEndpoint";
    private static final String RESOLVE_USE_FIPS_ENDPOINT = "ResolveUseFIPSEndpoint";

    public static final List<AddAwsConfigFields.AwsConfigField> AWS_CONFIG_FIELDS = ListUtils.of(
            AddAwsConfigFields.AwsConfigField.builder()
                    .name(USE_ARN_REGION_OPTION)
                    .type(getUniversalSymbol("bool"))
                    .generatedOnClient(false)
                    .servicePredicate(ResolveClientConfigFromSources::isS3SharedService)
                    .awsResolveFunction(SymbolUtils.createValueSymbolBuilder(USE_ARN_REGION_CONFIG_RESOLVER)
                            .build())
                    .build(),
        AddAwsConfigFields.AwsConfigField.builder()
                    .name(DISABLE_MRAP_OPTION)
                    .type(getUniversalSymbol("bool"))
                    .generatedOnClient(false)
                    .servicePredicate(ResolveClientConfigFromSources::isS3Service)
                    .awsResolveFunction(SymbolUtils.createValueSymbolBuilder(DISABLE_MRAP_CONFIG_RESOLVER)
                            .build())
                    .build(),
            AddAwsConfigFields.AwsConfigField.builder()
                    .name("DisableExpressAuth")
                    .type(getUniversalSymbol("*bool"))
                    .generatedOnClient(false)
                    .servicePredicate(ResolveClientConfigFromSources::isS3Service)
                    .awsResolveFunction(SymbolUtils.createValueSymbolBuilder("resolveDisableExpressAuth")
                            .build())
                    .build(),
            AddAwsConfigFields.AwsConfigField.builder()
                    .name(ENDPOINT_DISCOVERY_OPTION)
                    .type(ENDPOINT_DISCOVERY_OPTION_TYPE)
                    .generatedOnClient(false)
                    .servicePredicate(ResolveClientConfigFromSources::supportsEndpointDiscovery)
                    .awsResolveFunction(SymbolUtils.createValueSymbolBuilder(ENABLE_ENDPOINT_DISCOVERY_CONFIG_RESOLVER)
                            .build())
                    .build(),
            AddAwsConfigFields.AwsConfigField.builder()
                    .name("EndpointOptions.UseDualStackEndpoint")
                    .type(SymbolUtils.createPointableSymbolBuilder("bool").build())
                    .generatedOnClient(false)
                    .awsResolveFunction(SymbolUtils.createValueSymbolBuilder(DUAL_STACK_ENDPOINT_CONFIG_RESOLVER)
                            .build())
                    .build(),
            AddAwsConfigFields.AwsConfigField.builder()
                    .name("EndpointOptions.UseFIPSEndpoint")
                    .type(SymbolUtils.createPointableSymbolBuilder("bool").build())
                    .generatedOnClient(false)
                    .awsResolveFunction(SymbolUtils.createValueSymbolBuilder(USE_FIPS_ENDPOINT_CONFIG_RESOLVER)
                            .build())
                    .build(),
            AddAwsConfigFields.AwsConfigField.builder()
                    .name("BaseEndpoint")
                    .type(SymbolUtils.createPointableSymbolBuilder("string").build())
                    .generatedOnClient(false)
                    .awsResolveFunction(SymbolUtils.createValueSymbolBuilder(
                        AwsEndpointResolverInitializerGenerator.RESOLVE_BASE_ENDPOINT)
                            .build())
                    .build()
    );

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        LOGGER.info("generating client config resolver");
        ServiceShape serviceShape = settings.getService(model);
        goDelegator.useShapeWriter(serviceShape, writer -> {
            generateUseARNRegionResolver(model, serviceShape, writer);
            generateDisableMrapResolver(model, serviceShape, writer);
            generateEnableEndpointDiscoveryResolver(model, serviceShape, writer);
            generateUseUseDualStackResolver(model, serviceShape, writer);
            generateUseUseFIPSEndpointResolver(model, serviceShape, writer);
        });
    }

    private static void generatedResolverFunction(GoWriter writer, String name, String documentation, Runnable f) {
        if (documentation.length() > 0) {
            writer.writeDocs(documentation);
        }
        writer.addUseImports(AwsGoDependency.AWS_CORE);
        writer.openBlock("func $L(cfg aws.Config, o *Options) error {", "}", name, () -> {
            writer.openBlock("if len(cfg.$L) == 0 {", "}",
                    CONFIG_SOURCE_CONFIG_NAME,
                    () -> writer.write("return nil")
            );

            f.run();

            writer.write("return nil");
        });
    }

    private static void generateUseARNRegionResolver(Model model, ServiceShape serviceShape, GoWriter writer) {
        if (!isS3SharedService(model, serviceShape)) {
            return;
        }
        generatedResolverFunction(writer, USE_ARN_REGION_CONFIG_RESOLVER, "resolves UseARNRegion S3 configuration", () -> {
            writer.addUseImports(SmithyGoDependency.CONTEXT);
            Symbol resolverFunc = SymbolUtils.createValueSymbolBuilder(RESOLVE_USE_ARN_REGION,
                    AwsGoDependency.S3_SHARED_CONFIG).build();
            writer.write("value, found, err := $T(context.Background(), cfg.$L)", resolverFunc,
                    CONFIG_SOURCE_CONFIG_NAME);
            writer.write("if err != nil { return err }");
            writer.write("if found { o.$L = value }", USE_ARN_REGION_OPTION);
        });
        writer.write("");
    }

    private static void generateDisableMrapResolver(Model model, ServiceShape serviceShape, GoWriter writer) {
        if (!isS3Service(model, serviceShape)) {
            return;
        }
        generatedResolverFunction(writer, DISABLE_MRAP_CONFIG_RESOLVER, "resolves DisableMultiRegionAccessPoints S3 configuration", () -> {
            writer.addUseImports(SmithyGoDependency.CONTEXT);
            Symbol resolverFunc = SymbolUtils.createValueSymbolBuilder(RESOLVE_DISABLE_MRAP,
                    AwsGoDependency.S3_SHARED_CONFIG).build();
            writer.write("value, found, err := $T(context.Background(), cfg.$L)", resolverFunc,
                    CONFIG_SOURCE_CONFIG_NAME);
            writer.write("if err != nil { return err }");
            writer.write("if found { o.$L = value }", DISABLE_MRAP_OPTION);
        });
        writer.write("");
    }

    private static void generateEnableEndpointDiscoveryResolver(
            Model model,
            ServiceShape serviceShape,
            GoWriter writer
    ) {
        if (!supportsEndpointDiscovery(model, serviceShape)) {
            return;
        }

        generatedResolverFunction(writer, ENABLE_ENDPOINT_DISCOVERY_CONFIG_RESOLVER,
                "resolves EnableEndpointDiscovery configuration", () -> {
                    writer.addUseImports(SmithyGoDependency.CONTEXT);
                    Symbol resolverFunc = SymbolUtils.createValueSymbolBuilder(RESOLVE_ENABLE_ENDPOINT_DISCOVERY,
                            AwsGoDependency.SERVICE_INTERNAL_CONFIG).build();
                    writer.write("value, found, err := $T(context.Background(), cfg.$L)", resolverFunc,
                            CONFIG_SOURCE_CONFIG_NAME);
                    writer.write("if err != nil { return err }");
                    writer.write("if found { o.$L.$L = value }", ENDPOINT_DISCOVERY_OPTION, ENABLE_ENDPOINT_DISCOVERY_OPTION);
                });
        writer.write("");
    }

    private void generateUseUseDualStackResolver(Model model, ServiceShape serviceShape, GoWriter writer) {
        writer.addUseImports(AwsGoDependency.AWS_CORE);

        generatedResolverFunction(writer, DUAL_STACK_ENDPOINT_CONFIG_RESOLVER,
                "resolves dual-stack endpoint configuration", () -> {
                    writer.addUseImports(SmithyGoDependency.CONTEXT);
                    var resolverFunc = SymbolUtils.createValueSymbolBuilder(RESOLVE_USE_DUAL_STACK_ENDPOINT,
                            AwsGoDependency.SERVICE_INTERNAL_CONFIG).build();
                    writer.write("value, found, err := $T(context.Background(), cfg.$L)", resolverFunc,
                            CONFIG_SOURCE_CONFIG_NAME);
                    writer.write("if err != nil { return err }");

                    writer.openBlock("if found {", "}", () -> writer
                            .write("o.EndpointOptions.$L = value",
                                    EndpointGenerator.DUAL_STACK_ENDPOINT_OPTION));
                });
    }

    private void generateUseUseFIPSEndpointResolver(Model model, ServiceShape serviceShape, GoWriter writer) {
        writer.addUseImports(AwsGoDependency.AWS_CORE);

        generatedResolverFunction(writer, USE_FIPS_ENDPOINT_CONFIG_RESOLVER,
                "resolves FIPS endpoint configuration", () -> {
                    writer.addUseImports(SmithyGoDependency.CONTEXT);
                    var resolverFunc = SymbolUtils.createValueSymbolBuilder(RESOLVE_USE_FIPS_ENDPOINT,
                            AwsGoDependency.SERVICE_INTERNAL_CONFIG).build();
                    writer.write("value, found, err := $T(context.Background(), cfg.$L)", resolverFunc,
                            CONFIG_SOURCE_CONFIG_NAME);
                    writer.write("if err != nil { return err }");

                    writer.openBlock("if found {", "}", () -> {
                        writer.write("o.EndpointOptions.$L = value", EndpointGenerator.USE_FIPS_ENDPOINT_OPTION);
                    });
                });
    }

    private static Symbol getUniversalSymbol(String symbolName) {
        return SymbolUtils.createValueSymbolBuilder(symbolName)
                .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true).build();
    }

    private static boolean isS3SharedService(Model model, ServiceShape service) {
        return isS3Service(model, service) || isS3ControlService(model, service);
    }

    private static boolean isS3Service(Model model, ServiceShape service) {
        return service.expectTrait(ServiceTrait.class).getSdkId().equalsIgnoreCase("S3");
    }

    private static boolean isS3ControlService(Model model, ServiceShape service) {
        return service.expectTrait(ServiceTrait.class).getSdkId().equalsIgnoreCase("S3 Control");
    }

    private static boolean supportsEndpointDiscovery(Model model, ServiceShape service) {
        return service.hasTrait(ClientEndpointDiscoveryTrait.class);
    }
}
