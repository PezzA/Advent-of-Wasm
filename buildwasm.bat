set GOARCH=wasm
set GOOS=js
go build -o  .\docs\aoc.wasm .\aoc\main.go
go build -o  .\docs\snowfield.wasm .\snowfield\main.go





