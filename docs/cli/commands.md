# 命令參考

本頁面提供所有 gonovelmaker CLI 命令的完整參考文檔。

## 專案管理命令

### init

初始化新的小說專案結構。

```bash
novelmaker-obs init
```

**功能：**

在當前目錄建立完整的 Obsidian vault 結構，包括：

- `Config/project.md` - 專案設定檔
- `World/001_world_sample.md` - 世界觀範例
- `Character/character_sample.md` - 角色範例  
- `Story/001_prologue.md` - 故事序章範例
- `.obsidian/app.json` - Obsidian 設定
- `README.md` - 專案說明

**範例：**

```bash
mkdir my-novel
cd my-novel
novelmaker-obs init
```

**輸出：**

```
✅ 專案初始化成功！
📁 專案路徑: /path/to/my-novel
📖 已建立範例檔案
```

---

### scan

掃描並分析現有專案結構。

```bash
novelmaker-obs scan [flags]
```

**旗標：**

| 旗標 | 類型 | 說明 |
|------|------|------|
| `--json` | bool | 以 JSON 格式輸出 |

**功能：**

- 讀取專案設定
- 列出所有世界觀條目
- 列出所有角色檔案
- 列出所有章節（依順序排列）
- 顯示專案統計資訊

**範例：**

```bash
# 人類可讀格式
novelmaker-obs scan

# JSON 格式（適合程式處理）
novelmaker-obs scan --json
```

**JSON 輸出範例：**

```json
{
  "project": {
    "id": "my-novel-001",
    "name": "我的小說",
    "world": "奇幻世界"
  },
  "world_entries": [...],
  "characters": [...],
  "chapters": [...]
}
```

---

### config-check

檢查設定檔是否正確配置。

```bash
novelmaker-obs config-check
```

**功能：**

- 驗證設定檔格式
- 檢查 API 金鑰是否存在
- 顯示當前設定摘要
- 測試 API 連線（可選）

**範例：**

```bash
novelmaker-obs config-check
```

**輸出：**

```
✅ 設定檔驗證通過
📋 設定摘要：
  - 後端: openai
  - 模型: gpt-4o
  - 圖片模型: dall-e-3
  - API 端點: https://api.openai.com/v1
```

---

## 後端管理命令

### backend add

新增或編輯 LLM 後端設定。

```bash
novelmaker-obs backend add <name> [flags]
```

**必需參數：**

| 參數 | 類型 | 說明 |
|------|------|------|
| `<name>` | string | 後端名稱（唯一識別符） |

**可選參數：**

| 參數 | 類型 | 預設值 | 說明 |
|------|------|--------|------|
| `--type` | string | openai | 後端類型（openai, openrouter 等） |
| `--api_key` | string | - | API 金鑰 |
| `--base_url` | string | - | API 端點 URL |
| `--model` | string | - | 文字模型名稱 |
| `--image_model` | string | - | 圖片模型名稱 |
| `--timeout` | int | 0 | 請求超時（秒） |

**功能：**

- 新增新的 LLM 後端設定
- 更新現有後端設定（僅修改指定的欄位）
- 自動儲存到設定檔

**範例：**

```bash
# 新增 OpenAI 後端
novelmaker-obs backend add openai \
  --type openai \
  --api_key "sk-xxx" \
  --model "gpt-4o"

# 新增 OpenRouter 後端
novelmaker-obs backend add openrouter \
  --type openrouter \
  --api_key "sk-or-xxx" \
  --base_url "https://openrouter.ai/api/v1" \
  --model "anthropic/claude-3.5-sonnet"

# 新增自訂後端
novelmaker-obs backend add custom \
  --type openai \
  --api_key "xxx" \
  --base_url "https://api.example.com/v1" \
  --model "my-model" \
  --timeout 120

# 更新現有後端（僅修改模型）
novelmaker-obs backend add openai \
  --model "gpt-4o-mini"
```

**輸出：**

```
✓ Backend 'openai' added successfully
```

或

```
✓ Backend 'openai' updated successfully
```

