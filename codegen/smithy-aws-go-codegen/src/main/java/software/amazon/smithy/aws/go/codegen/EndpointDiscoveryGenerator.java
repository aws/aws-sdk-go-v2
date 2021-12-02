package software.amazon.smithy.aws.go.codegen;

import java.util.ArrayList;
import java.util.Collection;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.aws.traits.clientendpointdiscovery.ClientDiscoveredEndpointTrait;
import software.amazon.smithy.aws.traits.clientendpointdiscovery.ClientEndpointDiscoveryIdTrait;
import software.amazon.smithy.aws.traits.clientendpointdiscovery.ClientEndpointDiscoveryTrait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ClientMember;
import software.amazon.smithy.go.codegen.integration.ClientMemberResolver;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.EnumTrait;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.SetUtils;

public class EndpointDiscoveryGenerator implements GoIntegration {

    // endpoint cache
    private static final String CLIENT_MEMBER_ENDPOINT_CACHE = "endpointCache";
    private static final Integer DEFAULT_ENDPOINT_CACHE_SIZE = 10;
    private static final Symbol CLIENT_MEMBER_ENDPOINT_CACHE_TYPE = SymbolUtils.createPointableSymbolBuilder(
            "EndpointCache", AwsGoDependency.SERVICE_INTERNAL_ENDPOINT_DISCOVERY).build();
    private static final Symbol ENDPOINT_CACHE_RESOLVER = SymbolUtils.createValueSymbolBuilder(
            "resolveEndpointCache").build();
    private static final Symbol NEW_ENDPOINT_CACHE = SymbolUtils.createValueSymbolBuilder(
            "NewEndpointCache", AwsGoDependency.SERVICE_INTERNAL_ENDPOINT_DISCOVERY).build();

    // middleware
    private static final String ADD_DISCOVER_ENDPOINT_MIDDLEWARE_FUNC_NAME_FORMAT = "addOp%sDiscoverEndpointMiddleware";
    private static final Symbol DISCOVERY_MIDDLEWARE = SymbolUtils.createValueSymbolBuilder("DiscoverEndpoint",
            AwsGoDependency.SERVICE_INTERNAL_ENDPOINT_DISCOVERY).build();

    // internal endpoint discovery members
    private static final Symbol DISCOVERY_ENDPOINT_OPTIONS = SymbolUtils.createPointableSymbolBuilder(
            "DiscoverEndpointOptions", AwsGoDependency.SERVICE_INTERNAL_ENDPOINT_DISCOVERY).build();
    private static final Symbol DISCOVERY_ENDPOINT_WEIGHTED_ADDRESS = SymbolUtils.createPointableSymbolBuilder(
            "WeightedAddress", AwsGoDependency.SERVICE_INTERNAL_ENDPOINT_DISCOVERY).build();
    private static final Symbol DISCOVERY_ENDPOINT_TYPE = SymbolUtils.createValueSymbolBuilder(
            "Endpoint", AwsGoDependency.SERVICE_INTERNAL_ENDPOINT_DISCOVERY).build();

    // operation specific discovery function to fetch Endpoint
    private static final String FETCH_DISCOVERED_ENDPOINT_FUNC_NAME_FORMAT = "fetchOp%sDiscoverEndpoint";
    // client specific endpoint discovery handler function name
    private static final String DISCOVERY_ENDPOINT_HANDLER_NAME = "handleEndpointDiscoveryFromService";
    // Symbol for Endpoint resolver interface specific to the client.
    private static Symbol ENDPOINT_RESOLVER_INTERFACE_NAME= SymbolUtils.createValueSymbolBuilder(
            "EndpointResolver").build();

    // EndpointDiscovery options
    private static final String ENDPOINT_DISCOVERY_OPTION = "EndpointDiscovery";
    private static final Symbol ENDPOINT_DISCOVERY_OPTION_TYPE = SymbolUtils.createValueSymbolBuilder(
            "EndpointDiscoveryOptions").build();

