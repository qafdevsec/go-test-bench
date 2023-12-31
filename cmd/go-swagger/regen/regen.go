// Command regen generates go-swagger yaml and handlers from route data.
package main

import (
	"bytes"
	_ "embed"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/serveswagger/shared"
)

var (
	//go:embed data/swagger.goyaml
	ymlTmpl string

	//go:embed data/code.gogotmpl
	goTmpl string
)

func main() {
	shared.RegisterNewRoutes()

	rmap := common.PopulateRouteMap(common.AllRoutes)

	// do not generate code for input types currently unsupported with swagger.
	shared.FilterInputTypes(rmap)

	var rlist = make(common.Routes, 0, len(rmap))
	for _, r := range rmap {
		if len(r.Sinks) == 0 || len(r.Sinks[0].Name) == 0 {
			// skip
			continue
		}
		if len(r.Inputs) == 0 {
			continue
		}
		if len(r.Sinks) > 0 && len(r.Sinks[0].Name) > 0 {
			rlist = append(rlist, r)
		}
	}
	//sort so the generated code is stable
	sort.Sort(rlist)

	cmdDir, err := findSwagCmd()
	if err != nil {
		log.Fatal(err)
	}
	genYml, err := os.Create(filepath.Join(cmdDir, "swagger.yml"))
	if err != nil {
		log.Fatal(err)
	}
	defer genYml.Close()

	tdata := tmplData{
		GenNotice: "Code generated by cmd/go-swagger/regen; DO NOT EDIT.",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		GenCmd:    "go run ./cmd/go-swagger/regen/regen.go",
		Routes:    rlist,
	}
	tfuncs := template.FuncMap{
		"capital":         strings.Title,
		"routePkg":        routePkg,
		"routeIdentifier": routeIdentifier,
		"sinkName":        sinkName,
		"sinkFn":          sinkFn,
	}

	if err = generateYaml(tdata, tfuncs, genYml); err != nil {
		log.Fatal(err)
	}
	if err = runSwagger(cmdDir); err != nil {
		log.Fatal(err)
	}
	swagPkg, err := findSwagPkg()
	if err != nil {
		log.Fatal(err)
	}

	genGo := filepath.Join(swagPkg, "generatedInit.go")
	g, err := os.Create(genGo)
	if err != nil {
		log.Fatal(err)
	}
	defer goFmt(genGo)
	defer g.Close()
	if err = generateCode(tdata, tfuncs, g); err != nil {
		log.Fatal(err)
	}
}

func goFmt(path string) {
	gf := exec.Command("gofmt", "-w", path)
	gf.Stderr, gf.Stdout = os.Stderr, os.Stdout
	if err := gf.Run(); err != nil {
		log.Fatal(err)
	}
}

type tmplData struct {
	GenNotice string // kept separate from template so that automated stuff doesn't identify the template itself as generated.
	Timestamp string
	GenCmd    string
	Routes    common.Routes
}

func generateYaml(tdata tmplData, tfuncs template.FuncMap, w io.Writer) error {
	tmpl, err := template.New("yaml").Funcs(tfuncs).Parse(ymlTmpl)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, tdata)
}

func runSwagger(cmdDir string) error {
	//swagger generate server --target ../../go-swagger --name SwaggerBench --spec ../swagger.yml --principal interface{} --exclude-main
	swag := exec.Command("swagger")
	swag.Args = append(swag.Args,
		"generate", "server",
		"--target", cmdDir,
		"--name", "SwaggerBench",
		"--spec", filepath.Join(cmdDir, "swagger.yml"),
		"--principal", "interface{}",
		"--exclude-main",
	)
	swag.Stdout, swag.Stderr = os.Stdout, os.Stderr
	return swag.Run()
}

func generateCode(td tmplData, tfuncs template.FuncMap, w io.Writer) error {
	tmpl, err := template.New("go").Funcs(tfuncs).Parse(goTmpl)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, td)
}

// sink name, as used by swagger in url
// e.g. exec.Command
func sinkName(s *common.Sink) string {
	if len(s.URL) == 0 {
		return "Sink"
	}
	return s.URL
}

// sink name, for use in identifier
// e.g. Command (anything preceeding the dot is stripped)
func sinkIdent(s *common.Sink) string {
	fname := s.URL
	if len(fname) == 0 {
		return "Sink"
	}
	idx := strings.LastIndexByte(fname, '.')
	if idx < 0 {
		return exportIdentifier(fname)
	}
	return exportIdentifier(fname[idx+1:])
}

// name of wrapper around vulnerable function
// e.g. GetQueryCommand
func sinkFn(in string, s *common.Sink) string {
	return "Get" + exportIdentifier(in) + sinkIdent(s)
}

// generates package name swagger uses for route
// CmdInjection -> cmd_injection
func routePkg(r *common.Route) string {
	pkg := strings.ToLower(r.Base)
	// ignore 0th letter - start at 1
	j := 1
	for i := 1; i < len(r.Base); i++ {
		if pkg[j] != r.Base[i] {
			//case changed, insert underscore (and advance 1)
			pkg = pkg[:j] + "_" + pkg[j:]
			j++
		}
		j++
	}
	return strings.Trim(pkg, "-_ /")
}

//return an identifier for the route, suitable for use in an exported function name
func routeIdentifier(r *common.Route) string { return exportIdentifier(r.Base) }

//return an identifier suitable for use in an exported function name
func exportIdentifier(id string) string {
	id = strings.Trim(id, "-./ _")
	//swagger replaces lowercase initialisms that the linter would complain about
	// swag: https://github.com/go-openapi/swag/blob/e09cc4d/util.go#L41
	// upstream: https://github.com/golang/lint/blob/3390df4df2787994aea98de825b964ac7944b817/lint.go#L732-L769
	initialisms := []string{"xss", "xsrf", "sql", "json", "http"}
	for _, s := range initialisms {
		id = strings.ReplaceAll(id, s, strings.ToUpper(s))
	}
	id = capitalizeAfter(id, "-._")
	return strings.Title(id)
}

func findSwagCmd() (string, error) { return common.LocateDir("cmd/go-swagger", 5) }
func findSwagPkg() (string, error) { return common.LocateDir("pkg/serveswagger", 5) }

// for each occurrence of each character in 'special', remove that char and capitalize the following letter
func capitalizeAfter(s string, special string) string {
	in := []byte(s)
	for {
		idx := bytes.IndexAny(in, special)
		if idx == -1 {
			return string(in)
		}
		if idx == len(in)-1 {
			//last char
			return string(in[:idx])
		}
		in = append(in[:idx], in[idx+1:]...)
		if in[idx] >= 'a' && in[idx] <= 'z' {
			in[idx] -= 'a' - 'A'
		}
	}
}
