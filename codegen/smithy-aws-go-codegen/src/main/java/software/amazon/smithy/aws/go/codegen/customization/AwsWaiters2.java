package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.integration.Waiters2;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ShapeId;

import java.util.List;
import java.util.Set;

import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

public class AwsWaiters2 extends Waiters2 {
    public static final List<ShapeId> PHASED_ROLLOUT_SERVICES = List.of(
            ShapeId.from("com.amazonaws.ec2#AmazonEC2")
    );

    @Override
    public Set<Symbol> getAdditionalClientOptions() {
        return Set.of(buildPackageSymbol("addIsWaiterUserAgent"));
    }

    @Override
    public boolean enabledForService(Model model, ShapeId service) {
        return PHASED_ROLLOUT_SERVICES.contains(service);
    }
}
