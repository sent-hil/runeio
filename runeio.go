package runeio

import (
	"io"
	"io/ioutil"
	"unicode"
)

// RuneReader is the underlying interface Reader will use for its operations.
type RuneReader interface {
	ReadRune() (r rune, size int, err error)
	io.Reader
}

// Reader implements buffered manipulation of runes using RuneReader.
//
// NOTE: while calling Peek* methods won't affect result of ReadRunes(), it'll
// however read them from underlying RuneReader, but not put them back.
type Reader struct {
	RuneReader

	// Runes is temporary buffer used to store runes that are peeked, but not
	// yet read.
	Runes []rune
}

// NewReader is the required initializer for Reader.
func NewReader(r RuneReader) *Reader {
	return &Reader{r, []rune{}}
}

// Discard skips the given n runes, returning number of runes discarded.
//
// If given n is greater than amount of runes in the buffer, it'll discard all
// the runes and return `io.EOF`.
func (r *Reader) Discard(n uint) (uint, error) {
	runes, err := r.ReadRunes(n)
	return uint(len(runes)), err
}

// ReadRunes reads given n runes from buffers and returns slice of them.
//
// If given n is greater than amount of runes in the buffer, it'll return all
// the runes and `io.EOF` as error.
func (r *Reader) ReadRunes(n uint) (runes []rune, err error) {
	if err = r.readFromReader(n); err != nil {
		n = uint(len(r.Runes))
	}

	runes = r.Runes[0:n]
	r.Runes = r.Runes[n:]

	return runes, err
}

// ReadSingleRune reads a single rune from buffer and return it.
//
// If the are no runes left in the buffer, it'll return unicode.ReplacementChar
// and `io.EOF` error.
func (r *Reader) ReadSingleRune() (rune, error) {
	runes, err := r.ReadRunes(1)
	if err != nil {
		return unicode.ReplacementChar, err
	}

	return runes[0], nil
}

// ReadTill returns all the runes that matches the given matcher function.
func (r *Reader) ReadTill(matcherFn func(rune) bool) (runes []rune) {
	for {
		ru, err := r.PeekSingleRune()
		if err != nil || !matcherFn(ru) {
			break
		}

		r.ReadSingleRune()
		runes = append(runes, ru)
	}

	return runes
}

// PeekRunes peeks given n runes from buffers and returns slice of them. It does
// not however remove them the buffer and the same data will be returned on
// ReadRunes() operation.
//
// If given n is greater than amount of runes in the buffer, it'll peek all
// the runes and `io.EOF` as error.
func (r *Reader) PeekRunes(n uint) ([]rune, error) {
	if err := r.readFromReader(n); err != nil {
		return r.Runes, err
	}

	return r.Runes[0:n], nil
}

// PeekSingleRune peeks a single rune from buffer and return it.
//
// If the are no runes left in the buffer, it'll return unicode.ReplacementChar
// and `io.EOF` error.
func (r *Reader) PeekSingleRune() (rune, error) {
	runes, err := r.PeekRunes(1)
	if err != nil {
		return unicode.ReplacementChar, err
	}

	return runes[0], nil
}

// String returns ALL the unread runes in local buffer and underlying reader as
// a string.
//
// NOTE, it uses ioutil.ReadAll() to read runes from reader which may have
// performance issues depending on size of reader.
func (r *Reader) String() (string, error) {
	bites, err := ioutil.ReadAll(r.RuneReader)
	if err != nil {
		return "", err
	}
	return string(r.Runes) + string(bites), nil
}

// Reset replaces the underlying reader with the given reader.
func (r *Reader) Reset(bufReader RuneReader) {
	r.RuneReader = bufReader
}

// IsAtEnd returns if at the end of string, ie reading 1 more character
// would return `io.EOF` error.
func (r *Reader) IsAtEnd() bool {
	_, err := r.PeekSingleRune()
	return err == io.EOF
}

// readFromReader gets given x number of runes from underlying reader and stores
// it to make sure local buffer has n runes.
//
// If the are no runes left in the reader, it'll return `io.EOF` error.
func (r *Reader) readFromReader(n uint) error {
	l := int(n) - len(r.Runes)

	// check if we've already read enough runes
	if l <= 0 {
		return nil
	}

	// if not, read the remaining amount of runes
	for i := 0; i < l; i++ {
		ru, _, err := r.ReadRune()
		if err != nil {
			return err
		}
		r.Runes = append(r.Runes, ru)
	}

	return nil
}
