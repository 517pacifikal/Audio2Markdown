# Audio2Markdown

Audio2Markdown 是一个基于大语言模型(LLM)能力实现的文档生成方案，解决的是语音对话内容转化为 Markdown 格式的文档的问题。

Audio2Markdown 支持的场景包括但不限于：

- 面试录音 -> 面试复盘文档
- 会议录音 -> 会议纪要总结
- 电话录音 -> 电话录音摘要

每天都在海量面试，没有时间复盘？可以用 Audio2Markdown 帮你直接生成复盘文档！

公司/部门会议涉及内部术语，语音识别不准确？可以用 Audio2Markdown 帮你直接生成会议纪要！检索增强生成(RAG)能力保证生成结果的专业性！

欢迎大家广泛探索本项目的应用场景。这是一个简洁的实现方案，欢迎任何建议、贡献和反馈！

## 运行指南

### Step1. 依赖安装

```bash

# 创建虚拟环境并进入
python3 -m venv .audio2md
source .audio2md/bin/activate  

# 安装Python依赖
pip install -r requirements.txt

docker run -p 6379:6379 --name redismod \
-v /mydata/redismod/data:/data \
-d redislabs/redismod:preview

```

### Step2. 配置管理

1. 在项目根目录创建 `config.json` 文件，这是本项目的核心配置文件。文件具体字段参见 [config.json 配置文件说明](#configjson-配置文件说明)

## `config.json` 配置文件说明

本项目的 `config.json` 文件用于配置音频转写及模型调用的相关参数。结构示例如下：

```json
{
    "AUDIO_CONFIGS": {
        "MODEL_SRC": "BYTEDANCE",
        "BYTEDANCE": {
            "AUDIO_FILE": "./audio/input/xxx.mp3",
            "OUTPUT_FILE": "./audio/output/xxx.txt",
            "APP_KEY": "...",
            "ACCESS_KEY": "...",
            "TOS_BUCKET": "...",
            "TOS_REGION": "...",
            "TOS_ENDPOINT": "...",
            "TOS_ACCESS_KEY": "...",
            "TOS_SECRET_KEY": "..."
        }
    }
}
```

各字段含义说明：

- `AUDIO_CONFIGS`：音频转文本相关配置的总入口。
    - `MODEL_SRC`：指定当前使用的模型来源（如 "BYTEDANCE"）。
    - `BYTEDANCE`：具体模型来源的详细配置。
        - `AUDIO_FILE`：导入音频文件的路径。
        - `OUTPUT_FILE`：转写后文本的输出路径。
        - `APP_KEY`：模型服务的应用 App Key。
        - `ACCESS_KEY`：模型服务的访问密钥。
        - `TOS_BUCKET`：TOS（对象存储）桶名，用于上传本地音频文件。
        - `TOS_REGION`：TOS 区域。
        - `TOS_ENDPOINT`：TOS 服务 Endpoint。
        - `TOS_ACCESS_KEY`：TOS 访问密钥。
        - `TOS_SECRET_KEY`：TOS



