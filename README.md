# go-cp
`go-cp` is a command-line file copying tool, built in Go.

## Usage
go-cp takes two required arguments: a source file/directory and a destination file/directory. The tool copies the contents of the source file/directory to the destination file/directory.

`go-cp [flags] source destination`

Flags

`-i` Prompt before overwriting existing files.

`-n` Do not overwrite existing files.

`-r` Copy directories recursively.

## Examples
Copy a file:

`go-cp file.txt copy.txt`

Copy a directory recursively:

`go-cp -r dir/ copy/`

Prompt before overwriting an existing file:

`go-cp -i file.txt copy.txt`

Do not overwrite an existing file:

`go-cp -n file.txt copy.txt`

## License
go-cp is licensed under the Apache License, Version 2.0.