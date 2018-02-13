/*
author Vyacheslav Kryuchenko
*/
package web

import (
	"sort"
	"testing"
)

func TestAssetsLen(t *testing.T) {
	names := AssetNames()
	if len(names) == 0 {
		t.Error("assets empty")
	}
}

func TestAssetNames(t *testing.T) {
	sources := []string{
		"header.tmpl",
		"footer.tmpl",
		"login.tmpl",
		"main.tmpl",
		"static/css/bootstrap.min.css",
		"static/css/bootstrap-theme.min.css",
		"static/fonts/glyphicons-halflings-regular.svg",
		"static/fonts/glyphicons-halflings-regular.woff",
		"static/fonts/glyphicons-halflings-regular.woff2",
		"static/functions.js",
	}
	names := AssetNames()
	sourceLen := len(sources)
	namesLen := len(names)
	if sourceLen > namesLen {
		t.Errorf("elements count missmatch: %d > %d. May be generate assets fix this?", sourceLen, namesLen)
	}
	if sourceLen < namesLen {
		t.Errorf("elements count missmatch: %d < %d. Add source or re-generate assets to fix this", sourceLen, namesLen)
	}
	sort.Strings(sources)
	sort.Strings(names)
	for i := 0; i < namesLen; i++ {
		if sources[i] != names[i] {
			t.Errorf("%s != %s", sources[i], names[i])
		}
	}
}
