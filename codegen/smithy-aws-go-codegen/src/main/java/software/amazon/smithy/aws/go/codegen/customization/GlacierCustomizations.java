package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import java.util.stream.Collectors;
import software.amazon.smithy.aws.go.codegen.AwsSlotUtils;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.MiddlewareIdentifier;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.go.codegen.integration.StackSlotRegistrar;
import software.amazon.smithy.model.Model;
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
                                .resolvedFunction(customizationValue(TREE_HASH_ADDER))
                                .build())
                        .registerStackSlots(StackSlotRegistrar.builder()
                                .addFinalizeSlotMutators(AwsSlotUtils.addBefore(ListUtils.of(
                                        MiddlewareIdentifier.symbol(customizationValue("TreeHashMiddlewareID"))
                                )))
                                .build())
                        .build(),
                RuntimeClientPlugin.builder()
                        .servicePredicate(GlacierCustomizations::isGlacier)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(customizationValue(API_VERSION_ADDER))
                                .functionArguments(ListUtils.of(
                                        SymbolUtils.createValueSymbolBuilder("ServiceAPIVersion").build()))
                                .build())
                        .registerStackSlots(StackSlotRegistrar.builder()
                                .addSerializeSlotMutator(AwsSlotUtils.addBefore(ListUtils.of(
                                        MiddlewareIdentifier.symbol(customizationValue("APIVersionMiddlewareID"))
                                )))
                                .build())
                        .build(),
                RuntimeClientPlugin.builder()
                        .servicePredicate(GlacierCustomizations::isGlacier)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(customizationValue(ACCOUNT_ID_ADDER))
                                .functionArguments(ListUtils.of(
                                        SymbolUtils.createValueSymbolBuilder(SET_DEFAULT_ACCOUNT_ID).build()))
                                .build())
                        .registerStackSlots(StackSlotRegistrar.builder()
                                .addInitializeSlotMutator(AwsSlotUtils.addBefore(ListUtils.of(
                                        MiddlewareIdentifier.symbol(customizationValue("AccountIDMiddlewareID"))
                                )))
                                .build())
                        .build()
        );
    }

    private Symbol customizationValue(String name) {
        return SymbolUtils.createValueSymbolBuilder(name, AwsCustomGoDependency.GLACIER_CUSTOMIZATION).build();
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
                for (ShapeId operationId : service.getAllOperations()) {
                    OperationShape operation = model.expectShape(operationId, OperationShape.class);
                    StructureShape input = ProtocolUtils.expectInput(model, operation);

                    List<MemberShape> accountId = input.getAllMembers().values().stream()
                            .filter(m -> m.getMemberName().toLowerCase().equals("accountid"))
                            .collect(Collectors.toList());

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
