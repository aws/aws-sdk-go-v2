package software.amazon.smithy.aws.go.codegen;


import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.TriConsumer;
import software.amazon.smithy.go.codegen.endpoints.EndpointResolutionGenerator;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;
import software.amazon.smithy.utils.SetUtils;

import java.util.List;
import java.util.Map;
import java.util.function.Consumer;

public class AwsEndpointResolverAdapterGenerator implements GoIntegration {

    public static final String LEGACY_ADAPTER_TYPE = "legacyEndpointResolverAdapter";
    public static final String COMPATIBLE_ADAPTER_TYPE = "compatibleEndpointResolver";
    public static final String FINALIZE_ENDPOINT_RESOLVER_V2 = "finalize" + EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME;


    private Map<String, Object> commonCodegenArgs;

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            TriConsumer<String, String, Consumer<GoWriter>> writerFactory
    ) {
        this.commonCodegenArgs = MapUtils.of(
                "goContext", SymbolUtils.createValueSymbolBuilder("Context", SmithyGoDependency.CONTEXT).build(),
                "legacyAdapterType", SymbolUtils.createValueSymbolBuilder(LEGACY_ADAPTER_TYPE).build(),
                "legacyResolverType", SymbolUtils.createValueSymbolBuilder(EndpointGenerator.RESOLVER_INTERFACE_NAME).build(),
                "resolverType", SymbolUtils.createValueSymbolBuilder(EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME).build(),
                "paramsType", SymbolUtils.createValueSymbolBuilder(EndpointResolutionGenerator.PARAMETERS_TYPE_NAME).build(),
                "smithyEndpointType", SymbolUtils.createValueSymbolBuilder("Endpoint", SmithyGoDependency.SMITHY_ENDPOINTS).build(),
                "awsEndpointType", SymbolUtils.createValueSymbolBuilder("Endpoint", AwsGoDependency.AWS_CORE).build(),
                "fmtErrorf", SymbolUtils.createValueSymbolBuilder("Errorf", SmithyGoDependency.FMT).build());

        var content = new GoWriter.ChainWritable()
                .add(generateLegacyAdapter())
                .add(generateCompatibleAdapter())
                .add(generatePseudoRegionUtility())
                .add(generateFinalizeMethod())
                .compose();

        writerFactory.accept("endpoints.go", settings.getModuleName(), writer -> {
            writer.write("$W", content);
        });

    }

    private GoWriter.Writable generateLegacyAdapter() {
        return goTemplate("""
            type $legacyAdapterType:L struct {
                legacyResolver $legacyResolverType:T
                resolver       $resolverType:T
            }

            func (l *$legacyAdapterType:T) ResolveEndpoint(ctx $goContext:T, params $paramsType:T) (endpoint $smithyEndpointType:T, err error) {
                $legacyResolveMethodBody:W
            }
            """,
            commonCodegenArgs,
            MapUtils.of(
                    "legacyResolveMethodBody", generateLegacyResolveMethodBody()
            ));
    }


    private GoWriter.Writable generateLegacyResolveMethodBody() {
        return goTemplate(
            """
                $requiredOptions:W

                $resolveInvocation:W

                $hostnameImmutableCheck:W

                $pushdownToResolver:W

            """,
            commonCodegenArgs,
            MapUtils.of(
                    "requiredOptions", generateRequiredOptions(),
                    "resolveInvocation", generateResolveInvocation(),
                    "hostnameImmutableCheck", generateHostnameImmutableCheck(),
                    "pushdownToResolver", generatePushdownToResolver()));
    }

    private GoWriter.Writable generateRequiredOptions() {
        return goTemplate(
            """
            var fips $awsFipsEndpointStateType:T
            var dualStack $awsDualStackEndpointStateType:T

            if aws.ToBool(params.UseFIPS) {
                fips = $awsFipsEndpointStateEnabledType:T
            }
            if aws.ToBool(params.UseDualStack) {
                dualStack = $awsDualStackEndpointStateEnabledType:T
            }
            """,
            commonCodegenArgs,
            MapUtils.of(
                    "awsFipsEndpointStateType", SymbolUtils.createValueSymbolBuilder("FIPSEndpointState", AwsGoDependency.AWS_CORE).build(),
                    "awsFipsEndpointStateEnabledType", SymbolUtils.createValueSymbolBuilder("FIPSEndpointStateEnabled", AwsGoDependency.AWS_CORE).build(),
                    "awsDualStackEndpointStateEnabledType", SymbolUtils.createValueSymbolBuilder("DualStackEndpointStateEnabled", AwsGoDependency.AWS_CORE).build(),
                    "awsDualStackEndpointStateType", SymbolUtils.createValueSymbolBuilder("DualStackEndpointState", AwsGoDependency.AWS_CORE).build()
                    ));

    }

    private GoWriter.Writable generateResolveInvocation() {
        return goTemplate(
            """
                resolveEndpoint, err := l.legacyResolver.ResolveEndpoint($toString:T(params.Region), EndpointResolverOptions{
                    ResolvedRegion:       $toString:T(params.Region),
                    UseFIPSEndpoint:      fips,
                    UseDualStackEndpoint: dualStack,
                })
                if err != nil {
                    return endpoint, err
                }
            """,
            commonCodegenArgs,
            MapUtils.of(
                "toString", SymbolUtils.createValueSymbolBuilder("ToString", AwsGoDependency.AWS_CORE).build()
            ));
    }

    private GoWriter.Writable generateHostnameImmutableCheck() {
        return goTemplate(
            """
            if resolveEndpoint.HostnameImmutable {
                uriString := resolveEndpoint.URL
                uri, err := $parseUrl:T(uriString)
                if err != nil {
                    return endpoint, $fmtErrorf:T(\"Failed to parse uri: %s\", uriString)
                }

                return $smithyEndpointType:T{
                    URI: *uri,
                }, nil
            }
            """,
            commonCodegenArgs,
            MapUtils.of(
                "parseUrl", SymbolUtils.createValueSymbolBuilder("Parse", SmithyGoDependency.NET_URL).build()
            ));
    }

    private GoWriter.Writable generatePushdownToResolver() {
        return goTemplate(
            """
            if resolveEndpoint.Source == $endpointSourceMetadata:T {
                return l.resolver.ResolveEndpoint(ctx, params)
            }

            params = params.WithDefaults()
            params.Endpoint = &resolveEndpoint.URL

            return l.resolver.ResolveEndpoint(ctx, params)
            """,
            commonCodegenArgs,
            MapUtils.of(
                "endpointSourceMetadata", SymbolUtils.createValueSymbolBuilder("EndpointSourceServiceMetadata", AwsGoDependency.AWS_CORE).build()
            ));

    }


    private GoWriter.Writable generateCompatibleAdapter() {
        return goTemplate(
            """
                type isDefaultProvidedImplementation interface {
                    isDefaultProvidedImplementation()
                }

                type $compatibleResolverType:T struct {
                    EndpointResolverV2 $resolverType:T
                }

                func (n *$compatibleResolverType:T) isDefaultProvidedImplementation() {}

                func (n *$compatibleResolverType:T) ResolveEndpoint(region string, options EndpointResolverOptions) (endpoint $awsEndpointType:T, err error) {
                    $compatibleResolveMethodBody:W
                }
            """,
            commonCodegenArgs,
            MapUtils.of(
                    "compatibleResolverType", SymbolUtils.createValueSymbolBuilder(COMPATIBLE_ADAPTER_TYPE).build(),
                    "compatibleResolveMethodBody", generateCompatibleResolveMethodBody()
            ));
    }

    private GoWriter.Writable generateCompatibleResolveMethodBody() {
        return goTemplate(
            """
                reg := region
                fips := options.UseFIPSEndpoint
                if len(options.ResolvedRegion) > 0 {
                    reg = options.ResolvedRegion
                } else {
                    // $resolverInterfaceName:L needs to support pseudo-regions to maintain backwards-compatibility
                    // with the legacy $legacyResolverInterfaceName:L
                    reg, fips = mapPseudoRegion(region)
                }
                ctx := context.Background()
                resolved, err := n.EndpointResolverV2.ResolveEndpoint(ctx, $paramsType:T{
                    Region:       &reg,
                    UseFIPS:      $awsBoolType:T(fips == $awsFipsEndpointStateEnabledType:T),
                    UseDualStack: $awsBoolType:T(options.UseDualStackEndpoint == $awsDualStackEndpointStateEnabledType:T),
                })
                if err != nil {
                    return endpoint, err
                }

                endpoint = $awsEndpointType:T{
                    URL:               resolved.URI.String(),
                    HostnameImmutable: false,
                    Source:            $endpointSourceMetadata:T,
                }

                return endpoint, nil
            """,
            commonCodegenArgs,
            MapUtils.of(
                    "awsBoolType", SymbolUtils.createValueSymbolBuilder("Bool", AwsGoDependency.AWS_CORE).build(),
                    "awsFipsEndpointStateEnabledType", SymbolUtils.createValueSymbolBuilder("FIPSEndpointStateEnabled", AwsGoDependency.AWS_CORE).build(),
                    "awsDualStackEndpointStateEnabledType", SymbolUtils.createValueSymbolBuilder("DualStackEndpointStateEnabled", AwsGoDependency.AWS_CORE).build(),
                    "endpointSourceMetadata", SymbolUtils.createValueSymbolBuilder("EndpointSourceServiceMetadata", AwsGoDependency.AWS_CORE).build(),
                    "resolverInterfaceName", EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME,
                    "legacyResolverInterfaceName", EndpointGenerator.RESOLVER_INTERFACE_NAME
            ));
    }

    private GoWriter.Writable generatePseudoRegionUtility() {
        return goTemplate(
            """
            // Utility function to aid with translating pseudo-regions to classical regions
            // with the appropriate setting indicated by the pseudo-region
            func mapPseudoRegion(pr string) (region string, fips $awsFipsEndpointStateType:T) {
                const fipsInfix = \"-fips-\"
                const fipsPrefix = \"fips-\"
                const fipsSuffix = \"-fips\"

                if strings.Contains(pr, fipsInfix) ||
                    strings.Contains(pr, fipsPrefix) ||
                    strings.Contains(pr, fipsSuffix) {
                    region = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
                        pr, fipsInfix, "-"), fipsPrefix, ""), fipsSuffix, "")
                    fips = $awsFipsEndpointStateEnabledType:T
                } else {
                    region = pr
                }

                return region, fips
            }
            """,
            commonCodegenArgs,
            MapUtils.of(
                    "awsFipsEndpointStateEnabledType", SymbolUtils.createValueSymbolBuilder("FIPSEndpointStateEnabled", AwsGoDependency.AWS_CORE).build(),
                    "awsFipsEndpointStateType", SymbolUtils.createValueSymbolBuilder("FIPSEndpointState", AwsGoDependency.AWS_CORE).build()
            ));
    }


    private GoWriter.Writable generateFinalizeMethod() {
        return goTemplate(
            """
                func $finalizeMethodName:L(options *Options) {
                    // Check if the EndpointResolver was not user provided
                    // but is the SDK's default provided version.
                    _, ok := options.EndpointResolver.(isDefaultProvidedImplementation)
                    if options.EndpointResolverV2 == nil {
                        options.EndpointResolverV2 = $newResolverFuncName:L()
                    }
                    if ok {
                        // Nothing further to do
                        return
                    }

                    options.EndpointResolverV2 = &$legacyAdapterType:T{
                        legacyResolver: options.EndpointResolver,
                        resolver:       $newResolverFuncName:L(),
                    }
                }

            """,
            commonCodegenArgs,
            MapUtils.of(
                       "finalizeMethodName", FINALIZE_ENDPOINT_RESOLVER_V2,
                    "newResolverFuncName", EndpointResolutionGenerator.NEW_RESOLVER_FUNC_NAME
            ));

    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .configFields(SetUtils.of(
                                ConfigField.builder()
                                        .name(EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME)
                                        .type(SymbolUtils.createValueSymbolBuilder(EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME)
                                                .build())
                                        .documentation(String.format(
                                            """
                                            Resolves the endpoint used for a particular service. This should be used over the
                                            deprecated %s
                                            """,
                                            EndpointGenerator.RESOLVER_INTERFACE_NAME
                                        ))
                                        .withHelper(true)
                                        .build()
                        ))
                        .addConfigFieldResolver(ConfigFieldResolver.builder()
                                        .location(ConfigFieldResolver.Location.OPERATION)
                                        .target(ConfigFieldResolver.Target.FINALIZATION)
                                        .resolver(SymbolUtils.createValueSymbolBuilder(
                                                FINALIZE_ENDPOINT_RESOLVER_V2).build())
                                        .build())
                        .build());
    }
}
