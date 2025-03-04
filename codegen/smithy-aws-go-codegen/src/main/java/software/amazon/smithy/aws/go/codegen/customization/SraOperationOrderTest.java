package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.go.codegen.GoCodegenContext;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.knowledge.TopDownIndex;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

/**
 * Validates that service client operations are performed in the orders specified by the Smithy Reference Architecture (SRA).
 */
public class SraOperationOrderTest implements GoIntegration {
    @Override
    public void writeAdditionalFiles(GoCodegenContext ctx) {
        ctx.writerDelegator().useFileWriter("sra_operation_order_test.go", ctx.settings().getModuleName(), writer -> {
            writer.write(renderCommonTestSource());

            TopDownIndex.of(ctx.model())
                    .getContainedOperations(ctx.settings().getService(ctx.model()))
                    .forEach(it -> {
                        var operationName = ctx.symbolProvider().toSymbol(it).getName();
                        writer.write(renderTest(operationName));
                    });
        });
    }

    private GoWriter.Writable renderCommonTestSource() {
        return goTemplate("""
                $D $D
                var errTestReturnEarly = errors.New("errTestReturnEarly")

                func captureMiddlewareStack(stack *middleware.Stack) func(*middleware.Stack) error {
                	return func(inner *middleware.Stack) error {
                		*stack = *inner
                		return errTestReturnEarly
                	}
                }
                """, SmithyGoDependency.ERRORS, SmithyGoDependency.SMITHY_MIDDLEWARE);
    }

    private GoWriter.Writable renderTest(String operationName) {
        return goTemplate("""
                $1D $2D $3D $4D $5D $6D
                func TestOp$7LSRAOperationOrder(t *testing.T) {
                	expect := []string{
                		"OperationSerializer",
                		"Retry",
                		"ResolveAuthScheme",
                		"GetIdentity",
                		"ResolveEndpointV2",
                		"Signing",
                		"OperationDeserializer",
                	}

                	var captured middleware.Stack
                	svc := New(Options{
                		APIOptions: []func(*middleware.Stack) error{
                			captureMiddlewareStack(&captured),
                		},
                	})
                	_, err := svc.$7L(context.Background(), nil)
                	if err != nil && !errors.Is(err, errTestReturnEarly) {
                		t.Fatalf("unexpected error: %v", err)
                	}

                	var actual, all []string
                	for _, step := range strings.Split(captured.String(), "\\n") {
                		trimmed := strings.TrimSpace(step)
                		all = append(all, trimmed)
                		if slices.Contains(expect, trimmed) {
                			actual = append(actual, trimmed)
                		}
                	}

                	if !slices.Equal(expect, actual) {
                		t.Errorf("order mismatch:\\nexpect: %v\\nactual: %v\\nall: %v", expect, actual, all)
                	}
                }
                """,
                SmithyGoDependency.ERRORS,
                SmithyGoDependency.TESTING,
                SmithyGoDependency.CONTEXT,
                SmithyGoDependency.STRINGS,
                SmithyGoDependency.SLICES,
                SmithyGoDependency.SMITHY_MIDDLEWARE,
                operationName
        );
    }
}
