package draw

import "os"

var (
	DataPath = "data/"
	// 课题模式图标
	Challengemode string
	// 字体
	Font string
	// 评级
	Rank string
	// 曲绘
	Illustration string
	// 图标
	Icon string
	//战绩图
	Output string
	//头图
	Avatar string
)

// IsExist 文件/路径存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IsNotExist 文件/路径不存在
func IsNotExist(path string) bool {
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
}
