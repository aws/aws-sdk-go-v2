package software.amazon.smithy.aws.go.codegen;


import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.TriConsumer;
import software.amazon.smithy.go.codegen.endpoints.EndpointResolutionGenerator;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.MapUtils;


import java.util.Map;
import java.util.function.Consumer;

public class AwsEndpointResolverAdapterGenerator implements GoIntegration {

    private static final String LEGACY_ADAPTER_TYPE = "legacyEndpointResolverAdapter";

    private Map<String, Object> commonCodegenArgs;

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            TriConsumer<String, String, Consumer<GoWriter>> writerFactory
    ) {
        // String serviceId = settings.getService(model).expectTrait(ServiceTrait.class).getSdkId();

        this.commonCodegenArgs = MapUtils.of(
                "goContext", SymbolUtils.createValueSymbolBuilder("Context", SmithyGoDependency.CONTEXT).build(),
                "legacyAdapterType", SymbolUtils.createValueSymbolBuilder(LEGACY_ADAPTER_TYPE).build(),
                "legacyResolverType", SymbolUtils.createValueSymbolBuilder(EndpointGenerator.RESOLVER_INTERFACE_NAME).build(),
                "resolverType", SymbolUtils.createValueSymbolBuilder(EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME).build(),
                "paramsType", SymbolUtils.createValueSymbolBuilder(EndpointResolutionGenerator.PARAMETERS_TYPE_NAME).build(),
                "endpointType", SymbolUtils.createValueSymbolBuilder("Endpoint", SmithyGoDependency.SMITHY_ENDPOINTS).build(),
                "fmtErrorf", SymbolUtils.createValueSymbolBuilder("Errorf", SmithyGoDependency.FMT).build());

        // var content = new GoWriter.ChainWritable()
        //         .add(generateLegacyAdapter())
        //         .add(generateCompatibleAdapter())
        //         .add(generateFinalizeMethod(newResolverFuncName))
        //         .compose();

        var content = generateLegacyAdapter();
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

            func (l *$legacyAdapterType:T) ResolveEndpoint(ctx $goContext:T, params $paramsType:T) (endpoint $endpointType:T, err error) {
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
                    "awsDualStackEndpointStateType", SymbolUtils.createValueSymbolBuilder("DualStackEndpointState", AwsGoDependency.AWS_CORE).build(),
                    "awsFipsEndpointStateEnabledType", SymbolUtils.createValueSymbolBuilder("FIPSEndpointStateEnabled", AwsGoDependency.AWS_CORE).build(),
                    "awsDualStackEndpointStateEnabledType", SymbolUtils.createValueSymbolBuilder("DualStackEndpointStateEnabled", AwsGoDependency.AWS_CORE).build()
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

                return $endpointType:T{
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

            params.UseDualStack = nil
            params.UseFIPS = nil
            params.Endpoint = &resolveEndpoint.URL

            return l.resolver.ResolveEndpoint(ctx, params)
            """,
            commonCodegenArgs,
            MapUtils.of(
                "endpointSourceMetadata", SymbolUtils.createValueSymbolBuilder("EndpointSourceServiceMetadata", AwsGoDependency.AWS_CORE).build()
            ));
    }


    // private GoWriter.Writable generateCompatibleAdapter() {

    //     // generate type
    //     // generate isClientProvidedImplementation method
    //     // generate ResolveEndpointMethod
    //     //

    //     // external class: modify resolveDefaultEndpointConfiguration in EndpointGenerator.java


    //     return goTemplate(
    //         """
    //             type $clientProvidedImplType interface {
    //                 $clientProvidedImplMethod()
    //             }

    //             type $compatibleResolverType struct {
    //                 EndpointResolverV2 EndpointResolverV2
    //             }

    //             func (n *$compatibleResolverType) $clientProvidedImplMethod() {}

    //             func (n *$compatibleResolverType) ResolveEndpoint(region string, options EndpointResolverOptions) (endpoint aws.Endpoint, err error) {
    //                 $compatibleResolveMethodBody:W
    //             }
    //         """,
    //         commonCodegenArgs,
    //         MapUtils.of(
    //                 "compatibleResolveMethodBody", generateCompatibleResolveMethodBody()),
    //         overriddenArgs
    //     );
    // }

    // private GoWriter.Writable generateCompatibleResolveMethodBody() {
    // }




    // private GoWriter.Writable generateFinalizeMethod(String newResolverFuncName) {
    //     return goTemplate(
    //         """
    //             func finalizeEndpointResolverV2(options *Options) {
    //                 // check options.EndpointResolver is nil

    //                 // Check if the EndpointResolver was not user provided, but out default provided version
    //                 _, ok := options.EndpointResolver.(isClientProvidedImplementation)
    //                 if options.EndpointResolverV2 == nil {
    //                     options.EndpointResolverV2 = NewDefaultEndpointResolverV2()
    //                 }
    //                 if ok {
    //                     // Nothing further to do
    //                     return
    //                 }

    //                 options.EndpointResolverV2 = &legacyEndpointResolverAdaptor{
    //                     legacyResolver: options.EndpointResolver,
    //                     resolver:       NewDefaultEndpointResolverV2(),
    //                 }
    //             }

    //         """,
    //         commonCodegenArgs,
    //         MapUtils.of(
    //                 "", ),
    //         overriddenArgs

    //     )

    // }
}
