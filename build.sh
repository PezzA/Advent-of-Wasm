cd aoc
GOOS=js GOARCH=wasm go build -o ../docs/aoc.wasm 
cd ../snowfield
GOOS=js GOARCH=wasm go build -o ../docs/snowfield.wasm
