# Readme 
`

## Compiling the dynlib:

```sh
go build -buildmode=c-shared -o libgophersat.dll binding.go
```

## Cross-compilation

- All supported architectures:

```sh
go tool dist list
```

- Examples of compilations:

```
GOOS=windows GOARCH=386 go -o libgophersat.dll binding.go
GOOS=darwin GOARCH=amd64 go -o libgophersat.dylib binding.go
GOOS=linux GOARCH=arm go -o libgophersat.so binding.go
```