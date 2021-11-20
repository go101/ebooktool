package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"go101.org/ebooktool/internal/nstd"
)

// Better to fork the mamdContentrkdown lib?
// 1. change bullet list parsing
// 2. Add Node.[Get|Set]Attribute()
// 3. parser.AutoHeadingIDs and id in mannaul arrtibutes should not both present
// 5. support <!-- -->
// 4. use lists as tables
/*
	*-
		# AAA
		# BBB
		# CCC
	*-
		# bla bla
		* bla
		* bla
*/

type Markdown struct {
	doc        ast.Node
	filename   string
	renderInto string
	title      string
}

func ParseMarkdown(content []byte) (*Markdown, error) {
	extensions := 0 |
		parser.CommonExtensions |
		parser.SuperSubscript |
		parser.Attributes
	extensions &^= parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	parser.Opts.ParserHook = parseBlockHook
	docNode := markdown.Parse(content, parser)
	title := findTitle(docNode)
	return &Markdown{doc: docNode, title: title}, nil
}

func (md *Markdown) Title() string {
	return md.title
}

func (md *Markdown) Filename() string {
	return md.filename
}

func (md *Markdown) SetFilename(name string) {
	md.filename = name
}

func (md *Markdown) OutputFile() string {
	return md.renderInto
}

func (md *Markdown) SetOutputFile(outFile string) {
	md.renderInto = outFile
}

var commentStart, commentEnd = "<!--", []byte("-->")

func parseBlockHook(data []byte) (ast.Node, []byte, int) {
	if len(data) < len(commentStart) {
		return nil, nil, 0
	}
	if string(data[:len(commentStart)]) != commentStart {
		return nil, nil, 0
	}

	k := bytes.Index(data[4:], commentEnd)
	if k >= 0 {
		return nil, nil, k + len(commentStart) + len(commentEnd)
	}
	return nil, nil, len(data)
}

func findTitle(doc ast.Node) (title string) {
	ast.WalkFunc(
		doc,
		func(node ast.Node, entering bool) (status ast.WalkStatus) {
			if !entering {
				return
			}

			if heading, ok := node.(*ast.Heading); ok && len(heading.Children) > 0 {
				if text, ok := heading.Children[0].(*ast.Text); ok {
					title = string(text.Literal)
					return ast.Terminate
				}
			}

			return
		},
	)

	return
}

func findImages(doc ast.Node) (images []string) {
	// ToDo
	return nil
}

// Why this? Maybe *.mobi files require this.
func renderNodeHook(w io.Writer, node ast.Node, entering bool) (status ast.WalkStatus, handled bool) {
	if img, ok := node.(*ast.Image); ok {
		if entering {
			fmt.Fprintf(w, `<img src="%s">`, img.Destination)
		} else {
			fmt.Fprint(w, `</img>`)
		}
		return ast.SkipChildren, true
	}

	return
}

func (md *Markdown) Render() []byte {
	options := html.RendererOptions{
		Flags:          html.UseXHTML,
		RenderNodeHook: renderNodeHook,
	}
	renderer := html.NewRenderer(options)
	htmlBytes := markdown.Render(md.doc, renderer)

	// ...
	var buf bytes.Buffer
	fmt.Fprintf(&buf, `<p id="%s"></p>`, filename2ID(md.filename))
	fmt.Fprintln(&buf)
	buf.Write(htmlBytes)

	return buf.Bytes()
}

// ToDo: maybe this is too restricted. At least / is allowed in html5 IDs.
var invalidCharsInID = regexp.MustCompile(`[^0-9a-zA-Z\-\_\:\.]`)

func filename2ID(filename string) string {
	return "f-" + invalidCharsInID.ReplaceAllString(filename, "_")
}

