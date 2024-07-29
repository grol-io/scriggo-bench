module scriggo-bench

go 1.22.5

require (
	github.com/Shopify/go-lua v0.0.0-20210302141115-d8ac5566562d
	github.com/d5/tengo/v2 v2.8.0
	github.com/dop251/goja v0.0.0-20210904102640-6338b3246846
	github.com/open2b/scriggo v0.49.0
	github.com/traefik/yaegi v0.9.23
	github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9
	grol.io/grol v0.26.0
)

replace grol.io/grol => ../grol

require (
	fortio.org/log v1.16.0 // indirect
	fortio.org/struct2env v0.4.1 // indirect
	github.com/dlclark/regexp2 v1.4.1-0.20201116162257-a2a8dda75c91 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/kortschak/goroutine v1.1.2 // indirect
	golang.org/x/text v0.9.0 // indirect
)
