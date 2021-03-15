# Templated Markdown Renderer for Go

Package `github.com/wttw/markdownish` is a thing wrapper around [gomarkdown](https://github.com/gomarkdown/markdown)
that allows runtime configuration of the HTML output via a Go template. This makes it possible to, for example,
render to a subset of HTML that can be used in table-style layouts for HTML email.

All the documentation for [gomarkdown](https://github.com/gomarkdown/markdown) applies, and if run without
a template file it will behave similarly.

### Installation

Go to the [releases page](https://github.com/wttw/markdownish/releases/latest) and download the appropriate
zip file for your platform. Unzip it to get `mdtohtml` (or `mdtohtml.zip` if you're that way inclined).

### Usage

`mdtohtml -help` will display the supported flags.

`mdtohtml -template <templatefile> <input.md> <output.html>` will render input.md to output.html using the
configuration in the template file.

### Templates

A template file is a [go template](https://golang.org/pkg/text/template/) syntax file that defines one or more
blocks. Each block has a name, and when rendering markdown that would normally be rendered as, for example, a
`<hr>` tag then mdtohtml will look for a block named `hr` and use that template instead if it exists

A template file might look something like this:

```gotemplate
{{define "text"}}{{.Content}}{{end}}
{{define "cr"}}
{{end}}
{{define "br"}}<br>
{{end}}
{{define "nbsp"}}&nbsp;{{end}}
{{define "em"}}<em>{{.Content}}</em>{{end}}
{{define "strong"}}<strong>{{.Content}}</strong>{{end}}
{{define "del"}}<del>{{.Content}}</del>{{end}}
{{define "blockquote"}}<blockquote{{.Attrs}}>{{.Content}}</blockquote>
{{end}}
{{define "aside"}}<aside{{.Attrs}}>{{.Content}}</aside>
{{end}}
{{define "p"}}<p{{.Attrs}}>{{.Content}}</p>
{{end}}
{{define "hr"}}<hr>
{{end}}
{{define "th"}}<th>{{.Content}}</th>
{{end}}
{{define "tr"}}<tr>{{.Content}}</tr>
{{end}}
{{define "sub"}}<sub>{{.Content}}</sub>{{end}}
{{define "sup"}}<sup>{{.Content}}</sup>{{end}}
{{define "h1"}}<h1{{.Attrs}}>{{.Content}}</h1>
{{end}}
{{define "h2"}}<h2{{.Attrs}}>{{.Content}}</h2>
{{end}}
{{define "h3"}}<h3{{.Attrs}}>{{.Content}}</h3>
{{end}}
{{define "h4"}}<h4{{.Attrs}}>{{.Content}}</h4>
{{end}}
{{define "h5"}}<h5{{.Attrs}}>{{.Content}}</h5>
{{end}}
{{define "h6"}}<h6{{.Attrs}}>{{.Content}}</h6>
{{end}}

```

Each block starts with `{{define "name"}}` and ends with `{{end}}`. Inside the block text will be rendered as-is,
but `{{.Attrs}}` will be replaced with the HTML attributes for a tag, and {{.Content}} will be replaced with the
(appropriately HTML escaped, probably) content.

### Support

[Issue tracker](https://github.com/wttw/markdownish/issues). Patches welcome.