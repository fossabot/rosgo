package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	flag "github.com/ogier/pflag"
)

var (
	infile      string
	outfile     string
	packageName string
	dryRun      bool
)

//go:generate go-bindata -o tmpl.go msg.tmpl

type loopVarSetter interface {
	SetLoopVar(interface{})
}

func main() {
	log.SetFlags(0)

	flag.StringVarP(&infile, "in", "i", "", "input file")
	flag.StringVarP(&outfile, "out", "o", "", "output file; defaults to '<input file>.go'")
	flag.StringVarP(&packageName, "package", "p", "msgs", "package name for generated file")
	flag.BoolVar(&dryRun, "dry-run", false, "output the file that would be generated to stdout")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Printf("must provide type of generator")
		flag.PrintDefaults()
		os.Exit(1)
	}
	templateType := flag.Arg(0)
	basename := fmt.Sprintf("%s.tmpl", templateType)
	tmpl := template.New(basename)
	tmpl = tmpl.Funcs(map[string]interface{}{
		// HACK(ppg): Allow setting a loop variable a struct so we can use it
		"setloopvar": func(setter loopVarSetter, value interface{}) interface{} {
			setter.SetLoopVar(value)
			return setter
		},
	})
	data, err := Asset(basename)
	if err != nil {
		log.Printf("unrecognized generator template: %s (%s)", templateType, err)
		flag.PrintDefaults()
		os.Exit(1)
	}
	tmpl, err = tmpl.Parse(string(data))
	if err != nil {
		log.Printf("unable to template %s: %s", templateType, err)
		os.Exit(1)
	}

	if flag.NArg() > 1 {
		log.Printf("unrecognized arguments: %v", flag.Args()[1:])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if infile == "" {
		log.Printf("must provide input file")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if outfile == "" {
		outfile = infile + ".go"
	}

	typeParser, ok := parsers[templateType]
	if !ok {
		log.Fatalf("no parser configured: %s", templateType)
	}
	spec, err := typeParser(infile)
	if err != nil {
		log.Fatalf("failed to parse %s spec %s: %s", templateType, infile, err)
	}

	buf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buf, spec)
	if err != nil {
		log.Fatalf("failed to generate Go file: %s", err)
	}

	fset := token.NewFileSet()
	ast, err := parser.ParseFile(fset, "", buf.Bytes(), parser.ParseComments)
	if err != nil {
		log.Fatalf("bad Go source code was generated: %s\n%s", err, buf.String())
	}
	buf.Reset()
	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(buf, fset, ast)
	if err != nil {
		log.Fatalf("generated Go source code could not be reformatted: %s", err)
	}

	if dryRun {
		fmt.Println(buf.String())
		return
	}

	err = ioutil.WriteFile(outfile, buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("failed to write go file: %s", err)
	}
	log.Printf("Wrote %s from %s", outfile, infile)
}

var parsers = map[string]func(infile string) (interface{}, error){
	"msg": parseMsgSpec,
}

type msgSpec struct {
	InFile      string
	Raw         string
	PackageName string
	MD5Sum      string
	Name        string
	Fields      []*msgField
	Constants   []*msgConstant
	HasBuiltIn  bool
	HasSlice    bool
	HasArray    bool
}

var constMatcher = regexp.MustCompile(`^\s*([\w/]+)\s+(\w+)\s*=\s*(\d+)#?.*`)
var fieldMatcher = regexp.MustCompile(`^\s*([\w/]+)(\[(\d+)?\])?\s+(\w+)#?.*`)

