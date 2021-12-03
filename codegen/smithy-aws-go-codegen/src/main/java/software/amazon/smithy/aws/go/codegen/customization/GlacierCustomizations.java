package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import java.util.stream.Collectors;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.utils.ListUtils;

public class GlacierCustomizations implements GoIntegration {
    private static final String TREE_HASH_ADDER = "AddTreeHashMiddleware";
    private static final String API_VERSION_ADDER = "AddGlacierAPIVersionMiddleware";
    private static final String ACCOUNT_ID_ADDER = "AddDefaultAccountIDMiddleware";
    private static final String SET_DEFAULT_ACCOUNT_ID = "setDefaultAccountID";

    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(GlacierCustomizations::isGlacier)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(TREE_HASH_ADDER,
                                        AwsCustomGoDependency.GLACIER_CUSTOMIZATION).build())
                                .build())
                        .build(),
                RuntimeClientPlugin.builder()
                        .servicePredicate(GlacierCustomizations::isGlacier)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(API_VERSION_ADDER,
                                        AwsCustomGoDependency.GLACIER_CUSTOMIZATION).build())
                                .functionArguments(ListUtils.of(
                                        SymbolUtils.createValueSymbolBuilder("ServiceAPIVersion").build()))
                                .build())
                        .build(),
                RuntimeClientPlugin.builder()
                        .servicePredicate(GlacierCustomizations::isGlacier)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(ACCOUNT_ID_ADDER,
                                        AwsCustomGoDependency.GLACIER_CUSTOMIZATION).build())
                                .functionArguments(ListUtils.of(
                                        SymbolUtils.createValueSymbolBuilder(SET_DEFAULT_ACCOUNT_ID).build()))
                                .build())
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ServiceShape service = settings.getService(model);
        if (!isGlacier(model, service)) {
            return;
        }
        goDelegator.useShapeWriter(service, writer -> {
            writeAccountIdSetter(writer, model, symbolProvider, service);
        });
    }

    private void writeAccountIdSetter(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            ServiceShape service
    ) {
        writer.writeDocs("setDefaultAccountID sets the AccountID to the given value if the current value is nil");
        writer.openBlock("func setDefaultAccountID(input interface{}, accountID string) interface{} {", "}", () -> {
            writer.openBlock("switch i := input.(type) {", "}", () -> {
                for (OperationShape operation : TopDownIndex.of(model).getContainedOperations(service)) {
                    StructureShape input = ProtocolUtils.expectInput(model, operation);

                    List<MemberShape> accountId = input.getAllMembers().values().stream()
                            .filter(m -> m.getMemberName().toLowerCase().equals("accountid"))
                            .toList();

                    if (accountId.isEmpty()) {
                        continue;
                    }

                    writer.openBlock("case $P:", "", symbolProvider.toSymbol(input), () -> {
                        String memberName = symbolProvider.toMemberName(accountId.get(0));
                        writer.write("if i.$L == nil { i.$L = &accountID }", memberName, memberName);
                        writer.write("return i");
                    });
                }
                writer.write("default: return input");
            });
        });
    }

    private static boolean isGlacier(Model model, ServiceShape service) {
        return service.expectTrait(ServiceTrait.class).getSdkId().equalsIgnoreCase("Glacier");
    }
}
