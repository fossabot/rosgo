package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	flag "github.com/ogier/pflag"
)

var (
	infile      string
	outfile     string
	packageName string
)

//go:generate go-bindata -o tmpl.go msg.tmpl

func main() {
	log.SetFlags(0)

	flag.StringVarP(&infile, "in", "i", "", "input file")
	flag.StringVarP(&outfile, "out", "o", "", "output file; defaults to '<input file>.go'")
	flag.StringVarP(&packageName, "package", "p", "msgs", "package name for generated file")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Printf("must provide type of generator")
		flag.PrintDefaults()
		os.Exit(1)
	}
	templateType := flag.Arg(0)
	basename := fmt.Sprintf("%s.tmpl", templateType)
	tmpl := template.New(basename)
	//tmpl = tmpl.Funcs(map[string]interface{}{})
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

	parser, ok := parsers[templateType]
	if !ok {
		log.Fatalf("no parser configured: %s", templateType)
	}
	spec, err := parser(infile)
	if err != nil {
		log.Fatalf("failed to parse %s spec: %s", templateType, err)
	}

	buf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buf, spec)
	if err != nil {
		log.Fatalf("failed to generate go file: %s", err)
	}

	fmt.Println(buf.String())
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
	Fields      []msgField
}

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

		// Clean up the line and skip anything prefixed as comment
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Extract <type> <name> and create field specs
		items := strings.Split(line, " ")
		if len(items) != 2 {
			return nil, fmt.Errorf("unrecognized msg line: %s", line)
		}
		spec.Fields = append(spec.Fields, newMsgField(items[1], items[0]))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return spec, nil
}

type msgField struct {
	Name        string
	TypeName    string
	Info        goInfo
	BuiltIn     bool
	GoTypeName  string
	GoZeroValue string
}

type goInfo struct {
	TypeName  string
	ZeroValue string
}

// TODO(ppg): Handle arrays perhaps
func newMsgField(name, typeName string) (field msgField) {
	field.Name = snakeToCamel(name)
	field.TypeName = typeName

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
	items := strings.Split(name, "/")
	goMessageName := items[0]

	// If we have a namespace adjust goMessageName
	if len(items) == 2 {
		goPackageName := items[0]
		goMessageName = items[1]
		// if its not our package then prefix and return
		if goPackageName != packageName {
			field.GoTypeName = fmt.Sprintf("%s.%s", goPackageName, goMessageName)
			field.GoZeroValue = fmt.Sprintf("%s.%s{}", goPackageName, goMessageName)
			return
		}
	}

	// Returned non-namespaced
	field.GoTypeName = fmt.Sprintf("%s", goMessageName)
	field.GoZeroValue = fmt.Sprintf("%s{}", goMessageName)
	return
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
