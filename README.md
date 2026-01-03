# gonovelmaker

[![Go Report Card](https://goreportcard.com/badge/github.com/voilelab/gonovelmaker)](https://goreportcard.com/report/github.com/voilelab/gonovelmaker)
[![Go Tests](https://github.com/voilelab/gonovelmaker/actions/workflows/test.yml/badge.svg)](https://github.com/voilelab/gonovelmaker/actions/workflows/test.yml)

[Documentation](https://voilelab.github.io/gonovelmaker/) | [Release Notes](https://github.com/voilelab/gonovelmaker/wiki/Release-Notes)

一個用於管理 Obsidian vault 中小說專案的命令列工具，整合 OpenAI 功能。
目前專注於生成中文小說，所以預設模板和範例均為中文內容。

![](docs/image.png)

## 功能特色

- 初始化結構化的 Obsidian vault 小說專案
- 掃描和分析現有專案結構（支援 JSON 格式輸出）
- 使用 OpenAI GPT 模型生成章節
- 重新生成現有章節內容
- 建立空白章節範本
- 使用 OpenAI GPT 模型生成角色檔案
- 重新生成現有角色檔案
- 使用 OpenAI DALL-E 生成角色圖片
- 匯出小說為文字檔案
- 更新 Obsidian Novel Maker 外掛
- 組織角色（Character/）、世界觀（World/）和故事章節（Story/）
- 支援 YAML frontmatter 元數據
- 自訂提示詞模板系統

## 開始使用

### 先決條件

- Go 1.25.4 或更高版本
- OpenAI API 金鑰（用於章節生成）

### 安裝

```bash
brew tap voilelab/novelmaker https://github.com/voilelab/gonovelmaker
brew install voilelab/novelmaker/novelmaker-obs
```

## 設定

在 `~/.novelmaker/config.toml` 建立設定檔。

```toml
# 指定預設使用的後端
user_llm_backend = "openai"

# 定義多個命名後端
[llm_backend.openai]
type = "openai"
api_key = "sk-xxx"
base_url = "https://api.openai.com/v1"
model = "gpt-4o"
image_model = "dall-e-3"

[llm_backend.openrouter]
type = "openai"
api_key = "sk-or-xxx"
base_url = "https://openrouter.ai/api/v1"
model = "anthropic/claude-3.5-sonnet"
image_model = "openai/dall-e-3"

[llm_backend.custom]
type = "openai"
api_key = "your-api-key"
base_url = "https://your-custom-endpoint.com/v1"
model = "your-model-name"
image_model = "your-image-model"
```

**設定選項說明：**

- `user_llm_backend` - 預設使用的後端名稱（例如："openai", "openrouter"）
- `llm_backend.[名稱]` - 定義一個命名後端，可定義多個
  - `type` - 後端類型（目前支援 "openai"）
  - `api_key` - 該後端的 API 金鑰
  - `base_url` - API 端點 URL
  - `model` - 使用的文字生成模型
  - `image_model` - 使用的圖片生成模型

首次執行工具時，如果設定檔不存在，會自動在 `~/.novelmaker/` 建立範例設定檔。

## 開發說明

### 建置

```bash
go build ./cmd/novelmaker-obs
```

### 執行測試

```bash
go test ./...
```

## 常見問題

**Q: 如何更換 OpenAI 模型？**  
A: 在 `~/.novelmaker/config.toml` 中設定 `model` 欄位，例如 `model = "gpt-4o-mini"` 或 `model = "gpt-4"`。

**Q: 可以使用其他相容 OpenAI API 的服務嗎？**  
A: 可以，設定 `base_url` 欄位指向第三方服務端點，例如 Azure OpenAI 或其他本地部署的模型服務。

**Q: 生成的內容會自動加入專案嗎？**  
A: 是的，`gen-next` 和 `gen-char` 命令會自動建立對應的 Markdown 檔案並加入適當的 frontmatter。

**Q: 如何自訂生成的章節或角色格式？**  
A: 編輯 `~/.novelmaker/templates/` 中的提示詞範本檔案，調整提示詞內容和格式。

**Q: 如何更新 Obsidian 外掛？**  
A: 在 vault 根目錄執行 `novelmaker-obs update-plugin` 命令，會自動複製最新的外掛檔案到 `.obsidian/plugins/obsidian-novelmaker/` 目錄。

**Q: 生成圖片時遇到超時怎麼辦？**  
A: 使用 `--timeout` 參數增加超時時間，例如 `--timeout 120` 設定為 120 秒。DALL-E 生成圖片通常需要較長時間。

## 授權（License）

本專案 NovelMaker Core / CLI 以 BSD 3-Clause License（3 條款 BSD 授權） 釋出。
您可以自由使用、修改、重製與散布本專案的核心程式碼，
也可將其整合進您自己的專案（包含商業用途）。

### 此授權適用的範圍

BSD-3 授權 僅適用於本版本庫中的：

* 核心程式碼（cmd/, internal/, novelmaker/）
* 基礎範例與範本（cmd/novelmaker-obs/templates/）
* 設定檔格式與範例（internal/config/example_config.toml）
* 內建提示詞範本（internal/config/templates/）

以上內容均可依 BSD-3 使用，包含商用、重製與二次散布。

### 不包含在開源授權範圍內的內容

以下類型的內容並未包含於 BSD-3 授權範圍內，仍保留原作者完整權利：

* 作者實際使用的寫作提示語（Prompts）
* 任何 世界書（Worldbook）包含角色設定、劇情資料、背景知識、概念文檔等
* Obsidian Vault 或其他形式的小說原文、草稿、章節規劃

除非另行以書面方式授權，以上內容均不得視為本專案的開源部分。
