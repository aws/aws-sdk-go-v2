package software.amazon.smithy.aws.go.codegen;

import java.util.List;
import java.util.logging.Logger;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

/**
 * Registers additional client specific configuration fields
 */
public class ResolveClientConfig implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(AddAwsConfigFields.class.getName());

    private static final String USE_ARN_REGION_OPTION = "UseARNRegion";
    private static final String CLIENT_CONFIG_RESOLVER_FUNC = "resolveClientConfig";
    private static final String CONFIG_SOURCE_CONFIG_NAME = "ConfigSources";
    private static final String RESOLVE_USE_ARN_REGION= "ResolveUseARNRegion";

    public static final List<AddAwsConfigFields.AwsConfigField> AWS_CONFIG_FIELDS = ListUtils.of(
            AddAwsConfigFields.AwsConfigField.builder()
                    .name(USE_ARN_REGION_OPTION)
                    .type(getUniversalSymbol("boolean"))
                    .generatedOnClient(false)
                    .servicePredicate((model, serviceShape) -> {
                        return isS3SharedService(model, serviceShape);
                    })
                    .awsResolveFunction(SymbolUtils.createValueSymbolBuilder(CLIENT_CONFIG_RESOLVER_FUNC)
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
            if (!isS3SharedService(model, serviceShape)) {
                return;
            }

            writer.writeDocs("resolves client config");
            writer.addUseImports(AwsGoDependency.AWS_CORE);
            writer.openBlock("func $L(cfg aws.Config, o *Options) error {", "}",
                    CLIENT_CONFIG_RESOLVER_FUNC, () -> {
                writer.openBlock("if len(cfg.$L) == 0 {", "}",
                        CONFIG_SOURCE_CONFIG_NAME,
                        () -> writer.write("return nil")
                );

                writer.addUseImports(SmithyGoDependency.CONTEXT);
                Symbol resolverFunc = SymbolUtils.createValueSymbolBuilder(RESOLVE_USE_ARN_REGION,
                        AwsGoDependency.S3_SHARED_CONFIG).build();
                writer.write("value, found, err := $T(context.Background(), cfg.$L)", resolverFunc,
                        CONFIG_SOURCE_CONFIG_NAME);
                writer.write("if err != nil { return err }");
                writer.write("if found { o.$L = value }", USE_ARN_REGION_OPTION);

                writer.write("return nil");
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

}