    // enable endpoint discovery
    private static final String ENABLE_ENDPOINT_DISCOVERY_OPTION = "EnableEndpointDiscovery";
    private static final String ENDPOINT_RESOLVER_USED_FOR_DISCOVERY = "EndpointResolverUsedForDiscovery";

    private static final Symbol ENDPOINT_DISCOVERY_ENABLE_STATE_TYPE = SymbolUtils.createValueSymbolBuilder(
            "EndpointDiscoveryEnableState", AwsGoDependency.AWS_CORE).build();
    private static final Symbol ENDPOINT_DISCOVERY_ENABLE_STATE_UNSET = SymbolUtils.createValueSymbolBuilder(
            "EndpointDiscoveryUnset", AwsGoDependency.AWS_CORE).build();
    private static final Symbol ENDPOINT_DISCOVERY_ENABLE_STATE_AUTO = SymbolUtils.createValueSymbolBuilder(
            "EndpointDiscoveryAuto", AwsGoDependency.AWS_CORE).build();

    // disable-https option for endpoint discovery middleware option
    private static final String DISABLE_HTTPS = "DisableHTTPS";

    // resolver for enable endpoint discovery
    private static final String ENABLE_ENDPOINT_DISCOVERY_OPTION_RESOLVER = "resolveEnableEndpointDiscovery";

    List<RuntimeClientPlugin> runtimeClientPlugins = new ArrayList<>();

    /**
     * Generates endpoint discovery options struct that contains options used to configure endpoint discovery.
     *
     * @param model   is the generation model
     * @param writer  is the GoWriter
     * @param service is the service for which endpoint discovery is performed.
     */
    private static void generateEndpointDiscoveryOptions(Model model, GoWriter writer, ServiceShape service) {
        if (!serviceSupportsEndpointDiscovery(model, service)) {
            return;
        }

        writer.write("// $T used to configure endpoint discovery", ENDPOINT_DISCOVERY_OPTION_TYPE);
        writer.openBlock("type $T struct {", "}", ENDPOINT_DISCOVERY_OPTION_TYPE, () -> {
            writer.writeDocs("Enables endpoint discovery");
            writer.write("$L $T", ENABLE_ENDPOINT_DISCOVERY_OPTION, ENDPOINT_DISCOVERY_ENABLE_STATE_TYPE);

            if (serviceSupportsCustomDiscoveryEndpoint(model, service)) {
                writer.write("");
                writer.writeDocs("Allows configuring an endpoint resolver to use when attempting an endpoint discovery "
                        + "api request.");
                writer.write("$L $T", ENDPOINT_RESOLVER_USED_FOR_DISCOVERY, ENDPOINT_RESOLVER_INTERFACE_NAME);
            }
        });
        writer.write("");
    }

    /**
     * Generates resolver function to default EnableEndpointDiscovery to AUTO.
     *
     * @param model   is the generation model
     * @param writer  is the GoWriter
     * @param service is the service for which endpoint discovery is performed.
     */
    private static void generateEnableEndpointDiscoveryResolver(Model model, GoWriter writer, ServiceShape service) {
        if (!serviceSupportsEndpointDiscovery(model, service)) {
            return;
        }

        writer.openBlock("func $L(o *Options) {", "}",
                ENABLE_ENDPOINT_DISCOVERY_OPTION_RESOLVER, () -> {
                    writer.write("if o.$L.$L != $T { return }", ENDPOINT_DISCOVERY_OPTION,
                            ENABLE_ENDPOINT_DISCOVERY_OPTION, ENDPOINT_DISCOVERY_ENABLE_STATE_UNSET);
                    writer.write("o.$L.$L = $T", ENDPOINT_DISCOVERY_OPTION,
                            ENABLE_ENDPOINT_DISCOVERY_OPTION, ENDPOINT_DISCOVERY_ENABLE_STATE_AUTO);
                });
        writer.write("");
    }

