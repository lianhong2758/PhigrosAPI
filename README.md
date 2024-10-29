## PhigrosApi
Phigros的战绩查询SDK


## 接口如下
```
1.获取数据
GET localhost:8080/phigros/:session?n=num
示例 n为返回数据的数组长度,不写则全返回,session为key
GET localhost:8080/phigros/dwadsfawdsad?n=5

result
{
  "code": 200,
  "message": "", //err
  "data": {
     ......太长不写了
  }
}

```
2.绘图,获取数据同时进行
GET localhost:8080/phigros/:session?n=num&pic=bool&type=type
示例 n为返回数据的图片bn,pic指定输出为图片,true/false,type指定返回图片或者路径,pic/json
GET localhost:8080/phigros/dwadsfawdsad?n=21&pic=true&type=pic

result
图片数据

## 构建流程
- 1. 编译程序
  - 本地搭建Go环境,在项目目录下执行`go build`即可

## 在线体验
-  `http://106.54.63.95:8080`

## 鸣谢
- [PhigrosLibrary](https://github.com/7aGiven/PhigrosLibrary?tab=readme-ov-file)