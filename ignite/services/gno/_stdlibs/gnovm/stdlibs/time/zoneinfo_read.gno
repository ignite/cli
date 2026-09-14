// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Parse "zoneinfo" time zone file.
// This is a fairly standard file format used on OS X, Linux, BSD, Sun, and others.
// See tzfile(5), https://en.wikipedia.org/wiki/Zoneinfo,
// and ftp://munnari.oz.au/pub/oldtz/

package time

import (
	"errors"
)

// Simple I/O interface to binary blob of data.
type dataIO struct {
	p     []byte
	error bool
}

func (d *dataIO) read(n int) []byte {
	if len(d.p) < n {
		d.p = nil
		d.error = true
		return nil
	}
	p := d.p[0:n]
	d.p = d.p[n:]
	return p
}

func (d *dataIO) big4() (n uint32, ok bool) {
	p := d.read(4)
	if len(p) < 4 {
		d.error = true
		return 0, false
	}
	return uint32(p[3]) | uint32(p[2])<<8 | uint32(p[1])<<16 | uint32(p[0])<<24, true
}

func (d *dataIO) big8() (n uint64, ok bool) {
	n1, ok1 := d.big4()
	n2, ok2 := d.big4()
	if !ok1 || !ok2 {
		d.error = true
		return 0, false
	}
	return (uint64(n1) << 32) | uint64(n2), true
}

func (d *dataIO) byte() (n byte, ok bool) {
	p := d.read(1)
	if len(p) < 1 {
		d.error = true
		return 0, false
	}
	return p[0], true
}

// rest returns the rest of the data in the buffer.
func (d *dataIO) rest() []byte {
	r := d.p
	d.p = nil
	return r
}

// Make a string by stopping at the first NUL
func byteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

var errBadData = errors.New("malformed time zone information")

