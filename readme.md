好的！下面是基于你当前项目结构和多平台打包要求，完善后的 README.md 模板，包含安装指南、使用示例、及各平台下载链接示范。

⸻


# sm3sum

`sm3sum` 是一款兼容 `md5sum`、`sha1sum` 等工具的命令行程序，使用国密 SM3 哈希算法进行文件摘要计算和校验。

---

## 🚀 功能特性

- 支持 SM3 哈希计算和校验
- 完全兼容 md5sum 的常用命令行参数（`-c`, `--tag`, `-z` 等）
- 支持多平台交叉编译（Linux/macOS/Windows，amd64/arm64）
- 支持 Bash 和 Zsh 命令补全
- 提供 man 手册页

---

## 📦 安装指南

### 1. 从源码编译（需安装 Go 环境）

```bash
git clone https://github.com/yourname/sm3sum.git
cd sm3sum
go mod tidy
make build
sudo make install
```

### 2. 使用预编译二进制包（Release 下载）

请根据你的操作系统和 CPU 架构选择对应版本：

|平台|版本|示例|下载链接|
|---|---|---|---|
|Linux|amd64|v1.0.0|sm3sum-linux-amd64-v1.0.0.tar.gz|
|Linux|arm64|v1.0.0|sm3sum-linux-arm64-v1.0.0.tar.gz|
|macOS|amd64 (Intel)|v1.0.0|sm3sum-darwin-amd64-v1.0.0.tar.gz|
|macOS|arm64 (M1/M2)|v1.0.0|sm3sum-darwin-arm64-v1.0.0.tar.gz|
|Windows|amd64|v1.0.0|sm3sum-windows-amd64-v1.0.0.tar.gz|
|Windows|arm64|v1.0.0|sm3sum-windows-arm64-v1.0.0.tar.gz|

解压后将 sm3sum（Windows 下为 sm3sum.exe）复制到系统 PATH 目录即可使用。

⸻

📋 使用示例

计算文件 SM3 哈希值
```bash
sm3sum filename.txt
```
输出类似：
```bash
66c7f0f462eeedd9d1f2d46bdc10e4e24167c4875cf2f7a2297da02b8f4ba8e0  filename.txt
```
生成 BSD 风格哈希值（–tag）
```
sm3sum --tag filename.txt
```
输出：
```
SM3 (filename.txt) = 66c7f0f462eeedd9d1f2d46bdc10e4e24167c4875cf2f7a2297da02b8f4ba8e0
```
校验哈希文件
```
sm3sum filename.txt > filename.sm3
sm3sum -c filename.sm3
```
输出：
```
filename.txt: OK
```
使用零结尾（-z）支持
```
sm3sum -z filename.txt | xargs -0 -n1 echo
```

⸻

⚙️ 命令行参数
```
  -b, --binary         以二进制模式读取（默认）
  -c, --check FILE     从文件读取 SM3 校验和并验证
      --tag            以 BSD 风格输出校验和
  -t, --text           以文本模式读取（默认，忽略）
  -z, --zero           输出行以 NUL 字符结尾
      --ignore-missing 检查模式时忽略缺失文件
      --quiet          只报告错误，不输出成功信息
      --status         不输出任何内容，仅用退出状态表示结果
      --strict         格式错误时返回非零状态
  -w, --warn           格式错误时打印警告
  -h, --help           显示帮助信息
      --version        显示版本信息
```

⸻

📖 文档与支持
	•	查看手册页：
```
man sm3sum
```
	•	Bash 补全脚本：sm3sum_completion
	•	Zsh 补全脚本：_sm3sum

⸻

📝 许可证

MIT License，欢迎自由使用和贡献。

⸻

🙏 鸣谢

基于 [tjfoc/gmsm](https://github.com/tjfoc/gmsm) 的国密算法实现。

⸻

📬 联系方式

如需商业支持或定制开发，欢迎通过 GitHub Issue 或邮件联系。

⸻

提示
请根据实际版本号和 GitHub Release 地址替换上文中的 dist/ 下载链接。

