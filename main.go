package main

import (
	"encoding/json"
	"strconv"

	"github.com/lianhong2758/PhigrosAPI/phigros"

	"github.com/gin-gonic/gin"
)

func init() {
	_ = phigros.LoadDifficult("difficulty.tsv")
}
func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default() //初始化
	r.GET("/", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"msg": "hello"}) })
	r.GET("/phigros/:session", phi)
	r.Run("0.0.0.0:8080")
}

func phi(ctx *gin.Context) {
	session := ctx.Param("session")
	if session[len(session)-1] == '\\' {
		session = session[:len(session)-1]
	}
	sn, _ := ctx.GetQuery("n")
	bn := 0
	if sn != "" {
		bn, _ = strconv.Atoi(sn)
	}
	j := phigros.UserRecord{Session: session}
	data, err := phigros.GetDataFormTap(phigros.UserMeUrl, session) //获取id
	if err != nil {
		ctx.JSON(200, phigros.RespCode{Code: 200,Message: err.Error(),Data: nil})
		return
	}
	var um phigros.UserMe
	_ = json.Unmarshal(data, &um)
	j.PlayerInfo = phigros.PlayerInfo{
		Name:      um.Nickname,
		CreatedAt: um.CreatedAt,
		UpdatedAt: um.UpdatedAt,
		Avatar:    um.Avatar,
	}
	data, err = phigros.GetDataFormTap(phigros.SaveUrl, session) //获取存档链接
	if err != nil {
		ctx.JSON(200, phigros.RespCode{Code: 200,Message: err.Error(),Data: nil})
		return
	}
	var gs phigros.GameSave
	_ = json.Unmarshal(data, &gs)
	ScoreAcc, _ := phigros.ParseStatsByUrl(gs.Results[0].GameFile.URL)
	j.ScoreAcc = phigros.BN(ScoreAcc, bn)
	ctx.JSON(200, phigros.RespCode{Code: 200,Message: "",Data: j})
}
