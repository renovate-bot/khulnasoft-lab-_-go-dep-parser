module github.com/org/repo

go 1.17

require github.com/khulnasoft-lab/go-dep-parser v0.0.0-20231106001621-a21e9e52efce

require (
	golang.org/x/exp v0.0.0-20220407100705-7b9b53b0aca4 // indirect
	golang.org/x/mod v0.13.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace golang.org/x/xerrors v0.0.1 => ./xerrors
