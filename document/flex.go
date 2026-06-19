package document

// JustifyContent controls how children are distributed along the main axis
// of a horizontal (flex) container, mirroring CSS justify-content.
type JustifyContent int

const (
	// JustifyStart packs children to the start of the main axis (default).
	JustifyStart JustifyContent = iota
	// JustifyCenter packs children to the center of the main axis.
	JustifyCenter
	// JustifyEnd packs children to the end of the main axis.
	JustifyEnd
	// JustifyBetween distributes children so the first is at the start,
	// the last at the end, and equal space between the rest.
	JustifyBetween
	// JustifyEvenly distributes children with equal space between and
	// around them.
	JustifyEvenly
)

// AlignItems controls how children are aligned along the cross axis
// (vertical) of a horizontal container, mirroring CSS align-items.
type AlignItems int

const (
	// AlignItemsStretch stretches children to fill the container height.
	// This is the default (zero value), preserving backward compatibility.
	AlignItemsStretch AlignItems = iota
	// AlignItemsStart aligns children to the top of the container.
	AlignItemsStart
	// AlignItemsCenter centers children vertically within the container.
	AlignItemsCenter
	// AlignItemsEnd aligns children to the bottom of the container.
	AlignItemsEnd
)
