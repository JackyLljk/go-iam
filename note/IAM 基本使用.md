## apiserver 基本使用

### 使用 curl 进行操作 CURD 接口

#### 1. 使用管理员账号登录，获取 token

```bash
curl -s -XPOST -H'Content-Type: application/json' -d'{"username":"admin","password":"Admin@2021"}' http://127.0.0.1:8080/login | jq -r .token
	# 这里会输出 token xxx.xxx.xxx
```

#### 2. 使用 token 进行操作

**用户操作**

```bash
# 创建用户
curl -s -XPOST -H'Content-Type: application/json' -H'Authorization: Bearer {token}' -d'{"password":"User@2021","metadata":{"name":"colin"},"nickname":"colin","email":"colin@foxmail.com","phone":"1812884xxxx"}' http://127.0.0.1:8080/v1/users

# 列出用户
curl -s -XGET -H'Authorization: Bearer {token}' 'http://127.0.0.1:8080/v1/users?offset=0&limit=10'

# 获取 user_name 用户的详细信息
curl -s -XGET -H'Authorization: Bearer {token}' http://127.0.0.1:8080/v1/users/colin

# 修改用户
curl -s -XPUT -H'Content-Type: application/json' -H'Authorization: Bearer {token}' -d'{"nickname":"colin","email":"colin_modified@foxmail.com","phone":"1812884xxxx"}' http://127.0.0.1:8080/v1/users/colin

# 删除用户
curl -s -XDELETE -H'Authorization: Bearer {token}' http://127.0.0.1:8080/v1/users/colin

# 批量删除用户
curl -s -XDELETE -H'Authorization: Bearer {token}' 'http://127.0.0.1:8080/v1/users?name=colin&name=mark&name=john'
```







































