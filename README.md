# hub-mirror
加速国外镜像，例如gcr.io、registry.k8s.io、k8s.gcr.io、quay.io、ghcr.io等

为避免重复工作，使用前建议在issue中搜索是否已转换

示例：[issues搜索gcc:latest](https://github.com/zz-open/hub-mirror/issues?q=gcc%3Alatest)

## 原理
- [engine sdk](https://docs.docker.com/engine/api/sdk/)
- [client](https://pkg.go.dev/github.com/docker/docker/client)
- [multi-platform](https://docs.docker.com/build/building/multi-platform/)

调用docker sdk 进行 pull,tag,push等操作

## 如何使用

### 方案一：白嫖我的，点个 Star ，[提交issue](https://github.com/zz-open/hub-mirror/issues/new/choose)即可
要求：严格按照模板规范提交，参考： [成功案例](https://github.com/zz-open/hub-mirror/issues/1)

> 当任务失败时，可以查看失败原因并直接修改 issues 的内容，即可重新触发任务执行

限制：每次提交最多 20 个镜像地址，避免超时

本人 使用阿里云容器镜像服务免费存储，请勿滥用

### 方案二：Fork 本项目，绑定你自己的 DockerHub 账号或其他镜像服务账号
#### 1. 绑定账号
在 `Settings`-`Secrets and variables`-`Actions` 选择 `New repository secret` 新建

- 如果要使用默认的 hub.docker.com 镜像服务
以下secret必需创建
```text
DOCKER_LOGIN_USERNAME: 用户名
DOCKER_LOGIN_PASSWORD: 密码
```

- 如果需要使用其它镜像服务，例如腾讯云、阿里云等
以下secret必需创建
```text
DOCKER_LOGIN_USERNAME: 用户名
DOCKER_LOGIN_PASSWORD: 密码
DOCKER_LOGIN_SERVER: 仓库地址
```
DOCKER_LOGIN_SERVER 示例：
```text
腾讯云: `ccr.ccs.tencentyun.com/xxx`

阿里云: `registry.cn-hangzhou.aliyuncs.com/xxx`
```

#### 2. 在 Fork 的项目中开启 `Settings`-`General`-`Features` 中的 `Issues` 功能

#### 3. 在 Fork 的项目中修改 `Settings`-`Actions`-`General` 中的 `Workflow permissions` 为 `Read and write permissions`

#### 4. 在 `Issues`-`Labels` 选择 `New label` 依次添加三个 label ：`hub-mirror`、`success`、`failure`

#### 5. 在 `Actions` 里选择 `hub-mirror` ，在右边 `···` 菜单里选择 `Enable Workflow`

#### 6. 在 Fork 的项目中提交 issues
```json
[
    {
        "mirrors": [
            "格式：原始镜像$自定义镜像名:自定义标签名",
            "$ 为分隔符，左边为原始镜像，右边为自定义镜像",
            "以下是三个正确示例",
            "gcc:latest",
            "gcc:latest$mygcc",
            "gcc:latest$mygcc:v1.0.0",
            "注意：mirrors 标签必选，标题随意，默认每次最多支持转换 20 个镜像",
            "无效的镜像会忽略, 请确保 json 格式是正确的"
        ],
        "platform": ""
    }
]
```
例如:
```json
[
    {
        "mirrors": [
             "gcc:latest$mygcc-arm64:v1.0.0"
        ],
        "platform": "linux/arm64"
    }
    {
        "mirrors": [
            "gcc:latest$mygcc-amd64:v1.0.0"
        ],
        "platform": "linux/amd64"
    }
]
```
一般没有特殊情况，不需要platform参数

### 方案三：克隆代码到本地,（适合本地有'魔法'的同学使用，可以灵活修改）
修改conf.yaml文件，填入要转换的镜像

执行
```shell
go run main.go \
--username='xxx' \
--password='xxx' \
--server='registry.cn-hangzhou.aliyuncs.com/zzimage' \
--rawcontent='[{"mirrors": ["gcc:latest$mygcc:v1.0.0"],"platform": "linux/amd64"}]'
```
如果嫌弃json字符串太长，也可将rawcontent配置成yaml，例如[conf.yaml](./conf.yaml),然后指定-f参数即可
```shell
go run main.go \
--username='xxx' \
--password='xxx' \
--server='registry.cn-hangzhou.aliyuncs.com/zzimage' \
-f=./conf.yaml
```

## workflow 参考
- [github-script](https://github.com/marketplace/actions/github-script)
- [setup-go-environment](https://github.com/marketplace/actions/setup-go-environment)
- [checkout](https://github.com/marketplace/actions/checkout)
- [stale](https://github.com/marketplace/actions/close-stale-issues)