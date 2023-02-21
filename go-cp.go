package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// Define command-line flags
	interactiveFlag := flag.Bool("i", false, "prompt before overwriting existing files")
	noClobberFlag := flag.Bool("n", false, "do not overwrite existing files")
	recursiveFlag := flag.Bool("r", false, "copy directories recursively")
	// Parse command-line flags
	flag.Parse()

	// Check if there are two arguments (source and destination)
	if flag.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "error: missing source or destination file")
		os.Exit(1)
	}

	// Get source and destination paths
	src := flag.Arg(0)
	dst := flag.Arg(1)

	// Check if source and destination are the same file
	srcInfo, err := os.Stat(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	dstInfo, err := os.Stat(dst)
	if err == nil && os.SameFile(srcInfo, dstInfo) {
		fmt.Fprintln(os.Stderr, "error: source and destination files are the same")
		os.Exit(1)
	}

	// Check if source file exists
	if _, err := os.Stat(src); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: source file '%s' does not exist\n", src)
		os.Exit(1)
	}

	// Check if destination file already exists
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		if *noClobberFlag {
			fmt.Fprintf(os.Stderr, "error: file '%s' already exists\n", dst)
			os.Exit(1)
		} else if *interactiveFlag {
			fmt.Printf("overwrite '%s'? (y/n [n]) ", dst)
			var answer string
			fmt.Scanln(&answer)
			if answer != "y" && answer != "Y" {
				fmt.Fprintln(os.Stderr, "not overwritten")
				os.Exit(1)
			}
		}
	}

	// Check if source is a directory
	srcIsDir := false
	srcInfo, err = os.Stat(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if srcInfo.IsDir() {
		srcIsDir = true
		if !*recursiveFlag {
			fmt.Fprintln(os.Stderr, "error: cannot copy directory without recursive flag (-r)")
			os.Exit(1)
		}
	}

	// Copy file or directory
	if !srcIsDir {
		copyFile(src, dst)
	} else {
		copyDir(src, dst, *noClobberFlag, *interactiveFlag)
	}
}

// copyFile copies the contents of the file at src to the file at dst.
// copyFile copies the contents of the file at src to the file at dst.
func copyFile(src, dst string) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	srcFileInfo, err := os.Stat(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	err = os.Chmod(dst, srcFileInfo.Mode())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// copyDir copies the contents of the directory at src to the directory at dst.
func copyDir(src, dst string, noClobber, interactive bool) {
	// Get source directory contents
	srcFiles, err := os.ReadDir(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	// Create destination directory
	err = os.MkdirAll(dst, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Copy each file in the source directory to the destination directory
	for _, file := range srcFiles {
		srcFile := filepath.Join(src, file.Name())
		dstFile := filepath.Join(dst, file.Name())

		if file.IsDir() {
			copyDir(srcFile, dstFile, noClobber, interactive)
		} else {
			if _, err := os.Stat(dstFile); !os.IsNotExist(err) {
				if noClobber {
					fmt.Fprintf(os.Stderr, "error: file '%s' already exists\n", dstFile)
					os.Exit(1)
				} else if interactive {
					fmt.Printf("overwrite '%s'? (y/n [n]) ", dstFile)
					var answer string
					fmt.Scanln(&answer)
					if answer != "y" && answer != "Y" {
						fmt.Fprintln(os.Stderr, "not overwritten")
						continue
					}
				}
			}
			copyFile(srcFile, dstFile)
		}
	}
}
