package multi

// IntegerKey constrains grid and layer ids to unsigned integer key types.
//
// Use this for map-backed ids where callers may choose compact key widths
// such as uint8/uint16 or wider widths such as uint64.
type IntegerKey interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
