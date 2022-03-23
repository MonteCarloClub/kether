# kether
为区块链服务调度和管理 Docker 容器

## 1. 如何使用？
1.1. 拉取源代码，构建 `kether`。
```bash
$ git clone https://github.com/MonteCarloClub/kether.git
$ cd kether
$ make all # 或 make kether
```

1.2. 在主机 6379 端口部署 redis 并测试，期望输出 `ok`。
```bash
$ docker pull redis:6.2.6
$ docker run -d -p 6379:6379 redis:6.2.6
$ go test -run TestInitRedisClient github.com/MonteCarloClub/kether/registry
```

1.3. 运行和部署测试用例。
```bash
$ ./bin/kether deploy -f test/dao_2048.yml
```

1.4. 清理产物。
```bash
$ make clean
```

## 2. 如何开发？
- 请在自己的分支上开发，每个开发分支应该仅领先主分支1个提交。
- 请勿向主分支推送提交。
- 请新建 Pull Request，@KofClubs 将在 Review 通过后把你的分支合并到主分支。