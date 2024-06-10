package render

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type ThumbnailExtension struct {
}

// SetConfig implements renderer.Option.
func (e *ThumbnailExtension) SetConfig(c *renderer.Config) {
	c.NodeRenderers = append(c.NodeRenderers, util.PrioritizedValue{
		Priority: 0,
		Value:    &ThumbnailExtension{},
	})
}

func (e *ThumbnailExtension) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(e)
}

/**
 * Function copied from https://github.com/yuin/goldmark/blob/c15e394c2750495742ad03a4aa2674536664d273/renderer/html/html.go#L1011
 * Copyright (c) 2019 Yusuke Inuzuka
 * MIT License
 */
func nodeToHTMLText(n ast.Node, source []byte) []byte {
	var buf bytes.Buffer
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if s, ok := c.(*ast.String); ok && s.IsCode() {
			buf.Write(s.Text(source))
		} else if !c.HasChildren() {
			buf.Write(util.EscapeHTML(c.Text(source)))
			if t, ok := c.(*ast.Text); ok && t.SoftLineBreak() {
				buf.WriteByte('\n')
			}
		} else {
			buf.Write(nodeToHTMLText(c, source))
		}
	}
	return buf.Bytes()
}

/**
 * Simplified version of https://github.com/yuin/goldmark/blob/c15e394c2750495742ad03a4aa2674536664d273/renderer/html/html.go#L678
 * Copyright (c) 2019 Yusuke Inuzuka
 * MIT License
 */
func (e *ThumbnailExtension) renderImageWithoutThumbnail(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	_, _ = w.WriteString("<img src=\"")
	_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
	_, _ = w.WriteString(`" alt="`)
	_, _ = w.Write(nodeToHTMLText(n, source))
	_ = w.WriteByte('"')
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, html.ImageAttributeFilter)
	}
	_, _ = w.WriteString(">")
	return ast.WalkSkipChildren, nil
}

/**
 * Modified version of https://github.com/yuin/goldmark/blob/c15e394c2750495742ad03a4aa2674536664d273/renderer/html/html.go#L678
 * Copyright (c) 2019 Yusuke Inuzuka
 * MIT License
 */
func (e *ThumbnailExtension) renderImageWithThumbnail(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	newDest := strings.Replace(string(n.Destination), "/media/", "/thumbnail/", 1)
	_, _ = w.WriteString("<a href=\"")
	_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
	_, _ = w.WriteString("\">")
	_, _ = w.WriteString("<img src=\"")
	_, _ = w.Write(util.EscapeHTML(util.URLEscape([]byte(newDest), true)))
	_, _ = w.WriteString(`" alt="`)
	_, _ = w.Write(nodeToHTMLText(n, source))
	_ = w.WriteByte('"')
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, html.ImageAttributeFilter)
	}
	_, _ = w.WriteString("></a>")
	return ast.WalkSkipChildren, nil
}

func (e *ThumbnailExtension) NodeRendererFunc(
	writer util.BufWriter, source []byte, n ast.Node, entering bool,
) (ast.WalkStatus, error) {
	img := n.(*ast.Image)
	dest := string(img.Destination)
	// if dest starts with "/media/" it is a binary file
	// we can use a thumbnail
	if strings.HasPrefix(dest, "/media/") {
		return e.renderImageWithThumbnail(writer, source, n, entering)
	} else {
		// any other image will not use a thumbnail
		return e.renderImageWithoutThumbnail(writer, source, n, entering)
	}
}

func (e *ThumbnailExtension) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindImage, e.NodeRendererFunc)
}