    /**
     * Generates a helper method to add DiscoverEndpoint middleware into the middleware stack.
     *
     * @param model          is the generation model
     * @param symbolProvider is the symbol provider used in generation context
     * @param writer         is the GoWriter
     * @param service        is the service for which endpoint discovery is performed.
     * @param operation      is the operation for which discover endpoint middleware is added.
     */
    private static void generateAddDiscoverEndpointMiddleware(
            Model model, SymbolProvider symbolProvider, GoWriter writer, ServiceShape service, OperationShape operation
    ) {
        if (!operationUsesEndpointDiscovery(model, service, operation)) {
            return;
        }

        Symbol stackSymbol = SymbolUtils.createPointableSymbolBuilder("Stack", SmithyGoDependency.SMITHY_MIDDLEWARE)
                .build();

        String helperName = generateAddDiscoverEndpointMiddlewareName(service, operation);
        String discoverOperationName = generateFetchDiscoverEndpointFuncName(service, operation);

        // Generate Middleware Adder Helper
        writer.openBlock("func $L(stack $P, o Options, c *Client) error {", "}", helperName, stackSymbol,
                () -> {
                    writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);

                    String closeBlock = String.format("}, \"%s\", middleware.After)",
                            EndpointGenerator.MIDDLEWARE_NAME);

                    writer.openBlock("return stack.Serialize.Insert(&$T{", closeBlock, DISCOVERY_MIDDLEWARE, () -> {
                        writer.openBlock("Options: []func($P){", "},", DISCOVERY_ENDPOINT_OPTIONS, () -> {
                            writer.openBlock("func (opt $P) {", "},", DISCOVERY_ENDPOINT_OPTIONS, () -> {
                                writer.write("opt.$L = o.EndpointOptions.$L", DISABLE_HTTPS, DISABLE_HTTPS);
                                writer.write("opt.Logger = o.Logger");

                                if (serviceSupportsCustomDiscoveryEndpoint(model, service)) {
                                    writer.write("opt.$L = o.$L.$L", ENDPOINT_RESOLVER_USED_FOR_DISCOVERY,
                                            ENDPOINT_DISCOVERY_OPTION, ENDPOINT_RESOLVER_USED_FOR_DISCOVERY);
                                }

                            });
                        });
                        writer.write("DiscoverOperation: c.$L,", discoverOperationName);
                        writer.write("EndpointDiscoveryEnableState: o.$L.$L,", ENDPOINT_DISCOVERY_OPTION, ENABLE_ENDPOINT_DISCOVERY_OPTION);
                        writer.write("EndpointDiscoveryRequired: $L,",
                                operationRequiresEndpointDiscovery(model, service, operation));
                    });
                });
        writer.write("");
    }

    /**
     * Generates client scoped handler that attempts endpoint discovery from the service.
     *
     * @param model          is the generation model
     * @param symbolProvider is the symbol provider used in generation context
     * @param writer         is the GoWriter
     * @param service        is the service for which endpoint discovery is performed.
     */
    private static void generateEndpointDiscoveryHandler(
            Model model, SymbolProvider symbolProvider, GoWriter writer, ServiceShape service
    ) {
        if (!serviceSupportsEndpointDiscovery(model, service)) {
            return;
        }

        OperationShape discoveryOperation = model.expectShape(getEndpointDiscoveryOperationId(model, service),
                OperationShape.class);
        Symbol discoveryOperationSymbol = symbolProvider.toSymbol(discoveryOperation);
        Shape discoveryOperationInput = model.expectShape(discoveryOperation.getInput().get());
        Symbol discoveryOperationInputSymbol = symbolProvider.toSymbol(discoveryOperationInput);

        writer.addUseImports(SmithyGoDependency.CONTEXT);
        writer.openBlock("func (c *Client) $L(ctx context.Context, input $P, key string, opt $T) ($T, error) {", "}",
                DISCOVERY_ENDPOINT_HANDLER_NAME, discoveryOperationInputSymbol, DISCOVERY_ENDPOINT_OPTIONS,
                DISCOVERY_ENDPOINT_TYPE, () -> {

                    if (serviceSupportsCustomDiscoveryEndpoint(model, service)) {
                        // check if endpoint resolver for endpoint discovery is of service-specific endpoint resolver type
                        writer.writeDocs("assert endpoint resolver interface is of expected type.");
                        writer.write("endpointResolver, ok := opt.$L.($T)", ENDPOINT_RESOLVER_USED_FOR_DISCOVERY,
                                ENDPOINT_RESOLVER_INTERFACE_NAME);
                        writer.openBlock("if opt.$L != nil && !ok { return $T{},", "}",
                                ENDPOINT_RESOLVER_USED_FOR_DISCOVERY, DISCOVERY_ENDPOINT_TYPE, () -> {
                            writer.addUseImports(SmithyGoDependency.FMT);
                            writer.write("fmt.Errorf(\"Unexpected endpoint resolver type %T provided for endpoint "
                                            + "discovery api call\", opt.$L)",
                                    ENDPOINT_RESOLVER_USED_FOR_DISCOVERY);
                        });
                        writer.write("");
                    }

                    // fetch endpoint via making discovery call
                    writer.openBlock("output, err := c.$T(ctx, input, func(o *Options) {", "})",
                            discoveryOperationSymbol, () -> {
                                writer.write("o.EndpointOptions.$L = opt.$L", DISABLE_HTTPS, DISABLE_HTTPS);
                                writer.write("o.Logger = opt.Logger");

                                // if service supports custom discovery endpoint
                                if (serviceSupportsCustomDiscoveryEndpoint(model, service)) {
                                    writer.openBlock("if endpointResolver != nil {", "}", () -> {
                                        writer.write("o.EndpointResolver = endpointResolver");
                                    });
                                }
                            });
                    writer.write("if err != nil { return $T{}, err }", DISCOVERY_ENDPOINT_TYPE);
                    writer.write("");

                    // create an endpoint
                    writer.write("endpoint := $T{}", DISCOVERY_ENDPOINT_TYPE);
                    writer.write("endpoint.Key = key");
                    writer.write("");

                    writer.openBlock("for _, e := range output.Endpoints {", "}", () -> {
                        writer.write("if e.Address == nil { continue }");
                        writer.write("address := *e.Address");

                        writer.write("");
                        writer.write("var scheme string");
                        writer.addUseImports(SmithyGoDependency.STRINGS);
                        writer.openBlock("if idx := strings.Index(address, \"://\"); idx != -1 {", "}", () -> {
                            writer.write("scheme = address[:idx]");
                        });

                        writer.openBlock("if len(scheme) == 0 {", "}", () -> {
                            writer.write("scheme = \"https\"");
                            writer.openBlock("if opt.$L {", "}", DISABLE_HTTPS, () -> {
                                writer.write("scheme = \"http\"");
                            });
                            writer.addUseImports(SmithyGoDependency.FMT);
                            writer.write("address = fmt.Sprintf(\"%s://%s\", scheme, address)");
                        });

                        writer.write("");
                        writer.write("cachedInMinutes := e.CachePeriodInMinutes");
                        writer.addUseImports(SmithyGoDependency.NET_URL);
                        writer.write("u, err := url.Parse(address)");
                        writer.write("if err != nil { continue }");

                        writer.write("");
                        writer.openBlock("addr := $T {", "}", DISCOVERY_ENDPOINT_WEIGHTED_ADDRESS, () -> {
                            writer.write("URL: u,");
                            writer.addUseImports(SmithyGoDependency.TIME);
                            writer.write(
                                    "Expired: time.Now().Add(time.Duration(cachedInMinutes) * time.Minute).Round(0),");
                        });
                        writer.write("endpoint.Add(addr)");
                    });
                    writer.write("");
                    writer.write("c.$L.Add(endpoint)", CLIENT_MEMBER_ENDPOINT_CACHE);
                    writer.write("return endpoint, nil");
                });
    }

    //========================= Generation function ==============

    /**
     * Generates a function to fetch a discovered endpoint.
     *
     * @param model          is the generation model
     * @param symbolProvider is the symbol provider used in generation context
     * @param writer         is the GoWriter
     * @param service        is the service for which endpoint discovery is performed.
     * @param operation      is the operation for which discover endpoint middleware is added.
     */
    private static void generateFetchDiscoveredEndpointFunction(
            Model model, SymbolProvider symbolProvider, GoWriter writer, ServiceShape service, OperationShape operation
    ) {
        if (!operationUsesEndpointDiscovery(model, service, operation)) {
            return;
        }

        // Generate operation specific fetch discovered endpoint function
        String fetchDiscoveredEndpointFuncName = generateFetchDiscoverEndpointFuncName(service, operation);
        String operationName = operation.getId().getName(service);
        String sdkId = service.expectTrait(ServiceTrait.class).getSdkId();

        Shape inputShape = model.expectShape(operation.getInput().get());
        OperationShape discoveryOperation = model.expectShape(getEndpointDiscoveryOperationId(model, service),
                OperationShape.class);
        Shape discoveryOperationInput = model.expectShape(discoveryOperation.getInput().get());
        Symbol discoveryOperationInputSymbol = symbolProvider.toSymbol(discoveryOperationInput);
        String DISCOVERY_OPERATION_INPUT_NAME = "discoveryOperationInput";

        writer.addUseImports(SmithyGoDependency.CONTEXT);
        writer.openBlock(
                "func (c *Client) $L(ctx context.Context, input interface{}, optFns ...func($P)) ($T, error) {", "}",
                fetchDiscoveredEndpointFuncName, DISCOVERY_ENDPOINT_OPTIONS, DISCOVERY_ENDPOINT_WEIGHTED_ADDRESS,
                () -> {
                    writer.write("in, ok := input.($P)", symbolProvider.toSymbol(inputShape));
                    writer.openBlock("if !ok {", "}", () -> {
                        writer.addUseImports(SmithyGoDependency.FMT);
                        writer.write("return $T{}, fmt.Errorf(\"unknown input type %T\", input)",
                                DISCOVERY_ENDPOINT_WEIGHTED_ADDRESS);
                    });
                    writer.write("_ = in");
                    writer.write("");

                    // build identifier map
                    String IDENTIFIER_MAP = "identifierMap";
                    writer.write("$L := make(map[string]string, 0)", IDENTIFIER_MAP);

                    for (MemberShape member : getMembersUsedAsIdForDiscovery(model, service, operation)) {
                        String memberName = member.getMemberName();
                        Shape targetShape = model.expectShape(member.getTarget());
                        String assignFormat = String.format("$L[$S] = %s",
                                targetShape.hasTrait(EnumTrait.class) ? "string(in.$L)" : "in.$L");
                        writer.write(assignFormat, IDENTIFIER_MAP, memberName, memberName);
                    }

                    writer.write("");
                    writer.write("key := fmt.Sprintf(\"$L.%v\", $L)", sdkId, IDENTIFIER_MAP);

                    writer.write("");
                    writer.openBlock("if v, ok := c.$L.Get(key); ok {", "}", CLIENT_MEMBER_ENDPOINT_CACHE, () -> {
                        writer.write("return v, nil");
                    }).write("");

                    // build input for endpoint discovery call
                    writer.openBlock("$L := &$T{", "}", DISCOVERY_OPERATION_INPUT_NAME,
                            discoveryOperationInputSymbol, () -> {
                                for (MemberShape member : discoveryOperationInput.members()) {
                                    String memberName = member.getMemberName();
                                    if (memberName.equalsIgnoreCase("Operation")) {
                                        writer.write("$T : $S", member, operationName);
                                    }

                                    if (memberName.equalsIgnoreCase("Identifiers")) {
                                        writer.write("$T : $L", member, IDENTIFIER_MAP);
                                    }
                                }
                            }).write("");

                    // resolve options
                    writer.write("opt := $T{}", DISCOVERY_ENDPOINT_OPTIONS);
                    writer.openBlock("for _, fn := range optFns {", "}", () -> {
                        writer.write("fn(&opt)");
                    }).write("");

                    // if discovery not required, then spin up a unblocking go routine
                    if (!operationRequiresEndpointDiscovery(model, service, operation)) {
                        writer.write("go c.$L(ctx, $L, key, opt)", DISCOVERY_ENDPOINT_HANDLER_NAME,
                                DISCOVERY_OPERATION_INPUT_NAME);
                        writer.write("return $T{}, nil", DISCOVERY_ENDPOINT_WEIGHTED_ADDRESS);
                        return;
                    }

                    writer.write("endpoint, err := c.$L(ctx, $L, key, opt)", DISCOVERY_ENDPOINT_HANDLER_NAME,
                            DISCOVERY_OPERATION_INPUT_NAME);
                    writer.write("if err != nil { return $T{}, err }", DISCOVERY_ENDPOINT_WEIGHTED_ADDRESS);
                    writer.write("");

                    writer.write("weighted, ok := endpoint.GetValidAddress()");
                    writer.write(
                            "if !ok { return $T{}, fmt.Errorf(\"no valid endpoint address returned by the endpoint discovery api\")}",
                            DISCOVERY_ENDPOINT_WEIGHTED_ADDRESS);
                    writer.write("return weighted, nil");
                });
    }

    // generates an endpoint cache resolver function
    private static void generateEndpointCacheResolver(Model model, GoWriter writer, ServiceShape serviceShape) {
        if (!serviceSupportsEndpointDiscovery(model, serviceShape)) {
            return;
        }

        writer.writeDocs("resolves endpoint cache on client");
        writer.openBlock("func $T(c *Client) {", "}", ENDPOINT_CACHE_RESOLVER, () -> {
            writer.write("c.$L = $T($L)", CLIENT_MEMBER_ENDPOINT_CACHE, NEW_ENDPOINT_CACHE,
                    DEFAULT_ENDPOINT_CACHE_SIZE);
        });
        writer.write("");
    }

    // Returns name for operation specific function that fetches a discovered endpoint.
    private static final String generateFetchDiscoverEndpointFuncName(ServiceShape service, OperationShape operation) {
        return String.format(FETCH_DISCOVERED_ENDPOINT_FUNC_NAME_FORMAT, operation.getId().getName(service));
    }

    private static String generateAddDiscoverEndpointMiddlewareName(ServiceShape service, OperationShape operation) {
        return String.format(ADD_DISCOVER_ENDPOINT_MIDDLEWARE_FUNC_NAME_FORMAT, operation.getId().getName(service));
    }

    // Returns true if service supports endpoint discovery. This may or may not be true for each operation.
    public static boolean serviceSupportsEndpointDiscovery(Model model, ServiceShape service) {
        return service.hasTrait(ClientEndpointDiscoveryTrait.class);
    }

    // Returns true if operation may use endpoint discovery internally
    private static boolean operationUsesEndpointDiscovery(Model model, ServiceShape service, OperationShape operation) {
        return operation.hasTrait(ClientDiscoveredEndpointTrait.class);
    }

    // Returns true if operation requires endpoint discovery
    private static boolean operationRequiresEndpointDiscovery(
            Model model, ServiceShape service, OperationShape operation
    ) {
        if (!operation.hasTrait(ClientDiscoveredEndpointTrait.class)) return false;

        return operation.expectTrait(ClientDiscoveredEndpointTrait.class).isRequired();
    }


    // ============================Trait info retriever/analyzer==============

    // Return operation shape id of the operation to use for making endpoint discovery call.
    // eg. DescribeEndpoints
    private static ShapeId getEndpointDiscoveryOperationId(Model model, ServiceShape service) {
        return service.expectTrait(ClientEndpointDiscoveryTrait.class).getOperation();
    }

    // Returns input member shapes marked with clientEndpointDiscoveryId trait.
    private static Set<MemberShape> getMembersUsedAsIdForDiscovery(
            Model model, ServiceShape service, OperationShape operation
    ) {
        Set<MemberShape> shapes = new HashSet<>();
        StructureShape input = model.expectShape(operation.getInput().get(), StructureShape.class);
        for (MemberShape member : input.members()) {
            if (member.hasTrait(ClientEndpointDiscoveryIdTrait.class)) shapes.add(member);
        }
        return shapes;
    }

    // Returns true if service supports providing custom endpoint used with the endpoint discovery api call
    private static boolean serviceSupportsCustomDiscoveryEndpoint(Model model, ServiceShape service) {
        Set<String> ServicesSupportingCustomDiscoveryEndpoint = SetUtils.of(
                "Timestream Query", "Timestream Write");

        String sdkId = service.expectTrait(ServiceTrait.class).getSdkId();
        for (String id : ServicesSupportingCustomDiscoveryEndpoint) {
            if (sdkId.equalsIgnoreCase(id)) return true;
        }
        return false;
    }

    @Override
    public void processFinalizedModel(GoSettings settings, Model model) {
        ServiceShape service = settings.getService(model);

        // add client member field and resolver for endpoint cache
        runtimeClientPlugins.add(RuntimeClientPlugin.builder()
                .servicePredicate(EndpointDiscoveryGenerator::serviceSupportsEndpointDiscovery)
                .addClientMember(ClientMember.builder()
                        .name(CLIENT_MEMBER_ENDPOINT_CACHE)
                        .type(CLIENT_MEMBER_ENDPOINT_CACHE_TYPE)
                        .documentation("cache used to store discovered endpoints")
                        .build())
                .addClientMemberResolver(ClientMemberResolver.builder()
                        .resolver(ENDPOINT_CACHE_RESOLVER)
                        .build())
                .build());


        for (OperationShape operation : TopDownIndex.of(model).getContainedOperations(service)) {
            String helperFuncName = generateAddDiscoverEndpointMiddlewareName(service, operation);

            Collection<Symbol> middlewareArgs = ListUtils.of(
                    SymbolUtils.createValueSymbolBuilder("options").build(),
                    SymbolUtils.createValueSymbolBuilder("c").build());

            runtimeClientPlugins.add(
                    RuntimeClientPlugin.builder()
                            .operationPredicate((m, s, o) -> operation.equals(o)
                                    && EndpointDiscoveryGenerator.operationUsesEndpointDiscovery(m, s, o))
                            .registerMiddleware(MiddlewareRegistrar.builder()
                                    .resolvedFunction(SymbolUtils.createValueSymbolBuilder(helperFuncName).build())
                                    .functionArguments(middlewareArgs)
                                    .build())
                            .build());
        }
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        // Add enable endpoint discovery config option and resolver.
        runtimeClientPlugins.add(
                RuntimeClientPlugin.builder()
                        .servicePredicate(EndpointDiscoveryGenerator::serviceSupportsEndpointDiscovery)
                        .addConfigField(ConfigField.builder()
                                .name(ENDPOINT_DISCOVERY_OPTION)
                                .type(ENDPOINT_DISCOVERY_OPTION_TYPE)
                                .documentation("Allows configuring endpoint discovery")
                                .build())
                        .addConfigFieldResolver(ConfigFieldResolver.builder()
                                .location(ConfigFieldResolver.Location.CLIENT)
                                .target(ConfigFieldResolver.Target.INITIALIZATION)
                                .resolver(
                                        SymbolUtils.createValueSymbolBuilder(ENABLE_ENDPOINT_DISCOVERY_OPTION_RESOLVER)
                                                .build())
                                .build())
                        .build()
        );

        return runtimeClientPlugins;
    }

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
            generateEndpointCacheResolver(model, writer, service);
            generateEndpointDiscoveryOptions(model, writer, service);
            generateEnableEndpointDiscoveryResolver(model, writer, service);
            generateEndpointDiscoveryHandler(model, symbolProvider, writer, service);
        });

        // generate code specific to the operation
        for (OperationShape operation : TopDownIndex.of(model).getContainedOperations(service)) {
            goDelegator.useShapeWriter(operation, writer -> {
                generateAddDiscoverEndpointMiddleware(model, symbolProvider, writer, service, operation);
                generateFetchDiscoveredEndpointFunction(model, symbolProvider, writer, service, operation);
            });
        }
    }
}
