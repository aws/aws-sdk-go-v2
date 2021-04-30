package release

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/internal/repotools"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/git"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/semver"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// Manifest is a release description of changed modules and their associated tags to be released.
type Manifest struct {
	ID      string                    `json:"id"`
	Modules map[string]ModuleManifest `json:"modules"`
	Tags    []string                  `json:"tags"`
}

// ModuleManifest describes a changed module for release.
type ModuleManifest struct {
	ModulePath string `json:"module_path"`

	From string `json:"from,omitempty"`
	To   string `json:"to"`

	Changes ModuleChange `json:"changes,omitempty"`

	Annotations Annotations `json:"annotations,omitempty"`
}

func getNewModuleVersion(pathMajor string, increment changelog.SemVerIncrement, config repotools.ModuleConfig) (nextVersion string) {
	if len(pathMajor) == 0 {
		nextVersion = "v1.0.0"
	} else {
		nextVersion = pathMajor + ".0.0"
	}

	// Special case, by default new modules will have a pre-release tag, unless we have a corresponding change
	// annotation that marks the module for release.
	if increment == changelog.ReleaseBump {
		return nextVersion
	}

	if len(config.PreRelease) > 0 {
		nextVersion += "-" + config.PreRelease
	} else {
		nextVersion += "-preview"
	}

	return nextVersion
}

// CalculateNextVersion calculates the next version for the module. The provided set of annotations must be applicable
// for this specific module.
func CalculateNextVersion(modulePath string, latest string, config repotools.ModuleConfig, annotations []changelog.Annotation) (next string, err error) {
	_, pathMajor, ok := module.SplitPathVersion(modulePath)
	if !ok {
		return "", fmt.Errorf("invalid module path")
	}
	pathMajor = strings.TrimPrefix(pathMajor, "/")

	increment := changelog.GetVersionIncrement(annotations)

	if len(latest) == 0 {
		next = getNewModuleVersion(pathMajor, increment, config)
		return next, nil
	}

	parsed, ok := semver.Parse(semver.Canonical(latest))
	if !ok {
		return "", fmt.Errorf("failed to parse semver: %v, %v", latest, parsed.Err)
	}

	if increment == changelog.ReleaseBump {
		// Release Bumps are used to elevate pre-release tag versions to released versions
		// Examples:
		//   v1.4.0-preview   => v1.4.0
		//   v1.4.0-preview.1 => v1.4.0

		if len(parsed.Prerelease) == 0 {
			return "", fmt.Errorf("changelog annotation requests release bump, but latest tag is not a pre-release")
		}
		parsed.Prerelease = ""
	} else if len(parsed.Prerelease) > 0 {
		// The existing tag is a pre-release so just increment the pre-release tag number
		// Examples:
		//   v1.4.0-preview   => v1.4.0-preview.1
		//   v1.4.0-preview.2 => v1.4.0-preview.3
		//   v1.4.0-preview   => v1.4.0-rc (if different pre-release identifier is configured)

		if err := incrementPrerelease(&parsed.Prerelease, config.PreRelease); err != nil {
			return "", err
		}
	} else if len(parsed.Prerelease) == 0 && len(config.PreRelease) > 0 {
		// The latest tag was not a pre-release but module is configured for pre-release
		// It is assumed that the target final version is intended to be a minor bump, so we simulate that here
		// when constructing the pre-release tag.
		// Example: v1.3.5 => v1.4.0-preview

		if err := incrementStrInt(&parsed.Minor); err != nil {
			return "", err
		}

		parsed.Patch = "0"

		prerelease := config.PreRelease
		if !strings.HasPrefix(prerelease, "-") {
			prerelease = "-" + prerelease
		}

		parsed.Prerelease = prerelease

	} else if increment == changelog.MinorBump {
		// Module should be bumped by a minor version
		// Example: v1.2.3 => v1.3.0

		if err := incrementStrInt(&parsed.Minor); err != nil {
			return "", err
		}
		parsed.Patch = "0"
	} else {
		// Patch Bump
		// Example: v1.2.3 => v1.2.4
		if err := incrementStrInt(&parsed.Patch); err != nil {
			return "", err
		}
	}

	next = parsed.String()

	if semver.Compare(next, latest) <= 0 {
		return "", fmt.Errorf("computed next version %s is not higher then %s", next, latest)
	}

	return next, nil
}

