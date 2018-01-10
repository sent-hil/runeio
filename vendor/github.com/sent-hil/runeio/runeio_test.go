package runeio

import (
	"bytes"
	"io"
	"testing"
	"unicode"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRuneIo(t *testing.T) {
	Convey("RuneIo", t, func() {
		hw := NewReader(bytes.NewBufferString("Hello World"))
		om := NewReader(bytes.NewBufferString("H"))
		em := NewReader(bytes.NewBufferString(""))

		Convey("NewReader", func() {
			Convey("It returns initialized Reader", func() {
				So(hw, ShouldHaveSameTypeAs, &Reader{})
			})
		})

		Convey("Discard", func() {
			Convey("It discards given length of runes", func() {
				discarded, err := hw.Discard(1)
				So(err, ShouldEqual, nil)
				So(discarded, ShouldEqual, 1)

				str, err := hw.String()
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "ello World")
			})

			Convey("It discards all runes when given length is same length as reader", func() {
				discarded, err := hw.Discard(11)
				So(err, ShouldEqual, nil)
				So(discarded, ShouldEqual, 11)

				str, err := hw.String()
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "")
			})

			Convey("It returns io.EOF when given length is greater than length in reader", func() {
				discarded, err := hw.Discard(12)
				So(err, ShouldEqual, io.EOF)
				So(discarded, ShouldEqual, 11)

				str, err := hw.String()
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "")

				discarded, err = em.Discard(1)
				So(err, ShouldEqual, io.EOF)
				So(discarded, ShouldEqual, 0)
			})
		})

		Convey("ReadRunes", func() {
			Convey("It discards given length of runes", func() {
				runes, err := hw.ReadRunes(1)
				So(err, ShouldBeNil)
				So(runes, ShouldHaveSameTypeAs, []rune{})
				So(string(runes), ShouldResemble, "H")
			})

			Convey("It returns all runes when given length is same length as reader", func() {
				runes, err := hw.ReadRunes(11)
				So(err, ShouldBeNil)
				So(string(runes), ShouldEqual, "Hello World")
			})

			Convey("It returns io.EOF when given length is greater than length in reader", func() {
				runes, err := hw.ReadRunes(12)
				So(err, ShouldEqual, io.EOF)
				So(string(runes), ShouldEqual, "Hello World")
			})

			Convey("It removes runes from reader", func() {
				_, err := hw.ReadRunes(11)
				So(err, ShouldBeNil)

				str, err := hw.String()
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "")
			})
		})

		Convey("ReadSingleRune", func() {
			Convey("It returns single rune from reader", func() {
				h, err := hw.ReadSingleRune()
				So(err, ShouldBeNil)
				So(string(h), ShouldEqual, "H")
			})

			Convey("It returns io.EOF if at end of reader", func() {
				_, err := em.ReadSingleRune()
				So(err, ShouldEqual, io.EOF)
			})

			Convey("It returns last char when index is at end of reader", func() {
				h, err := om.ReadSingleRune()
				So(err, ShouldBeNil)
				So(string(h), ShouldEqual, "H")
			})
		})

		Convey("ReadTill", func() {
			Convey("It returns no runes if 1st rune does not match", func() {
				runes := hw.ReadTill(func(r rune) bool { return r == 'o' })
				So(len(runes), ShouldEqual, 0)
			})

			Convey("It returns single rune if only 1st rune matches", func() {
				runes := hw.ReadTill(func(r rune) bool { return r == 'H' })
				So(len(runes), ShouldEqual, 1)
				So(string(runes[0]), ShouldEqual, "H")
			})

			Convey("It returns all runes till io.EOF if all runes matches", func() {
				runes := hw.ReadTill(func(r rune) bool {
					return unicode.IsLetter(r) || unicode.IsSpace(r)
				})
				So(len(runes), ShouldEqual, 11)
				So(string(runes), ShouldEqual, "Hello World")
			})

			Convey("It unreads rune if it doesn't match", func() {
				runes := hw.ReadTill(func(r rune) bool {
					return unicode.IsLetter(r)
				})
				So(len(runes), ShouldEqual, 5)

				str, err := hw.String()
				So(err, ShouldBeNil)
				So(str, ShouldEqual, " World")
			})
		})

		Convey("PeekRunes", func() {
			Convey("It returns given length of runes", func() {
				runes, err := hw.PeekRunes(1)
				So(err, ShouldBeNil)
				So(runes, ShouldHaveSameTypeAs, []rune{})
				So(string(runes), ShouldResemble, "H")
			})

			Convey("It returns io.EOF when given length is greater than length in reader", func() {
				runes, err := hw.PeekRunes(12)
				So(err, ShouldEqual, io.EOF)
				So(string(runes), ShouldEqual, "Hello World")
			})

			Convey("It does not remove runes from reader", func() {
				_, err := hw.PeekRunes(1)
				So(err, ShouldBeNil)

				str, err := hw.String()
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "Hello World")
			})
		})

		Convey("PeekSingleRune", func() {
			Convey("It returns single rune from reader", func() {
				h, err := hw.PeekSingleRune()
				So(err, ShouldBeNil)
				So(string(h), ShouldEqual, "H")
			})

			Convey("It does not remove runes from reader", func() {
				_, err := hw.PeekSingleRune()
				So(err, ShouldBeNil)

				str, err := hw.String()
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "Hello World")
			})

			Convey("It returns io.EOF if at end of reader", func() {
				_, err := em.PeekSingleRune()
				So(err, ShouldEqual, io.EOF)
			})

			Convey("It returns last char when index is at end of reader", func() {
				h, err := om.PeekSingleRune()
				So(err, ShouldBeNil)
				So(string(h), ShouldEqual, "H")
			})
		})

		Convey("Reset", func() {
			Convey("It resets reader to given reader", func() {
				nr := bytes.NewBufferString("New")
				hw.Reset(nr)

				runes, err := hw.PeekRunes(3)
				So(err, ShouldBeNil)
				So(string(runes), ShouldEqual, "New")
			})
		})

		Convey("IsAtEnd", func() {
			Convey("It returns if there are more chars left to read", func() {
				So(om.IsAtEnd(), ShouldEqual, false)
			})

			Convey("It returns if no more chars left to be read", func() {
				So(em.IsAtEnd(), ShouldEqual, true)
			})
		})
	})
}
