# runeio

runeio is a library for manipulating runes from an underlying io.Reader.

[![Documentation](https://godoc.org/github.com/sent-hil/runeio?status.svg)](https://godoc.org/github.com/sent-hil/runeio)

## Getting started

    // See https://godoc.org/github.com/sent-hil/runeio#RuneReader
    // for `RuneReader` interface.
    //
    // `bufio.Reader`, `bytes.Reader` and `strings.Reader` all implement the
    // interface  and can be used here.
    buf := bufio.NewStringReader("Hello World")
    runeio.NewRuneio(buf)

## Install

    go get -u github.com/sent-hil/runeio
