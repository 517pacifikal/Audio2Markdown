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

# 安装并运行支持 RediSearch 的 Redis 服务，作为默认的 RAG 向量数据库
docker run -d --name redis-stack-server -p 6379:6379 redis/redis-stack-server:latest

# 启动 redis 客户端
audio2md % redis-cli -h 127.0.0.1 -p 6379

# 创建索引，注意维度与实际使用的 embedding 模型对齐
FT.CREATE rag ON HASH PREFIX 1 rag. SCHEMA content TEXT vector VECTOR FLAT 6 TYPE FLOAT32 DIM 1024 DISTANCE_METRIC COSINE

```

### Step2. 配置管理

1. 在项目根目录创建 `config.json` 文件，这是本项目的核心配置文件。文件具体字段参见 [config.json 配置文件说明](#configjson-配置文件说明)

## `config.json` 配置文件说明

本项目的 `config.json` 文件用于配置音频转写、文档加载、向量化、索引等各环节的参数。结构示例如下：

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
    },
    "INDEXING": {
        "FILE_PATH": "./files/",
        "EMBEDDING": {
            "BASE_URL": "https://dashscope.aliyuncs.com/compatible-mode/v1",
            "API_KEY": "sk-xxx",
            "MODEL": "text-embedding-v4"
        },
        "INDEXER": {
            "TYPE": "REDIS",
            "REDIS": {
                "ADDR": "localhost:6379",
                "KEY_PREFIX": "rag.",
                "BATCH_SIZE": 1
            },
            "FAISS": {
                "INDEX_PATH": "/tmp/faiss.index",
                "BATCH_SIZE": 1
            }
        },
        "LOADER": {},
        "TRANSFORMER": {
            "HEADERS": {
                "#": "title"
            },
            "TRIM_HEADERS": false
        }
    }
}
```

各字段含义说明：

- `AUDIO_CONFIGS`：音频转文本相关配置。
    - `MODEL_SRC`：指定当前使用的语音转写模型来源（如 "BYTEDANCE"）。
    - `BYTEDANCE`：字节大模型语音转写相关参数。
        - `AUDIO_FILE`：输入音频文件路径。
        - `OUTPUT_FILE`：转写后文本输出路径。
        - `APP_KEY`、`ACCESS_KEY`：API 访问凭证。
        - `TOS_BUCKET`、`TOS_REGION`、`TOS_ENDPOINT`、`TOS_ACCESS_KEY`、`TOS_SECRET_KEY`：TOS 对象存储相关配置，用于音频文件上传。

- `INDEXING`：RAG 检索增强生成相关配置。
    - `FILE_PATH`：待加载文档的目录或文件路径，支持批量加载目录下所有文件。
    - `EMBEDDING`：向量化模型相关配置。
        - `BASE_URL`：向量化 API 服务地址。
        - `API_KEY`：API 访问密钥。
        - `MODEL`：使用的向量模型名称。
    - `INDEXER`：向量数据库相关配置。
        - `TYPE`：索引类型，支持 "REDIS" 或 "FAISS"。
        - `REDIS`：Redis 相关参数（如地址、key 前缀、批处理大小）。
        - `FAISS`：FAISS 相关参数（如索引文件路径、批处理大小）。
    - `LOADER`：文档加载器相关配置（如有特殊参数可在此扩展）。
    - `TRANSFORMER`：文档切分与预处理相关配置。
        - `HEADERS`：用于 Markdown 等文档的标题层级映射。
        - `TRIM_HEADERS`：是否去除标题前后的空白字符。

请根据实际需求填写和调整上述参数，确保各环节配置正确，系统即可顺利完成音频转写、文档加载、向量化、索引与检索等全