---

### backend list

列出所有已設定的 LLM 後端。

```bash
novelmaker-obs backend list [flags]
```

**可選參數：**

| 參數 | 類型 | 說明 |
|------|------|------|
| `--json` | bool | 以 JSON 格式輸出 |

**功能：**

- 顯示所有後端設定
- 標記預設後端
- 隱藏 API 金鑰（僅顯示最後 4 碼）
- 支援 JSON 格式輸出

**範例：**

```bash
# 人類可讀格式
novelmaker-obs backend list

# JSON 格式
novelmaker-obs backend list --json
```

**輸出範例（人類可讀）：**

```
Configured backends:

• openai (default)
  Type:        openai
  API Key:     ...1234
  Model:       gpt-4o
  Image Model: dall-e-3

• openrouter
  Type:        openrouter
  API Key:     ...abcd
  Base URL:    https://openrouter.ai/api/v1
  Model:       anthropic/claude-3.5-sonnet
  Timeout:     120s
```

**輸出範例（JSON）：**

```json
[
  {
    "name": "openai",
    "type": "openai",
    "api_key": "...1234",
    "base_url": "",
    "model": "gpt-4o",
    "image_model": "dall-e-3",
    "timeout": 0,
    "is_default": true
  },
  {
    "name": "openrouter",
    "type": "openrouter",
    "api_key": "...abcd",
    "base_url": "https://openrouter.ai/api/v1",
    "model": "anthropic/claude-3.5-sonnet",
    "image_model": "",
    "timeout": 120,
    "is_default": false
  }
]
```

---

### backend check

檢查 LLM 後端連線是否正常。

```bash
novelmaker-obs backend check <name> [flags]
```

**必需參數：**

| 參數 | 類型 | 說明 |
|------|------|------|
| `<name>` | string | 要檢查的後端名稱 |

**可選參數：**

| 參數 | 類型 | 預設值 | 說明 |
|------|------|--------|------|
| `--timeout` | int | 30 | 請求超時（秒） |
| `--json` | bool | false | 以 JSON 格式輸出 |

**功能：**

- 向 LLM 後端發送測試請求
- 驗證 API 金鑰和設定
- 測量回應時間
- 顯示 token 使用情況

**範例：**

```bash
# 檢查後端連線
novelmaker-obs backend check openai

# 使用較長超時
novelmaker-obs backend check openrouter --timeout 60

# JSON 輸出
novelmaker-obs backend check openai --json
```

**輸出範例（成功）：**

```
Checking backend 'openai'...
  Type:     openai
  Model:    gpt-4o

Sending test request... ✅ SUCCESS

Response time: 1.234s
Response: OK

Token usage:
  Input tokens:  15
  Output tokens: 2
  Total tokens:  17
```

**輸出範例（失敗）：**

```
Checking backend 'openai'...
  Type:     openai
  Model:    gpt-4o

Sending test request... ❌ FAILED
Error: check failed: API key is invalid
```

**JSON 輸出範例（成功）：**

```json
{
  "success": true,
  "backend": "openai",
  "backend_type": "openai",
  "model": "gpt-4o",
  "base_url": "",
  "response_time": 1234,
  "response": "OK",
  "token_usage": {
    "input_tokens": 15,
    "output_tokens": 2,
    "total_tokens": 17
  }
}
```

---

### backend use

設定預設使用的 LLM 後端。

```bash
novelmaker-obs backend use <name>
```

**必需參數：**

| 參數 | 類型 | 說明 |
|------|------|------|
| `<name>` | string | 要設為預設的後端名稱 |

**功能：**

- 將指定後端設為預設
- 所有生成命令將使用此後端
- 自動儲存設定

**範例：**

```bash
# 設定 OpenAI 為預設
novelmaker-obs backend use openai

# 切換到 OpenRouter
novelmaker-obs backend use openrouter
```

**輸出：**

```
✓ Default backend set to 'openai'
```

---

### backend remove

移除 LLM 後端設定。

```bash
novelmaker-obs backend remove <name>
```

