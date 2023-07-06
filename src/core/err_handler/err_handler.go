package err_handler

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/we7coreteam/w7-rangine-go-support/src/facade"
	"golang.org/x/xerrors"
	"math"
	"os"
	"runtime"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

type CustomHandler func(err error)

var DefaultHandler = func(err error) {
	logger, logErr := facade.GetLoggerFactory().Channel("default")
	if logErr == nil {
		if _, ok := err.(xerrors.Wrapper); ok {
			logger.Error(fmt.Sprintf("%s \n %s", err.Error(), string(Stack(3, 0))))
		} else {
			logger.Error(fmt.Sprintf("%s \n %s", err.Error(), string(Stack(3, 0))))
		}
	}
}

func SetHandler(handler CustomHandler) {
	DefaultHandler = handler
}

func Throw(message string, previous error) error {
	var err error
	if previous == nil {
		err = errors.New(message)
	} else {
		err = errors.Wrap(previous, message)
	}

	return err
}

func Found(err error) bool {
	if err != nil {
		return true
	}
	return false
}

func Handle(err error) {
	DefaultHandler(err)
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func Stack(skip int, until int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	if until <= 0 {
		until = math.MaxInt16
	}

	for i := skip; i < until; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contain dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