func incrementStrInt(v *string) error {
	if v == nil {
		return fmt.Errorf("must be a non-nil pointer")
	}

	i, err := strconv.Atoi(*v)
	if err != nil {
		return err
	}
	*v = strconv.Itoa(i + 1)

	return nil
}

func incrementPrerelease(prerelease *string, identifier string) error {
	if prerelease == nil {
		return fmt.Errorf("must be non-nil pointer")
	}

	if !strings.HasSuffix(identifier, "-") {
		identifier = "-" + identifier
	}

	if len(identifier) > 0 && !strings.HasPrefix(*prerelease, identifier) {
		*prerelease = identifier
		return nil
	}

	index := strings.LastIndex(*prerelease, ".")
	if index == -1 {
		*prerelease += ".1"
		return nil
	}

	i, err := strconv.Atoi((*prerelease)[index+1:])
	if err != nil {
		return fmt.Errorf("failed to parse pre-release version number: %v", err)
	}
	*prerelease = (*prerelease)[:index+1] + strconv.Itoa(i+1)

	return nil
}

// BuildReleaseManifest given a mapping of Go module paths to their Module descriptions, returns a summarized
// manifest for release.
func BuildReleaseManifest(id string, modules map[string]*Module) (rm Manifest, err error) {
	rm.ID = id

	rm.Modules = make(map[string]ModuleManifest)

	for modulePath, mod := range modules {
		if mod.Changes == 0 || mod.ModuleConfig.NoTag {
			continue
		}

		nextVersion, err := CalculateNextVersion(modulePath, mod.Latest, mod.ModuleConfig, mod.ChangeAnnotations)
		if err != nil {
			return Manifest{}, err
		}

		mm := ModuleManifest{
			ModulePath:  modulePath,
			From:        mod.Latest,
			To:          nextVersion,
			Changes:     mod.Changes,
			Annotations: annotationsToIDs(mod.ChangeAnnotations),
		}

		rm.Modules[mod.RelativeRepoPath] = mm

		moduleTag, err := git.ToModuleTag(mod.RelativeRepoPath, nextVersion)
		if err != nil {
			return Manifest{}, err
		}

		rm.Tags = append(rm.Tags, moduleTag)
	}

	sort.Strings(rm.Tags)

	return rm, nil
}

// Annotations is a type alias for changelog.Annotation to control how annotations
// are marshaled in a release manifest.
type Annotations []string

func annotationsToIDs(annotations []changelog.Annotation) []string {
	var ids []string

	for _, annotation := range annotations {
		ids = append(ids, annotation.ID)
	}

	return ids
}

// Module is a description of a repository Go module and knowledge about it's current release state.
type Module struct {
	// The parsed go.mod file
	File *modfile.File

	// The modules relative path from the repository root
	RelativeRepoPath string

	// The most recent semver tagged release
	Latest string

	// The next semver tag to release
	Next string

	// The changes for the module
	Changes ModuleChange

	// The change note identifiers applicable for this module
	ChangeAnnotations []changelog.Annotation

	// The release configuration for this module
	ModuleConfig repotools.ModuleConfig
}

// ModuleChange is a bit field to describe the changes for a module
type ModuleChange uint64

// MarshalJSON marshals the chnage bits into a structure JSON object.
func (m ModuleChange) MarshalJSON() ([]byte, error) {
	j := moduleChangeJSON{
		SourceChange:     m&SourceChange != 0,
		NewModule:        m&NewModule != 0,
		DependencyUpdate: m&DependencyUpdate != 0,
	}

	return json.Marshal(j)
}

// UnmarshalJSON unmarshals the JSON object bytes into the ModuleChange bit-field representation.
func (m *ModuleChange) UnmarshalJSON(bytes []byte) error {
	var j moduleChangeJSON

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	if j.SourceChange {
		*m |= SourceChange
	}

	if j.NewModule {
		*m |= NewModule
	}

	if j.DependencyUpdate {
		*m |= DependencyUpdate
	}

	return nil
}