**必需參數：**

| 參數 | 類型 | 說明 |
|------|------|------|
| `<name>` | string | 要移除的後端名稱 |

**功能：**

- 刪除指定的後端設定
- 如果移除的是預設後端，會提示設定新的預設
- 自動儲存設定

**範例：**

```bash
# 移除後端
novelmaker-obs backend remove custom
```

**輸出：**

```
✓ Backend 'custom' removed successfully
```

如果移除的是預設後端：

```
✓ Backend 'openai' removed successfully
Note: No default backend set. Use 'backend use <name>' to set one.
```

---

## 章節生成命令

### gen-next

生成下一個章節。

```bash
novelmaker-obs gen-next --title "章節標題" [flags]
```

**必需參數：**

| 參數 | 簡寫 | 類型 | 說明 |
|------|------|------|------|
| `--title` | `-t` | string | 章節標題 |

**可選參數：**

| 參數 | 簡寫 | 類型 | 預設值 | 說明 |
|------|------|------|--------|------|
| `--prompt` | `-p` | string | - | 自訂生成提示 |
| `--json` | `-j` | bool | false | JSON 輸出 |
| `--model` | - | string | - | 覆蓋文字模型 |
| `--timeout` | - | int | 0 | API 超時（秒） |

**功能：**

1. 分析專案背景（世界觀、角色）
2. 讀取前面的章節作為上下文
3. 使用 AI 生成新章節內容
4. 自動儲存為 Markdown 檔案
5. 更新 frontmatter 元數據

**範例：**

```bash
# 基本使用
novelmaker-obs gen-next --title "第二章：冒險開始"

# 使用自訂提示
novelmaker-obs gen-next \
  --title "第三章：決戰" \
  --prompt "描寫一場激烈的魔法對決"

# 使用不同模型
novelmaker-obs gen-next \
  --title "第四章" \
  --model gpt-4o-mini
```

---

### gen-curr

重新生成現有章節。

```bash
novelmaker-obs gen-curr --filepath "章節路徑" [flags]
```

**必需參數：**

| 參數 | 簡寫 | 類型 | 說明 |
|------|------|------|------|
| `--filepath` | `-f` | string | 章節檔案路徑（相對於 vault 根目錄） |

**可選參數：**

| 參數 | 簡寫 | 類型 | 預設值 | 說明 |
|------|------|------|--------|------|
| `--prev-chapters` | `-c` | int | 3 | 用作上下文的前面章節數 |
| `--json` | `-j` | bool | false | JSON 輸出 |
| `--model` | - | string | - | 覆蓋文字模型 |
| `--timeout` | - | int | 0 | API 超時（秒） |

**功能：**

- 讀取章節 frontmatter 中的 `prompt` 欄位
- 基於該提示重新生成章節內容
- 保留 frontmatter 元數據
- 更新章節檔案

**範例：**

```bash
# 重新生成章節
novelmaker-obs gen-curr --filepath "Story/002_ch2.md"

# 使用更多前文上下文
novelmaker-obs gen-curr \
  --filepath "Story/005_ch5.md" \
  --prev-chapters 5

# JSON 輸出
novelmaker-obs gen-curr \
  --filepath "Story/003_ch3.md" \
  --json
```

---

### gen-next-empty

建立空白章節範本。

```bash
novelmaker-obs gen-next-empty --title "章節標題" [flags]
```

**必需參數：**

| 參數 | 簡寫 | 類型 | 說明 |
|------|------|------|------|
| `--title` | `-t` | string | 章節標題 |

**可選參數：**

| 參數 | 簡寫 | 類型 | 預設值 | 說明 |
|------|------|------|--------|------|
| `--prompt` | `-p` | string | - | 額外的提示或說明 |
| `--json` | `-j` | bool | false | JSON 輸出 |

**功能：**

- 建立章節檔案框架
- 包含完整的 frontmatter
- 內容區域為空白
- 適合手動撰寫或稍後生成

**範例：**

