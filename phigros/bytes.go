package phigros

type Bytes struct {
	Data []byte
	ptr  int
}

func NewBytesReader(b []byte) *Bytes {
	return &Bytes{Data: b, ptr: 0}
}

func (b *Bytes) ReadShort() byte {
	num := b.Data[b.ptr]
	if num < 128 {
		b.ptr++
	} else {
		b.ptr += 2
	}
	return num
}

func (b *Bytes) ReadBool() byte {
	return b.Data[b.ptr]
}

func (b *Bytes) ReadNext() {
	b.ptr++
}

func (b *Bytes) ReadString() string {
	length := b.Data[b.ptr]
	b.ptr += int(length + 1)
	return BytesToString(b.Data[b.ptr-int(length) : b.ptr])
}

func (b *Bytes) ReadScoreAcc() ScoreAcc {

	return ScoreAcc{Score: int(b.ReadInt32()), Acc: b.ReadFloat32()}
}

func (b *Bytes) ReadInt32() int32 {
	b.ptr += 4
	return BytesToInt(b.Data[b.ptr-4 : b.ptr])
}

func (b *Bytes) ReadFloat32() float32 {
	b.ptr += 4
	return ByteToFloat32(b.Data[b.ptr-4 : b.ptr])
}

func GetBool(num byte, index int) bool {
	return num&(1<<index) == 1
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
			scoreAcc := b.ReadScoreAcc()
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