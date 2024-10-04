package phigros_test

import (
	"PhigrosApi/phigros"
	"encoding/json"
	"fmt"
	"testing"
)

var token = "zp4fenkahnb7639x9t6v38jmn"

func TestSave(t *testing.T) {
	//data, _ := phigros.GetDataFormTap(phigros.UserMeUrl, token) //获取id
	data, _ := phigros.GetDataFormTap(phigros.SaveUrl, token) //获取存档链接
	var us phigros.GameSave
	_ = json.Unmarshal(data, &us)
	phigros.SaveGameData(us.Results[0].GameFile.URL, "../data/gamesave/"+token+".zip")
	_ = phigros.LoadDifficult("../data/difficulty.tsv")
	jsons, _ := phigros.ParseSave("../data/gamesave/" + token + ".zip")
	fmt.Println(jsons)
}
