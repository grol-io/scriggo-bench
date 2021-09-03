package scriggo_bench

import (
	"embed"
	_ "embed"
	"io/fs"
	"testing"

	"github.com/Shopify/go-lua"
	"github.com/d5/tengo/script"
	"github.com/open2b/scriggo"
	"github.com/traefik/yaegi/interp"
	gophlua "github.com/yuin/gopher-lua"
)

//go:embed programs/*
var programsFolder embed.FS

func BenchmarkScriggo(b *testing.B) {
	programs, err := loadPrograms("go")
	if err != nil {
		b.Fatal(err)
	}

	for _, program := range programs {
		s, err := scriggo.Build(scriggo.Files{"main.go": program.Code}, nil)
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
	yaegiInterpreter := interp.New(interp.Options{})
	for _, program := range programs {
		b.Run(program.Name, func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				_, err = yaegiInterpreter.Eval(string(program.Code))
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
		tengoScript := script.New(program.Code)
		b.Run(program.Name, func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				_, err := tengoScript.Run()
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
				_ = lua.DoString(L, string(program.Code))
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
				_ = L.DoString(string(program.Code))
			}
		})
	}
}

type testFile struct {
	Name string
	Code []byte
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
