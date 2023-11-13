# my-cubing-robot

魔方QQ-微信机器人交流

安装时有几个依赖：

```
qsign签名服务器:
docker run -d --restart=always --name qsign -p 5709:8080 -e ANDROID_ID="3ea26c64f3914bc3"  xzhouqd/qsign:8.9.63


html转image工具
docker run -d --restart=always --name doctron -p 8080:8080 lampnick/doctron 
```
