package phigros

type Bytes struct {
	Data []byte
	ptr  int
	bit  int
}

func NewBytesReader(b []byte) *Bytes {
	return &Bytes{Data: b, ptr: 0, bit: 0}
}

func (b *Bytes) Alignment() {
	if b.bit > 0 {
		b.bit = 0
		b.ptr++
	}
}

// 也叫VarInt,这里用来对应uint16,因为内置的uint8==byte无法共存
func (b *Bytes) ReadVarShort() byte {
	b.Alignment()
	num := b.Data[b.ptr]
	if num < 128 {
		b.ptr++
	} else {
		num = num&0b01111111 ^ b.Data[b.ptr+1]<<7
		b.ptr += 2
	}
	return num
}

func (b *Bytes) ReadShort() int16 {
	b.ptr += 2
	return int16(b.Data[b.ptr-2]) + int16(b.Data[b.ptr-1])<<8
}

func (b *Bytes) ReadByte1() byte {
	b.Alignment()
	b.ptr++
	return b.Data[b.ptr-1]
}

func (b *Bytes) ReadBool() (tb bool) {
	if b.bit >= 4 {
		b.bit = 0
		b.ptr++
	}
	t := b.Data[b.ptr]
	tb = GetBool(t, b.bit)
	b.bit++
	return
}

func (b *Bytes) ReadNext() {
	b.ptr++
}

func (b *Bytes) ReadString() string {
	b.Alignment()
	length := b.ReadVarShort()
	b.ptr += int(length)
	return BytesToString(b.Data[b.ptr-int(length) : b.ptr])
}

func (b *Bytes) ReadInt32() int32 {
	b.Alignment()
	b.ptr += 4
	return BytesToInt(b.Data[b.ptr-4 : b.ptr])
}

func (b *Bytes) ReadFloat32() float32 {
	b.Alignment()
	b.ptr += 4
	return ByteToFloat32(b.Data[b.ptr-4 : b.ptr])
}

func GetBool(num byte, index int) bool {
	return (num>>index)&1 == 1
}

func (b *Bytes) ReadRecord(songId string) []ScoreAcc {
	endPosition := b.ptr + int(b.Data[b.ptr]) + 1
	b.ptr += 1
	exists := b.Data[b.ptr]
	b.ptr += 1
	fc := b.Data[b.ptr]
	b.ptr += 1
	diff := difficulty[songId]
	records := []ScoreAcc{}
	for level := range len(diff) {
		if GetBool(exists, level) {
			scoreAcc := ScoreAcc{}
			scoreAcc.Score = int(b.ReadInt32())
			scoreAcc.Acc = b.ReadFloat32()
			scoreAcc.Level = levels[level]
			scoreAcc.Fc = GetBool(fc, level)
			scoreAcc.SongId = songId
			scoreAcc.Difficulty = diff[level]
			scoreAcc.Rks = (scoreAcc.Acc - 55) / 45
			scoreAcc.Rks = scoreAcc.Rks * scoreAcc.Rks * scoreAcc.Difficulty
			records = append(records, scoreAcc)
		}
	}
	b.ptr = endPosition
	return records
}
