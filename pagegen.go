package main

import (
	"html/template"
	"os"
	"text/template/parse"
)

func ListTemplFields(t *template.Template) []string {
	return listNodeFields(t.Tree.Root, nil)
}

func listNodeFields(node parse.Node, res []string) []string {
	if node.Type() == parse.NodeAction {
		res = append(res, node.String())
	}

	if ln, ok := node.(*parse.ListNode); ok {
		for _, n := range ln.Nodes {
			res = listNodeFields(n, res)
		}
	}
	return res
}

func readTemplate(path *string, name *string) string {
	tpl, err := os.ReadFile(*path)
	if err != nil {
		return ""
	}
	return string(tpl)

}

func parseTemplate(tObject *template.Template, tContent *string) []string {
	t := template.Must(tObject.Parse(*tContent))
	l := ListTemplFields(t)
	return l
}
