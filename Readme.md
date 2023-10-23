### Gogpt是什么？

---------------

**Gogpt** 是一个非常简洁直观的基于[TUI](https://github.com/charmbracelet/bubbletea)的ChatGPT客户端.

### 功能描述

比[j178](https://github.com/j178/chatgpt)更好用。
---------------
- 支持本地代理配置(http或者socks5)。
- 可以在TUI界面进行配置，无需手动编辑json文件或者设置环境变量等。
- 更简洁直观的界面，无冗余功能。
- 更多的Prompt选择，支持一百多项选择。

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


---------------

### What's gogpt?

---------------

**Gogpt** is a simple client for chatgpt based on [TUI](https://github.com/charmbracelet/bubbletea).

### Features

Easies to use than [j178](https://github.com/j178/chatgpt).
---------------
- Local proxy settings.
- Configurations in TUI.
- More simple and intuitive Interface.
- More chatgpt prompt choices.

### Thanks to
- [go-openai](https://github.com/sashabaranov/go-openai)
- [bubbletea](https://github.com/charmbracelet/bubbletea)
- [chatgpt](https://github.com/j178/chatgpt)

### Gallaries

![config](https://github.com/moqsien/gogpt/blob/main/docs/gogpt_conf.png)

![help](https://github.com/moqsien/gogpt/blob/main/docs/gogpt_help_info.png)

![main_tab](https://github.com/moqsien/gogpt/blob/main/docs/gogpt_main.png)

![quicksort1](https://github.com/moqsien/gogpt/blob/main/docs/gpgpt_quicksort_1.png)

![quicksort2](https://github.com/moqsien/gogpt/blob/main/docs/gogpt_quicksort_2.png)