const (
	// SourceChange indicates that the module has source changes since the last tagged release
	SourceChange ModuleChange = 1 << (64 - 1 - iota)

	// NewModule indicates that the module is new and has not been tagged previously
	NewModule

	// DependencyUpdate indicates the module has changes due to a dependency bump
	DependencyUpdate
)

type moduleChangeJSON struct {
	SourceChange     bool `json:"source_change,omitempty"`
	NewModule        bool `json:"new_module,omitempty"`
	DependencyUpdate bool `json:"dependency_update,omitempty"`
}

// buildInverseDependencyGraph builds an inverse dependency graphs mapping a module path to a slice of
// dependents.
func buildInverseDependencyGraph(modules map[string]*Module) (reverseDepGraph map[string][]string) {
	reverseDepGraph = make(map[string][]string)

	for modulePath, mod := range modules {
		for _, require := range mod.File.Require {
			requireModPath := require.Mod.Path
			_, ok := modules[requireModPath]
			if !ok {
				continue
			}
			reverseDepGraph[requireModPath] = append(reverseDepGraph[requireModPath], modulePath)
		}
	}

	return reverseDepGraph
}

// CalculateDependencyUpdates determines which modules require a dependency update bump
// due to one or more of its direct or indirect dependencies being bumped. This will set
// the DependencyUpdate bit flag on the modules set of changes.
func CalculateDependencyUpdates(modules map[string]*Module) error {
	reverseDepGraph := buildInverseDependencyGraph(modules)

	var toVisit []string
	for modulePath := range reverseDepGraph {
		toVisit = append(toVisit, modulePath)
	}
	sort.Strings(toVisit)

	var current string
	for len(toVisit) > 0 {
		current, toVisit = toVisit[0], toVisit[1:]

		m := modules[current]

		if m.Changes == 0 {
			continue
		}

		dependents := reverseDepGraph[current]

		if m.ModuleConfig.NoTag && len(dependents) > 0 {
			return fmt.Errorf("module %v is configured for no releases, but has %d dependents", current,
				len(dependents))
		} else if m.ModuleConfig.NoTag {
			continue
		}

		for _, dependent := range dependents {
			dependentModule := modules[dependent]
			if dependentModule.Changes&DependencyUpdate != 0 {
				continue
			}
			dependentModule.Changes |= DependencyUpdate
			if _, ok := reverseDepGraph[dependent]; ok {
				toVisit = repotools.AppendIfNotPresent(toVisit, dependent)
			}
		}
	}

	return nil
}

var nowTime = time.Now

// NextReleaseID returns the next release identifier based on current YYYY-MM-DD and whether there are multiple tags
// for the given date.
// For example:
//   First Release           => YYYY-MM-DD
//   Second Same-Day Release => YYYY-MM-DD.2
func NextReleaseID(tags []string) (next string) {
	const releaseTagPrefix = "release-"
	const dt = "2006-01-02"

	ct := nowTime().UTC()

	nextTime := time.Date(ct.Year(), ct.Month(), ct.Day(), 0, 0, 0, 0, time.UTC)

	latestNum := 0

	for _, tag := range tags {
		if !strings.HasPrefix(tag, releaseTagPrefix) {
			continue
		}
		tag = strings.TrimPrefix(tag, releaseTagPrefix)
		split := strings.SplitN(tag, ".", 2)

		t, err := time.Parse(dt, split[0])
		if err != nil {
			continue
		}

		if !t.Equal(nextTime) {
			continue
		}

		if len(split) != 2 {
			if latestNum == 0 {
				latestNum = 1
			}
			continue
		}

		i, err := strconv.Atoi(split[1])
		if err != nil {
			continue
		}

		if i > latestNum {
			latestNum = i
		}
	}

	if latestNum == 0 {
		return nextTime.Format(dt)
	}

	latestNum++
	return nextTime.Format(dt) + "." + strconv.Itoa(latestNum)
}
