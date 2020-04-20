package software.amazon.smithy.aws.go.codegen;

import software.amazon.smithy.go.codegen.CodegenUtils;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.ShapeCloner;
import software.amazon.smithy.go.codegen.SyntheticClone;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.NeighborProviderIndex;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.neighbor.Relationship;
import software.amazon.smithy.model.shapes.ListShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.SetShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.shapes.UnionShape;
import software.amazon.smithy.utils.CaseUtils;
import software.amazon.smithy.utils.ListUtils;

import java.util.Optional;
import java.util.TreeSet;
import java.util.logging.Logger;

/**
 * Ensures that each operation has a unique input and output shape
 */
public class AddOperationShapes implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(AddOperationShapes.class.getName());

    @Override
    public byte getOrder() {
        return (byte) 127;
    }

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        TopDownIndex topDownIndex = model.getKnowledge(TopDownIndex.class);
        TreeSet<OperationShape> operations = new TreeSet<>(topDownIndex.getContainedOperations(settings.getService()));

        Model.Builder modelBuilder = model.toBuilder();

        for (OperationShape operation : operations) {
            ShapeId operationId = operation.getId();
            LOGGER.info(() -> "building unique input/output shapes for " + operationId);

            StructureShape newInputShape = operation.getInput()
                    .map(shapeId -> cloneOperationShape(operationId, (StructureShape) model.expectShape(shapeId),
                            "Input"))
                    .orElseGet(() -> emptyOperationStructure(operationId, "Input"));

            StructureShape newOutputShape = operation.getInput()
                    .map(shapeId -> cloneOperationShape(operationId, (StructureShape) model.expectShape(shapeId),
                            "Output"))
                    .orElseGet(() -> emptyOperationStructure(operationId, "Output"));

            // Add new input/output to model
            modelBuilder.addShape(newInputShape);
            modelBuilder.addShape(newOutputShape);

            // Update operation model with the input/output shape ids
            modelBuilder.addShape(operation.toBuilder()
                    .input(newInputShape.toShapeId())
                    .output(newOutputShape.toShapeId())
                    .build());
        }

        return modelBuilder.build();
    }

    private static StructureShape emptyOperationStructure(ShapeId opShapeId, String suffix) {
        return StructureShape.builder()
                .id(ShapeId.fromParts(CodegenUtils.getSyntheticTypeNamespace(), opShapeId.getName() + suffix))
                .build();
    }

    private static StructureShape cloneOperationShape(ShapeId opShapeId, StructureShape structureShape, String suffix) {
        return ((StructureShape) new ShapeCloner(shapeId -> ShapeId.fromParts(CodegenUtils.getSyntheticTypeNamespace(),
                opShapeId.getName() + suffix))
                .structureShape(structureShape)).toBuilder()
                .build();
    }
}
