##浏览器

###使用前准备


1. 数据库备份位置 db/backup/

2. 还原数据库publicChain

     在此之前保证数据库,能够正常连通,mongodb 版本version v4.0.9
     
     首先进入db/backck/文件夹（）
     使用命令：
     
     mongorestore -h 127.0.0.1:27017 -d publicChain --dir ./publicChain

3. 打开项目工程运行 main.go即可 



注：如有包失效请下载 

如：web3.go/common/hexutil

请下载web3go 的引用包


