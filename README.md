# tiktok-demo
## 0.青训营相关
![img.png](_images/青训营奖项.jpg)
## 1.安装部署
1. 安装mysql，执行sql文件
2. 安装nginx，根据项目中conf/nginx.conf修改nginx.conf配置文件，重新启动
3. 安装minio，建议默认9000端口，创建两个桶image和video，公开访问权限
![img.png](_images/minio.png)
4. 安装ffmpeg，linux上将可执行文件配置到/usr/bin中，windows配置环境变量
5. 安装redis
6. 安装kafka，新建三个topic：video、favorite、chat
7. 下载代码，go mod tidy拉取第三方库
8. 修改conf/config.yaml参数适配

## 2.数据库设计
![img.png](_images/数据库表设计.png)

## 3.详细设计
[https://zvrkpwe3jg.feishu.cn/docx/T2g7de8ONoGLZax9IxBc2PBRnUf](https://zvrkpwe3jg.feishu.cn/docx/T2g7de8ONoGLZax9IxBc2PBRnUf)
