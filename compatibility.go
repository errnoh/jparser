package jparser

// TODO: null voi olla mikä tahansa muu asia.

type Compatibility uint8

const (
	Equal  = 0
	Subset = Compatibility(1 << iota)
	Superset
	Incompatible
)

func (c Compatibility) String() string {
	// incompatible and superset could be both thought as being incompatible
	// since you'll lose fields if you unmarshal 'a' into 'b'
	switch {
	case c == Equal:
		return "equal"
	case c&Incompatible == Incompatible:
		return "incompatible"
	case c == Subset|Superset:
		return "subset|superset"
	case c == Subset:
		return "subset"
	case c == Superset:
		return "superset"
	}
	return "???"
}

// Compare checks if two JSON blops are equal or compatible.
// It can answer questions like "Is 'a' subset of 'b'?"
func Compare(a, b *Token) (compatible Compatibility) {
	// Assume that NULL values are supported.
	// TODO: maybe have boolean switch to disable null values: (a.tokenType != NULL && b.tokenType != NULL) -blocks
	if a.tokenType != b.tokenType && (a.tokenType != NULL && b.tokenType != NULL) {
		compatible = compatible | Incompatible
		return
	}

	if a.arrayType != b.arrayType && b.arrayType != MIXED && (a.arrayType != NULL && b.arrayType != NULL) {
		compatible = compatible | Incompatible
		return
	}

	if a.tokenType == ARRAY {
		if a.arrayType == OBJECT {
			if len(a.array) > 0 && len(b.array) > 0 {
				compatible = compatible | Compare(a.array[0], b.array[0])
			}
		}
	}

	if a.tokenType == OBJECT {
		for k, va := range a.object {
			vb, exists := b.object[k]
			if !exists {
				compatible = compatible | Superset
			} else {
				compatible = compatible | Compare(va, vb)
			}
		}
		for k, _ := range b.object {
			_, exists := a.object[k]
			if !exists {
				compatible = compatible | Subset
			}
		}
	}
	return
}
