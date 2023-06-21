/*
 * Copyright 2021 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;

/**
 * Collection of static utilities for working with S3 services.
 */
public final class S3ModelUtils {
    public static String SERVICE_S3_MODEL_FILE = "s3.json";
    public static String SERVICE_S3_CONTROL_MODEL_FILE = "s3-control.json";

    private S3ModelUtils() {}

    /**
     * Return true if service is S3.
     *
     * @param model the model used for generation.
     * @param service the service shape for which default HTTP Client is generated.
     * @return true if service is S3
     */
    public static boolean isServiceS3(Model model, ServiceShape service) {
        return service.expectTrait(ServiceTrait.class).getSdkId().equalsIgnoreCase("S3");
    }

    /**
     * Return true if service is S3.
     *
     * @param model the model used for generation.
     * @param service the service shape for which default HTTP Client is generated.
     * @return true if service is S3
     */
    public static boolean isServiceS3Control(Model model, ServiceShape service) {
        return service.expectTrait(ServiceTrait.class).getSdkId().equalsIgnoreCase("S3 Control");
    }
}
