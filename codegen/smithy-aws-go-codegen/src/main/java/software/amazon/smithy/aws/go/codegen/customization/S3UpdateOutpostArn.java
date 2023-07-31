package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
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
import software.amazon.smithy.utils.SetUtils;

/**
 * Implements the arnable interface on all relevant S3/S3Control outposts operations.
 */
public class S3UpdateOutpostArn implements GoIntegration {

    private final Set<String> LIST_ACCESSPOINT_ARN_INPUT = SetUtils.of(
            "GetAccessPoint", "DeleteAccessPoint", "PutAccessPointPolicy",
            "GetAccessPointPolicy", "DeleteAccessPointPolicy"
    );

    private final MiddlewareRegistrar middlewareAdder =
        MiddlewareRegistrar.builder()
                .resolvedFunction(SymbolUtils.createValueSymbolBuilder("AddUpdateOutpostARN", AwsCustomGoDependency.S3CONTROL_CUSTOMIZATION).build())
                .build();

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3ModelUtils::isServiceS3Control)
                        .registerMiddleware(middlewareAdder)
                        .build()
                );
    }

    public void writeAdditionalFiles(GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator) {
        ServiceShape service = settings.getService(model);
        if (!S3ModelUtils.isServiceS3Control(model, service)) return;

        for (OperationShape operation : TopDownIndex.of(model).getContainedOperations(service)) {
            goDelegator.useShapeWriter(operation, writer -> {
                writeARNHelper(settings, writer, model, symbolProvider, operation);
            });
        }
    }
    
    private void writeARNHelper(
            GoSettings settings, GoWriter writer,
            Model model, SymbolProvider symbolProvider, OperationShape operation
    ) {
        ServiceShape service = settings.getService(model);

        String arnType = LIST_ACCESSPOINT_ARN_INPUT.contains(
                operation.getId().getName(service)
        ) ? "AccessPointName" : "BucketName";

        StructureShape input = model.expectShape(operation.getInput().get(), StructureShape.class);
        List<MemberShape> listOfARNMembers = input.getAllMembers().values().stream()
                .filter(m -> m.getTarget().getName(service).equals(arnType))
                .collect(Collectors.toList());
        if (listOfARNMembers.size() > 1) {
            throw new CodegenException(arnType + " shape should be targeted by only one input member, found " +
                                        listOfARNMembers.size() + " for Input shape: " + input.getId());
        }

        if (listOfARNMembers.isEmpty()) {
            return;
        }

        String inputName = symbolProvider.toSymbol(input).getName();
        String memberName = listOfARNMembers.get(0).getMemberName();

        writer.write(
            """
                func (m *$1L) GetARNMember() (*string, bool) {
                    if m.$2L == nil {
                        return nil, false
                    }
                    return m.$2L, true
                }
            """,
            inputName,
            memberName
        );


        writer.write(
            """
                func (m *$1L) SetARNMember(v string) error {
                    m.$2L = &v
                    return nil
                }
            """,
            inputName,
            memberName
        );
    }



}
