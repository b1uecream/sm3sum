package main

import (
	"bufio"
	"bytes"
	"crypto/subtle"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tjfoc/gmsm/sm3"
)

const version = "1.0.0"

var (
	checkFile       = flag.String("c", "", "Read hash sums from the FILE and check them")
	binaryMode      = flag.Bool("b", false, "Read files in binary mode (default)")
	textMode        = flag.Bool("t", false, "Read files in text mode (ignored on Unix)")
	tagOutput       = flag.Bool("tag", false, "Create a BSD-style checksum")
	zeroMode        = flag.Bool("z", false, "End each output line with NUL, not newline")
	ignoreMissing   = flag.Bool("ignore-missing", false, "Ignore missing files in check mode")
	quietMode       = flag.Bool("quiet", false, "Don't print OK for successfully verified files")
	statusMode      = flag.Bool("status", false, "Don't output anything, status code shows success")
	strictMode      = flag.Bool("strict", false, "Exit non-zero for improperly formatted lines")
	warnMode        = flag.Bool("w", false, "Warn about improperly formatted lines")
	help            = flag.Bool("h", false, "Show help")
	showVersion     = flag.Bool("version", false, "Show version")
)

func usage() {
	fmt.Fprintf(os.Stderr, `sm3sum - SM3 hash calculator and verifier

Usage:
  sm3sum [OPTION]... [FILE]...
  sm3sum -c [FILE]

Options:
  -b, --binary         Read in binary mode (default)
  -c, --check FILE     Read SM3 sums from the FILE and check them
      --tag            Create a BSD-style checksum
  -t, --text           Read in text mode (ignored)
  -z, --zero           End lines with NUL instead of newline
      --ignore-missing Ignore missing files in check mode
      --quiet          Don't print OK for successfully verified files
      --status         Don't output anything, status code shows success
      --strict         Exit non-zero for improperly formatted lines
  -w, --warn           Warn about improperly formatted lines
  -h, --help           Display this help and exit
      --version        Display version information and exit
`)
}

func versionInfo() {
	fmt.Printf("sm3sum version %s\n", version)
}

func hashFile(filename string) ([]byte, error) {
	var reader io.Reader
	if filename == "-" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader = bufio.NewReader(f)
	}

	h := sm3.New()
	_, err := io.Copy(h, reader)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func outputHash(hash []byte, filename string) {
	var buf bytes.Buffer
	if *tagOutput {
		fmt.Fprintf(&buf, "SM3 (%s) = %x", filename, hash)
	} else {
		fmt.Fprintf(&buf, "%x  %s", hash, filename)
	}
	if *zeroMode {
		buf.WriteByte(0)
	} else {
		buf.WriteByte('\n')
	}
	os.Stdout.Write(buf.Bytes())
}

func checkFromFile(listFile string) {
	file, err := os.Open(listFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sm3sum: cannot open %s: %v\n", listFile, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	allOk := true
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		expected := ""
		filename := ""
		isValid := true

		if strings.HasPrefix(line, "SM3 (") && *tagOutput {
			parts := strings.SplitN(line, ") = ", 2)
			if len(parts) != 2 {
				isValid = false
			} else {
				filename = strings.TrimPrefix(parts[0], "SM3 (")
				expected = parts[1]
			}
		} else {
			parts := strings.SplitN(line, "  ", 2)
			if len(parts) != 2 {
				isValid = false
			} else {
				expected = parts[0]
				filename = parts[1]
			}
		}

		if !isValid {
			if *strictMode {
				fmt.Fprintf(os.Stderr, "sm3sum: %d: improperly formatted line\n", lineNum)
				os.Exit(1)
			}
			if *warnMode {
				fmt.Fprintf(os.Stderr, "sm3sum: %d: warning: invalid format\n", lineNum)
			}
			continue
		}

		hash, err := hashFile(filename)
		if err != nil {
			if *ignoreMissing && os.IsNotExist(err) {
				continue
			}
			if !*statusMode {
				fmt.Printf("%s: FAILED open or read: %v\n", filename, err)
			}
			allOk = false
			continue
		}

		actual := fmt.Sprintf("%x", hash)
		if subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) == 1 {
			if !*quietMode && !*statusMode {
				fmt.Printf("%s: OK\n", filename)
			}
		} else {
			if !*statusMode {
				fmt.Printf("%s: FAILED\n", filename)
			}
			allOk = false
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "sm3sum: error reading list: %v\n", err)
		os.Exit(1)
	}
	if !allOk {
		os.Exit(1)
	}
}

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "--help" {
			usage()
			return
		}
		if arg == "--version" {
			versionInfo()
			return
		}
		if arg == "--warn" {
			*warnMode = true
		}
	}

	flag.Parse()

	if *help {
		usage()
		return
	}
	if *showVersion {
		versionInfo()
		return
	}

	if *checkFile != "" {
		checkFromFile(*checkFile)
		return
	}

	files := flag.Args()
	if len(files) == 0 {
		usage()
		os.Exit(1)
	}

	for _, fname := range files {
		hash, err := hashFile(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sm3sum: %s: %v\n", fname, err)
			continue
		}
		outputHash(hash, fname)
	}
}
