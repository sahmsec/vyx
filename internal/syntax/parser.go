// Package syntax evaluates an embedded pure-Go parser candidate. Syntax
// validation does not establish authentication, authorization, or safety.
package syntax

import (
	"fmt"
	"path/filepath"
	"strings"

	gts "github.com/odvcencio/gotreesitter"
	blobs "github.com/odvcencio/gotreesitter/grammars/grammar_blobs"
	gruntime "github.com/odvcencio/gotreesitter/grammars/runtime"
)

const Version = "v0.56.0"

type Parser struct{ parsers map[string]*gts.Parser }

type Facts struct {
	Language        string
	ClientDirective bool
	ServerDirective bool
	Imports         int
}

func New() *Parser { return &Parser{parsers: make(map[string]*gts.Parser)} }

func Language(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".js", ".jsx", ".mjs", ".cjs":
		return "javascript"
	case ".ts", ".mts", ".cts":
		return "typescript"
	case ".tsx":
		return "tsx"
	default:
		return ""
	}
}

func (p *Parser) parser(name string) (*gts.Parser, error) {
	if cached := p.parsers[name]; cached != nil {
		return cached, nil
	}
	var data []byte
	switch name {
	case "javascript":
		gruntime.RegisterJavascriptSupport()
		data = blobs.Javascript()
	case "typescript":
		gruntime.RegisterTypescriptSupport()
		data = blobs.Typescript()
	case "tsx":
		gruntime.RegisterTsxSupport()
		data = blobs.Tsx()
	default:
		return nil, fmt.Errorf("unsupported source language")
	}
	// Load supplied embedded bytes directly, bypassing environment-based
	// external grammar selection and the broad language registry.
	lang, err := gruntime.LoadLanguage(name, data)
	if err != nil {
		return nil, fmt.Errorf("load embedded %s grammar", name)
	}
	parser := gts.NewParser(lang)
	parser.SetTimeoutMicros(500_000)
	parser.SetMemoryBudgetBytes(64 << 20)
	p.parsers[name] = parser
	return parser, nil
}

func (p *Parser) Parse(path string, source []byte) (facts Facts, err error) {
	defer func() {
		if recover() != nil {
			facts, err = Facts{}, fmt.Errorf("parser failed unexpectedly; source was not fully analyzed")
		}
	}()
	return p.parse(path, source)
}

func (p *Parser) parse(path string, source []byte) (Facts, error) {
	name := Language(path)
	parser, err := p.parser(name)
	if err != nil {
		return Facts{}, err
	}
	tree, err := parser.ParseStrict(source)
	if tree != nil {
		defer tree.Release()
	}
	if err != nil {
		return Facts{}, fmt.Errorf("parser failed or stopped before completion")
	}
	root := tree.RootNode()
	if root == nil || root.HasErrorOrMissing() {
		return Facts{}, fmt.Errorf("syntax contains an error or an unsupported construct")
	}
	facts := Facts{Language: name}
	inDirectives := true
	for i := 0; i < root.NamedChildCount(); i++ {
		node := root.NamedChild(i)
		kind := node.Type(tree.Language())
		if kind == "comment" || kind == "hash_bang_line" {
			continue
		}
		if kind == "import_statement" {
			facts.Imports++
		}
		if inDirectives && kind == "expression_statement" && node.NamedChildCount() == 1 {
			literal := node.NamedChild(0)
			if literal.Type(tree.Language()) == "string" {
				switch literal.Text(source) {
				case "'use client'", "\"use client\"":
					facts.ClientDirective = true
				case "'use server'", "\"use server\"":
					facts.ServerDirective = true
				}
				continue
			}
		}
		inDirectives = false
	}
	return facts, nil
}
