# 测试 #

1. 获取状态
    `curl -XGET "localhost:8888/status/" -v`
2. 添加
    `curl -XPOST "localhost:8888/cache/name/" -d "kk" -v`
3. 获取
    `curl -XGET "localhost:8888/cache/name/" -v`
4. 删除
    `curl -XDELETE "localhost:8888/cache/name/" -v`