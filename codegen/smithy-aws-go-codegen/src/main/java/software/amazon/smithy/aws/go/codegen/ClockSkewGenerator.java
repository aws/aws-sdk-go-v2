package software.amazon.smithy.aws.go.codegen;

import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.*;
import software.amazon.smithy.go.codegen.integration.*;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

import java.util.List;

import static software.amazon.smithy.go.codegen.SmithyGoDependency.ATOMIC;

/**
 * Class to handle clock skew, the discrepancy of time between the client and the server
 * that can cause SDK calls to fail
 */
public class ClockSkewGenerator implements GoIntegration {
    private static final String TIME_OFFSET = "timeOffset";
    private static final String ADD_CLOCK_SKEW_BUILD = "addTimeOffsetBuild";
    private static final String ADD_CLOCK_SKEW_BUILD_MIDDLEWARE = "AddTimeOffsetBuildMiddleware";
    private static final String ADD_CLOCK_SKEW_DESERIALIZER = "addTimeOffsetDeserializer";
    private static final String ADD_CLOCK_SKEW_DESERIALIZE_MIDDLEWARE = "AddTimeOffsetDeserializeMiddleware";
    private static final Symbol TIME_OFFSET_RESOLVER = SymbolUtils.createValueSymbolBuilder(
            "initializeTimeOffsetResolver").build();

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {

        ServiceShape service = settings.getService(model);

        // generate code specific to service client
        goDelegator.useShapeWriter(service, writer -> {
            generateClockSkewInsertMiddleware(writer);
            generateClockSkewDeserializeMiddleware(writer);
            generateTimeOffsetResolver(writer);
        });
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        ClientMember timeOffset = ClientMember.builder()
                .name(TIME_OFFSET)
                .type(ATOMIC.struct("Int64"))
                .documentation("Difference between the time reported by the server and the client")
                .build();
        ClientMemberResolver resolver = ClientMemberResolver.builder()
                .resolver(TIME_OFFSET_RESOLVER)
                .build();
        MiddlewareRegistrar initializeMiddleware = MiddlewareRegistrar.builder()
                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(ADD_CLOCK_SKEW_BUILD).build())
                .functionArguments(ListUtils.of(
                        SymbolUtils.createValueSymbolBuilder("c").build()
                )).build();
        MiddlewareRegistrar finalizeMiddleware = MiddlewareRegistrar.builder()
                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(ADD_CLOCK_SKEW_DESERIALIZER).build())
                .functionArguments(ListUtils.of(
                        SymbolUtils.createValueSymbolBuilder("c").build()
                )).build();
        return List.of(
                RuntimeClientPlugin.builder()
                        .addClientMember(timeOffset)
                        .addClientMemberResolver(resolver)
                        .registerMiddleware(initializeMiddleware)
                        .build(),
                RuntimeClientPlugin.builder()
                        .registerMiddleware(finalizeMiddleware)
                        .build()
        );
    }

    private void generateClockSkewInsertMiddleware(GoWriter writer) {
        Symbol stackSymbol = SymbolUtils.createPointableSymbolBuilder("Stack", SmithyGoDependency.SMITHY_MIDDLEWARE)
                .build();
        writer.openBlock("func $L(stack $P, c *Client) error{", "}", ADD_CLOCK_SKEW_BUILD, stackSymbol, () -> {
            writer.write("return stack.Build.Add(&awsmiddleware.$L{Offset: c.$L}, middleware.After)",
                    ADD_CLOCK_SKEW_BUILD_MIDDLEWARE, TIME_OFFSET);
        });
    }

    private void generateClockSkewDeserializeMiddleware(GoWriter writer) {
        Symbol stackSymbol = SymbolUtils.createPointableSymbolBuilder("Stack", SmithyGoDependency.SMITHY_MIDDLEWARE)
                .build();
        writer.openBlock("func $L(stack $P, c *Client) error{", "}", ADD_CLOCK_SKEW_DESERIALIZER, stackSymbol, () ->
                writer.write("return stack.Deserialize.Insert(&awsmiddleware.$L{Offset: c.$L}, \"RecordResponseTiming\", middleware.Before)",
                        ADD_CLOCK_SKEW_DESERIALIZE_MIDDLEWARE, TIME_OFFSET)
        );
    }

    private void generateTimeOffsetResolver(GoWriter writer) {
        writer.openBlock("func $L(c *Client) {", "}", TIME_OFFSET_RESOLVER, () -> {
            Symbol atomic = SymbolUtils.createValueSymbolBuilder("Int64", ATOMIC).build();
            writer.write("c.$L = new($P)", TIME_OFFSET, atomic);
        });
    }
}