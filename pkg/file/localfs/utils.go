package localfs

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/my-network/fsutil/pkg/file"
)

func localToPathBytes(b []byte) file.Path {
	destinationPathParts := bytes.Split(b, []byte{filepath.Separator})

	result := make(file.Path, 0, len(destinationPathParts))
	for _, destinationPathPart := range destinationPathParts {
		result = append(result, string(destinationPathPart))
	}
	return result
}

func localToPath(s string) file.Path {
	return strings.Split(s, string(filepath.Separator))
}
