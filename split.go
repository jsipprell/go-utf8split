/*
 * Copyright (c) 2014-2015 Jesse Sipprell <jessesipprell@gmail.com>
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 */

// utf-8 splitting utility to split a valid utf8 string into a number
// of substrings as delimited by some artibrary number of utf-8
// delimiters.

package utf8split // "github.com/jsipprell/go-utf8split"

import (
	"bytes"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Splitter struct {
	ranges []*unicode.RangeTable
}

func insertRune(rt *unicode.RangeTable, r rune) *unicode.RangeTable {
	if uint32(r) > uint32(0xffff) {
		r32 := uint32(r)
		for i, _ := range rt.R32 {
			if rt.R32[i].Stride != 0 {
				range32 := &rt.R32[i]
				if range32.Lo > 0 && r32 == range32.Lo-1 {
					range32.Lo--
					return rt
				}
				if range32.Hi < uint32(0xffffffff) && r32 == range32.Hi+1 {
					range32.Hi++
					return rt
				}
			}
		}
		rt.R32 = append(rt.R32, unicode.Range32{Lo: r32, Hi: r32, Stride: 1})
	} else {
		r16 := uint16(r)
		if uint32(r) <= unicode.MaxLatin1 {
			rt.LatinOffset++
		}
		for i, _ := range rt.R16 {
			if rt.R16[i].Stride != 0 {
				range16 := &rt.R16[i]
				if range16.Lo > 0 && r16 == range16.Lo-1 {
					range16.Lo--
					return rt
				}
				if range16.Hi < uint16(0xffff) && r16 == range16.Hi+1 {
					range16.Hi++
					return rt
				}
			}
		}
		rt.R16 = append(rt.R16, unicode.Range16{Lo: r16, Hi: r16, Stride: 1})
	}

	return rt
}

func insertRunes(rt *unicode.RangeTable, runes ...rune) *unicode.RangeTable {
	for _, r := range runes {
		rt = insertRune(rt, r)
	}
	return rt
}

func addRangeTable(rt *unicode.RangeTable, seps ...[]byte) *unicode.RangeTable {
	for _, b := range seps {
		if !utf8.FullRune(b) {
			log.Panicf("invalid utf8 encoding '%v'", b)
		}
		runes := make([]rune, 0, len(b))
		for len(b) > 0 {
			r, sz := utf8.DecodeRune(b)
			if !utf8.ValidRune(r) {
				b = b[1:]
				continue
			}
			runes = append(runes, r)
			b = b[sz:]
		}
		if len(runes) > 0 {
			rt = insertRunes(rt, runes...)
		}
	}

	return rt
}

func makeStringRangeTable(seps ...string) (rt *unicode.RangeTable) {
	rt = new(unicode.RangeTable)
	for _, s := range seps {
		if !utf8.ValidString(s) {
			log.Panicf("invalid utf8 encoding '%v'", []byte(s))
		}
		rt = makeRangeTable([]byte(s))
	}
	return
}

func makeRangeTable(seps ...[]byte) *unicode.RangeTable {
	return addRangeTable(new(unicode.RangeTable), seps...)
}

// Returns a new Splitter which will split strings or byte
// slices by utf8 delimiters, specified as []byte slices.
func New(seps ...[]byte) *Splitter {
	return &Splitter{[]*unicode.RangeTable{makeRangeTable(seps...)}}
}

// Returns a new Splitter which will split strings or byte
// slices by utf8 deliiters, specified as strings
func WithDelimiters(seps ...string) *Splitter {
	return &Splitter{[]*unicode.RangeTable{makeStringRangeTable(seps...)}}
}

// Returns true if a rune is one of the separators handled by the splitter
func (s *Splitter) In(r rune) bool {
	return unicode.In(r, s.ranges...)
}

// Returns true if any rune passed is one of the separators handled by the spliter
func (s *Splitter) AnyIn(runes ...rune) bool {
	for _, r := range runes {
		if unicode.In(r, s.ranges...) {
			return true
		}
	}
	return false
}

// Returns true if all runes passed are in one of the separators handled by the splitter
func (s *Splitter) AllIn(runes ...rune) bool {
	var c int
	for _, r := range runes {
		if !unicode.In(r, s.ranges...) {
			return false
		}
		c++
	}

	return c > 0
}

// Returns a slice of []byte slices as delimited by the utf8 characters the splitter
// was initialized with. Repeated delimiters are concatenated and thus trimmed.
func (s *Splitter) Split(src []byte) [][]byte {
	return bytes.FieldsFunc(src, s.In)
}

// Returns a slice of strings as delimited by the utf8 characters the splitter
// was initialized with. Repeated delimiters are concatenated and thus trimmed.
func (s *Splitter) SplitString(src string) []string {
	return strings.FieldsFunc(src, s.In)
}

// Given a source byte slice, split it into slices as delimited by
// any arbitrary number of utf8 delimiters. Repeated delimiters are
// concatenated and thus trimmed.
func Bytes(src []byte, delims []byte, addl ...[]byte) [][]byte {
	rangeTable := makeRangeTable(delims)
	if len(addl) > 0 {
		rangeTable = addRangeTable(rangeTable, addl...)
	}
	if len(rangeTable.R16)+len(rangeTable.R32) > 0 {
		return bytes.FieldsFunc(src, func(r rune) bool {
			return unicode.In(r, rangeTable)
		})
	}
	return nil
}

// Given a source string, split it into slices as delimited by
// any arbitrary number of utf8 delimiters. Repeated delimiters
// are concatenated and thus trimmed.
func Strings(src string, delims string, addl ...string) []string {
	rangeTable := makeStringRangeTable(delims)
	for _, s := range addl {
		rangeTable = addRangeTable(rangeTable, []byte(s))
	}
	if len(rangeTable.R16)+len(rangeTable.R32) > 0 {
		return strings.FieldsFunc(src, func(r rune) bool {
			return unicode.In(r, rangeTable)
		})
	}
	return nil
}
