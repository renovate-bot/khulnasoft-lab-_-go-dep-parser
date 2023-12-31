package conan_test

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft-labs/go-dep-parser/pkg/c/conan"
	"github.com/khulnasoft-labs/go-dep-parser/pkg/types"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string // Test input file
		wantLibs  []types.Library
		wantDeps  []types.Dependency
	}{
		{
			name:      "happy path",
			inputFile: "testdata/happy.lock",
			wantLibs: []types.Library{
				{
					ID:      "pkga/0.0.1",
					Name:    "pkga",
					Version: "0.0.1",
				},
				{
					ID:       "pkgb/system",
					Name:     "pkgb",
					Version:  "system",
					Indirect: true,
				},
				{
					ID:      "pkgc/0.1.1",
					Name:    "pkgc",
					Version: "0.1.1",
				},
			},
			wantDeps: []types.Dependency{
				{
					ID: "pkga/0.0.1",
					DependsOn: []string{
						"pkgb/system",
					},
				},
			},
		},
		{
			name:      "happy path. lock file with revisions support",
			inputFile: "testdata/happy2.lock",
			wantLibs: []types.Library{
				{
					ID:      "openssl/3.0.3",
					Name:    "openssl",
					Version: "3.0.3",
				},
				{
					ID:       "zlib/1.2.12",
					Name:     "zlib",
					Version:  "1.2.12",
					Indirect: true,
				},
			},
			wantDeps: []types.Dependency{
				{
					ID: "openssl/3.0.3",
					DependsOn: []string{
						"zlib/1.2.12",
					},
				},
			},
		},
		{
			name:      "happy path. lock file without dependencies",
			inputFile: "testdata/empty.lock",
		},
		{
			name:      "sad path. wrong ref format",
			inputFile: "testdata/sad.lock",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.inputFile)
			require.NoError(t, err)
			defer f.Close()

			gotLibs, gotDeps, err := conan.NewParser().Parse(f)
			require.NoError(t, err)

			sort.Slice(gotLibs, func(i, j int) bool {
				ret := strings.Compare(gotLibs[i].Name, gotLibs[j].Name)
				if ret != 0 {
					return ret < 0
				}
				return gotLibs[i].Version < gotLibs[j].Version
			})

			assert.Equal(t, tt.wantLibs, gotLibs)
			assert.Equal(t, tt.wantDeps, gotDeps)
		})
	}
}
