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

import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.TreeMap;
import java.util.function.Consumer;
import java.util.stream.Collectors;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.TriConsumer;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.node.Node;
import software.amazon.smithy.model.node.ObjectNode;
import software.amazon.smithy.model.node.StringNode;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.IoUtils;
import software.amazon.smithy.utils.ListUtils;

/**
 * Writes out a file that resolves endpoints using endpoints.json, but the
 * created resolver resolves endpoints for a single service.
 */
final class EndpointGenerator implements Runnable {
    public static final String MIDDLEWARE_NAME = "ResolveEndpoint";
    public static final String ADD_MIDDLEWARE_HELPER_NAME = String.format("Add%sMiddleware", MIDDLEWARE_NAME);
    public static final String RESOLVER_INTERFACE_NAME = "EndpointResolver";
    public static final String RESOLVER_FUNC_NAME = "EndpointResolverFunc";
    public static final String RESOLVER_OPTIONS = "ResolverOptions";
    public static final String CLIENT_CONFIG_RESOLVER = "resolveDefaultEndpointConfiguration";

    private static final int ENDPOINT_MODEL_VERSION = 3;
    private static final String RESOLVER_CONSTRUCTOR_NAME = "NewDefaultEndpointResolver";
    private static final String INTERNAL_ENDPOINT_PACKAGE = "internal/endpoints";
    private static final String INTERNAL_RESOLVER_NAME = "Resolver";
    private static final String INTERNAL_RESOLVER_OPTIONS_NAME = "Options";
    private static final String INTERNAL_ENDPOINTS_DATA_NAME = "defaultPartitions";
    private static final List<ResolveConfigField> resolveConfigFields = ListUtils.of(
            ResolveConfigField.builder()
                    .name("DisableHTTPS")
                    .type(SymbolUtils.createValueSymbolBuilder("bool").build())
                    .shared(true)
                    .build()
    );

    private final GoSettings settings;
    private final Model model;
    private final TriConsumer<String, String, Consumer<GoWriter>> writerFactory;
    private final ServiceShape serviceShape;
    private final ObjectNode endpointData;
    private final String endpointPrefix;
    private final Map<String, Partition> partitions = new TreeMap<>();

    EndpointGenerator(
            GoSettings settings,
            Model model,
            TriConsumer<String, String, Consumer<GoWriter>> writerFactory
    ) {
        this.settings = settings;
        this.model = model;
        this.writerFactory = writerFactory;
        serviceShape = settings.getService(model);
        this.endpointPrefix = getEndpointPrefix(serviceShape);
        this.endpointData = Node.parse(IoUtils.readUtf8Resource(getClass(), "endpoints.json")).expectObjectNode();
        validateVersion();
        loadPartitions();
    }

    private void validateVersion() {
        int version = endpointData.expectNumberMember("version").getValue().intValue();
        if (version != ENDPOINT_MODEL_VERSION) {
            throw new CodegenException("Invalid endpoints.json version. Expected version 3, found " + version);
        }
    }

    // Get service's endpoint prefix from a known list. If not found, fallback to ArnNamespace
    private String getEndpointPrefix(ServiceShape service) {
        ObjectNode endpointPrefixData = Node.parse(IoUtils.readUtf8Resource(getClass(), "endpoint-prefix.json"))
                .expectObjectNode();
        ServiceTrait serviceTrait = service.getTrait(ServiceTrait.class)
                .orElseThrow(() -> new CodegenException("No service trait found on " + service.getId()));
        return endpointPrefixData.getStringMemberOrDefault(serviceTrait.getSdkId(), serviceTrait.getArnNamespace());
    }

    private void loadPartitions() {
        List<ObjectNode> partitionObjects = endpointData
                .expectArrayMember("partitions")
                .getElementsAs(Node::expectObjectNode);

        for (ObjectNode partition : partitionObjects) {
            String partitionName = partition.expectStringMember("partition").getValue();
            partitions.put(partitionName, new Partition(partition, partitionName));
        }
    }

    @Override
    public void run() {
        writerFactory.accept(INTERNAL_ENDPOINT_PACKAGE + "/endpoints.go", getInternalEndpointImportPath(), (writer) -> {
            generateInternalResolverImplementation(writer);
            generateInternalEndpointsModel(writer);
        });
        writerFactory.accept("endpoints.go", settings.getModuleName(), writer -> {
            generatePublicResolverTypes(writer);
            generateMiddleware(writer);
        });
        writerFactory.accept(INTERNAL_ENDPOINT_PACKAGE + "/endpoints_test.go",
                getInternalEndpointImportPath(), (writer) -> {
                    writer.addUseImports(SmithyGoDependency.TESTING);
                    writer.openBlock("func TestRegexCompile(t *testing.T) {", "}", () -> {
                        writer.write("_ = $T", getInternalEndpointsSymbol(INTERNAL_ENDPOINTS_DATA_NAME, false).build());
                    });
                });
    }

