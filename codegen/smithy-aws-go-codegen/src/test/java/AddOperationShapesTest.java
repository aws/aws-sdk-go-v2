import org.commonmark.node.Code;
import org.hamcrest.Matcher;
import org.hamcrest.MatcherAssert;
import org.hamcrest.Matchers;
import org.junit.jupiter.api.Test;
import software.amazon.smithy.aws.go.codegen.AddOperationShapes;
import software.amazon.smithy.go.codegen.CodegenUtils;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.SyntheticClone;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StringShape;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.utils.ListUtils;

import java.util.Optional;
import java.util.logging.Logger;

public class AddOperationShapesTest {
    private static final Logger LOGGER = Logger.getLogger(AddOperationShapesTest.class.getName());

    private static final String NAMESPACE = "go.codege.test";
    private static final ServiceShape SERVICE = ServiceShape.builder()
            .id(ShapeId.fromParts(NAMESPACE, "TestService"))
            .version("1.0")
            .build();

    @Test
    public void testOperationWithoutInputOutput() {
        GoSettings settings = new GoSettings();
        settings.setService(SERVICE.toShapeId());

        OperationShape op = OperationShape.builder()
                .id(ShapeId.fromParts(NAMESPACE, "TestOperation"))
                .build();

        ServiceShape service = SERVICE.toBuilder()
                .addOperation(op.getId())
                .build();

        Model model = Model.builder()
                .addShapes(service, op)
                .build();

        AddOperationShapes addOperationShapes = new AddOperationShapes();
        Model processedModel = addOperationShapes.preprocessModel(model, settings);

        ListUtils.of(
                ShapeId.fromParts(CodegenUtils.getSyntheticTypeNamespace(), "TestOperationInput"),
                ShapeId.fromParts(CodegenUtils.getSyntheticTypeNamespace(), "TestOperationOutput")
        ).forEach(shapeId -> {
            Optional<Shape> shape = processedModel.getShape(shapeId);
            MatcherAssert.assertThat(shapeId + " shape must be present in model",
                    shape.isPresent(), Matchers.is(true));
            MatcherAssert.assertThat(shapeId + " shape must have no members",
                    shape.get().members().size(), Matchers.equalTo(0));
            if (shape.get().getTrait(SyntheticClone.class).isPresent()) {
                MatcherAssert.assertThat("shape is not synthetic clone", false);
            }
        });
    }

    @Test
    public void testOperationWithExistingInputOutput() {
        GoSettings settings = new GoSettings();
        settings.setService(SERVICE.toShapeId());

        StringShape stringShape = StringShape.builder().id(ShapeId.fromParts(NAMESPACE, "String")).build();

        StructureShape inputShape = StructureShape.builder()
                .id(ShapeId.fromParts(NAMESPACE, "TestOperationRequest"))
                .addMember("foo", stringShape.getId())
                .build();
        StructureShape outputShape = StructureShape.builder()
                .id(ShapeId.fromParts(NAMESPACE, "TestOperationResponse"))
                .addMember("foo", stringShape.getId())
                .build();

        OperationShape op = OperationShape.builder()
                .id(ShapeId.fromParts(NAMESPACE, "TestOperation"))
                .input(inputShape)
                .output(outputShape)
                .build();

        ServiceShape service = SERVICE.toBuilder()
                .addOperation(op.getId())
                .build();

        Model model = Model.builder()
                .addShapes(stringShape, inputShape, outputShape, op, service)
                .build();

        AddOperationShapes addOperationShapes = new AddOperationShapes();
        Model processedModel = addOperationShapes.preprocessModel(model, settings);

        ListUtils.of(
                ShapeId.fromParts(CodegenUtils.getSyntheticTypeNamespace(), "TestOperationInput"),
                ShapeId.fromParts(CodegenUtils.getSyntheticTypeNamespace(), "TestOperationOutput")
        ).forEach(shapeId -> {
            Optional<Shape> shape = processedModel.getShape(shapeId);
            MatcherAssert.assertThat(shapeId + " shape must be present in model",
                    shape.isPresent(), Matchers.is(true));

            StructureShape structureShape = shape.get().asStructureShape().get();
            Optional<MemberShape> fooMember = structureShape.getMember("foo");
            MatcherAssert.assertThat("foo member present", fooMember.isPresent(), Matchers.is(true));

            ShapeId id = fooMember.get().getId();
            MatcherAssert.assertThat("foo is correct namespace", id.getNamespace(), Matchers.equalTo(shapeId.getNamespace()));
            MatcherAssert.assertThat("foo is correct parent", id.getName(), Matchers.equalTo(shapeId.getName()));

            Optional<SyntheticClone> syntheticClone = shape.get().getTrait(SyntheticClone.class);
            if (!syntheticClone.isPresent()) {
                MatcherAssert.assertThat("shape must be marked as synthetic clone", false);
            } else {
                MatcherAssert.assertThat(syntheticClone.get().getArchetype().toString(),
                        Matchers.is(Matchers.oneOf(NAMESPACE + "#TestOperationRequest",
                                NAMESPACE + "#TestOperationResponse")));
            }
        });
    }

    @Test
    public void testOperationWithExistingInputOutputWithConflitcs() {
        GoSettings settings = new GoSettings();
        settings.setService(SERVICE.toShapeId());

        StructureShape inputConflict = StructureShape.builder()
                .id(ShapeId.fromParts(NAMESPACE, "TestOperationInput"))
                .build();
        StructureShape outputConflict = StructureShape.builder()
                .id(ShapeId.fromParts(NAMESPACE, "TestOperationOutput"))
                .build();

        OperationShape op = OperationShape.builder()
                .id(ShapeId.fromParts(NAMESPACE, "TestOperation"))
                .build();

        ServiceShape service = SERVICE.toBuilder()
                .addOperation(op.getId())
                .build();

        Model model = Model.builder()
                .addShapes(inputConflict, outputConflict, op, service)
                .build();

        AddOperationShapes addOperationShapes = new AddOperationShapes();
        Model processedModel = addOperationShapes.preprocessModel(model, settings);

        ShapeId expInputRename = ShapeId.fromParts(CodegenUtils.getSyntheticTypeNamespace(), "TestOperationInput");
        ShapeId expOutputRename = ShapeId.fromParts(CodegenUtils.getSyntheticTypeNamespace(), "TestOperationOutput");

        ListUtils.of(expInputRename, expOutputRename, inputConflict.getId(), outputConflict.getId())
                .forEach(shapeId -> {
                    Optional<Shape> shape = processedModel.getShape(shapeId);
                    MatcherAssert.assertThat(shapeId + " shape must be present in model",
                            shape.isPresent(), Matchers.is(true));

                    if (shape.get().getTrait(SyntheticClone.class).isPresent()) {
                        MatcherAssert.assertThat("shape must not be marked as synthetic clone", false);
                    }
                });
    }
}
