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

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

import java.util.List;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoUniverseTypes;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SmithyGoTypes;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;

public class AwsClientUserAgent implements GoIntegration {
    public static final String MIDDLEWARE_RESOLVER = "addClientUserAgent";

    public static final String SDK_UA_APP_ID = "AppID";

    @Override
    public byte getOrder() {
        return -49;
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ServiceTrait serviceTrait = settings.getService(model).expectTrait(ServiceTrait.class);
        String serviceId = serviceTrait.getSdkId()
                .replace("-", "")
                .replace(" ", "")
                .toLowerCase();
        goDelegator.useFileWriter("api_client.go", settings.getModuleName(), addMiddleware(serviceId));
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .configFields(ListUtils.of(
                                ConfigField.builder()
                                        .name(SDK_UA_APP_ID)
                                        .type(GoUniverseTypes.String)
                                        .documentation("The optional application specific identifier appended to the User-Agent header.")
                                        .build()
                        ))
                        .build()
        );
    }

    private GoWriter.Writable addMiddleware(String serviceId) {
        return goTemplate("""
                func addClientUserAgent(stack $middlewareStack:P, options Options) error {
                    ua, err := getOrAddRequestUserAgent(stack)
                    if err != nil {
                        return err
                    }

                    ua.AddSDKAgentKeyValue(awsmiddleware.APIMetadata, $service:S, goModuleVersion)
                    if len(options.AppID) > 0 {
                        ua.AddSDKAgentKey(awsmiddleware.ApplicationIdentifier, options.AppID)
                    }

                    return nil
                }

                func getOrAddRequestUserAgent(stack *middleware.Stack) (*awsmiddleware.RequestUserAgent, error) {
                    id := (*awsmiddleware.RequestUserAgent)(nil).ID()
                    mw, ok := stack.Build.Get(id)
                    if !ok {
                        mw = awsmiddleware.NewRequestUserAgent()
                        if err := stack.Build.Add(mw, middleware.After); err != nil {
                            return nil, err
                        }
                    }

                    ua, ok := mw.(*awsmiddleware.RequestUserAgent)
                    if !ok {
                        return nil, fmt.Errorf("%T for %s middleware did not match expected type", mw, id)
                    }

                    return ua, nil
                }
                """,
                MapUtils.of(
                        "service", serviceId,
                        "middlewareStack", SmithyGoTypes.Middleware.Stack
                ));
    }
}
