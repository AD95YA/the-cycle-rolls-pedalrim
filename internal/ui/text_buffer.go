package ui

import (
	"strconv"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// FastBuffer is a fixed-size byte buffer that formats integers, floats, and strings
// without causing heap allocations.
type FastBuffer struct {
	data   [1024]byte
	length int
}

// Reset clears the buffer.
func (b *FastBuffer) Reset() {
	b.length = 0
}

// AppendString appends raw ASCII text to the buffer.
func (b *FastBuffer) AppendString(s string) {
	n := copy(b.data[b.length:], s)
	b.length += n
}

// AppendInt appends an integer to the buffer.
func (b *FastBuffer) AppendInt(val int64) {
	start := b.length
	b.data[start] = '0'
	buf := strconv.AppendInt(b.data[:start], val, 10)
	b.length = len(buf)
}

// AppendFloat appends a float to the buffer with specified decimal precision.
func (b *FastBuffer) AppendFloat(val float64, precision int) {
	start := b.length
	buf := strconv.AppendFloat(b.data[:start], val, 'f', precision, 64)
	b.length = len(buf)
}

// String returns a zero-allocation string view of the buffer contents.
func (b *FastBuffer) String() string {
	if b.length == 0 {
		return ""
	}
	return unsafe.String(&b.data[0], b.length)
}

// DrawTextAt draws the buffer text at (x, y) using ebitenutil.DebugPrintAt.
func (b *FastBuffer) DrawTextAt(screen *ebiten.Image, x, y int) {
	if b.length > 0 {
		ebitenutil.DebugPrintAt(screen, b.String(), x, y)
	}
}

