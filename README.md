# 1. gost-subscribe

这是一个基于`ss`订阅生成`gost`配置文件的`go`程序，可一次性将`ss`链接的所有订阅转化为`gost`程序所需要的node节点

支持过滤指定关键字和包含指定关键字，方便对同一个地区进行代理，如生成所有香港节点的`config.yaml`代理文件

使用方法

```bash
gost-subscribe -u "https://example.com/ss" -o "config.yaml"
```

运行成功后，检查`config.yaml`文件，并用`gost`读取运行即可

可用的命令行选项：

- `-h`：显示帮助信息
- `-u`：订阅链接
- `-p`：TCP 代理端口号（默认为 11080）
- `-s`：代理策略（round|rand|fifo|hash，默认为 fifo）
- `-t`：代理失败超时时间（以秒为单位，默认为 600）
- `-m`：代理最大失败次数（默认为 1）
- `-V`：显示版本信息
- `-f`：过滤包含关键字的订阅（默认为 "套餐|重置|剩余|更新"）
- `-k`：仅获取包含关键字的订阅（默认为空，即包含所有）
- `-o`：输出文件路径（默认为 config.yml）


## 1.1. 开发环境
采用`vscode`进行开发
- Go Version: go version go1.21.2 linux/amd64

环境配置信息`.vscode/launch.json`参考

```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Go Subscribe",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/main.go",
            "console": "internalConsole",
            "args": [
                "-u=https://example.com/api/v1/client/subscribe?token=mytoken&flag=shadowsocks",
                "-k=韩国"
            ]
        }
    ]
}
```


## 1.2. 贡献

欢迎对该项目进行贡献！如果您发现问题或有改进建议，请提出新的 Issue 或提交 Pull Request