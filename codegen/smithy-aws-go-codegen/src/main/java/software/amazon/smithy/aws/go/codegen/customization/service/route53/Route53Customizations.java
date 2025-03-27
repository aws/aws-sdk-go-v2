package software.amazon.smithy.aws.go.codegen.customization.service.route53;

import java.util.List;
import java.util.stream.Collectors;

import software.amazon.smithy.aws.go.codegen.customization.AwsCustomGoDependency;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.utils.ListUtils;

public class Route53Customizations implements GoIntegration {
    private static final String ADD_ERROR_HANDLER_INTERNAL = "HandleCustomErrorDeserialization";
    private static final String URL_SANITIZE_ADDER = "addSanitizeURLMiddleware";
    private static final String URL_SANITIZE_INTERNAL_ADDER= "AddSanitizeURLMiddleware";
    private static final String SANITIZE_URL_INPUT = "sanitizeURLInput";

    @Override
    public byte getOrder() {
        // The associated customization ordering is relative to operation deserializers
        // and thus the integration should be added at the end.
        return 127;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .operationPredicate(Route53Customizations::supportsCustomError)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(ADD_ERROR_HANDLER_INTERNAL,
                                        AwsCustomGoDependency.ROUTE53_CUSTOMIZATION).build())
                                .build())
                        .build(),
                RuntimeClientPlugin.builder()
                        .operationPredicate(Route53Customizations::supportsSanitizeURLValue)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(URL_SANITIZE_ADDER).build())
                                .build())
                        .build()
        );
    }

    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        if (!isRoute53Service(model, settings.getService(model))) {
            return;
        }

        ServiceShape service = settings.getService(model);
        goDelegator.useShapeWriter(service, this::writeMiddlewareHelper);

        goDelegator.useShapeWriter(service, writer -> {
                writeURLInputSanitizer(writer, model, symbolProvider, service);
            });
    }

    private void writeMiddlewareHelper(GoWriter writer) {
        writer.openBlock("func $L(stack *middleware.Stack) error {", "}", URL_SANITIZE_ADDER, () -> {
            writer.write("return $T(stack, $T{SanitizeURLInput: $L})",
                    SymbolUtils.createValueSymbolBuilder(URL_SANITIZE_INTERNAL_ADDER,
                            AwsCustomGoDependency.ROUTE53_CUSTOMIZATION).build(),
                    SymbolUtils.createValueSymbolBuilder(URL_SANITIZE_INTERNAL_ADDER + "Options",
                            AwsCustomGoDependency.ROUTE53_CUSTOMIZATION).build(),
                    SANITIZE_URL_INPUT
            );
        });
        writer.insertTrailingNewline();
    }

    private void writeURLInputSanitizer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            ServiceShape service
    ) {

        writer.writeDocs("Check for and split apart Route53 resource IDs, setting only the last piece. " +
                "This allows the output of one operation e.g. foo/1234 to be used as input in another operation " +
                "(e.g. it expects just '1234')");

        writer.openBlock("func sanitizeURLInput(input interface{}) error {", "}", () -> {
            writer.openBlock("switch i:= input.(type) {", "}", () -> {
                TopDownIndex.of(model).getContainedOperations(service).forEach((operation)-> {
                    StructureShape input = model.expectShape(operation.getInput().get(), StructureShape.class);
                    List<MemberShape> resourceMembers = getResourceMembers(service, input);

                    if (!resourceMembers.isEmpty()){
                        writer.openBlock("case $P :", "", symbolProvider.toSymbol(input), () -> {
                            writer.addUseImports(SmithyGoDependency.STRINGS);
                            for (MemberShape member : resourceMembers) {
                                String memberName = member.getMemberName();
                               writer.openBlock("if i.$L != nil {", "}", memberName, () -> {
                                writer.write("idx := strings.LastIndex(*i.$L, `/`)", memberName);
                                writer.write("v := (*i.$L)[idx+1:]", memberName);
                                writer.write("i.$L = &v", memberName);
                               });
                            }
                        });
                    }
                });
                writer.write("default: break");
            });
            writer.write("return nil");
        });
    }

    private static List<MemberShape> getResourceMembers(ServiceShape service, StructureShape input) {
        return input.getAllMembers().values().stream()
                .filter(m -> m.getTarget().getName(service).equalsIgnoreCase("ResourceId")
                        || m.getTarget().getName(service).equalsIgnoreCase("DelegationSetId")
                        || m.getTarget().getName(service).equalsIgnoreCase("ChangeId"))
                .collect(Collectors.toList());
    }

    // returns true if the operation supports custom route53 error response
    private static boolean supportsCustomError(Model model, ServiceShape service, OperationShape operation) {
        if (!isRoute53Service(model, service)) {
            return false;
        }

        return operation.getId().getName(service).equalsIgnoreCase("ChangeResourceRecordSets");
    }

    // return true if the operation takes input that supports resource sanitization.
    //
    // For Route53, HostedZoneID, and changeID is supported by member shapes targeting `ResourceId` or `DelegationSetId`.
    private static boolean supportsSanitizeURLValue(Model model, ServiceShape service, OperationShape operation) {
        if (!isRoute53Service(model, service)) {
            return false;
        }

        StructureShape input = model.expectShape(operation.getInput().get(), StructureShape.class);
        List<MemberShape> targetMembers = getResourceMembers(service, input);

        if (targetMembers.size() >1 ){
            throw new CodegenException(String.format("Route53 service has ResourceId, DelegationSetId members " +
                            "modeled on %s shape", input.getId().getName(service)));
        }

        return targetMembers.size() != 0;
    }

    // returns true if service is route53
    private static boolean isRoute53Service(Model model, ServiceShape service) {
        String serviceId= service.expectTrait(ServiceTrait.class).getSdkId();
        return serviceId.equalsIgnoreCase("Route 53");
    }
}
