//go:generate msgp

package file

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/my-network/fsutil/pkg/mathutils"
	"github.com/tinylib/msgp/msgp"
)

// Path is OS-agnostic representation of a path.
//
// Each nesting level is in a separate string. Any characters are allowed in
// a name (including '/', '\\' and '\000').
type Path []string

func ParseLocalPath(pathString string) Path {
	return strings.Split(pathString, string(filepath.Separator))
}

func (p Path) LocalPath() string {
	result := filepath.Join([]string(p)...)
	if len(p) > 0 && p[0] == "" { // absolute path in UNIX-like systems
		return string(filepath.Separator) + result
	}
	return result
}

// Key is OS-agnostic key which represents this path.
//
// It is supposed to be used as a key for `map`-s.
func (p Path) Key() string {
	var buf bytes.Buffer
	w := msgp.NewWriter(&buf)
	err := p.EncodeMsg(w)
	if err != nil {
		panic(err)
	}
	err = w.Flush()
	if err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}

func (p Path) Append(appendie ...string) Path {
	p = append(p, appendie...)
	return p
}

func (p Path) Up() Path {
	if len(p) == 0 {
		return nil
	}

	return p[:len(p)-1]
}

func (p Path) RelativeTo(root Path) Path {
	minLen := mathutils.Min(len(p), len(root))
	var matchUntilIdx int
	for matchUntilIdx = 0; matchUntilIdx < minLen; matchUntilIdx++ {
		if p[matchUntilIdx] != root[matchUntilIdx] {
			break
		}
	}
	pRest := p[matchUntilIdx:]
	if matchUntilIdx == len(root) {
		return pRest
	}

	rootRest := root[matchUntilIdx:]
	result := make(Path, len(pRest)+len(rootRest))
	for i := 0; i < len(rootRest); i++ {
		result[i] = ".."
	}
	copy(result[len(rootRest):], pRest)
	return result
}
