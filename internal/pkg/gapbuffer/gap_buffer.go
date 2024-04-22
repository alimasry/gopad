package gapbuffer

type GapBuffer struct {
	buffer []rune
	gap    int
	end    int
}

func NewGapBuffer(capacity int) *GapBuffer {
	return &GapBuffer{
		buffer: make([]rune, capacity),
		gap:    0,
		end:    capacity - 1,
	}
}

func NewGapBufferWithContent(content string) *GapBuffer {
	gapBuffer := NewGapBuffer(len(content) + 1)
	gapBuffer.Insert(content)
	return gapBuffer
}

func (gb *GapBuffer) MoveCursor(position int) {
	if position < gb.gap {
		gb.moveCursorLeft(position)
	} else {
		gb.moveCursorRight(position)
	}
}

func (gb *GapBuffer) Insert(text string) {
	for gb.end-gb.gap < len(text) {
		gb.grow()
	}
	copy(gb.buffer[gb.gap:], []rune(text))
	gb.gap += len(text)
}

func (gb *GapBuffer) InsertAt(position int, text string) {
	gb.MoveCursor(position)
	gb.Insert(text)
}

func (gb *GapBuffer) Delete(count int) {
	gb.gap = min(gb.gap, max(0, gb.gap-count))
}

func (gb *GapBuffer) DeleteAt(position int, count int) {
	gb.MoveCursor(position + 1)
	gb.Delete(count)
}

func (gb *GapBuffer) String() string {
	return string(gb.buffer[:gb.gap]) + string(gb.buffer[gb.end+1:])
}

func (gb *GapBuffer) grow() {
	newEnd := gb.end + len(gb.buffer)
	newBuffer := make([]rune, len(gb.buffer)*2)
	copy(newBuffer, gb.buffer[:gb.gap])
	copy(newBuffer[newEnd+1:], gb.buffer[gb.end+1:])
	gb.end = newEnd
	gb.buffer = newBuffer
}

func (gb *GapBuffer) moveCursorLeft(position int) {
	diff := min(gb.gap, max(0, gb.gap-position))
	copy(gb.buffer[gb.end-diff+1:gb.end+1], gb.buffer[position:gb.gap])
	gb.gap -= diff
	gb.end -= diff
}

func (gb *GapBuffer) moveCursorRight(position int) {
	diff := min(len(gb.buffer)-gb.end-1, max(0, position-gb.gap))
	copy(gb.buffer[gb.gap:gb.gap+diff], gb.buffer[gb.end+1:gb.end+diff+1])
	gb.gap += diff
	gb.end += diff
}