func RenderMarkdownFiles(mdFiles []Markdown) map[string][]File {
	// stat merged file count for each rendered file
	allFiles := make([]File, len(mdFiles))
	renderIntoFiles := make(map[string][]File, len(mdFiles))
	filesRenderInto := make(map[string]string, len(mdFiles))
	for i := range mdFiles {
		md := &mdFiles[i]
		k := md.renderInto
		filesRenderInto[md.filename] = k
		files := renderIntoFiles[k]
		if files == nil {
			renderIntoFiles[k] = allFiles[:1]
		} else {
			renderIntoFiles[k] = allFiles[:len(files)+1]
		}
	}

	nstd.Assertf(
		len(renderIntoFiles) == 1 || len(renderIntoFiles) == len(mdFiles),
		"the count of render-into files must be 1 or %d, but it is %d",
		len(mdFiles), len(renderIntoFiles),
	)

	// modify ast nodes as needed
	filenamesAsID := make(map[string]string, len(mdFiles))
	for i := range mdFiles {
		md := &mdFiles[i]
		filenamesAsID[md.filename] = filename2ID(md.filename)
	}
	for i := range mdFiles {
		md := &mdFiles[i]
		ast.WalkFunc(
			md.doc,
			makeWalkFunc(
				md.filename,
				filenamesAsID,
				filesRenderInto,
			),
		)
	}

	// render markdown files and collect html files for each rendered file
	var sum = 0
	for k, files := range renderIntoFiles {
		renderIntoFiles[k] = allFiles[:sum]
		sum += len(files)
	}
	for k, files := range renderIntoFiles {
		start := len(files)
		renderIntoFiles[k] = allFiles[start:start]
	}
	for i := range mdFiles {
		md := &mdFiles[i]
		k := md.renderInto
		files := renderIntoFiles[k]
		files = append(files, File{
			Name:    md.filename,
			Content: md.Render(),
		})
		renderIntoFiles[k] = files
	}
	for k, files := range renderIntoFiles {
		renderIntoFiles[k] = files[:len(files):len(files)]
	}

	return renderIntoFiles
}

func makeWalkFunc(filename string, filenamesAsID, filesRenderInto map[string]string,
) func(node ast.Node, entering bool) (status ast.WalkStatus) {
	var renderInto, ok = filesRenderInto[filename]
	if !ok {
		nstd.Panicf("file %s has not render target", filename)
	}

	var (
		// ToDo. Always false now.
		prefixFilenameInIDs = false && len(filenamesAsID) > 1
		filenamePrefix      = nstd.Bytes(filenamesAsID[filename] + "-").Decap()

		modifyID = func(id []byte) []byte {
			if !prefixFilenameInIDs {
				return id
			}
			return append(filenamePrefix, id...)
		}

		modifyHref = func(href []byte) []byte {
			if len(href) == 0 {
				log.Printf("warning: blank href in file %s", filename)
				return href
			}
			if nstd.Bytes(href).Index([]byte{'/', '/'}) >= 0 {
				return href
			}

			// "file#anchor"
			// "file"
			// "#anchor"
			// "#"

			tokens := nstd.String(href).SplitN("#", 2)
			targetFilename, targetAnchor := tokens[0], ""
			if len(tokens) > 1 { // == 2
				targetAnchor = tokens[1]
			}
			if targetFilename == "" {
				targetFilename = filename
			}

			if prefixFilenameInIDs && targetAnchor != "" {
				targetAnchor = filenamesAsID[targetFilename] + "-" + targetAnchor
			}

			var targetRenderInto string
			targetRenderInto, ok = filesRenderInto[targetFilename]
			if !ok {
				nstd.Panicf("target %s has not render target", targetFilename)
			}

			if targetRenderInto == renderInto {
				if targetAnchor != "" {
					return []byte("#" + targetAnchor)
				}
				return []byte("#" + filenamesAsID[targetFilename])
			}

			if targetAnchor != "" {
				return []byte(targetRenderInto + "#" + targetAnchor)
			}
			return []byte(targetRenderInto + "#" + filenamesAsID[targetFilename])
		}

		getAttr = func(node ast.Node) *ast.Attribute {
			var attr *ast.Attribute
			if c := node.AsContainer(); c != nil && c.Attribute != nil {
				attr = c.Attribute
			}
			if l := node.AsLeaf(); l != nil && l.Attribute != nil {
				attr = l.Attribute
			}

			return attr
		}
		tryToSetAttributeID = func(node ast.Node, id string) {
			var attr *ast.Attribute
			if c := node.AsContainer(); c != nil && c.Attribute == nil {
				c.Attribute = &ast.Attribute{}
				attr = c.Attribute
			}
			if l := node.AsLeaf(); l != nil && l.Attribute != nil {
				l.Attribute = &ast.Attribute{}
				attr = l.Attribute
			}
			attr.ID = []byte(id)
		}
	)

	return func(node ast.Node, entering bool) (status ast.WalkStatus) {
		if !entering {
			return
		}

		if link, ok := node.(*ast.Link); ok {
			link.Destination = modifyHref(link.Destination)
		} else if heading, ok := node.(*ast.Heading); ok {
			if heading.HeadingID != "" {
				// avoid duplication
				// https://github.com/gomarkdown/markdown/issues/206
				tryToSetAttributeID(node, heading.HeadingID)
				heading.HeadingID = ""
			}
		}

		if attr := getAttr(node); attr != nil && len(attr.ID) > 0 {
			attr.ID = modifyID(attr.ID)
		}

		return
	}
}
