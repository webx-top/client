package upload

import (
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
	svg "github.com/h2non/go-is-svg"
	"github.com/hamochi/safesvg"
)

var ReadHeadSizeBytes = 261
var SVGMaxSizeBytes int64 = 5 * 1024 * 1024
var ErrIncorrectFileFormat = errors.New(`file format is incorrect`)
var defaultSafeSVGValidator = safesvg.NewValidator()
var svgHeadRegex = regexp.MustCompile(`(?i)^\s*(?:<\?xml[^>]*>\s*)?(?:<!doctype svg[^>]*>\s*)?<svg[^>]*>`)

func ReadHeadBytes(r io.Reader, readSizes ...int) ([]byte, error) {
	readSize := ReadHeadSizeBytes
	if len(readSizes) > 0 && readSizes[0] > 0 {
		readSize = readSizes[0]
	}
	// We only have to pass the file header = first 261 bytes
	head := make([]byte, readSize)
	_, err := r.Read(head)
	return head, err
}

func IsImage(b []byte) bool {
	return filetype.IsImage(b)
}

func IsSVGImage(b []byte) bool {
	return svg.Is(b)
}

func ValidateSVGImage(b []byte) error {
	return defaultSafeSVGValidator.Validate(b)
}

func IsVideo(b []byte) bool {
	return filetype.IsVideo(b)
}

func IsAudio(b []byte) bool {
	return filetype.IsAudio(b)
}

func IsFont(b []byte) bool {
	return filetype.IsFont(b)
}

func IsArchive(b []byte) bool {
	return filetype.IsArchive(b)
}

func IsDocument(b []byte) bool {
	return filetype.IsDocument(b)
}

func IsApplication(b []byte) bool {
	return filetype.IsApplication(b)
}

func IsType(b []byte, expected types.Type) bool {
	kind, _ := filetype.Match(b)
	return kind == expected
}

func IsTypeStringByReader(rd io.Reader, expected string) error {
	head, err := ReadHeadBytes(rd)
	if err != nil {
		return err
	}
	if sk, ok := rd.(io.ReadSeeker); ok {
		sk.Seek(0, 0)
	}
	if IsTypeString(head, expected) {
		return nil
	}
	if expected != `image` {
		return ErrIncorrectFileFormat
	}
	if !svgHeadRegex.Match(head) {
		return ErrIncorrectFileFormat
	}
	// SVG
	buf, err := io.ReadAll(io.LimitReader(rd, SVGMaxSizeBytes))
	if err != nil {
		return err
	}
	if !IsSVGImage(buf) {
		return ErrIncorrectFileFormat
	}
	err = ValidateSVGImage(buf)
	if err != nil {
		return err
	}
	if sk, ok := rd.(io.ReadSeeker); ok {
		sk.Seek(0, 0)
	}
	return nil
}

func IsTypeString(b []byte, expected string) bool {
	switch expected {
	case `image`:
		return IsImage(b)
	case `video`:
		return IsVideo(b)
	case `audio`:
		return IsAudio(b)
	case `archive`:
		return IsArchive(b)
	case `document`, `office`:
		return IsDocument(b)
	case `doc`:
		return IsType(b, matchers.TypeDoc) || IsType(b, matchers.TypeDocx)
	case `ppt`:
		return IsType(b, matchers.TypePpt) || IsType(b, matchers.TypePptx)
	case `xls`:
		return IsType(b, matchers.TypeXls) || IsType(b, matchers.TypeXlsx)
	case `file`:
		return IsApplication(b)
	case `font`:
		return IsFont(b)
	case `pdf`:
		return IsType(b, matchers.TypePdf)
	case `photoshop`:
		return IsType(b, matchers.TypePsd)
	default:
		return false
	}
}

func IsSupported(extension string) bool {
	extension = strings.TrimPrefix(extension, `.`)
	return filetype.IsSupported(extension)
}

func IsMIMESupported(mime string) bool {
	return filetype.IsMIMESupported(mime)
}
