package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.utils.ListUtils;

import java.util.Comparator;
import java.util.List;

/**
 * Implements a bucket() method on applicable S3 input structures, which returns the principal bucket name from the input.
 */
public class S3Bucketer implements GoIntegration {
    // when deriving the principal input bucket for an operation, if more than one member is targeting the BucketName
    // shape, we take one that has one of the following names
    private final List<String> MATCHED_BUCKET_MEMBERS = ListUtils.of("Bucket", "BucketName");

    public void writeAdditionalFiles(GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator) {
        ServiceShape service = settings.getService(model);
        if (!S3ModelUtils.isServiceS3(model, service)) return;

        for (OperationShape operation : TopDownIndex.of(model).getContainedOperations(service)) {
            goDelegator.useShapeWriter(operation, writer -> {
                writeBucketer(writer, model, symbolProvider, operation);
            });
        }
    }
    
    private void writeBucketer(GoWriter writer, Model model, SymbolProvider symbolProvider, OperationShape operation) {
        StructureShape input = model.expectShape(operation.getInput().get(), StructureShape.class);
        MemberShape bucket = getBucketMember(input);
        if (bucket == null) return;

        String inputName = symbolProvider.toSymbol(input).getName();
        String memberName = bucket.getMemberName();
        writer.write("""
            func (v *$1L) bucket() (string, bool) {
                if v.$2L == nil {
                    return "", false
                }
                return *v.$2L, true
            }""",
            inputName,
            memberName
        );
    }

    private MemberShape getBucketMember(StructureShape input) {
        List<MemberShape> members = input.members().stream()
                .filter(it -> it.getTarget().getName().equals("BucketName"))
                .toList();
        if (members.isEmpty()) return null;
        if (members.size() == 1) return members.get(0);

        members = members.stream()
                .filter(it -> MATCHED_BUCKET_MEMBERS.contains(it.getMemberName()))
                .sorted(Comparator.comparing(MemberShape::getMemberName))
                .toList();
        if (members.size() > 0) return members.get(0);

        throw new CodegenException("could not determine principal bucket input: " + input.getId());
    }
}
