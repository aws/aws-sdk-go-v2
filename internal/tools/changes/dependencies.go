package changes

import "path/filepath"

// ModuleGraph is a mapping between modules in a repository and a list of modules within the same repository that depend on that module.
type ModuleGraph map[string][]string

// moduleGraph returns a map between each given module in modules and a slice of modules within the repository that depend on that module.
func moduleGraph(repoRoot string, modules []string) (ModuleGraph, error) {
	deps := map[string][]string{}

	for _, m := range modules {
		mDeps, err := listDependencies(filepath.Join(repoRoot, m))
		if err != nil {
			return nil, err
		}

		for _, d := range mDeps {
			if depList, ok := deps[d]; ok {
				deps[d] = append(depList, m)
			} else {
				deps[d] = []string{m}
			}
		}
	}

	return deps, nil
}

// dependencyUpdates traverses the given module dependency graph, returning a mapping between each module that needs to have
// its dependencies updated and a list of the dependency modules that must be updated.
func (graph ModuleGraph) dependencyUpdates(updatedModules []string) map[string][]string {
	seen := make(map[string]struct{}) // keeps track of which modules have been visited
	updates := map[string][]string{}  // updates stores a list of necessary dependency updates to return.

	// perform a breadth first search through module dependency graph, updating each module that depends on an updated module.
	for len(updatedModules) > 0 {
		m := updatedModules[0]
		if _, ok := seen[m]; ok {
			updatedModules = updatedModules[1:]
			continue
		}

		seen[m] = struct{}{}

		for _, d := range graph[m] {
			updates[d] = append(updates[d], m)

			if _, ok := seen[d]; !ok {
				// add the dependency module to the recursive search if we haven't already encountered it.
				updatedModules = append(updatedModules, d)
			}
		}

		updatedModules = updatedModules[1:]
	}

	return updates
}
