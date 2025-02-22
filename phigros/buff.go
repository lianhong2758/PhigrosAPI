package phigros

import "bytes"

type Buff struct {
	Bytes   bytes.Buffer
	bit     int
	tempbit byte
}

func (b *Buff) Alignment() {
	if b.bit > 0 {
		b.Bytes.WriteByte(b.tempbit)
		b.bit = 0
		b.tempbit = 0
	}
}

func (b *Buff) SaveBool(bt bool) {
	b.tempbit += byte(1 << b.bit)
	b.bit++
}

func (b *Buff) SaveString(s string) {
 
}