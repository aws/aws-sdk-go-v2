package software.amazon.smithy.aws.go.codegen.customization;

import java.util.Map;
import java.util.Optional;
import java.util.Set;
import java.util.logging.Logger;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.DocumentationTrait;
import software.amazon.smithy.utils.MapUtils;
import software.amazon.smithy.utils.Pair;
import software.amazon.smithy.utils.SetUtils;

public class S3AddPutObjectUnseekableBodyDoc implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(S3AddPutObjectUnseekableBodyDoc.class.getName());

    private static final Map<ShapeId, Set<Pair<ShapeId, String>>> SERVICE_TO_SHAPE_MAP = MapUtils.of(
            ShapeId.from("com.amazonaws.s3#AmazonS3"), SetUtils.of(
                    new Pair(ShapeId.from("com.amazonaws.s3#PutObjectRequest"), "Body"),
                    new Pair(ShapeId.from("com.amazonaws.s3#UploadPartRequest"), "Body")
            )
    );

    @Override
    public byte getOrder() {
        // This integration should happen before other integrations that rely on the presence of this trait
        return -60;
    }

    @Override
    public Model preprocessModel(
            Model model, GoSettings settings
    ) {
        ShapeId serviceId = settings.getService();
        if (!SERVICE_TO_SHAPE_MAP.containsKey(serviceId)) {
            return model;
        }

        Set<Pair<ShapeId, String>> shapeIds = SERVICE_TO_SHAPE_MAP.get(serviceId);

        Model.Builder builder = model.toBuilder();
        for (Pair<ShapeId, String> pair : shapeIds) {
            ShapeId shapeId = pair.getLeft();
            String memberName = pair.getRight();
            StructureShape parent = model.expectShape(shapeId, StructureShape.class);

            Optional<MemberShape> memberOpt = parent.getMember(memberName);
            if (!memberOpt.isPresent()) {
                // Throw in case member is not present, bad things must of happened.
                throw new CodegenException("expect to find " + memberName + " member in shape " + parent.getId());
            }

            MemberShape member = memberOpt.get();
            Shape target = model.expectShape(member.getTarget());

            Optional<DocumentationTrait> docTrait = member.getTrait(DocumentationTrait.class);
            String currentDocs = "";
            if (docTrait.isPresent()) {
                currentDocs = docTrait.get().getValue();
            }
            if (currentDocs.length() != 0) {
                currentDocs += "<br/><br/>";
            }

            final String finalCurrentDocs = currentDocs;
            StructureShape.Builder parentBuilder = parent.toBuilder();
            parentBuilder.removeMember(memberName);
            parentBuilder.addMember(memberName, target.getId(), (memberBuilder) -> {
                memberBuilder
                        .addTraits(member.getAllTraits().values())
                        .addTrait(new DocumentationTrait(finalCurrentDocs +
                                "For using values that are not seekable (io.Seeker) see, " +
                                "https://aws.github.io/aws-sdk-go-v2/docs/sdk-utilities/s3/#unseekable-streaming-input"));
            });


            builder.addShape(parentBuilder.build());
        }

        return builder.build();
    }
}
