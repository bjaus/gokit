// Package caller provides an abstraction on top of the `runtime.Caller`
// function to provide
package caller

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Caller holds information pertaining to the runtime caller
// in a stack trace.
type Caller struct {
	fn   string
	pkg  string
	path string
	line int
}

// Packge provides the package of the caller.
func (c Caller) Package() string { return c.pkg }

// Function provides the function/method of the caller.
func (c Caller) Function() string { return c.fn }

// LineNumber provides the line number within the file of the caller.
func (c Caller) LineNumber() int { return c.line }

// FilePath provides the location of the caller file on disk.
//
// The variadic basedir parameter provides the ability to
// indicate the directory name for which the filepath
// should begin. It is variadic in order to make this
// functionality optional. If nothing is provided or
// the value provided is not in the filepath then the
// the whole system filepath will be returned.
func (c Caller) FilePath(basedir ...string) string {

	// If the path is empty or no basedir parameter is provided
	// then just return the system filepath.
	if c.path == "" || len(basedir) == 0 {
		return c.path
	}

	// Strip the first provided basedir value and if it is
	// an empty string then just return the system filepath.
	directory := strings.TrimSpace(basedir[0])
	if directory == "" {
		return c.path
	}

	// If the filepath does not contain the dir value extracted
	// from the basedir value then there's no point in rebuilding.
	if !strings.Contains(c.path, directory) {
		return c.path
	}

	// Use the operating system's file path seperator value.
	// This keeps this method OS independent.
	sep := fmt.Sprintf("%c", os.PathSeparator)

	var parts []string

	// Rebuild the path looking for the dir value.
	for _, part := range strings.Split(c.path, sep) {
		if part == directory {
			// Wipe the parts slice to start from this directory.
			parts = []string{}
		}
		parts = append(parts, part)
	}

	return filepath.Join(parts...)
}

// replacer removes the noise which the runtime package adds.
var replacer = *strings.NewReplacer("(*", "", ")", "", ".go", "")

// New creates a Caller value via the runtime package.
func New(skip int) Caller {

	// Setup invalid data to indicate unsucessful call.
	unk := "unknown"
	c := Caller{
		pkg:  unk,
		fn:   unk,
		path: unk,
		line: -1,
	}

	if skip < 0 {
		skip = 0
	}
	// Add one to account for this function in the call stack.
	skip++

	pc, path, line, ok := runtime.Caller(skip)
	if !ok {
		return c
	}

	pcfn := runtime.FuncForPC(pc)
	_, fp := filepath.Split(replacer.Replace(pcfn.Name()))

	// Remove the `.func` part if provided by runtime package.
	parts := strings.Split(fp, ".func")
	// Take only the first value since we don't want `.func`.
	caller := parts[0]

	// Seperate the package from what remains.
	parts = strings.SplitN(caller, ".", 2)

	switch len(parts) {
	case 0:
		return c
	case 1:
		// Somehow there wasn't a function and/or method...?
		c.pkg = parts[0]
	default:
		c.pkg = parts[0]
		c.fn = parts[1]
	}

	// Add the file path and line number now that we've
	// successfully made it to this point.
	c.path = path
	c.line = line

	return c
}
