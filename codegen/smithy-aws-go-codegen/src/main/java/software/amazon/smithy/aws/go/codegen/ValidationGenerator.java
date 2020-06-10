/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *  http://aws.amazon.com/apache2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *  http://aws.amazon.com/apache2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *  http://aws.amazon.com/apache2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *  http://aws.amazon.com/apache2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package software.amazon.smithy.aws.go.codegen;

import java.util.Map;
import java.util.Set;
import java.util.TreeMap;
import java.util.function.Consumer;
import software.amazon.smithy.aws.go.codegen.knowledge.ValidationIndex;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDependency;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.TriConsumer;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.CollectionShape;
import software.amazon.smithy.model.shapes.MapShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.traits.EnumTrait;
import software.amazon.smithy.utils.StringUtils;

/**
 * Generates Go validation middleware and shape helpers.
 */
public class ValidationGenerator implements GoIntegration {

    private void execute(GoWriter writer, Model model, SymbolProvider symbolProvider, ServiceShape service) {
        ValidationIndex validationIndex = model.getKnowledge(ValidationIndex.class);
        Map<Shape, OperationShape> inputShapeToOperation = new TreeMap<>();
        validationIndex.getOperationsRequiringValidation(service).forEach(operationShape -> {
            Shape inputShape = model.expectShape(operationShape.getInput().get());
            inputShapeToOperation.put(inputShape, operationShape);
        });
        Set<Shape> shapesWithHelpers = validationIndex.getShapesRequiringValidationHelpers(service);

        generateOperationValidationMiddleware(writer, symbolProvider, inputShapeToOperation);
        generateShapeValidationFunctions(writer, symbolProvider, model, inputShapeToOperation.keySet(),
                shapesWithHelpers);
    }

    private void generateOperationValidationMiddleware(
            GoWriter writer,
            SymbolProvider symbolProvider,
            Map<Shape, OperationShape> operationShapeMap
    ) {
        for (Map.Entry<Shape, OperationShape> entry : operationShapeMap.entrySet()) {
            GoStackStepMiddlewareGenerator generator = GoStackStepMiddlewareGenerator.createBuildStepMiddleware(
                    getOperationValidationMiddlewareName(entry.getValue()));
            String helperName = getShapeValidationFunctionName(entry.getKey(), true);
            Symbol inputSymbol = symbolProvider.toSymbol(entry.getKey());
            generator.writeMiddleware(writer, (g, w) -> {
                writer.addUseImports(GoDependency.FMT);
                // cast input parameters type to the input type of the operation
                writer.write("input, ok := in.Parameters.($P)", inputSymbol);
                writer.openBlock("if !ok {", "}", () -> {
                    writer.write("return out, metadata, "
                            + "fmt.Errorf(\"unknown input parameters type %T\", in.Parameters)");
                });
                writer.openBlock("if err := $L(input); err != nil {", "}", helperName,
                        () -> writer.write("return err"));
                writer.write("return next.$L(ctx, in)", g.getHandleMethodName());
            });
            writer.write("");
        }
    }

    private void generateShapeValidationFunctions(
            GoWriter writer,
            SymbolProvider symbolProvider,
            Model model,
            Set<Shape> operationInputShapes,
            Set<Shape> shapesWithHelpers
    ) {
        for (Shape shape : shapesWithHelpers) {
            boolean topLevelShape = operationInputShapes.contains(shape);
            String functionName = getShapeValidationFunctionName(shape, topLevelShape);
            Symbol shapeSymbol = symbolProvider.toSymbol(shape);
            writer.openBlock("func $L(v $P) error {", "}", functionName, shapeSymbol, () -> {
                writer.addUseImports(GoDependency.SMITHY);
                writer.openBlock("if v == nil {", "}", () -> writer.write("return nil"));
                writer.write("invalidParams := smithy.InvalidParamsError{Context: $S}", shapeSymbol.getName());
                switch (shape.getType()) {
                    case STRUCTURE:
                        shape.members().forEach(memberShape -> {
                            String memberName = symbolProvider.toMemberName(memberShape);
                            Shape targetShape = model.expectShape(memberShape.getTarget());
                            boolean required = ValidationIndex.isRequiredParameter(model, memberShape, topLevelShape);
                            boolean hasHelper = shapesWithHelpers.contains(targetShape);
                            boolean isEnum = targetShape.getTrait(EnumTrait.class).isPresent();
                            if (required) {
                                if (isEnum) {
                                    writer.write("if len(v.$L) == 0 {", memberName);
                                } else {
                                    writer.write("if v.$L == nil {", memberName);
                                }
                                writer.write("invalidParams.Add(smithy.NewErrParamRequired($S))", memberName);
                                if (hasHelper) {
                                    writer.writeInline("} else ");
                                } else {
                                    writer.write("}");
                                }
                            }
                            if (hasHelper) {
                                Runnable runnable = () -> {
                                    String helperName = getShapeValidationFunctionName(targetShape, false);
                                    writer.openBlock("if err := $L(v.$L); err != nil {", "}", helperName, memberName,
                                            () -> writer.write("invalidParams.AddNested($S, err)", memberName));
                                };
                                if (isEnum) {
                                    writer.openBlock("if len(v.$L) > 0 {", "}", memberName, runnable);
                                } else {
                                    writer.openBlock("if v.$L != nil {", "}", memberName, runnable);
                                }

                            }
                        });
                        break;
                    case LIST:
                    case SET:
                        String helperName = getShapeValidationFunctionName(model.expectShape(((CollectionShape) shape)
                                .getMember().getTarget()), false);
                        writer.openBlock("for i := range v {", "}", () -> {
                            writer.openBlock("if err := $L(v[i]); err != nil {", "}", helperName, () -> {
                                writer.write("invalidParams.AddNested(fmt.Sprintf(\"[%d]\", i), err)");
                            });
                        });
                        break;
                    case MAP:
                        helperName = getShapeValidationFunctionName(model.expectShape(((MapShape) shape).getValue()
                                .getTarget()), false);
                        writer.openBlock("for key := range v {", "}", () -> {
                            writer.openBlock("if err := $L(v[key]); err != nil {", "}", helperName, () -> {
                                writer.write("invalidParams.AddNested(fmt.Sprintf(\"[%q]\", key), err)");
                            });
                        });
                        break;
                    case UNION:
                        // TODO: Implement Union support
                    default:
                        throw new CodegenException("Unexpected validation helper shape type " + shape.getType());
                }

                writer.write("if invalidParams.Len() > 0 {");
                writer.write("return invalidParams");
                writer.write("} else {");
                writer.write("return nil");
                writer.write("}");
            });
            writer.write("");
        }
    }

    private static String getOperationValidationMiddlewareName(OperationShape operationShape) {
        return "validateOp"
                + StringUtils.capitalize(operationShape.getId().getName());
    }

    private static String getShapeValidationFunctionName(Shape shape, boolean topLevelOpShape) {
        StringBuilder builder = new StringBuilder();
        builder.append("validate");
        if (topLevelOpShape) {
            builder.append("Op");
        }
        builder.append(StringUtils.capitalize(shape.getId().getName()));
        return builder.toString();
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            TriConsumer<String, String, Consumer<GoWriter>> writerFactory
    ) {
        writerFactory.accept("validators.go", settings.getModuleName(), writer -> {
            execute(writer, model, symbolProvider, settings.getService(model));
        });
    }
}
