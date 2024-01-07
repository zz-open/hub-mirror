---
name: hub-mirror issue template
about: 用于执行 hub-mirror workflow 的 issue 模板
title: "[hub-mirror] 请求执行任务"
labels: ["hub-mirror"]
---

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
