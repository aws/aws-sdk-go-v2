# Description

`changelog` is a utility for creating and managing change annotations for a repository. Change annotations are used in
the repository to refine the [semver] increment behavior for modules that have pending changes to be released.
Additionally, these annotations include descriptions that are used to produce `CHANGELOG.md` entries for each module
that has been changed as part of a release.

# Usage

```
changelog create [-c <tree-ish> | (-cs <tree-ish> -ce <tree-ish>)] [-t <change-type>] [-d <description>] [<module>...]

Options:
-c  <tree-ish>   A commit or tag to generate a change annotation for
-cs <tree-ish>   A starting commit or tag for a change annotation, must be used with -ce to compare changes between two trees
-ce <tree-ish>   An ending commit or tag for a change annotation, must be used with -cs to compare changes between two trees
-r               Declare that the annotation description should be rolled up as a summary when producing summarized CHANGELOG digests
-t <change-type> The change annotation type (release, feature, bugfix, dependency, announcement)
-d <description> The description of the change annotation, must be a string or a valid markdown list block
-ni              Non-Interactive mode

changelog ls

changelog edit <id>

changelog view <id>
```

# Examples

## Create an annotation for modules that were changed in a specific commit

1. Determine the git commit ID for the change you wish to annotate.
    ```
    $ git log --oneline
    e22f8f0948 Update API clients from latest models (#1250)
    ```
1. Use the changelog CLI's `create` verb to
   ```
   changelog create -c e22f8f0948
   ```
1. By default, the CLI will prompt you interactively via text editor (vim by default)
1. Adjust the `type`, `description`, and `modules` fields by populating them into the provided TOML template.
1. Once editing is completed save the file and exit the editor

## Create an annotation for modules that were changed over a commit range

1. Determine the git commit ID for the change you wish to annotate.
    ```
    $ git log --oneline
    e22f8f0948 Update API clients from latest models (#1250)
    9b93441d7f service/ec2: Fix generation of number and bool struct members to be pointers (#1195)
    ```
1. Use the changelog CLI's `create` verb to
   ```
   $ changelog create -cs 9b93441d7f -ce e22f8f0948
   ```
1. By default, the CLI will prompt you interactively via text editor (vim by default)
1. Adjust the `type`, `description`, and `modules` fields by populating them into the provided TOML template.
1. Once editing is completed save the file and exit the editor

## Create an annotation for modules that were changed over a commit range

1. Determine the git commit ID for the change you wish to annotate.
    ```
    $ git log --oneline
    e22f8f0948 Update API clients from latest models (#1250)
    9b93441d7f service/ec2: Fix generation of number and bool struct members to be pointers (#1195)
    ```
1. Use the changelog CLI's `create` verb to
   ```
   $ changelog create -cs 9b93441d7f -ce e22f8f0948
   ```
1. By default, the CLI will prompt you interactively via text editor (vim by default)
1. Adjust the `type`, `description`, and `modules` fields by populating them into the provided TOML template.
1. Once editing is completed save the file and exit the editor

## Create an annotation (non-interactive)

1. By passing the required annotation parameters and the `-ni` flag to the CLI you can create an annotation without
   being prompted interactively using a text editor.
   ```
   $ changelog create -ni -type feature -description "addewd new feature foo" service/s3 feature/s3/manager
   ```

## List Change Annotations

```
$ changelog ls
+--------------------------------------+--------+---------+----------+----------------------+
|                  ID                  |  TYPE  | MODULES | COLLAPSE |     DESCRIPTION      |
+--------------------------------------+--------+---------+----------+----------------------+
| 0ba0c6bf-d697-49d1-ac8f-1f6c7f29663e | bugfix |       1 | false    | a change description |
+--------------------------------------+--------+---------+----------+----------------------+
```

## View Change Annotation

```
$ changelog view 0ba0c6bf-d697-49d1-ac8f-1f6c7f29663e
{
    "id": "0ba0c6bf-d697-49d1-ac8f-1f6c7f29663e",
    "type": "bugfix",
    "description": "a change description",
    "modules": [
        "internal/repotools"
    ]
}
```

## Remove one or more annotations

1. Provide one or more annotations identifiers as position arguments the `changelog`
   ```
   $ changelog rm <id1> <id2> <id3>
   ```


## Remove ALL annotations

```
$ changelog rm -all
```

[semver]: https://semver.org
