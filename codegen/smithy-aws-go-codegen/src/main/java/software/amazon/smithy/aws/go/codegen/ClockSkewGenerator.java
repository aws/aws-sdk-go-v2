package software.amazon.smithy.aws.go.codegen;

import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.*;
import software.amazon.smithy.go.codegen.integration.*;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

import java.util.List;
import java.util.Map;

import static software.amazon.smithy.aws.go.codegen.AwsGoDependency.INTERNAL_MIDDLEWARE;
import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SmithyGoDependency.ATOMIC;

/**
 * Class to handle clock skew, the discrepancy of time between the client and the server
 * that can cause SDK calls to fail
 */
public class ClockSkewGenerator implements GoIntegration {
    private static final String TIME_OFFSET = "timeOffset";

    private static final Symbol TIME_OFFSET_RESOLVER = SymbolUtils.createValueSymbolBuilder(
            "initializeTimeOffsetResolver").build();

    private static final Writable TIME_OFFSET_RESOLVER_TEMPLATE = goTemplate(
            """
                    $import:D
                    func $fn:L(c *Client) {
                        c.$off:L = new(atomic.Int64)
                    }
                    """,
            Map.of(
                    "import", ATOMIC,
                    "fn", TIME_OFFSET_RESOLVER,
                    "off", TIME_OFFSET
            )
    );

    private static final ClientMember TIME_OFFSET_MEMBER = ClientMember.builder()
            .name(TIME_OFFSET)
            .type(ATOMIC.struct("Int64"))
            .documentation("Difference between the time reported by the server and the client")
            .build();
    private static final ClientMemberResolver TIME_OFFSET_MEMBER_RESOLVER = ClientMemberResolver.builder()
            .resolver(TIME_OFFSET_RESOLVER)
            .build();
    private static final List<RuntimeClientPlugin> CLIENT_PLUGINS = List.of(
            RuntimeClientPlugin.builder()
                    .addClientMember(TIME_OFFSET_MEMBER)
                    .addClientMemberResolver(TIME_OFFSET_MEMBER_RESOLVER)
                    .build()
    );

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
            writer.write(TIME_OFFSET_RESOLVER_TEMPLATE);
        });
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return CLIENT_PLUGINS;
    }
}