    private void generateMiddleware(GoWriter writer) {
        // Generate middleware definition
        GoStackStepMiddlewareGenerator middleware = GoStackStepMiddlewareGenerator.createSerializeStepMiddleware(
                MIDDLEWARE_NAME, MIDDLEWARE_NAME);
        middleware.writeMiddleware(writer, this::generateMiddlewareResolverBody,
                this::generateMiddlewareStructureMembers);

        Symbol stackSymbol = SymbolUtils.createPointableSymbolBuilder("Stack", SmithyGoDependency.SMITHY_MIDDLEWARE)
                .build();
        Symbol optionsSymbol = SymbolUtils.createValueSymbolBuilder(String.format("%sMiddlewareOptions",
                MIDDLEWARE_NAME)).build();

        // Generate Middleware options interface
        writer.openBlock("type $T interface {", "}", optionsSymbol, () -> {
            writer.write("GetEndpointResolver() $L", RESOLVER_INTERFACE_NAME);
            writer.write("GetEndpointOptions() $L", RESOLVER_OPTIONS);
        });
        writer.write("");

        // Generate Middleware Adder Helper
        writer.openBlock("func $L(stack $P, options $T) {", "}", ADD_MIDDLEWARE_HELPER_NAME, stackSymbol,
                optionsSymbol, () -> {
                    writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);
                    writer.openBlock("stack.Serialize.Add(&$T{", "}, middleware.After)",
                            middleware.getMiddlewareSymbol(), () -> {
                                writer.write("Resolver: options.GetEndpointResolver(),");
                                writer.write("Options: options.GetEndpointOptions(),");
                            });
                });
        writer.write("");
        // Generate Middleware Remover Helper
        writer.openBlock("func Remove$LMiddleware(stack $P) error {", "}", middleware.getMiddlewareSymbol(),
                stackSymbol, () -> {
                    writer.write("return stack.Serialize.Remove((&$T{}).ID())", middleware.getMiddlewareSymbol());
                });
    }

    private void generateMiddlewareResolverBody(GoStackStepMiddlewareGenerator g, GoWriter w) {
        w.addUseImports(SmithyGoDependency.FMT);
        w.addUseImports(SmithyGoDependency.NET_URL);
        w.addUseImports(AwsGoDependency.AWS_MIDDLEWARE);
        w.addUseImports(SmithyGoDependency.SMITHY_HTTP_TRANSPORT);

        w.write("req, ok := in.Request.(*smithyhttp.Request)");
        w.openBlock("if !ok {", "}", () -> {
            w.write("return out, metadata, fmt.Errorf(\"unknown transport type %T\", in.Request)");
        });
        w.write("");
        w.openBlock("if m.Resolver == nil {", "}", () -> {
            w.write("return out, metadata, fmt.Errorf(\"expected endpoint resolver to not be nil\")");
        });
        w.write("");
        w.write("var endpoint $T", SymbolUtils.createValueSymbolBuilder("Endpoint", AwsGoDependency.AWS_CORE)
                .build());
        w.write("endpoint, err = m.Resolver.ResolveEndpoint(awsmiddleware.GetRegion(ctx), m.Options)");
        w.openBlock("if err != nil {", "}", () -> {
            w.write("return out, metadata, fmt.Errorf(\"failed to resolve service endpoint\")");
        });
        w.write("");
        w.write("req.URL, err = url.Parse(endpoint.URL)");
        w.openBlock("if err != nil {", "}", () -> {
            w.write("return out, metadata, fmt.Errorf(\"failed to parse endpoint URL: %w\", err)");
        });
        w.write("");
        w.openBlock("if len(awsmiddleware.GetSigningName(ctx)) == 0 {", "}", () -> {
            w.write("signingName := endpoint.SigningName");
            w.openBlock("if len(signingName) == 0 {", "}", () -> {
                w.write("signingName = $S", serviceShape.expectTrait(ServiceTrait.class).getArnNamespace());
            });
            w.write("ctx = awsmiddleware.SetSigningName(ctx, signingName)");
        });
        w.write("ctx = awsmiddleware.SetSigningRegion(ctx, endpoint.SigningRegion)");
        w.write("");
        w.write("return next.HandleSerialize(ctx, in)");
    }

    private void generateMiddlewareStructureMembers(GoStackStepMiddlewareGenerator g, GoWriter w) {
        w.write("Resolver $L", RESOLVER_INTERFACE_NAME);
        w.write("Options $L", RESOLVER_OPTIONS);
    }

    private Symbol.Builder getInternalEndpointsSymbol(String symbolName, boolean pointable) {
        Symbol.Builder builder;
        if (pointable) {
            builder = SymbolUtils.createPointableSymbolBuilder(symbolName);
        } else {
            builder = SymbolUtils.createValueSymbolBuilder(symbolName);
        }
        return builder.namespace(getInternalEndpointImportPath(), "/")
                .putProperty(SymbolUtils.NAMESPACE_ALIAS, "internalendpoints");
    }

    private String getInternalEndpointImportPath() {
        return settings.getModuleName() + "/" + INTERNAL_ENDPOINT_PACKAGE;
    }

    private void generatePublicResolverTypes(GoWriter writer) {
        Symbol awsEndpointSymbol = SymbolUtils.createValueSymbolBuilder("Endpoint", AwsGoDependency.AWS_CORE).build();
        Symbol internalEndpointsSymbol = getInternalEndpointsSymbol(INTERNAL_RESOLVER_NAME, true).build();

        Symbol resolverOptionsSymbol = SymbolUtils.createPointableSymbolBuilder(RESOLVER_OPTIONS).build();
        writer.writeDocs(String.format("%s is the service endpoint resolver options",
                resolverOptionsSymbol.getName()));
        writer.write("type $T = $T", resolverOptionsSymbol, getInternalEndpointsSymbol(INTERNAL_RESOLVER_OPTIONS_NAME,
                false).build());
        writer.write("");

        // Generate Resolver Interface
        writer.writeDocs(String.format("%s interface for resolving service endpoints.", RESOLVER_INTERFACE_NAME));
        writer.openBlock("type $L interface {", "}", RESOLVER_INTERFACE_NAME, () -> {
            writer.write("ResolveEndpoint(region string, options $T) ($T, error)", resolverOptionsSymbol,
                    awsEndpointSymbol);
        });
        writer.write("var _ $L = &$T{}", RESOLVER_INTERFACE_NAME, internalEndpointsSymbol);
        writer.write("");

        // Resolver Constructor
        writer.writeDocs(String.format("%s constructs a new service endpoint resolver", RESOLVER_CONSTRUCTOR_NAME));
        writer.openBlock("func $L() $P {", "}", RESOLVER_CONSTRUCTOR_NAME, internalEndpointsSymbol, () -> {
            writer.write("return $T()", getInternalEndpointsSymbol("New", false)
                    .build());
        });

        // Generate resolver function creator
        writer.writeDocs(String.format("%s is a helper utility that wraps a function so it satisfies the %s "
                + "interface. This is useful when you want to add additional endpoint resolving logic, or stub out "
                + "specific endpoints with custom values.", RESOLVER_FUNC_NAME, RESOLVER_INTERFACE_NAME));
        writer.write("type $L func(region string, options $T) ($T, error)",
                RESOLVER_FUNC_NAME, resolverOptionsSymbol, awsEndpointSymbol);

        writer.openBlock("func (fn $L) ResolveEndpoint(region string, options $T) ($T, error) {", "}",
                RESOLVER_FUNC_NAME, resolverOptionsSymbol, awsEndpointSymbol, () -> {
            writer.write("return fn(region, options)");
        }).write("");

        // Generate Client Options Configuration Resolver
        writer.openBlock("func $L(o $P) {", "}", CLIENT_CONFIG_RESOLVER,
                SymbolUtils.createPointableSymbolBuilder("Options").build(), () -> {
                    writer.openBlock("if o.EndpointResolver != nil {", "}", () -> writer.write("return"));
                    writer.write("o.EndpointResolver = $L()", RESOLVER_CONSTRUCTOR_NAME);
                });
    }

    private void generateInternalResolverImplementation(GoWriter writer) {
        Symbol awsEndpointSymbol = SymbolUtils.createValueSymbolBuilder("Endpoint", AwsGoDependency.AWS_CORE).build();

        // Options
        Symbol resolverOptionsSymbol = SymbolUtils.createPointableSymbolBuilder(INTERNAL_RESOLVER_OPTIONS_NAME).build();
        writer.writeDocs(String.format("%s is the endpoint resolver configuration options",
                resolverOptionsSymbol.getName()));
        writer.openBlock("type $T struct {", "}", resolverOptionsSymbol, () -> {
            resolveConfigFields.forEach(field -> {
                writer.write("$L $T", field.getName(), field.getType());
            });
        });
        writer.write("");

        // Resolver
        Symbol resolverImplSymbol = SymbolUtils.createPointableSymbolBuilder(INTERNAL_RESOLVER_NAME).build();
        writer.writeDocs(String.format("%s %s endpoint resolver", resolverImplSymbol.getName(),
                serviceShape.expectTrait(ServiceTrait.class).getSdkId()));
        writer.openBlock("type $T struct {", "}", resolverImplSymbol, () -> {
            writer.write("partitions $T", SymbolUtils.createValueSymbolBuilder("Partitions",
                    AwsGoDependency.AWS_ENDPOINTS).build());
        });
        writer.write("");
        writer.writeDocs("ResolveEndpoint resolves the service endpoint for the given region and options");
        writer.openBlock("func (r $P) ResolveEndpoint(region string, options $T) ($T, error) {", "}",
                resolverImplSymbol, resolverOptionsSymbol, awsEndpointSymbol, () -> {
                    Symbol sharedOptions = SymbolUtils.createPointableSymbolBuilder("Options",
                            AwsGoDependency.AWS_ENDPOINTS).build();
                    writer.openBlock("opt := $T{", "}", sharedOptions, () -> {
                        resolveConfigFields.stream().filter(ResolveConfigField::isShared).forEach(field -> {
                            writer.write("$L: options.$L,", field.getName(), field.getName());
                        });
                    });
                    writer.write("return r.partitions.ResolveEndpoint(region, opt)");
                });
        writer.write("");
        writer.writeDocs(String.format("New returns a new %s", resolverImplSymbol.getName()));
        writer.openBlock("func New() *$T {", "}", resolverImplSymbol, () -> writer.openBlock("return &$T{", "}",
                resolverImplSymbol, () -> {
                    writer.write("partitions: $L,", INTERNAL_ENDPOINTS_DATA_NAME);
                }));
    }

    private void generateInternalEndpointsModel(GoWriter writer) {
        writer.addUseImports(AwsGoDependency.AWS_ENDPOINTS);

        Symbol partitionsSymbol = SymbolUtils.createPointableSymbolBuilder("Partitions", AwsGoDependency.AWS_ENDPOINTS)
                .build();
        writer.openBlock("var $L = $T{", "}", INTERNAL_ENDPOINTS_DATA_NAME, partitionsSymbol, () -> {
            List<Partition> entries = partitions.entrySet().stream()
                    .sorted((x, y) -> {
                        // Always sort standard aws partition first
                        if (x.getKey().equals("aws")) {
                            return -1;
                        }
                        return x.getKey().compareTo(y.getKey());
                    }).map(Map.Entry::getValue).collect(Collectors.toList());

            entries.forEach(entry -> {
                writer.openBlock("{", "},", () -> writePartition(writer, entry));
            });
        });
    }

    private void writePartition(GoWriter writer, Partition partition) {
        writer.write("ID: $S,", partition.getId());
        Symbol endpointSymbol = SymbolUtils.createValueSymbolBuilder("Endpoint",
                AwsGoDependency.AWS_ENDPOINTS).build();
        writer.openBlock("Defaults: $T{", "},", endpointSymbol,
                () -> writeEndpoint(writer, partition.getDefaults()));

        writer.addUseImports(AwsGoDependency.REGEXP);
        writer.write("RegionRegex: regexp.MustCompile($S),", partition.getConfig().expectStringMember("regionRegex")
                .getValue());

        Optional<String> optionalPartitionEndpoint = partition.getPartitionEndpoint();
        Symbol isRegionalizedValue = SymbolUtils.createValueSymbolBuilder(optionalPartitionEndpoint.isPresent()
                ? "false" : "true").build();
        writer.write("IsRegionalized: $T,", isRegionalizedValue);
        optionalPartitionEndpoint.ifPresent(s -> writer.write("PartitionEndpoint: $S,", s));

        Map<StringNode, Node> endpoints = partition.getEndpoints().getMembers();
        if (endpoints.size() > 0) {
            Symbol endpointsSymbol = SymbolUtils.createPointableSymbolBuilder("Endpoints",
                    AwsGoDependency.AWS_ENDPOINTS)
                    .build();
            writer.openBlock("Endpoints: $T{", "},", endpointsSymbol, () -> {
                endpoints.forEach((s, n) -> {
                    writer.openBlock("$S: $T{", "},", s, endpointSymbol,
                            () -> writeEndpoint(writer, n.expectObjectNode()));
                });
            });
        }
    }

    private void writeEndpoint(GoWriter writer, ObjectNode node) {
        node.getStringMember("hostname").ifPresent(n -> {
            writer.write("Hostname: $S,", n.getValue());
        });
        node.getArrayMember("protocols").ifPresent(nodes -> {
            writer.writeInline("Protocols: []string{");
            nodes.forEach(n -> {
                writer.writeInline("$S, ", n.expectStringNode().getValue());
            });
            writer.write("},");
        });
        node.getArrayMember("signatureVersions").ifPresent(nodes -> {
            writer.writeInline("SignatureVersions: []string{");
            nodes.forEach(n -> writer.writeInline("$S, ", n.expectStringNode().getValue()));
            writer.write("},");
        });
        node.getMember("credentialScope").ifPresent(n -> {
            ObjectNode credentialScope = n.expectObjectNode();
            Symbol credentialScopeSymbol = SymbolUtils.createValueSymbolBuilder("CredentialScope",
                    AwsGoDependency.AWS_ENDPOINTS)
                    .build();
            writer.openBlock("CredentialScope: $T{", "},", credentialScopeSymbol, () -> {
                credentialScope.getStringMember("region").ifPresent(nn -> {
                    writer.write("Region: $S,", nn.getValue());
                });
                credentialScope.getStringMember("service").ifPresent(nn -> {
                    writer.write("Service: $S,", nn.getValue());
                });
            });
        });
    }

    private final class Partition {
        private final String id;
        private final ObjectNode defaults;
        private String dnsSuffix;
        private final ObjectNode config;

        private Partition(ObjectNode config, String partition) {
            id = partition;
            this.config = config;

            // Resolve the partition defaults + the service defaults.
            ObjectNode serviceDefaults = config.expectObjectMember("defaults").merge(getService()
                    .getObjectMember("defaults")
                    .orElse(Node.objectNode()));

            // Resolve the hostnameTemplate to use for this service in this partition.
            String hostnameTemplate = serviceDefaults.expectStringMember("hostname").getValue();
            hostnameTemplate = hostnameTemplate.replace("{service}", endpointPrefix);
            hostnameTemplate = hostnameTemplate.replace("{dnsSuffix}",
                    config.expectStringMember("dnsSuffix").getValue());

            this.defaults = serviceDefaults.withMember("hostname", hostnameTemplate);

            dnsSuffix = config.expectStringMember("dnsSuffix").getValue();
        }

        /**
         * @return the partition defaults merged with the service defaults
         */
        ObjectNode getDefaults() {
            return defaults;
        }

        ObjectNode getService() {
            ObjectNode services = config.getObjectMember("services").orElse(Node.objectNode());
            return services.getObjectMember(endpointPrefix).orElse(Node.objectNode());
        }

        ObjectNode getEndpoints() {
            return getService().getObjectMember("endpoints").orElse(Node.objectNode());
        }

        Optional<String> getPartitionEndpoint() {
            ObjectNode service = getService();
            // Note: regionalized services always use regionalized endpoints.
            return service.getBooleanMemberOrDefault("isRegionalized", true)
                    ? Optional.empty()
                    : service.getStringMember("partitionEndpoint").map(StringNode::getValue);
        }

        public String getId() {
            return id;
        }

        public ObjectNode getConfig() {
            return config;
        }
    }

    private static class ResolveConfigField extends ConfigField {
        private final boolean shared;

        public ResolveConfigField(Builder builder) {
            super(builder);
            this.shared = builder.shared;
        }

        public boolean isShared() {
            return shared;
        }

        public static Builder builder() {
            return new Builder();
        }

        private static class Builder extends ConfigField.Builder {
            private boolean shared;

            public Builder() {
                super();
            }

            /**
             * Set the resolver config field to be shared common parameter
             *
             * @param shared whether the resolver config field is shared
             * @return the builder
             */
            public Builder shared(boolean shared) {
                this.shared = shared;
                return this;
            }

            @Override
            public ResolveConfigField build() {
                return new ResolveConfigField(this);
            }

            @Override
            public Builder name(String name) {
                super.name(name);
                return this;
            }

            @Override
            public Builder type(Symbol type) {
                super.type(type);
                return this;
            }

            @Override
            public Builder documentation(String documentation) {
                super.documentation(documentation);
                return this;
            }
        }
    }
}
