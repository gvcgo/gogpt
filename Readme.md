## Demo
[![asciicast](https://asciinema.org/a/617183.svg)](https://asciinema.org/a/617183)

### Gogpt是什么？

---------------

**Gogpt** 是一个非常简洁直观的基于[TUI](https://github.com/charmbracelet/bubbletea)的GPT客户端.
支持ChatGPT(3.5, 4.0)和讯飞星火(1.1, 2.1, 3.1)。

### 安装使用

---------------

- 从releases下载，解压可得到可执行文件，支持Win/Mac/Linux。
- 编译安装。
```bash
go install github.com/moqsien/gogpt/tree/main/pkgs/cmd/gogptm@latest
```
- 你也可以安装(下载或者自行编译)[gvc](https://github.com/moqsien/gvc), gvc是一个超级强大的本地开发环境管理工具，已经集成了gogpt。

```text
"→" 切换到下一个Tab

"←" 切换到上一个Tab

"ctrl+w" 在ChatGPT和讯飞星火之间切换。
```

### 配置(Configuration Tab，使用左右箭头切换Tab)
<img src="https://github.com/moqsien/gogpt/blob/main/docs/gogpt_config.png" width="85%">

### 功能描述

---------------
- 比[j178](https://github.com/j178/chatgpt)更好用。
---------------
- 支持本地代理配置(http或者socks5)。
- 可以在TUI界面进行配置，无需手动编辑json文件或者设置环境变量等。
- 更简洁直观的界面，无冗余功能。
- 更多的Prompt选择，支持170+项选择。也可以自行在Configuration页面定制。

### 感谢
- [go-openai](https://github.com/sashabaranov/go-openai)
- [bubbletea](https://github.com/charmbracelet/bubbletea)
- [chatgpt](https://github.com/j178/chatgpt)

### 特别说明

本项目参考了[chatgpt](https://github.com/j178/chatgpt)。绝大部分代码进行了重新设计。剔除了没有太多用的功能。
增加了界面配置，本地代理支持。
最开始，本来想将[chatgpt](https://github.com/j178/chatgpt)集成到[gvc](https://github.com/moqsien/gvc)中，但是发现一些代码属于包内私有，给作者提过issue，但不幸遭拒。
后来想想，可以使用其中一部分，其余自行实现。结果在还未完全实现之前，被作者拉黑了，导致无法fork该项目。原因在于当时功能尚未完成，所以没有在[gvc](https://github.com/moqsien/gvc)的感谢中添加[chatgpt](https://github.com/j178/chatgpt)项目。这，确实有点尴尬……所以，后来把[gvc](https://github.com/moqsien/gvc)中，关于[chatgpt](https://github.com/j178/chatgpt)的引用都删了。
直到最近，有空了，所以，重新另起炉灶，做一个满足自己功能需求的基于TUI的Chatgpt客户端。

最后，特别感谢[chatgpt](https://github.com/j178/chatgpt)项目。没有这个项目，就没有gogpt。

<br/>
<br/>
<br/>
<br/>

---------------

### What's gogpt?

---------------

**Gogpt** is a simple client for GPT based on [TUI](https://github.com/charmbracelet/bubbletea).
Openai chatgpt(3.5, 4.0) and Iflytek spark(1.1, 2.1, 3.1) are supported.

### Install

---------------

- Download zip file from releases, decompress, use the binaries for Win/Mac/Linux.
- Compile and install.
```bash
go install github.com/moqsien/gogpt/tree/main/pkgs/cmd/gogptm@latest
```
- You can also download binaries or complie from source code from [gvc](https://github.com/moqsien/gvc). GVC has [gogpt](https://github.com/moqsien/gvc#subcommand-gpt) integrated.
```text
"→" Switch to next Tab

"←" Switch to previous Tab

"ctrl+w" Switch between Chatgpt and Spark.
```

### Features

---------------
- Easier to use than [j178](https://github.com/j178/chatgpt).
---------------
- Local proxy settings.
- Configurations in TUI.
- More simple and intuitive Interface.
- More chatgpt prompt choices.

### Thanks to
- [go-openai](https://github.com/sashabaranov/go-openai)
- [bubbletea](https://github.com/charmbracelet/bubbletea)
- [chatgpt](https://github.com/j178/chatgpt)
