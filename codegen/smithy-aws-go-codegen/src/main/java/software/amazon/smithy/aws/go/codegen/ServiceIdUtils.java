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

/**
 * Provides utilities for normalizing a ServiceId.
 */
public final class ServiceIdUtils {
    /**
     * Converts beginning of words to title case with spaces removed.
     * Does not modify the case of characters used within a word.
     *
     * @param id service id
     * @return the title cased string
     */
    public static String toTitleCase(String id) {
        char[] charArray = id.toCharArray();
        char prev = ' ';
        for (int i = 0; i < charArray.length; i++) {
            if (prev != ' ') {
                continue;
            }
            charArray[i] = Character.toTitleCase(charArray[i]);
            prev = charArray[i];
        }
        return new String(charArray);
    }
}
