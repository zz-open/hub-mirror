# hub-mirror
加速国外镜像，例如gcr.io、registry.k8s.io、k8s.gcr.io、quay.io、ghcr.io等

为避免重复工作，使用前建议在issue中搜索是否已转换

示例：[issues搜索gcc:latest](https://github.com/zz-open/hub-mirror/issues?q=registry.k8s.io%2Fkube-apiserver%3Av1.28.2)

## 原理
- [engine sdk](https://docs.docker.com/engine/api/sdk/)
- [client](https://pkg.go.dev/github.com/docker/docker/client)
- [multi-platform](https://docs.docker.com/build/building/multi-platform/)

调用docker sdk 进行 pull,tag,push等操作

## 如何使用

### 场景一：github issue(适用于开放给其他同学使用，不太灵活)

#### 1. 绑定账号
在 `Settings`-`Secrets and variables`-`Actions` 选择 `New repository secret` 新建

#### 如果要使用默认的 hub.docker.com 镜像服务
以下参数必填
```text
DOCKER_LOGIN_USERNAME: 用户名
DOCKER_LOGIN_PASSWORD: 密码
```

#### 如果需要使用其它镜像服务，例如腾讯云、阿里云等
以下参数必填
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
            "格式：原始镜像[分隔符]自定义镜像全名",
            "其中 [分隔符]自定义镜像全名 是可选的",
            "以下是三个正确示例",
            "registry.k8s.io/kube-apiserver:v1.28.2",
            "registry.k8s.io/kube-apiserver:v1.28.2$demo",
            "registry.k8s.io/kube-apiserver:v1.28.2$demo:mytag",
            "要求：mirrors 标签是必选的，标题随意，内容严格按照该 json 格式，默认每次最多支持转换 20 个镜像",
            "错误的镜像都会被跳过, 请确保 json 格式是正确的",
            "注意最后一项没有逗号"
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
            "k8s.gcr.io/kube-proxy:v1.20.13"
        ],
        "platform": "linux/arm64"
    }
    {
        "mirrors": [
            "k8s.gcr.io/kube-proxy:v1.20.13"
        ],
        "platform": "linux/amd64"
    }
]
```

### 场景二：克隆代码到本地（适合本地有魔法的同学使用，可以灵活修改）
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