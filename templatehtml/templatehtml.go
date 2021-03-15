package templatehtml

import (
	"bytes"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"os"
	"sort"
	"strings"
	"text/template"
	"io"
)

type TemplateRenderer struct {
	Template *template.Template
	Fallback markdown.Renderer
}

type Attributes struct {
	Content string
	Attr map[string]string
	Classes []string
	ID string
	Attrs string
}

var _ markdown.Renderer = &TemplateRenderer{}

func NewRenderer(tpl *template.Template, opts html.RendererOptions) (markdown.Renderer, error) {
	r := &TemplateRenderer{
		Template: tpl,
		Fallback: html.NewRenderer(opts),
	}
	return r, nil
}

func (t TemplateRenderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	if t.Template == nil {
		return t.Fallback.RenderNode(w, node, entering)
	}
	switch node.(type) {
	case *ast.Text:
		return t.Text("text", w, node)
	case *ast.Softbreak:
		return t.Closed("cr", w, node)
	case *ast.Hardbreak:
		return t.Closed("br", w, node)
	case *ast.NonBlockingSpace:
		return t.Closed("nbsp", w, node)
	case *ast.Emph:
		return t.Span("em", w, node, entering)
	case *ast.Strong:
		return t.Span("strong", w, node, entering)
	case *ast.Del:
		return t.Span("del", w, node, entering)
	case *ast.BlockQuote:
		return t.Block("blockquote", w, node, entering)
	case *ast.Aside:
		return t.Block("aside", w, node, entering)
	case *ast.Link:
		break
	case *ast.CrossReference:
		break
	case *ast.Citation:
		break
	case *ast.Image:
		break
	case *ast.Code:
		break
	case *ast.Caption:
		break
	case *ast.CaptionFigure:
		break
	case *ast.Document:
		break
	case *ast.Paragraph:
		return t.Block("p", w, node, entering)
	case *ast.HTMLSpan:
		break
	case *ast.HTMLBlock:
		break
	case *ast.Heading:
		break
	case *ast.HorizontalRule:
		return t.Closed("hr", w, node)
	case *ast.List:
		break
	case *ast.ListItem:
		break
	case *ast.Table:
		return t.Block("table", w, node, entering)
	case *ast.TableCell:
		break
	case *ast.TableHeader:
		return t.Span("th", w, node, entering)
	case *ast.TableBody:
		break
	case *ast.TableRow:
		return t.Span("tr", w, node, entering)
	case *ast.Math:
		break
	case *ast.MathBlock:
		break
	case *ast.DocumentMatter:
		break
	case *ast.Callout:
		break
	case *ast.Index:
		break
	case *ast.Subscript:
		return t.Span("sub", w, node, entering)
	case *ast.Superscript:
		return t.Span("sup", w, node, entering)
	default:
		break;
	}
	return t.Fallback.RenderNode(w, node, entering)
}

func (t TemplateRenderer) RenderHeader(w io.Writer, ast ast.Node) {
	t.Fallback.RenderHeader(w, ast)
}

func (t TemplateRenderer) RenderFooter(w io.Writer, ast ast.Node) {
	t.Fallback.RenderFooter(w, ast)
}

func (t TemplateRenderer) Block(name string, w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	return t.Span(name, w, node, entering)
}

func (t TemplateRenderer) Span(name string, w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	tpl := t.Template.Lookup(name)
	if tpl == nil {
		return t.Fallback.RenderNode(w, node, entering)
	}
	container := node.AsContainer()
	if container == nil {
		return t.Fallback.RenderNode(w, node, entering)
	}
	if !entering {
		return ast.SkipChildren
	}
	var content bytes.Buffer
	for _, child := range container.Children {
		status := ast.Walk(child, ast.NodeVisitorFunc(func(n ast.Node, entering bool) ast.WalkStatus {
			return t.RenderNode(&content, n, entering)
		}))
		if status == ast.Terminate {
			return status
		}
	}
	attrs := t.NodeAttributes(node)
	attrs.Content = content.String()
	err := tpl.Execute(w, attrs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "template rendering error for %s: %v", name, err)
		return ast.Terminate
	}
	return ast.SkipChildren
}

func (t TemplateRenderer) Text(name string, w io.Writer, node ast.Node) ast.WalkStatus {
	var content bytes.Buffer
	status := t.Fallback.RenderNode(&content, node, true)
	if status == ast.Terminate {
		return status
	}
	tpl := t.Template.Lookup(name)
	if tpl == nil {
		return t.Fallback.RenderNode(w, node, true)
	}

	err := tpl.Execute(w, struct{
		Content string
	}{
		Content: content.String(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "template rendering error for %s: %v", name, err)
		return ast.Terminate
	}
	return ast.GoToNext
}

func (t TemplateRenderer) NodeAttributes(node ast.Node) Attributes {
	var attr *ast.Attribute
	switch node := node.(type) {
		case *ast.Container:
			attr = node.Attribute
	case *ast.Leaf:
		attr = node.Attribute
	default:
		return Attributes{}
	}
	ret := Attributes{
		Attr:    map[string]string{},
		Classes: []string{},
		ID:      string(attr.ID),
	}
	for k, v := range attr.Attrs {
		ret.Attr[k] = string(v)
	}
	if attr.ID != nil {
		ret.Attr["id"] = ret.ID
	}
	for _, c := range attr.Classes {
		ret.Classes = append(ret.Classes, string(c))
	}
	sort.Strings(ret.Classes)
	ret.Attr["class"] = strings.Join(ret.Classes, " ")

	keys := []string{}
	for k := range ret.Attr {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	s := []string{}
	for _, k := range keys {
		s = append(s, fmt.Sprintf(` %s="%s"`, k, ret.Attr[k]))
	}
	ret.Attrs = strings.Join(s, "")
	return ret
}

func (t TemplateRenderer) Closed(name string, w io.Writer, node ast.Node) ast.WalkStatus {
	tpl := t.Template.Lookup(name)
	if tpl == nil {
		return t.Fallback.RenderNode(w, node, true)
	}
	err := tpl.Execute(w, struct{}{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "template rendering error for %s: %v", name, err)
		return ast.Terminate
	}
	return ast.GoToNext
}

func (t TemplateRenderer) Heading(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	heading, ok := node.(*ast.Heading)
	if !ok {
		return t.Fallback.RenderNode(w, node, entering)
	}
	level := heading.Level
	if level < 1 || level > 6 {
		level = 6
	}
	return t.Span(fmt.Sprintf("h%d", level), w, node, entering)
}