```bash
# 建立空白章節
novelmaker-obs gen-next-empty --title "第六章：待續"

# 附帶提示說明
novelmaker-obs gen-next-empty \
  --title "第七章：高潮" \
  --prompt "描寫最終決戰，需要融入前面所有伏筆"
```

---

## 角色生成命令

### gen-char

生成新角色檔案。

```bash
novelmaker-obs gen-char --prompt "角色描述" [flags]
```

**必需參數：**

| 參數 | 簡寫 | 類型 | 說明 |
|------|------|------|------|
| `--prompt` | `-p` | string | 角色描述或生成提示 |

**可選參數：**

| 參數 | 簡寫 | 類型 | 預設值 | 說明 |
|------|------|------|--------|------|
| `--name` | `-n` | string | - | 指定角色名稱 |
| `--json` | `-j` | bool | false | JSON 輸出 |
| `--model` | - | string | - | 覆蓋文字模型 |
| `--timeout` | - | int | 0 | API 超時（秒） |

**功能：**

1. 基於專案背景和提示生成角色
2. 自動提取或使用指定的角色名稱
3. 生成詳細的角色描述
4. 儲存為 Markdown 檔案

**範例：**

```bash
# 基本使用
novelmaker-obs gen-char \
  --prompt "一位神秘的魔法師，擅長時間魔法"

# 指定角色名稱
novelmaker-obs gen-char \
  --prompt "主角的摯友，勇敢但衝動" \
  --name "艾倫"

# 詳細描述
novelmaker-obs gen-char \
  --prompt "反派角色：野心勃勃的貴族，掌握黑魔法，曾經是主角的導師"
```

---

### gen-char-curr

重新生成現有角色。

```bash
novelmaker-obs gen-char-curr --filepath "角色路徑" [flags]
```

**必需參數：**

| 參數 | 簡寫 | 類型 | 說明 |
|------|------|------|------|
| `--filepath` | `-f` | string | 角色檔案路徑（相對於 vault 根目錄） |

**可選參數：**

| 參數 | 簡寫 | 類型 | 預設值 | 說明 |
|------|------|------|--------|------|
| `--json` | `-j` | bool | false | JSON 輸出 |
| `--model` | - | string | - | 覆蓋文字模型 |
| `--timeout` | - | int | 0 | API 超時（秒） |

**功能：**

- 讀取角色 frontmatter 中的 `prompt` 欄位
- 基於該提示重新生成角色描述
- 保留角色名稱和其他元數據
- 更新角色檔案

**範例：**

```bash
# 重新生成角色
novelmaker-obs gen-char-curr \
  --filepath "Character/protagonist.md"

# JSON 輸出
novelmaker-obs gen-char-curr \
  --filepath "Character/villain.md" \
  --json
```

---

### gen-char-img

生成角色圖片。

```bash
novelmaker-obs gen-char-img --name "角色名稱" [flags]
```

**必需參數：**

| 參數 | 簡寫 | 類型 | 說明 |
|------|------|------|------|
| `--name` | `-n` | string | 角色名稱 |

**可選參數：**

| 參數 | 簡寫 | 類型 | 預設值 | 說明 |
|------|------|------|--------|------|
| `--prompt` | `-p` | string | - | 自訂圖片生成提示 |
| `--output-dir` | - | string | Character/ | 圖片輸出目錄 |
| `--json` | `-j` | bool | false | JSON 輸出 |
| `--image-model` | - | string | - | 覆蓋圖片模型 |
| `--timeout` | - | int | 60 | API 超時（秒） |

**功能：**

1. 查找角色檔案
2. 讀取角色描述
3. 使用 DALL-E 生成圖片
4. 下載並儲存圖片
5. 支援自訂提示詞

**範例：**

```bash
# 基本使用
novelmaker-obs gen-char-img --name "愛麗絲"

# 自訂提示
novelmaker-obs gen-char-img \
  --name "魔法師梅林" \
  --prompt "Portrait of a wise old wizard with long white beard, wearing blue robes with stars, holding a magical staff, fantasy art style"

# 指定輸出目錄
novelmaker-obs gen-char-img \
  --name "主角" \
  --output-dir "assets/characters/"

# 增加超時時間
novelmaker-obs gen-char-img \
  --name "反派" \
  --timeout 120
```

