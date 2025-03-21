## PhigrosApi
Phigros的战绩查询SDK


## 接口如下
```
1.获取数据
GET localhost:8080/phigros/:session?n=num
示例 n为返回数据的数组长度,不写则全返回,session为key
GET localhost:8080/phigros/nkyjch88ydrg4js83bea9jyiw?n=5

result
{
  "code": 200,
  "message": "", //err
  "data": {
     ......太长不写了
  }
}

```
```
2.绘图,获取数据同时进行
GET localhost:8080/phigros/:session?n=num&pic=bool&type=type
示例 n为返回数据的图片bn,pic指定输出为图片,true/false,type指定返回图片或者路径,pic/json
GET localhost:8080/phigros/nkyjch88ydrg4js83bea9jyiw?n=21&pic=true&type=pic

result
图片数据
```
## 构建流程
- 1. 编译程序
  - 本地搭建Go环境,在项目目录下执行`go build`即可
- 2. 图片服务需要额外下载内容
  - 下载这个文件夹的内容下载到data文件夹[RosmBot-Data-phi](https://github.com/lianhong2758/RosmBot-Data/tree/main/phi)
    - 若此文件版本更新不及时,可以前往这里[Phigros_Resource-illustration](https://github.com/7aGiven/Phigros_Resource/tree/illustration)下载最新曲绘
    - 最新定数表[difficulty.tsv](https://github.com/7aGiven/Phigros_Resource/blob/info/difficulty.tsv)
  - 下载这个文件夹里面的MaokenZhuyuanTi.ttf,放到data文件夹[RosmBot-Data-font](https://github.com/lianhong2758/RosmBot-Data/tree/main/font)
- 3. qr扫码登录部分可自由设计, `r.Data.QrcodeURL`也可直接浏览器打开登录

## 在线体验
-  `http://117.72.123.235:8080/phigros/nkyjch88ydrg4js83bea9jyiw?n=5`

## 鸣谢
- [PhigrosLibrary](https://github.com/7aGiven/PhigrosLibrary?tab=readme-ov-file)