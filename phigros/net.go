package phigros

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	UserMeUrl = "https://rak3ffdi.cloud.tds1.tapapis.cn/1.1/users/me"
	SaveUrl   = "https://rak3ffdi.cloud.tds1.tapapis.cn/1.1/classes/_GameSave"
)

func GetDataFormTap(url, token string) (data []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-LC-Id", "rAK3FfdieFob2Nn8Am")
	req.Header.Add("X-LC-Key", "Qr9AEqtuoSVS3zeD6iVbM4ZC0AtkJcQ89tywVyi0")
	req.Header.Add("X-LC-Session", token)
	req.Header.Add("User-Agent", "LeanCloud-CSharp-SDK/1.0.3")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "rak3ffdi.cloud.tds1.tapapis.cn")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}

func SaveGameData(url, path string) error {
	rsp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error while downloading: %s", err.Error())

	}
	defer rsp.Body.Close()
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %s", err.Error())

	}
	defer file.Close()

	// 将下载的内容保存到文件
	_, err = io.Copy(file, rsp.Body)
	if err != nil {
		return fmt.Errorf("error saving file: %s", err.Error())
	}
	return nil
}
