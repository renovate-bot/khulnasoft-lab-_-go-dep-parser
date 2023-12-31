package conan

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/xerrors"

	dio "github.com/khulnasoft-labs/go-dep-parser/pkg/io"
	"github.com/khulnasoft-labs/go-dep-parser/pkg/log"
	"github.com/khulnasoft-labs/go-dep-parser/pkg/types"
)

type LockFile struct {
	GraphLock GraphLock `json:"graph_lock"`
}

type GraphLock struct {
	Nodes map[string]Node `json:"nodes"`
}

type Node struct {
	Ref      string   `json:"ref"`
	Requires []string `json:"requires"`
}

type Parser struct{}

func NewParser() types.Parser {
	return &Parser{}
}

func (p *Parser) Parse(r dio.ReadSeekerAt) ([]types.Library, []types.Dependency, error) {
	var lock LockFile
	if err := json.NewDecoder(r).Decode(&lock); err != nil {
		return nil, nil, xerrors.Errorf("failed to decode conan.lock file: %s", err.Error())
	}

	// Get a list of direct dependencies
	var directDeps []string
	if root, ok := lock.GraphLock.Nodes["0"]; ok {
		directDeps = root.Requires
	}

	// Parse packages
	parsed := map[string]types.Library{}
	for i, node := range lock.GraphLock.Nodes {
		if node.Ref == "" {
			continue
		}
		lib, err := parseRef(node.Ref)
		if err != nil {
			log.Logger.Debug(err)
			continue
		}

		// Determine if the package is a direct dependency or not
		direct := slices.Contains(directDeps, i)
		lib.Indirect = !direct

		parsed[i] = lib
	}

	// Parse dependency graph
	var libs []types.Library
	var deps []types.Dependency
	for i, node := range lock.GraphLock.Nodes {
		lib, ok := parsed[i]
		if !ok {
			continue
		}

		var childDeps []string
		for _, req := range node.Requires {
			if child, ok := parsed[req]; ok {
				childDeps = append(childDeps, child.ID)
			}
		}
		if len(childDeps) != 0 {
			deps = append(deps, types.Dependency{
				ID:        lib.ID,
				DependsOn: childDeps,
			})
		}

		libs = append(libs, lib)
	}
	return libs, deps, nil
}

func parseRef(ref string) (types.Library, error) {
	// full ref format: package/version@user/channel#rrev:package_id#prev
	// various examples:
	// 'pkga/0.1@user/testing'
	// 'pkgb/0.1.0'
	// 'pkgc/system'
	// 'pkgd/0.1.0#7dcb50c43a5a50d984c2e8fa5898bf18'
	ss := strings.Split(strings.Split(strings.Split(ref, "@")[0], "#")[0], "/")
	if len(ss) != 2 {
		return types.Library{}, xerrors.Errorf("Unable to determine conan dependency: %q", ref)
	}
	return types.Library{
		ID:      fmt.Sprintf("%s/%s", ss[0], ss[1]),
		Name:    ss[0],
		Version: ss[1],
	}, nil
}
