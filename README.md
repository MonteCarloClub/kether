# kether
为区块链服务调度和管理 Docker 容器

## 1. 如何使用？
1.1. 拉取源代码，构建 `kether`。
```bash
git clone https://github.com/MonteCarloClub/kether.git
cd kether
make all # 或 make kether
```

1.2. 在主机 6379 端口部署 redis 并测试，期望输出 `ok`。
```bash
docker pull redis:6.2.6
docker run -d -p 6379:6379 redis:6.2.6
go test -run TestInitRedisClient github.com/MonteCarloClub/kether/registry
```

1.3. 运行和部署测试用例，对内发布 HTTP 服务，对外发布 HTTPS 服务。

1.3.1. 创建 `kether-net` 网络，查询网关 IP，填充 `test/http_https_echo_server.yml` 的 `network_list` 字段的 `*` 处，部署 `http-https-echo-server`。
```bash
docker network create --driver bridge kether-net
docker network inspect kether-net
./bin/kether deploy -f test/http_https_echo_server.yml
```
1.3.2. 在主机 8443 端口访问 HTTPS 服务。
```bash
curl -k -X PUT -H "Arbitrary:Header" -d aaa=bbb https://localhost:8443/hello-world
```
1.3.3. 构建 `http-echo-client` 镜像，填充 `test/http_echo_client.yml` 的 `network_list` 和 `volume_list` 字段的 `*` 处，部署 `http-echo-client`。
```bash
cd test/http_echo_client
docker build -t kofclubs/http-echo-client:testing .
cd ../..
./bin/kether deploy -f test/http_echo_client.yml
```
1.3.4. 打开 `test/http_echo_client.yml` 的 `volume_list` 字段指定的主机文件，验证文件 I/O。
```bash
cat /*/response.txt
```

1.4. 清理产物。
```bash
make clean
```

## 2. 如何开发？
- 请在自己的分支上开发，每个开发分支应该仅领先主分支1个提交。
- 请勿向主分支推送提交。
- 请新建 Pull Request，@KofClubs 将在 Review 通过后把你的分支合并到主分支。