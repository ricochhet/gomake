# gomake
A Makefile alternative / task runner.

## Requirements
- Go 1.22 or later

## Build
- Run `go build -o gomake-bootstrap.exe`
- Run `gomake-bootstrap all` to build `gomake.exe` using the `make.gomake` file.

## Usage
- See [make.gomake](./make.gomake) for an example of a `gomake` gomake file.
- Run `gomake function_name` to execute the specified function.
- Optionally specify the file path: `gomake ./make.gomake function_name`
- Prepend `-dump` before additional commands to view the parsed function blocks.
    - Example: `gomake -dump ...`

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

# License
See LICENSE file.