//uint8 status
//uint8 PENDING         = 0   # The goal has yet to be processed by the action server
func parseMsgSpec(infile string) (interface{}, error) {
	data, err := ioutil.ReadFile(infile)
	if err != nil {
		return nil, err
	}

	spec := new(msgSpec)
	spec.PackageName = packageName
	basename := filepath.Base(infile)
	spec.Name = strings.TrimSuffix(basename, filepath.Ext(basename))
	spec.InFile = infile
	spec.Raw = string(data)
	spec.MD5Sum = fmt.Sprintf("%x", md5.Sum(data))

	// Read the data line by line
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Check constant first
		items := constMatcher.FindStringSubmatch(line)
		if items != nil {
			constant := newMsgConstant(items[2], items[1], items[3])
			spec.Constants = append(spec.Constants, constant)
			continue
		}

		// Now check field
		items = fieldMatcher.FindStringSubmatch(line)
		if items != nil {
			//log.Println(line)
			//log.Printf("items(%d): %+v", len(items), items)
			field := newMsgField(items[4], items[1], len(items[2]) > 0, items[3])
			if field.BuiltIn {
				spec.HasBuiltIn = true
			}
			if field.IsArray {
				spec.HasSlice = true
				if field.ArraySize > 0 {
					spec.HasArray = true
				}
			}
			//log.Printf("field: %+v", field)
			spec.Fields = append(spec.Fields, field)
			continue
		}

		if len(line) != 0 && !strings.HasPrefix(line, "#") {
			return nil, fmt.Errorf("unrecognized msg line: %s", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return spec, nil
}

type msgField struct {
	IsArray     bool
	ArraySize   int
	Name        string
	TypeName    string
	BuiltIn     bool
	GoTypeName  string
	GoZeroValue string

	LoopVar interface{}
}

func (m *msgField) SetLoopVar(i interface{}) {
	m.LoopVar = i
}

func newMsgField(name, typeName string, isArray bool, arraySize string) (field *msgField) {
	//log.Printf("name: %s", name)
	//log.Printf("typeName: %s", typeName)
	field = new(msgField)
	field.Name = snakeToCamel(name)
	field.TypeName = typeName
	field.IsArray = isArray
	if isArray && len(arraySize) > 0 {
		var err error
		field.ArraySize, err = strconv.Atoi(arraySize)
		if err != nil {
			panic(err)
		}
	}

	// Try to get field info from built-in first
	info, ok := builtInInfo[typeName]
	if ok {
		field.BuiltIn = true
		field.GoTypeName = info.TypeName
		field.GoZeroValue = info.ZeroValue
		return
	}

	// Not built-in, deconstruct
	field.BuiltIn = false

	// TODO(ppg): Can we set all complex types to nil instead?

	// Split on / for namespaced messages
	items := strings.Split(typeName, "/")
	goMessageTypeName := snakeToCamel(items[0])

	// If we have a namespace adjust goMessageName
	if len(items) == 2 {
		goPackageName := items[0]
		goMessageTypeName = snakeToCamel(items[1])
		// if its not our package then prefix and return
		if goPackageName != packageName {
			field.GoTypeName = fmt.Sprintf("%s.%s", goPackageName, goMessageTypeName)
			field.GoZeroValue = fmt.Sprintf("%s.%s{}", goPackageName, goMessageTypeName)
			return
		}
	}

	// Returned non-namespaced
	field.GoTypeName = fmt.Sprintf("%s", goMessageTypeName)
	field.GoZeroValue = fmt.Sprintf("%s{}", goMessageTypeName)
	return
}

type msgConstant struct {
	Name     string
	TypeName string
	Value    string
}

func newMsgConstant(name, typeName, value string) (constant *msgConstant) {
	constant = new(msgConstant)
	constant.Name = snakeToCamel(name)
	constant.TypeName = typeName
	constant.Value = value
	return
}

type goInfo struct {
	TypeName  string
	ZeroValue string
}

var builtInInfo = map[string]goInfo{
	"bool": goInfo{"bool", "false"},
	"byte": goInfo{"int8", "0"},
	"char": goInfo{"uint8", "0"},
	// TODO(ppg): Change to time.Duration, handle conversion to ROS type during serialization/deserialization
	"duration": goInfo{"ros.Duration", "ros.Duration{}"},
	"float32":  goInfo{"float32", "0.0"},
	"float64":  goInfo{"float64", "0.0"},
	"int8":     goInfo{"int8", "0"},
	"int16":    goInfo{"int16", "0"},
	"int32":    goInfo{"int32", "0"},
	"int64":    goInfo{"int64", "0"},
	"string":   goInfo{"string", `""`},
	// TODO(ppg): Change to time.Time, handle conversion to ROS type during serialization/deserialization
	"time":   goInfo{"ros.Time", "ros.Time{}"},
	"uint8":  goInfo{"uint8", "0"},
	"uint16": goInfo{"uint16", "0"},
	"uint32": goInfo{"uint32", "0"},
	"uint64": goInfo{"uint64", "0"},
}

// snakeToCamel returns a string converted from snake case to uppercase
func snakeToCamel(s string) string {
	var result string

	words := strings.Split(s, "_")

	for _, word := range words {
		if upper := strings.ToUpper(word); commonInitialisms[upper] {
			result += upper
			continue
		}

		if len(word) > 0 {
			w := []rune(word)
			w[0] = unicode.ToUpper(w[0])
			result += string(w)
		}
	}

	return result
}

// commonInitialisms, taken from
// https://github.com/golang/lint/blob/32a87160691b3c96046c0c678fe57c5bef761456/lint.go#L702
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}
