package main

import (
	"errors"
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/git"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/gomod"
)

const createHelpDoc = `changelog create [-c <tree-ish> | (-cs <tree-ish> -ce <tree-ish>)] [-t <change-type>] [-d <description>] [<module>...]

Options:
-c  <tree-ish>   A commit or tag to generate a change annotation for
-cs <tree-ish>   A starting commit or tag for a change annotation, must be used with -ce to compare changes between two trees
-ce <tree-ish>   An ending commit or tag for a change annotation, must be used with -cs to compare changes between two trees
-r               Declare that the annotation description should be rolled up as a summary when producing summarized CHANGELOG digests
-t <change-type> The change annotation type (release, feature, bugfix, dependency, announcement)
-d <description> The description of the change annotation, must be a string or a valid markdown list block
-ni              Non-Interactive mode
`

var createCommand = struct {
	Commit string

	CommitStart string
	CommitEnd   string

	Type ChangeType

	Description string

	Collapse bool

	NonInteractive bool
}{}

var createFlagSet = func() *flag.FlagSet {
	fs := flag.NewFlagSet("create", flag.ExitOnError)

	fs.Usage = func() {
		fmt.Fprint(fs.Output(), createHelpDoc)
	}

	fs.StringVar(&createCommand.Commit, "c", "", "")
	fs.StringVar(&createCommand.CommitStart, "cs", "", "")
	fs.StringVar(&createCommand.CommitEnd, "ce", "", "")
	fs.Var(&createCommand.Type, "t", "")
	fs.StringVar(&createCommand.Description, "d", "", "")
	fs.BoolVar(&createCommand.Collapse, "r", false, "")
	fs.BoolVar(&createCommand.NonInteractive, "ni", false, "")

	return fs
}()

func runCreateCommand(args []string, repoRoot string) error {
	if err := createFlagSet.Parse(args); err != nil {
		return err
	}

	discoverer := gomod.NewDiscoverer(repoRoot)

	if err := discoverer.Discover(); err != nil {
		return err
	}

	modules, err := discoverer.ModulesRel()
	if err != nil {
		return fmt.Errorf("failed to discover repository go modules: %w", err)
	}

	if err := validateCreateCommandArguments(createFlagSet.Args(), modules); err != nil {
		return fmt.Errorf("invalid arguments: %w", err)
	}

	modulesToAnnotate := make(map[string]struct{})

	for _, arg := range createFlagSet.Args() {
		modulesToAnnotate[arg] = struct{}{}
	}

	var commitChanges []string

	if createCommand.Commit != "" {
		commitChanges, err = git.Changed(repoRoot, createCommand.Commit)
		if err != nil {
			return fmt.Errorf("failed to get changed files for commit: %v", err)
		}
	} else if createCommand.CommitStart != "" && createCommand.CommitEnd != "" {
		commitChanges, err = git.Changes(repoRoot, createCommand.CommitStart, createCommand.CommitEnd)
		if err != nil {
			return fmt.Errorf("failed to get changed files for commit: %v", err)
		}
	}

	if len(commitChanges) > 0 {
		for moduleDir, submodules := range modules {
			if isChanged, err := gomod.IsModuleChanged(moduleDir, submodules, commitChanges); err != nil {
				return err
			} else if isChanged {
				modulesToAnnotate[moduleDir] = struct{}{}
			}
		}
	}

	annotation, err := changelog.NewAnnotation()
	if err != nil {
		return err
	}

	annotation.Type = changelog.ChangeType(createCommand.Type)
	if annotation.Type == changelog.UnknownChangeType {
		annotation.Type = changelog.BugFixChangeType
	}

	annotation.Description = createCommand.Description
	annotation.Collapse = createCommand.Collapse

	if len(modulesToAnnotate) > 0 {
		for moduleDir := range modulesToAnnotate {
			annotation.Modules = append(annotation.Modules, moduleDir)
		}
	}

	sort.Strings(annotation.Modules)

	if createCommand.NonInteractive {
		if invalid := validateModules(annotation.Modules, modules); len(invalid) > 0 {
			return fmt.Errorf("invalid modules: %v", invalid)
		}

		if err = changelog.Validate(annotation); err != nil {
			return err
		}
	} else {
		if err = interactiveEdit(&annotation, modules); err != nil {
			return err
		}
	}

	return changelog.WriteAnnotation(repoRoot, annotation)
}

func interactiveEdit(annotation *changelog.Annotation, modules map[string][]string) error {
	var issues []string

	template, err := changelog.AnnotationToTemplate(*annotation)
	if err != nil {
		return err
	}

	filledTemplate, err := editTemplate(template)
	if err != nil {
		return fmt.Errorf("failed to create change: %v", err)
	}

	filledAnnotation, err := changelog.TemplateToAnnotation(filledTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse annotation: %v", err)
	}

	if err = changelog.Validate(filledAnnotation); err != nil {
		var ve *changelog.ValidationError
		if !errors.As(err, &ve) {
			return err
		}
		issues = append(issues, ve.Issues...)
	}

	if invalidModules := validateModules(filledAnnotation.Modules, modules); len(invalidModules) > 0 {
		for _, module := range invalidModules {
			issues = append(issues, fmt.Sprintf("unknown module: %s", module))
		}
	}

	if len(issues) > 0 {
		var sb strings.Builder
		sb.WriteString("Invalid Template:\n")
		for _, issue := range issues {
			sb.WriteRune('\t')
			sb.WriteString(issue)
			sb.WriteRune('\n')
		}
		return fmt.Errorf(sb.String())
	}

	// Ensure this didn't get swapped / changed
	filledAnnotation.ID = annotation.ID

	*annotation = filledAnnotation

	return nil
}

func validateModules(input []string, modules map[string][]string) (invalid []string) {
	for _, module := range input {
		if _, ok := modules[module]; !ok {
			invalid = append(invalid, module)
		}
	}
	return invalid
}

func validateCreateCommandArguments(args []string, modules map[string][]string) error {
	if createCommand.Commit != "" && (createCommand.CommitStart != "" || createCommand.CommitEnd != "") {
		return fmt.Errorf("only -c can not be specified with -cs and -ce")
	}

	if (createCommand.CommitStart != "" && createCommand.CommitEnd == "") ||
		(createCommand.CommitEnd != "" && createCommand.CommitStart == "") {
		return fmt.Errorf("-cs must be specified with -ce")
	}

	var unknown []string
	for _, moduleDir := range args {
		if _, ok := modules[moduleDir]; !ok {
			unknown = append(unknown, moduleDir)
		}
	}

	if len(unknown) > 0 {
		return fmt.Errorf("unknown modules: %v", unknown)
	}

	return nil
}
