package main

// Derived from https://github.com/gomarkdown/mdtohtml

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/gomarkdown/markdown"
	"github.com/wttw/markdownish/templatehtml"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"text/template"
)

func main() {
	var smartypants, latexdashes, fractions bool
	var templateFile string
	flag.BoolVar(&smartypants, "smartypants", true,
		"Apply smartypants-style substitutions")
	flag.BoolVar(&latexdashes, "latexdashes", true,
		"Use LaTeX-style dash rules for smartypants")
	flag.BoolVar(&fractions, "fractions", true,
		"Use improved fraction rules for smartypants")
	flag.StringVar(&templateFile, "template", "",
		"Load this template file to define output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Markdown Processor "+
			"\nAvailable at http://github.com/wttw/markdownish/\n\n"+
			"Copyright © 2011 Russ Ross <russ@russross.com>\n"+
			"Copyright © 2018 Krzysztof Kowalczyk <https://blog.kowalczyk.info>\n"+
			"Copyright © 2021 Turscar <https://turscar.ie/>\n" +
			"Distributed under the Simplified BSD License\n"+
			"Usage:\n"+
			"  %s [options] [inputfile [outputfile]]\n\n"+
			"Options:\n",
			os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()


	var err error

	// Load the template
	var tpl *template.Template
	if templateFile != "" {
		tpl, err = template.ParseFiles(templateFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading template %s: %v\n", templateFile, err)
			os.Exit(-1)
		}
	}

	// read the input
	var input []byte

	args := flag.Args()
	switch len(args) {
	case 0:
		if input, err = ioutil.ReadAll(os.Stdin); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from Stdin:", err)
			os.Exit(-1)
		}
	case 1, 2:
		if input, err = ioutil.ReadFile(args[0]); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from", args[0], ":", err)
			os.Exit(-1)
		}
	default:
		flag.Usage()
		os.Exit(-1)
	}

	// set up options
	var extensions = parser.NoIntraEmphasis |
		parser.Tables |
		parser.FencedCode |
		parser.Autolink |
		parser.Strikethrough |
		parser.SpaceHeadings

	var htmlFlags html.Flags

	if smartypants {
		htmlFlags |= html.Smartypants
	}
	if fractions {
		htmlFlags |= html.SmartypantsFractions
	}
	if latexdashes {
		htmlFlags |= html.SmartypantsLatexDashes
	}

	params := html.RendererOptions{
		Flags: htmlFlags,
	}
	renderer, _ := templatehtml.NewRenderer(tpl, params)

	// parse and render
		parser := parser.NewWithExtensions(extensions)
		output := markdown.ToHTML(input, parser, renderer)

	// output the result
	var out *os.File
	if len(args) == 2 {
		if out, err = os.Create(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s: %v", args[1], err)
			os.Exit(-1)
		}
		defer out.Close()
	} else {
		out = os.Stdout
	}

	if _, err = out.Write(output); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing output:", err)
		os.Exit(-1)
	}
}

