## hub-mirror
使用国内云镜像服务加速gcr.io、registry.k8s.io、k8s.gcr.io、quay.io、ghcr.io等国外镜像

为减少重复请求，合理利用资源，建议提前在 issues 中搜索镜像是否已转换过，可以直接复用

示例：[issues搜索registry.k8s.io/kube-apiserver:v1.27.4](https://github.com/zz-open/hub-mirror/issues?q=registry.k8s.io%2Fkube-apiserver%3Av1.27.4)

## 原理
利用github actions 拉取镜像上传到自己的云镜像服务

## 开始使用

### 绑定自己的 DockerHub 账号或其他镜像服务账号

1. 绑定账号

    - 如果要使用默认的 hub.docker.com 镜像服务

      在 `Settings`-`Secrets and variables`-`Actions` 选择 `New repository secret` 新建 `DOCKER_USERNAME`（你的 Docker
      用户名） 和 `DOCKER_TOKEN`（你的 Docker 密码） 两个 Secrets

    - 如果需要使用其它镜像服务，例如腾讯云、阿里云等

      在 `Settings`-`Secrets and variables`-`Actions` 选择 `New repository secret` 新建 `DOCKER_USERNAME`（你的其它镜像服务用户名）
      和 `DOCKER_TOKEN`（你的其它镜像服务密码）以及 `DOCKER_REPOSITORY` 三个 Secrets

      其中 `DOCKER_REPOSITORY` 配置例子：

        - 腾讯云: `ccr.ccs.tencentyun.com/xxxxxx`
        - 阿里云: `registry.cn-hangzhou.aliyuncs.com/xxxxxx`

2. 在 Fork 的项目中开启 `Settings`-`General`-`Features` 中的 `Issues` 功能

3. 在 Fork 的项目中修改 `Settings`-`Actions`-`General` 中的 `Workflow permissions` 为 `Read and write permissions`

4. 在 `Issues`-`Labels` 选择 `New label` 依次添加三个 label ：`hub-mirror`、`success`、`failure`

5. 在 `Actions` 里选择 `hub-mirror` ，在右边 `···` 菜单里选择 `Enable Workflow`

6. 在 Fork 的项目中提交 issues


## docker sdk
- [sdk](https://docs.docker.com/engine/api/sdk/)
- [client](https://pkg.go.dev/github.com/docker/docker/client)