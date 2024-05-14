fmt() {
    gofumpt -l -w .
}

lint() {
    golangci-lint run  
}

test() {
    go test ./...
}

build() {
    go build -o gomake.exe
}

prebuild() {
    @fmt
    @lint
    @test
}

all() {
    @fmt
    @lint
    @test
    @build
}