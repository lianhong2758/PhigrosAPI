package phigros

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"reflect"
	"sort"
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

func DecoderWithStruct[T PhigrosStruct](in []byte) *T {
	var ps T
	v := reflect.ValueOf(&ps).Elem()
	bit := 0
	reader := NewBytesReader(in)

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Bool {
			v.Field(i).SetBool(GetBool(reader.ReadBool(), bit))
			continue
		}
		if bit > 0 {
			bit = 0
			reader.ReadNext()
		}
		switch v.Field(i).Kind() {
		case reflect.String:
			v.Field(i).SetString(reader.ReadString())
		case reflect.Float32:
			v.Field(i).SetFloat(float64(reader.ReadFloat32()))
		}
		if bit > 0 {
			reader.ReadNext()
		}

	}
	return &ps
}

func DecoderGameRecord(in []byte) []ScoreAcc {
	records := []ScoreAcc{}
	reader := NewBytesReader(in)
	for range reader.ReadShort() {
		t := reader.ReadString()
		songId := t[:len(t)-2]
		record := reader.ReadRecord(songId)
		records = append(records, record...)
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Rks < records[j].Rks
	})
	var maxRecord ScoreAcc
	for _, r := range records {
		if r.Score == 1000000 {
			if r.Difficulty > maxRecord.Difficulty {
				maxRecord = r
			}
		}
	}
	b19 := []ScoreAcc{maxRecord}
	// 将records中的前19个记录加入b19
	if len(records) >= 19 {
		b19 = append(b19, records[:19]...)
	} else {
		b19 = append(b19, records...)
	}
	return b19

}

func ParseSave(path string) (map[string]any, error) {
	m, err := ReadZip(path)
	if err != nil {
		return nil,err
	}
	for k, v := range m {
		out, err := Decrypt(v)
		if err != nil {
			return nil,fmt.Errorf("Decrypt file %s Error %s", k, err.Error())
		}
		m[k] = out
	}
	if m["gameRecord"][:1][0] != byte(0x01) {
		return nil,errors.New("版本号不正确，可能协议已更新。")
	}
	//json
	jsons := make(map[string]any)
	jsons["gameRecord"] = DecoderGameRecord(m["gameRecord"][1:])
	// jsons["settings"] = *DecoderWithStruct[Settings](m["settings"][1:])
	// jsons["user"] = *DecoderWithStruct[User](m["user"][1:])
	// fmt.Println(jsons)
	return jsons,nil
}
