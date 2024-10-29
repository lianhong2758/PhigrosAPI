package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/lianhong2758/PhigrosAPI/draw"
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
	pic, havepic := ctx.GetQuery("pic")
	//不需要pic
	if !havepic || pic == "false" {
		havepic = false
	}
	pic, havePath := ctx.GetQuery("type")
	if !havePath || pic != "json" {
		havePath = false
	}
	j := phigros.UserRecord{Session: session}
	data, err := phigros.GetDataFormTap(phigros.UserMeUrl, session) //获取id
	if err != nil {
		ctx.JSON(200, phigros.RespCode{Code: 200, Message: err.Error(), Data: nil})
		return
	}
	var um phigros.UserMe
	_ = json.Unmarshal(data, &um)
	j.PlayerInfo = &phigros.PlayerInfo{
		Name:      um.Nickname,
		CreatedAt: um.CreatedAt,
		UpdatedAt: um.UpdatedAt,
		Avatar:    um.Avatar,
	}
	data, err = phigros.GetDataFormTap(phigros.SaveUrl, session) //获取存档链接
	if err != nil {
		ctx.JSON(200, phigros.RespCode{Code: 200, Message: err.Error(), Data: nil})
		return
	}
	var gs phigros.GameSave
	_ = json.Unmarshal(data, &gs)
	ScoreAcc, _ := phigros.ParseStatsByUrl(gs.Results[0].GameFile.URL)
	j.ScoreAcc = phigros.BN(ScoreAcc, bn)
	j.Summary = phigros.ProcessSummary(gs.Results[0].Summary)
	if !havepic {
		ctx.JSON(200, phigros.RespCode{Code: 200, Message: "", Data: j})
		return
	}
	draw.DownloadAvatar(j.PlayerInfo.Avatar, session)
	err = draw.DrawPic(0.5, j, strconv.FormatFloat(float64(j.Summary.Rks), 'f', 6, 64), draw.Challengemoderank[(j.Summary.ChallengeModeRank-(j.Summary.ChallengeModeRank%100))/100], strconv.Itoa(int(j.Summary.ChallengeModeRank%100)), session)
	if err != nil {
		ctx.JSON(200, phigros.RespCode{Code: 200, Message: err.Error(), Data: nil})
		return
	}
	if havePath {
		ctx.JSON(200, phigros.RespCode{Code: 200, Message: "", Data: map[string]string{"file": Pwd() + "/" + draw.Output + session + ".png"}})
		return
	}
	ctx.File(Pwd() + "/" + draw.Output + session + ".png")
}
func Pwd() string {
	path, _ := os.Getwd()
	return strings.ReplaceAll(path, "\\", "/")
}
