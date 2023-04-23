package software.amazon.smithy.aws.go.codegen;

import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.endpoints.FnProvider;


public class AwsFnProvider implements FnProvider {
    @Override
    public Symbol fnFor(String name) {
        return switch (name) {
            case "aws.partition" -> SymbolUtils.createValueSymbolBuilder("GetPartition",
                    AwsGoDependency.AWS_ENDPOINT_RULESFN).build();
            case "aws.parseArn" -> SymbolUtils.createValueSymbolBuilder("ParseARN",
                    AwsGoDependency.AWS_ENDPOINT_RULESFN).build();
            case "aws.isVirtualHostableS3Bucket" ->
                    SymbolUtils.createValueSymbolBuilder("IsVirtualHostableS3Bucket",
                    AwsGoDependency.AWS_ENDPOINT_RULESFN).build();

            default -> null;
        };
    }
}
