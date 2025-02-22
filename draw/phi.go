package draw

import (
	"os"

	"github.com/lianhong2758/PhigrosAPI/phigros"
)

var Challengemoderank = []string{"white", "green", "blue", "red", "gold", "rainbow"}
var fontsd []byte

func init() {
	Challengemode = DataPath + "challengemode/"
	// 默认字体,如需修改直接调用LoadFont
	Font = DataPath + "MaokenZhuyuanTi.ttf"
	// 评级
	Rank = DataPath + "rank/"
	// 曲绘
	Illustration = DataPath + "illustration/"
	//战绩图
	Output = DataPath + "output/"
	//头图
	Avatar = DataPath + "avatar/"
	//字体
	LoadFont(Font)
	if IsNotExist(DataPath + "output") {
		_ = os.MkdirAll(DataPath+"output", 0755)
	}
	if IsNotExist(DataPath + "avatar") {
		_ = os.MkdirAll(DataPath+"avatar", 0755)
	}
}

// 判断是否为有效session
func SessionIsEfficient(session string) error {
	_, err := phigros.GetDataFormTap(phigros.UserMeUrl, session) //获取id
	return err
}
func DownloadAvatar(url, session string) error {
	data, err := phigros.GetData(url)
	if err != nil {
		return err
	}
	f, err := os.Create(Avatar + session + ".png")
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	f.Close()
	return err
}

func LoadFont(path string) {
	fontsd, _ = os.ReadFile(Font)
}
