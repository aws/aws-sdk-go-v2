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
     *
     * @param id service id
     * @return the title cased string
     */
    public static String toTitleCase(String id) {
        return toTitleCase(id, true);
    }

    /**
     * Converts begining of words to title case with or without the space separator removed.
     * Does not modify the case of characters used within a word.
     *
     * @param id service id
     * @param removeSeparator whether the separator character should be removed between words
     * @return the title cases string
     */
    public static String toTitleCase(String id, boolean removeSeparator) {
        StringBuilder builder = new StringBuilder();
        char[] charArray = id.toCharArray();
        char prev = ' ';
        for (int i = 0; i < charArray.length; i++) {
            char c = charArray[i];

            if (isSeparator(prev)) {
                c = Character.toTitleCase(c);
            }

            if (!removeSeparator || !isSeparator(c)) {
                builder.append(c);
            }

            prev = c;
        }
        return builder.toString();
    }

    private static boolean isSeparator(char c) {
        return c == ' ';
    }
}
