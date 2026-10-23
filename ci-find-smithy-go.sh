#!/bin/bash

# looks for (and modreplaces if existing) a smithy-go branch matching the
# current branch name
#
# the loop will unfurl -*s off of the branch, e.g. sdk branch
# 'feat-foo-bar-baz' will match any of the following (in order):
#  - feat-foo-bar-baz
#  - feat-foo-bar
#  - feat-foo

if [ -z "$SMITHY_GO_REPOSITORY" ]; then
    SMITHY_GO_REPOSITORY=aws/smithy-go
fi
if [ "$SMITHY_GO_REPOSITORY" == /smithy-go ]; then
    SMITHY_GO_REPOSITORY=aws/smithy-go
fi

if [ -z "$RUNNER_TMPDIR" ]; then
    echo env RUNNER_TMPDIR is required
    exit 1
fi

# TRUST_BRANCH_MATCH gates whether we honor a same-named branch in smithy-go for
# co-development. It is set upstream in codegen.yml: true for pushes to main and
# non-fork PRs (author has write access to canonical, so is trusted), false for
# fork PRs. When not trusted, skip branch matching entirely and check out the
# pinned canonical smithy-go, matching what prod codegen uses. Default to the
# untrusted path when unset so a missing value never widens trust.
if [ "$TRUST_BRANCH_MATCH" != "true" ]; then
    codegen_version=$(tr -d '[:space:]' < SMITHY_GO_CODEGEN_VERSION)
    echo "branch matching not trusted; checking out canonical smithy-go at SMITHY_GO_CODEGEN_VERSION=${codegen_version}"
    git clone https://github.com/aws/smithy-go "$RUNNER_TMPDIR"/smithy-go
    git -C "$RUNNER_TMPDIR"/smithy-go checkout "$codegen_version"
    exit 0
fi

if [ -n "$GIT_PAT" ]; then
    repository=https://$GIT_PAT@github.com/$SMITHY_GO_REPOSITORY
else
    repository=https://github.com/$SMITHY_GO_REPOSITORY
fi

branch=$(git branch --show-current)
if [ "$branch" == main ]; then
    echo aws-sdk-go-v2 is on branch main
    git clone "$repository" "$RUNNER_TMPDIR"/smithy-go
    exit 0
fi

# For PR workflows, only the triggering ref is checked out, which in isolation
# is not recognized as a branch by git. Use the specific workflow env instead.
if [ -z "$branch" ]; then
    branch=$GITHUB_HEAD_REF
fi

echo on branch \""$branch"\"
while [ -n "$branch" ] && [[ "$branch" == *-* ]]; do
    echo looking for "$branch"...
    if git ls-remote --exit-code --heads "$repository" refs/heads/"$branch"; then
        echo found "$branch"
        matched_branch=$branch
        break
    fi

    branch=${branch%-*}
done

if [ -z "$matched_branch" ]; then
    # default to SMITHY_GO_CODEGEN_VERSION so CI uses the same smithy-go as prod
    # strip any trailing whitespace/newlines that editors may add to the file
    codegen_version=$(tr -d '[:space:]' < SMITHY_GO_CODEGEN_VERSION)
    echo "no matching branch, checking out canonical smithy-go at SMITHY_GO_CODEGEN_VERSION=${codegen_version}"
    # Always clone canonical, public aws/smithy-go here rather than $repository:
    # the pinned codegen commit only exists upstream, and $repository may point
    # at a fork with no smithy-go repo (or no matching branch).
    git clone https://github.com/aws/smithy-go "$RUNNER_TMPDIR"/smithy-go
    git -C "$RUNNER_TMPDIR"/smithy-go checkout "$codegen_version"
    exit 0
fi

git clone -b "$matched_branch" "$repository" "$RUNNER_TMPDIR"/smithy-go
SMITHY_GO_SRC=$RUNNER_TMPDIR/smithy-go make gen-mod-replace-smithy-.
