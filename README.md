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

```
## This is a comment

example_function() {
    echo Hello, World!
}

example_caller_in_function() {
    @example_function
}

example_directory_switch_function() {
    mkdir mydir1
    ~./mydir1
    mkdir subdir1
}

example_function_with_params({param1}, {param2}) {
    echo {param1}, {param2}!
}
```

...

```
gomake example_function

gomake example_caller_in_function

gomake example_directory_switch_function

gomake example_function_with_params Hello World
```

# License
See LICENSE file.