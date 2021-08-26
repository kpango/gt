package gt

const (
	// H represents horizontal separator
	H int = iota
	// V represents vertical separator
	V
	// VH represents vertical&horizontal separator
	VH
	// HU represents horizontal and up separator
	HU
	// HD represents horizontal and down separator
	HD
	// VL represents vertical and left separator
	VL
	// VR represents vertical and right separator
	VR
	// DL represents down and left separator
	DL
	// DR represents down and right separator
	DR
	// UL represents up and left separator
	UL
	// UR represents up and right separator
	UR
	length
)

var boxLine = [length]rune{'─', '│', '┼', '┴', '┬', '┤', '├', '┐', '┌', '┘', '└'}
