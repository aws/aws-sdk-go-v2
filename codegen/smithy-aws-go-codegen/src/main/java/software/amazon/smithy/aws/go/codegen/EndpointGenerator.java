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
import software.amazon.smithy.go.codegen.GoDependency;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.TriConsumer;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.node.Node;
import software.amazon.smithy.model.node.ObjectNode;
import software.amazon.smithy.model.node.StringNode;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.IoUtils;

/**
 * Writes out a file that resolves endpoints using endpoints.json, but the
 * created resolver resolves endpoints for a single service.
 */
final class EndpointGenerator implements Runnable {
    private static final int VERSION = 3;

    private final GoSettings settings;
    private final Model model;
    private final TriConsumer<String, String, Consumer<GoWriter>> writerFactory;
    private final ServiceShape serviceShape;
    private final ObjectNode endpointData;
    private final String endpointPrefix;
    private final Map<String, Partition> partitions = new TreeMap<>();
    private final Map<String, ObjectNode> endpoints = new TreeMap<>();

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
        loadServiceEndpoints();
    }

    private void validateVersion() {
        int version = endpointData.expectNumberMember("version").getValue().intValue();
        if (version != VERSION) {
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

    private void loadServiceEndpoints() {
        for (Partition partition : partitions.values()) {
            ObjectNode serviceData = partition.getService();
            ObjectNode endpointMap = serviceData.getObjectMember("endpoints").orElse(Node.objectNode());

            for (Map.Entry<String, Node> entry : endpointMap.getStringMap().entrySet()) {
                ObjectNode config = entry.getValue().expectObjectNode();
                endpoints.put(entry.getKey(), config);
            }
        }
    }

    private String getInternalEndpointsPath() {
        return "internal/endpoints";
    }

    @Override
    public void run() {
        writerFactory.accept(getInternalEndpointsPath() + "/endpoints.go", "endpoints", (writer) -> {
            generateResolverImplementation(writer);
            generateInternalEndpointsModel(writer);
        });
        writerFactory.accept("endpoints.go", settings.getModuleName(), writer -> {
            generatePublicResolverTypes(writer);
            generateMiddleware(writer);
        });
    }

    private void generateMiddleware(GoWriter writer) {
        // Generate middleware definition
        GoStackStepMiddlewareGenerator middleware = GoStackStepMiddlewareGenerator.createSerializeStepMiddleware(
                getMiddlewareName(), getMiddlewareName());
        middleware.writeMiddleware(writer, this::generateMiddlewareResolverBody,
                this::generateMiddlewareStructureMembers);

        Symbol stackSymbol = SymbolUtils.createPointableSymbolBuilder("Stack", SmithyGoDependency.SMITHY_MIDDLEWARE)
                .build();
        Symbol optionsSymbol = SymbolUtils.createValueSymbolBuilder(String.format("%sMiddlewareOptions",
                getMiddlewareName())).build();

        // Generate Middleware options interface
        writer.openBlock("type $T interface {", "}", optionsSymbol, () -> {
            writer.write("GetEndpointResolver() $L", getResolverInterfaceName());
        });
        writer.write("");
        // Generate Middleware Adder Helper
        writer.openBlock("func $L(stack $P, options $T) {", "}", getAddMiddlewareHelperName(), stackSymbol,
                optionsSymbol, () -> {
                    Symbol afterSymbol = SymbolUtils.createPointableSymbolBuilder("After",
                            SmithyGoDependency.SMITHY_MIDDLEWARE).build();
                    writer.write("resolver := options.GetEndpointResolver()");
                    writer.openBlock("if resolver == nil {", "}", () -> {
                        writer.write("resolver = $L()", getPublicResolverConstructorName());
                    });
                    writer.write("stack.Serialize.Add(&$T{Resolver: resolver}, $T)", middleware.getMiddlewareSymbol(),
                            afterSymbol);
                });
        writer.write("");
        // Generate Middleware Remover Helper
        writer.openBlock("func Remove$LMiddleware(stack $P) error {", "}", middleware.getMiddlewareSymbol(),
                stackSymbol, () -> {
                    writer.write("return stack.Serialize.Remove((&$T{}).ID())", middleware.getMiddlewareSymbol());
                });
    }

    /**
     * @return the add middleware helper name
     */
    public static String getAddMiddlewareHelperName() {
        return String.format("Add%sMiddleware", getMiddlewareName());
    }

    private static String getMiddlewareName() {
        return "ResolveEndpoint";
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
        w.write("endpoint, err = m.Resolver.ResolveEndpoint(awsmiddleware.GetRegion(ctx))");
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
        w.write("Resolver $L", getResolverInterfaceName());
    }

    private String getPublicResolverConstructorName() {
        return "NewDefaultEndpointResolver";
    }

    private Symbol.Builder getInternalEndpointsSymbol(String symbolName, boolean pointable) {
        Symbol.Builder builder;
        if (pointable) {
            builder = SymbolUtils.createPointableSymbolBuilder(symbolName);
        } else {
            builder = SymbolUtils.createValueSymbolBuilder(symbolName);
        }
        return builder.namespace(settings.getModuleName() + "/" + getInternalEndpointsPath(), "/");
    }

    private Symbol.Builder getInternalEndpointsSymbol(String symbolName, boolean pointable, String moduleAlias) {
        return getInternalEndpointsSymbol(symbolName, pointable)
                .putProperty(SymbolUtils.NAMESPACE_ALIAS, moduleAlias);
    }

    private void generatePublicResolverTypes(GoWriter writer) {
        final String internalAlias = "serviceEndpoints";

        Symbol awsEndpointSymbol = SymbolUtils.createValueSymbolBuilder("Endpoint", AwsGoDependency.AWS_CORE).build();
        Symbol internalEndpointsSymbol = getInternalEndpointsSymbol(getResolverImplementationName(), true,
                internalAlias).build();

        // Generate Resolver Interface
        writer.openBlock("type $L interface {", "}", getResolverInterfaceName(), () -> {
            writer.write("ResolveEndpoint(region string) ($T, error)", awsEndpointSymbol);
        });
        writer.write("var _ $L = &$T{}", getResolverInterfaceName(), internalEndpointsSymbol);
        writer.write("");

        Symbol resolverOptionsSymbol = SymbolUtils.createPointableSymbolBuilder("ResolveOptions").build();
        writer.writeDocs(String.format("%s is the set endpoint resolver configuration options",
                resolverOptionsSymbol.getName()));
        writer.write("type $T = $T", resolverOptionsSymbol, SymbolUtils.createValueSymbolBuilder("ResolveOptions",
                AwsGoDependency.AWS_ENDPOINTS).build());

        // Resolver Constructor
        writer.openBlock("func $L(options ... func($P)) $P {", "}", getPublicResolverConstructorName(),
                resolverOptionsSymbol, internalEndpointsSymbol, () -> {
                    writer.write("o := &$T{}", resolverOptionsSymbol);
                    writer.openBlock("for _, fn := range options {", "}", () -> {
                        writer.write("fn(o)");
                    });
                    writer.write("return $T(*o)", getInternalEndpointsSymbol("NewResolver", false, internalAlias)
                            .build());
                });
    }

    private void generateResolverImplementation(GoWriter writer) {
        Symbol awsEndpointSymbol = SymbolUtils.createValueSymbolBuilder("Endpoint", AwsGoDependency.AWS_CORE).build();

        Symbol resolverOptionsSymbol = SymbolUtils.createPointableSymbolBuilder("ResolveOptions",
                AwsGoDependency.AWS_ENDPOINTS).build();

        // Resolver Implementation
        Symbol resolverSymbol = SymbolUtils.createPointableSymbolBuilder(getResolverImplementationName()).build();
        writer.openBlock("type $T struct {", "}", resolverSymbol, () -> {
            writer.write("options $T", resolverOptionsSymbol);
            writer.write("partitions $T", SymbolUtils.createValueSymbolBuilder("Partitions",
                    AwsGoDependency.AWS_ENDPOINTS).build());
        });
        writer.write("");
        writer.openBlock("func (r $P) ResolveEndpoint(region string) ($T, error) {", "}", resolverSymbol, awsEndpointSymbol,
                () -> {
                    writer.write("return r.partitions.EndpointFor(region, r.options)");
                });
        writer.write("");
        writer.openBlock("func NewResolver(o $T) *$T {", "}", resolverOptionsSymbol, resolverSymbol,
                () -> {
                    writer.openBlock("return &$T{", "}", resolverSymbol, () -> {
                        writer.write("options: o,");
                        writer.write("partitions: $L,", getGeneratedEndpointsVariable());
                    });
                });
    }

    private String getGeneratedEndpointsVariable() {
        return "DefaultPartitions";
    }

    /**
     * @return the endpoint resolver interface name
     */
    public static String getResolverInterfaceName() {
        return "EndpointResolver";
    }

    private String getResolverImplementationName() {
        return "EndpointResolver";
    }

    private void generateInternalEndpointsModel(GoWriter writer) {
        writer.addUseImports(AwsGoDependency.AWS_ENDPOINTS);

        Symbol partitionsSymbol = SymbolUtils.createPointableSymbolBuilder("Partitions", AwsGoDependency.AWS_ENDPOINTS)
                .build();
        writer.openBlock("var $L = $T{", "}", getGeneratedEndpointsVariable(), partitionsSymbol, () -> {
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

        writer.openBlock("RegionRegex: func() $P {", "}(),",
                SymbolUtils.createPointableSymbolBuilder("Regexp", AwsGoDependency.REGEXP).build(), () -> {
                    writer.write("r, _ := regexp.Compile($S)", partition.getConfig().expectStringMember("regionRegex")
                            .getValue());
                    writer.write("return r");
                });

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
            hostnameTemplate = hostnameTemplate.replace("{dnsSuffix}", config.expectStringMember("dnsSuffix").getValue());

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
}
