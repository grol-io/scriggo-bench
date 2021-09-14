// Copyright 2021 The Scriggo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scriggo_bench

import (
	"bytes"
	"embed"
	_ "embed"
	"io/fs"
	"strings"
	"testing"

	"github.com/Shopify/go-lua"

	"github.com/d5/tengo/v2"
	tengoDtdlib "github.com/d5/tengo/v2/stdlib"

	"github.com/dop251/goja"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/native"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"

	gophlua "github.com/yuin/gopher-lua"
)

//go:embed programs/*
var programsFolder embed.FS

func BenchmarkScriggo(b *testing.B) {
	programs, err := loadPrograms("go")
	if err != nil {
		b.Fatal(err)
	}

	opts := &scriggo.BuildOptions{
		Packages: native.Packages{
			"strings": native.Package{
				Name: "strings",
				Declarations: native.Declarations{
					"ContainsRune": strings.ContainsRune,
				},
			},
		},
	}
	for _, program := range programs {
		s, err := scriggo.Build(scriggo.Files{"main.go": program.Code}, opts)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(program.Name, func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				err := s.Run(nil)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkYaegi(b *testing.B) {
	programs, err := loadPrograms("go")
	if err != nil {
		b.Fatal(err)
	}
	for _, program := range programs {
		i := interp.New(interp.Options{})
		err = i.Use(stdlib.Symbols)
		if err != nil {
			b.Fatal(err)
		}
		code := string(program.Code)
		if index := strings.LastIndex(code, "// YAEGI\n"); index > 0 {
			// This assumes that the programs do not contain multibyte utf8 characters
			code = code[index+9:]

			_, err = i.Eval(`import "strings"`)
			if err != nil {
				b.Fatal(err)
			}
		}

		b.Run(program.Name, func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				_, err = i.Eval(code)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkTengo(b *testing.B) {
	programs, err := loadPrograms("tengo")
	if err != nil {
		b.Fatal(err)
	}
	for _, program := range programs {
		tengoScript := tengo.NewScript(program.Code)
		var imports []string
		if bytes.Contains(program.Code, []byte(`import("text")`)) {
			imports = append(imports, "text")
		}
		if bytes.Contains(program.Code, []byte(`import("fmt")`)) {
			imports = append(imports, "fmt")
		}
		tengoScript.SetImports(tengoDtdlib.GetModuleMap(imports...))
		b.Run(program.Name, func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				_, err = tengoScript.Run()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkGoLua(b *testing.B) {
	programs, err := loadPrograms("lua")
	if err != nil {
		b.Fatal(err)
	}
	L := lua.NewState()
	lua.OpenLibraries(L)
	for _, program := range programs {
		b.Run(program.Name, func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				err = lua.DoString(L, string(program.Code))
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkGopherLua(b *testing.B) {
	programs, err := loadPrograms("lua")
	if err != nil {
		b.Fatal(err)
	}
	L := gophlua.NewState()
	defer L.Close()
	for _, program := range programs {
		b.Run(program.Name, func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				err = L.DoString(string(program.Code))
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkGoja(b *testing.B) {
	programs, err := loadPrograms("javascript")
	if err != nil {
		b.Fatal(err)
	}
	for _, program := range programs {
		b.Run(program.Name, func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				vm := goja.New()
				_, err = vm.RunString(string(program.Code))
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

type testFile struct {
	Name string
	Code []byte
}

var extensionOf = map[string]string{
	"go":         "go",
	"javascript": "js",
	"lua":        "lua",
	"tengo":      "tengo",
}

func loadPrograms(language string) ([]testFile, error) {
	var tests []testFile
	err := fs.WalkDir(programsFolder, "programs", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && (path != "programs" && path != "programs/"+language) {
			return fs.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), "."+extensionOf[language]) {
			return nil
		}
		var test = testFile{Name: d.Name()}
		test.Code, err = programsFolder.ReadFile(path)
		tests = append(tests, test)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tests, nil
}
