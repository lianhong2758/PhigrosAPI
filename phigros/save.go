package phigros

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"github.com/tidwall/gjson"
)

func ReadZip(path string) (m map[string][]byte, err error) {
	m = map[string][]byte{}
	// 打开 zip 文件
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// 遍历 zip 文件中的文件
	for _, file := range reader.File {
		// 打开文件
		f, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()
		// 读取文件内容
		buf := make([]byte, file.FileInfo().Size())
		_, _ = f.Read(buf)
		m[file.Name] = buf
	}
	return m, nil
}

func Decrypt(in []byte) (out []byte, err error) {
	// CBCDecrypt AES-CBC 解密
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(in) < aes.BlockSize {
		return nil, fmt.Errorf("cipherText too short")
	}
	out = make([]byte, len(in))

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(out, in[1:])

	return append(in[0:1], unpad(out)...), nil
	//return out, nil
}

// 去填充函数的示例实现
func unpad(data []byte) []byte {
	padding := data[len(data)-1]
	return data[:len(data)-int(padding)]
}

func Unmarshal[T PhigrosStruct](in []byte) *T {
	var ps T
	v := reflect.ValueOf(&ps).Elem()
	t := reflect.TypeOf(&ps).Elem()
	reader := NewBytesReader(in)
	for i := range v.NumField() {
		if t.Field(i).Tag.Get("phi") != "-" {
			set(v.Field(i), reader)
		}
	}
	return &ps
}

// 嵌套结构不检查tag
func set(rv reflect.Value, reader *Bytes) {
	switch rv.Kind() {
	case reflect.Bool:
		rv.SetBool(reader.ReadBool())
	case reflect.String:
		rv.SetString(reader.ReadString())
	case reflect.Float32:
		rv.SetFloat(float64(reader.ReadFloat32()))
	case reflect.Int16:
		rv.SetInt(int64(reader.ReadShort()))
	case reflect.Uint8:
		rv.SetUint(uint64(reader.ReadByte1()))
	case reflect.Uint16:
		rv.SetUint(uint64(reader.ReadVarShort()))
	case reflect.Array:
		for i := range rv.Len() {
			if rv.Index(i).Kind() == reflect.Struct {
				for ii := range rv.Index(i).NumField() {
					set(rv.Index(i).Field(ii), reader)
				}
			} else {
				//非结构体数组
				set(rv.Index(i), reader)
			}
		}
	default:
	}
}

// 未来实现
func Marshal[T PhigrosStruct](v *T) []byte {
	var buff bytes.Buffer
	buff.WriteByte('1')
	return buff.Bytes()
}

func UnmarshalGameRecord(in []byte) []ScoreAcc {
	records := []ScoreAcc{}
	reader := NewBytesReader(in)
	for range reader.ReadVarShort() {
		t := reader.ReadString()
		songId := t[:len(t)-2]
		record := reader.ReadRecord(songId)
		records = append(records, record...)
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Rks > records[j].Rks
	})
	return records

}

// 前19成绩,取最高成绩放第一位
func B19(records []ScoreAcc) []ScoreAcc {
	return BN(records, 19)
}

// 取前n成绩,取最高成绩放第一位
func BN(records []ScoreAcc, n int) []ScoreAcc {
	var maxRecord ScoreAcc
	for _, r := range records {
		if r.Score == 1000000 {
			if r.Difficulty > maxRecord.Difficulty {
				maxRecord = r
			}
		}
	}
	bn := []ScoreAcc{maxRecord}
	if n <= 0 {
		return append(bn, records...)
	}
	// 将records中的前19个记录加入b19
	if len(records) >= n {
		bn = append(bn, records[:n]...)
	} else {
		bn = append(bn, records...)
	}
	return bn
}

// 通过zip文件读取所有云端内容
func ParseSave(path string) (map[string]any, error) {
	m, err := ReadZip(path)
	if err != nil {
		return nil, err
	}
	for k, v := range m {
		out, err := Decrypt(v)
		if err != nil {
			return nil, fmt.Errorf("Decrypt file %s Error %s", k, err.Error())
		}
		m[k] = out
	}
	if m["gameRecord"][0] != byte(0x01) {
		return nil, errors.New("版本号不正确，可能协议已更新。")
	}
	//json
	jsons := make(map[string]any)
	jsons["gameRecord"] = B19(UnmarshalGameRecord(m["gameRecord"][1:]))
	jsons["settings"] = *Unmarshal[Settings](m["settings"][1:])
	jsons["user"] = *Unmarshal[User](m["user"][1:])
	jsons["gameProgress"] =*Unmarshal[GameProgress](m["gameProgress"][1:])
	return jsons, nil
}

// 通过url获取战绩,其余内容丢弃
func ParseStatsByUrl(url string) ([]ScoreAcc, error) {
	d, err := GetGameRecordData(url)
	if err != nil {
		return nil, err
	}
	d, err = Decrypt(d)
	if err != nil {
		return nil, fmt.Errorf("Decrypt file gameRecord Error %s", err.Error())
	}
	if d[0] != byte(0x01) {
		return nil, errors.New("版本号不正确，可能协议已更新。")
	}
	return UnmarshalGameRecord(d[1:]), nil
}

func ProcessSummary(sum string) (s *Summary) {
	if sum == "" {
		return nil
	}
	b, err := base64.StdEncoding.DecodeString(sum)
	if err != nil {
		return nil
	}
	s = Unmarshal[Summary](b)
	s.ChalID = (s.ChallengeModeRank - (s.ChallengeModeRank % 100)) / 100
	s.Chalnum = strconv.Itoa(int(s.ChallengeModeRank % 100))
	return
}
func GetUserRecordQuickly(session string) (*UserRecord, error) {
	j := UserRecord{Session: session}
	data, err := GetDataFormTap(UserMeUrl, session) //获取id
	if err != nil {
		return nil, err
	}
	um := gjson.Parse(BytesToString(data))
	j.PlayerInfo = &PlayerInfo{
		Name:      um.Get("nickname").String(),
		CreatedAt: um.Get("createdAt").Time(),
		UpdatedAt: um.Get("updatedAt").Time(),
		Avatar:    um.Get("avatar").String(),
	}
	data, err = GetDataFormTap(SaveUrl, session) //获取存档链接
	if err != nil {
		return nil, err
	}
	gs := gjson.Parse(BytesToString(data))
	j.ScoreAcc, _ = ParseStatsByUrl(gs.Get("results.0.gameFile.url").String()) //gs.Results[0].GameFile.URL
	j.Summary = ProcessSummary(gs.Get("results.0.summary").String())           //gs.Results[0].Summary
	return &j, nil
}
