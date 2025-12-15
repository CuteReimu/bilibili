module github.com/CuteReimu/bilibili/v2

go 1.23.0

require (
	github.com/Baozisoftware/qrcode-terminal-go v0.0.0-20170407111555-c0650d8dff0f
	github.com/go-resty/resty/v2 v2.17.1
	github.com/pkg/errors v0.9.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/spf13/cast v1.10.0
	golang.org/x/sync v0.16.0
)

require (
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
)

retract v2.3.0 // 接口调整
