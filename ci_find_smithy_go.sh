#!/bin/bash

# looks for (and modreplaces if existing) a smithy-go branch matching the
# current branch name
# e.g. sdk branch 'feat-foo-bar' will match either of the following:
#  - feat-foo-bar
#  - feat-foo
# only one instance of -* will be stripped from the base branch name, so
# sdk branch 'feat-foo-bar-baz' wll match either:
#  - feat-foo-bar-baz
#  - feat-foo-bar

if [ -z "$SMITHY_GO_REPOSITORY" ]; then
    echo env SMITHY_GO_REPOSITORY is required
    exit 1
fi
if [ -z "$RUNNER_TMPDIR" ]; then
    echo env RUNNER_TMPDIR is required
    exit 1
fi

branch=`git branch --show-current`
if [ "$branch" == main ]; then
    echo aws-sdk-go-v2 is on branch main, stop
    exit 0
fi

if [ -n "$CLONE_PAT" ]; then
    repository=https://$CLONE_PAT@github.com/$SMITHY_GO_REPOSITORY
else
    repository=https://github.com/$SMITHY_GO_REPOSITORY
fi

echo looking for smithy-go branch $branch...
git ls-remote --exit-code --heads $repository refs/heads/$branch
if [ "$?" == 0 ]; then
    echo found $branch
    matched_branch=$branch
fi

branch_trimmed=${branch%-*}
echo looking for smithy-go branch $branch_trimmed...
git ls-remote --exit-code --heads $repository refs/heads/$branch_trimmed
if [ "$?" == 0 ]; then
    echo found $branch_trimmed
    matched_branch=$branch_trimmed
fi

if [ -z "$matched_branch" ]; then
    echo found no matching smithy-go branch, stop
    exit 0
fi

git clone -b $matched_branch $repository $RUNNER_TMPDIR/smithy-go
SMITHY_GO_SRC=$RUNNER_TMPDIR/smithy-go make gen-mod-replace-smithy-.
