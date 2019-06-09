# 测试 #

## http ##
1. 获取状态
    `curl -XGET "localhost:8888/status/" -v`
2. 添加
    `curl -XPOST "localhost:8888/cache/name/" -d "kk" -v`
3. 获取
    `curl -XGET "localhost:8888/cache/name/" -v`
4. 删除
    `curl -XDELETE "localhost:8888/cache/name/" -v`


## tcp ##

1. 添加
    ```
    import socket
    c = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    c.connect(('localhost', 8889))
    c.send(b"S4 2 namekk")
    c.recv(256)
    ```
2. 获取
    ```
    import socket
    c = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    c.connect(('localhost', 8889))
    c.send(b"G4 name")
    c.recv(256)
    ```
3. 删除
    ```
    import socket
    c = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    c.connect(('localhost', 8889))
    c.send(b"D4 name")
    c.recv(256)
    ```