**支援的圖片模型：**

- `dall-e-2` - 較快，成本較低
- `dall-e-3` - 更高品質（預設）

---

## 匯出與外掛命令

### export

匯出完整小說為文字檔。

```bash
novelmaker-obs export --output "檔案路徑" --type txt [flags]
```

**必需參數：**

| 參數 | 簡寫 | 類型 | 說明 |
|------|------|------|------|
| `--output` | `-o` | string | 輸出檔案路徑 |
| `--type` | `-t` | string | 匯出格式（目前僅支援 txt） |

**功能：**

- 讀取所有章節（依順序）
- 組合為單一文件
- 包含專案名稱和章節標題
- 移除 frontmatter
- 格式化輸出

**範例：**

```bash
# 匯出為文字檔
novelmaker-obs export \
  --output "my-novel.txt" \
  --type txt

# 匯出到特定目錄
novelmaker-obs export \
  --output "../exports/novel-v1.txt" \
  --type txt
```

**輸出格式：**

```
# 我的小說

## 第一章：序章

章節內容...

## 第二章：冒險開始

章節內容...
```

---

### update-plugin

更新 Obsidian Novel Maker 外掛。

```bash
novelmaker-obs update-plugin
```

**功能：**

- 複製最新的外掛檔案到 vault
- 更新 main.js
- 更新 manifest.json
- 更新 styles.css
- 建立外掛目錄（如不存在）

**注意事項：**

- 必須在 vault 根目錄執行
- 需要有 `.obsidian` 目錄
- 會覆蓋現有的外掛檔案

**範例：**

```bash
cd /path/to/your/vault
novelmaker-obs update-plugin
```

**輸出：**

```
✅ 外掛更新成功！
📦 已更新檔案：
  - .obsidian/plugins/obsidian-novelmaker/main.js
  - .obsidian/plugins/obsidian-novelmaker/manifest.json
  - .obsidian/plugins/obsidian-novelmaker/styles.css
💡 請重新啟動 Obsidian 以載入更新
```

---

## 進階用法

### 組合多個命令

使用 shell 腳本自動化工作流程：

```bash
#!/bin/bash

# 初始化專案
novelmaker-obs init

# 生成三個角色
for char in "主角" "配角" "反派"; do
  novelmaker-obs gen-char --name "$char" --prompt "待設定"
done

# 生成前五章
for i in {1..5}; do
  novelmaker-obs gen-next --title "第${i}章"
done

# 匯出
novelmaker-obs export --output "draft.txt" --type txt
```

### JSON 處理

使用 `jq` 處理 JSON 輸出：

```bash
# 取得所有章節標題
novelmaker-obs scan --json | jq -r '.chapters[].title'

# 統計角色數量
novelmaker-obs scan --json | jq '.characters | length'

# 取得最新章節
novelmaker-obs scan --json | jq -r '.chapters[-1].filepath'
```

### 參數覆蓋

臨時使用不同設定：

```bash
# 使用不同後端
novelmaker-obs gen-next \
  --title "測試章節" \
  --model "anthropic/claude-3.5-sonnet"
```

## 疑難排解

### 常見錯誤

**錯誤：找不到專案設定檔**

```bash
# 確認當前目錄
pwd

# 掃描專案結構
ls -la Config/

# 重新初始化（如果需要）
novelmaker-obs init
```

**錯誤：API 請求超時**

```bash
# 增加超時時間
novelmaker-obs gen-next \
  --title "長章節" \
  --timeout 300
```

**錯誤：權限不足**

```bash
# 檢查檔案權限
ls -l ~/.novelmaker/

# 修正權限
chmod 600 ~/.novelmaker/config.toml
```

## 下一步

- 💡 [使用範例](examples.md) - 查看實際使用案例
- 🏗️ [世界書結構](../worldbook/schema.md) - 了解檔案格式