// LoadLocationFromTZData returns a Location with the given name
// initialized from the IANA Time Zone database-formatted data.
// The data should be in the format of a standard IANA time zone file
// (for example, the content of /etc/localtime on Unix systems).
func LoadLocationFromTZData(name string, data []byte) (*Location, error) {
	d := dataIO{data, false}

	// 4-byte magic "TZif"
	if magic := d.read(4); string(magic) != "TZif" {
		return nil, errBadData
	}

	// 1-byte version, then 15 bytes of padding
	var version int
	var p []byte
	if p = d.read(16); len(p) != 16 {
		return nil, errBadData
	} else {
		switch p[0] {
		case 0:
			version = 1
		case '2':
			version = 2
		case '3':
			version = 3
		default:
			return nil, errBadData
		}
	}

	// six big-endian 32-bit integers:
	//	number of UTC/local indicators
	//	number of standard/wall indicators
	//	number of leap seconds
	//	number of transition times
	//	number of local time zones
	//	number of characters of time zone abbrev strings
	const (
		NUTCLocal = iota
		NStdWall
		NLeap
		NTime
		NZone
		NChar
	)
	var n [6]int
	for i := 0; i < 6; i++ {
		nn, ok := d.big4()
		if !ok {
			return nil, errBadData
		}
		if uint32(int(nn)) != nn {
			return nil, errBadData
		}
		n[i] = int(nn)
	}

	// If we have version 2 or 3, then the data is first written out
	// in a 32-bit format, then written out again in a 64-bit format.
	// Skip the 32-bit format and read the 64-bit one, as it can
	// describe a broader range of dates.

	is64 := false
	if version > 1 {
		// Skip the 32-bit data.
		skip := n[NTime]*4 +
			n[NTime] +
			n[NZone]*6 +
			n[NChar] +
			n[NLeap]*8 +
			n[NStdWall] +
			n[NUTCLocal]
		// Skip the version 2 header that we just read.
		skip += 4 + 16
		d.read(skip)

		is64 = true

		// Read the counts again, they can differ.
		for i := 0; i < 6; i++ {
			nn, ok := d.big4()
			if !ok {
				return nil, errBadData
			}
			if uint32(int(nn)) != nn {
				return nil, errBadData
			}
			n[i] = int(nn)
		}
	}

	size := 4
	if is64 {
		size = 8
	}

	// Transition times.
	txtimes := dataIO{d.read(n[NTime] * size), false}

	// Time zone indices for transition times.
	txzones := d.read(n[NTime])

	// Zone info structures
	zonedata := dataIO{d.read(n[NZone] * 6), false}

	// Time zone abbreviations.
	abbrev := d.read(n[NChar])

	// Leap-second time pairs
	d.read(n[NLeap] * (size + 4))

	// Whether tx times associated with local time types
	// are specified as standard time or wall time.
	isstd := d.read(n[NStdWall])

	// Whether tx times associated with local time types
	// are specified as UTC or local time.
	isutc := d.read(n[NUTCLocal])

	if d.error { // ran out of data
		return nil, errBadData
	}

	var extend string
	rest := d.rest()
	if len(rest) > 2 && rest[0] == '\n' && rest[len(rest)-1] == '\n' {
		extend = string(rest[1 : len(rest)-1])
	}

	// Now we can build up a useful data structure.
	// First the zone information.
	//	utcoff[4] isdst[1] nameindex[1]
	nzone := n[NZone]
	if nzone == 0 {
		// Reject tzdata files with no zones. There's nothing useful in them.
		// This also avoids a panic later when we add and then use a fake transition (golang.org/issue/29437).
		return nil, errBadData
	}
	zones := make([]zone, nzone)
	for i := range zones {
		var ok bool
		var n uint32
		if n, ok = zonedata.big4(); !ok {
			return nil, errBadData
		}
		if uint32(int(n)) != n {
			return nil, errBadData
		}
		zones[i].offset = int(int32(n))
		var b byte
		if b, ok = zonedata.byte(); !ok {
			return nil, errBadData
		}
		zones[i].isDST = b != 0
		if b, ok = zonedata.byte(); !ok || int(b) >= len(abbrev) {
			return nil, errBadData
		}
		zones[i].name = byteString(abbrev[b:])
	}

	// Now the transition time info.
	tx := make([]zoneTrans, n[NTime])
	for i := range tx {
		var n int64
		if !is64 {
			if n4, ok := txtimes.big4(); !ok {
				return nil, errBadData
			} else {
				n = int64(int32(n4))
			}
		} else {
			if n8, ok := txtimes.big8(); !ok {
				return nil, errBadData
			} else {
				n = int64(n8)
			}
		}
		tx[i].when = n
		if int(txzones[i]) >= len(zones) {
			return nil, errBadData
		}
		tx[i].index = txzones[i]
		if i < len(isstd) {
			tx[i].isstd = isstd[i] != 0
		}
		if i < len(isutc) {
			tx[i].isutc = isutc[i] != 0
		}
	}

	if len(tx) == 0 {
		// Build fake transition to cover all time.
		// This happens in fixed locations like "Etc/GMT0".
		tx = append(tx, zoneTrans{when: alpha, index: 0})
	}

	// Committed to succeed.
	l := &Location{zone: zones, tx: tx, name: name, extend: extend}

	// Fill in the cache with information about right now,
	// since that will be the most common lookup.
	sec, _, _ := now()
	for i := range tx {
		if tx[i].when <= sec && (i+1 == len(tx) || sec < tx[i+1].when) {
			l.cacheStart = tx[i].when
			l.cacheEnd = omega
			l.cacheZone = &l.zone[tx[i].index]
			if i+1 < len(tx) {
				l.cacheEnd = tx[i+1].when
			} else if l.extend != "" {
				// If we're at the end of the known zone transitions,
				// try the extend string.
				if name, offset, estart, eend, isDST, ok := tzset(l.extend, l.cacheStart, sec); ok {
					l.cacheStart = estart
					l.cacheEnd = eend
					// Find the zone that is returned by tzset to avoid allocation if possible.
					if zoneIdx := findZone(l.zone, name, offset, isDST); zoneIdx != -1 {
						l.cacheZone = &l.zone[zoneIdx]
					} else {
						l.cacheZone = &zone{
							name:   name,
							offset: offset,
							isDST:  isDST,
						}
					}
				}
			}
			break
		}
	}

	return l, nil
}

func findZone(zones []zone, name string, offset int, isDST bool) int {
	for i, z := range zones {
		if z.name == name && z.offset == offset && z.isDST == isDST {
			return i
		}
	}
	return -1
}

// get4 returns the little-endian 32-bit value in b.
func get4(b []byte) int {
	if len(b) < 4 {
		return 0
	}
	return int(b[0]) | int(b[1])<<8 | int(b[2])<<16 | int(b[3])<<24
}

// get2 returns the little-endian 16-bit value in b.
func get2(b []byte) int {
	if len(b) < 2 {
		return 0
	}
	return int(b[0]) | int(b[1])<<8
}

// loadLocation returns the Location with the given name from the
// embedded data.The first timezone data matching the given name
// that is successfully loaded and parsed is returned as a Location.
func loadLocation(name string) (*Location, error) {
	zoneData, ok := loadFromEmbeddedTZData(name)
	if !ok {
		return nil, errors.New("unknown time zone " + name)
	}

	if loc, err := LoadLocationFromTZData(name, zoneData); err == nil {
		return loc, nil
	}

	return nil, errors.New("unknown time zone " + name)
}
