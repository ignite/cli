package time

// locationDefinitions are loaded during initialization using modified code from:
// https://cs.opensource.google/go/go/+/refs/tags/go1.23.0:src/time/tzdata/tzdata.go
var locationDefinitions = make(map[string][]byte)

func init() {
	const (
		zecheader = 0x06054b50
		zcheader  = 0x02014b50
		ztailsize = 22

		zheadersize = 30
		zheader     = 0x04034b50
	)

	z := zipdata

	idx := len(z) - ztailsize
	n := get2s(z[idx+10:])
	idx = get4s(z[idx+16:])

	for range n {
		// See time.loadTzinfoFromZip for zip entry layout.
		if get4s(z[idx:]) != zcheader {
			break
		}
		meth := get2s(z[idx+10:])
		size := get4s(z[idx+24:])
		namelen := get2s(z[idx+28:])
		xlen := get2s(z[idx+30:])
		fclen := get2s(z[idx+32:])
		off := get4s(z[idx+42:])
		zname := z[idx+46 : idx+46+namelen]
		idx += 46 + namelen + xlen + fclen

		if meth != 0 {
			panic("unsupported compression for " + zname + " in embedded tzdata")
		}

		// See time.loadTzinfoFromZip for zip per-file header layout.
		nidx := off
		if get4s(z[nidx:]) != zheader ||
			get2s(z[nidx+8:]) != meth ||
			get2s(z[nidx+26:]) != namelen {
			panic("corrupt embedded tzdata")
		}

		nxlen := get2s(z[nidx+28:])
		nidx += 30 + namelen + nxlen
		locationDefinitions[zname] = []byte(z[nidx : nidx+size])
	}
}

// get4s returns the little-endian 32-bit value at the start of s.
func get4s(s string) int {
	if len(s) < 4 {
		return 0
	}
	return int(s[0]) | int(s[1])<<8 | int(s[2])<<16 | int(s[3])<<24
}

// get2s returns the little-endian 16-bit value at the start of s.
func get2s(s string) int {
	if len(s) < 2 {
		return 0
	}
	return int(s[0]) | int(s[1])<<8
}

func X_loadFromEmbeddedTZData(name string) ([]byte, bool) {
	definition, ok := locationDefinitions[name]
	return definition, ok
}
