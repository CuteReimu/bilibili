<div align="center">

# 哔哩哔哩-API-Go版本

![](https://img.shields.io/github/go-mod/go-version/CuteReimu/bilibili?filename=go.mod "语言")
[![](https://img.shields.io/github/workflow/status/CuteReimu/bilibili/Go)](https://github.com/CuteReimu/bilibili/actions/workflows/golangci-lint.yml "代码分析")
[![](https://img.shields.io/github/contributors/CuteReimu/bilibili)](https://github.com/CuteReimu/bilibili/graphs/contributors "贡献者")
[![](https://img.shields.io/github/license/CuteReimu/bilibili)](https://github.com/CuteReimu/bilibili/blob/master/LICENSE "许可协议")
</div>

本项目是基于Go语言编写的哔哩哔哩API调用

**声明**：

1. 本项目遵守 AGPL 开源协议。
2. 本项目基于 [SocialSisterYi/bilibili-API-collect](https://github.com/SocialSisterYi/bilibili-API-collect)
   中描述的接口编写。请尊重该项目作者的努力，遵循该项目的开源要求，禁止一切商业使用。
3. **请勿滥用，本项目仅用于学习和测试！利用本项目提供的接口、文档等造成不良影响及后果与本人无关。**
4. 由于本项目的特殊性，可能随时停止开发或删档
5. 本项目为开源项目，不接受任何形式的催单和索取行为，更不容许存在付费内容

PS：因为B站接口同时支持`http`和`https`，为了数据安全，本项目调用接口时统一使用 `https`

## 快速开始

### 安装

```bash
go get -u github.com/CuteReimu/bilibili
```

因为项目正在不断更新中，请经常使用`go get -u`更新依赖，确保处于最新版本。
