---
name: hub-mirror issue template
about: 用于执行 hub-mirror workflow 的 issue 模板
title: "[hub-mirror] 请求执行任务"
labels: ["hub-mirror"]
---

[
    {
        "mirrors": [
            "格式：原始镜像 | 自定义镜像名:自定义标签名",
            "其中 |自定义镜像名:自定义标签名 是可选的",
            "以下是三个正确示例",
            "registry.k8s.io/kube-apiserver:v1.27.4",
            "registry.k8s.io/kube-apiserver:v1.27.4$demo",
            "registry.k8s.io/kube-apiserver:v1.27.4$demo:mytag",
            "要求：mirrors 标签是必选的，标题随意，内容严格按照该 json 格式，默认每次最多支持转换 20 个镜像",
            "错误的镜像都会被跳过, 请确保 json 格式是正确的"
        ],
        "platform": ""
    }
]
