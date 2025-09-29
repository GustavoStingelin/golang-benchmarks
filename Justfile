# Use `just` to run common tasks.

set dotenv-load := false

# Defaults
count   := env_var_or_default('COUNT', '1')
pattern := env_var_or_default('PATTERN', '.')
out     := env_var_or_default('OUT', 'reports/vizb.html')
name    := env_var_or_default('NAME', 'Benchmarks')
group   := env_var_or_default('GROUP', 'n/w/s')

# Shared command fragments (DRY)
_test_base := "go test ./... -run '^$'"
_vizb_cmd  := "$(command -v vizb || echo $(go env GOBIN)/vizb || echo $(go env GOPATH)/bin/vizb)"
_benchstat := "$(command -v benchstat || echo $(go env GOBIN)/benchstat || echo $(go env GOPATH)/bin/benchstat)"

help:
    @echo "Tasks:"
    @echo "  just deps                      # go mod tidy"
    @echo "  just tools                     # install benchstat + vizb"
    @echo "  just bench [PATTERN=. COUNT=10]# run benchmarks"
    @echo "  just save FILE [PATTERN=. COUNT=10]  # run and tee to file"
    @echo "  just stat BEFORE AFTER         # benchstat compare"
    @echo "  just viz  [PATTERN=. OUT=vizb.html NAME=... DESC=... GROUP=n/w/s]"

deps:
    go mod tidy

tools:
    go install golang.org/x/perf/cmd/benchstat@latest
    go install github.com/goptics/vizb@latest

bench:
    {{_test_base}} -bench "{{pattern}}" -benchmem -count {{count}}

save file:
    mkdir -p reports
    {{_test_base}} -bench "{{pattern}}" -benchmem -count {{count}} | tee reports/{{file}}

stat before after:
    {{_benchstat}} reports/{{before}} reports/{{after}}

viz:
    mkdir -p reports
    {{_test_base}} -bench "{{pattern}}" -benchmem -json | {{_vizb_cmd}} -o {{out}} -n {{name}} -p {{group}}
