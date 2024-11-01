package main

import (
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
	bn := 21
	if sn != "" {
		bn, _ = strconv.Atoi(sn)
	}
	if bn > 200 {
		bn = 200
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
	j, err := phigros.GetUserRecordQuickly(session)
	if err != nil {
		ctx.JSON(200, phigros.RespCode{Code: 200, Message: err.Error(), Data: nil})
		return
	}
	j.ScoreAcc = phigros.BN(j.ScoreAcc, bn)
	if !havepic {
		ctx.JSON(200, phigros.RespCode{Code: 200, Message: "", Data: j})
		return
	}
	draw.DownloadAvatar(j.PlayerInfo.Avatar, session)
	err = draw.DrawPic(0.5, j, strconv.FormatFloat(float64(j.Summary.Rks), 'f', 6, 64), draw.Challengemoderank[j.Summary.ChalID], j.Summary.Chalnum, session)
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
