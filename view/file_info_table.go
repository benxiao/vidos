package view

import (
	"github.com/kirillrdy/nadeshiko/html"
	"github.com/kirillrdy/vidos/path"
	"github.com/sparkymat/webdsl/css"
	"github.com/sparkymat/webdsl/css/size"
	"os"
	"path/filepath"
)

//FilesTable display table of files returned by ioutil.ReadDir()
//TODO make each table own type, so that basePath doesn't need to be passed in
func FilesTable(files []os.FileInfo, basePath string) html.Node {

	if len(files) == 0 {
		return html.H4().Text("No files have been added")
	}

	//TODO this is duplicated
	style := tableClass.Style(
		css.Width(size.Percent(100)),
	)

	//TODO use layout
	page := html.Div().Children(
		html.Style().Text(
			style.String(),
		),

		html.Table().Class(tableClass).Children(
			html.Thead().Children(
				html.Tr().Children(
					html.Th().Text("Name"),
					html.Th().Text(""),
				),
			),

			html.Tbody().Children(
				filesTrs(files, basePath)...,
			),
		),
	)

	return page
}

func filesTrs(files []os.FileInfo, basePath string) []html.Node {
	var nodes []html.Node
	for index := range files {
		nodes = append(nodes, fileTr(files[index], basePath))
	}
	return nodes
}

func canBeEncoded(file os.FileInfo) bool {
	if file.IsDir() {
		return false
	}
	ext := filepath.Ext(file.Name())
	if ext == ".mp4" || ext == ".avi" || ext == ".mkv" {
		return true
	}
	return false
}

//TODO only encode link if it can be encoded
func actionsLinksForFile(file os.FileInfo, basePath string) html.Node {
	div := html.Div()
	if canBeEncoded(file) {
		div.Append(
			html.A().Href(path.AddFileForEncodingPath(basePath + file.Name())).Text("Encode"),
		)
	}
	div.Append(
		html.A().Href(path.DeleteFileOrDirectoryPath(basePath + file.Name())).Text("Delete"),
	)
	return div
}

func fileTr(file os.FileInfo, basePath string) html.Node {
	var name html.Node
	if !file.IsDir() {
		name = html.Span().Text(file.Name())
	} else {
		path := path.ViewFilesPath(basePath + file.Name() + "/")
		name = html.A().Href(path).Text(file.Name())
	}

	return html.Tr().Children(
		html.Td().Children(
			name,
		),
		html.Td().Children(
			actionsLinksForFile(file, basePath),
		),
	)
}
