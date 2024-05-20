# gomake
A Makefile alternative / task runner.

## Requirements
- Go 1.22 or later

## Build
- Run `go build -o gomake-bootstrap.exe`
- Run `gomake-bootstrap -run all` to build `gomake.exe` using the `make.gomake` file.

## Usage
- See [make.gomake](./make.gomake) for an example of a `gomake` file.
- Run `gomake -h` for a list of all commands.

### Functions

```
# This is a comment.
task() {
    ...
}
```

### Comparison
```
@(eq:"aaa","bbb")
# aaa == bbb = false

@(neq:"aaa","bbb")
# aaa != bbb = true
```

### Directory
```
@(cd:"./path/to/directory/")
```

### Operating System
```
# Command only runs on windows
@(os:"windows")

# Command runs on all platforms
@(os:"all")
```

### Environment Variables
```
echo %{GOPATH}
# Echo ./path/to/go/bin

@(env:"GOPATH=./other/path/to/go/bin")
# Set GOPATH to ./other/path/to/go/bin during runtime.
```

# License
See LICENSE file.