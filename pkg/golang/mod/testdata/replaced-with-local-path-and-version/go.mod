module github.com/org/repo

go 1.17

require github.com/khulnasoft-lab/go-dep-parser v1.0.0

require (
	golang.org/x/exp v0.0.0-20230811145659-89c5cff77bcb // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 => ./xerrors
