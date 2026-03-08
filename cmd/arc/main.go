package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/pafthang/arc"
	"github.com/pafthang/bms/internal/app"
	"github.com/pafthang/bms/internal/config"
)

func main() {
	var (
		outPath          string
		format           string
		stdout           bool
		withSystem       bool
		validateQuality  bool
		requireTags      bool
		requireServers   bool
		requireExamples  bool
		requireSchemesCS string
	)
	flag.StringVar(&outPath, "out", "openapi/openapi.json", "output file path")
	flag.StringVar(&format, "format", "json", "output format: json|yaml")
	flag.BoolVar(&stdout, "stdout", false, "write spec to stdout instead of file")
	flag.BoolVar(&withSystem, "with-system", false, "include system routes in generated spec")
	flag.BoolVar(&validateQuality, "validate-quality", false, "validate OpenAPI quality gates and exit non-zero on violations")
	flag.BoolVar(&requireTags, "require-tags", true, "quality gate: require root tags list")
	flag.BoolVar(&requireServers, "require-servers", true, "quality gate: require root servers list")
	flag.BoolVar(&requireExamples, "require-examples", true, "quality gate: require at least one operation example")
	flag.StringVar(&requireSchemesCS, "require-security-schemes", "BearerAuth", "quality gate: comma-separated required security scheme names")
	flag.Parse()

	cfg := config.Load()
	engine := app.BuildEngine(cfg, nil, app.BuildOptions{
		IncludeSystemRoutes: withSystem,
		IncludeHealthRoutes: withSystem,
	})
	spec := engine.OpenAPISpec()

	var (
		data []byte
		err  error
	)
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		data, err = engine.MarshalOpenAPIJSON()
	case "yaml", "yml":
		data, err = engine.MarshalOpenAPIYAML()
	default:
		fatalf("unsupported format %q (use json|yaml)", format)
	}
	if err != nil {
		fatalf("generate openapi: %v", err)
	}
	if validateQuality {
		issues := arc.ValidateOpenAPIQuality(spec, arc.OpenAPIQualityGates{
			RequireRootTags:         requireTags,
			RequireServers:          requireServers,
			RequireExamples:         requireExamples,
			RequiredSecuritySchemes: parseCSV(requireSchemesCS),
		})
		if len(issues) > 0 {
			for _, item := range issues {
				fmt.Fprintf(os.Stderr, "bms openapi quality: %s\n", item)
			}
			os.Exit(2)
		}
	}

	if stdout {
		_, _ = os.Stdout.Write(data)
		if len(data) == 0 || data[len(data)-1] != '\n' {
			_, _ = os.Stdout.Write([]byte("\n"))
		}
		return
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		fatalf("create output dir: %v", err)
	}
	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		fatalf("write output: %v", err)
	}
	fmt.Fprintf(os.Stderr, "openapi written to %s\n", outPath)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "bms arc: "+format+"\n", args...)
	os.Exit(1)
}

func parseCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	out = slices.Compact(out)
	return out
}
