# 安裝指南

本指南將協助您安裝和配置 gonovelmaker。

## 系統需求

### 必要條件

- **Go 1.25.4 或更高版本**
- **OpenAI API 金鑰**（用於 AI 生成功能）
- **作業系統**：macOS、Linux 或 Windows

### 可選條件

- **Obsidian**：若要使用 Obsidian 外掛功能（[下載 Obsidian](https://obsidian.md/)）
- **Git**：用於從原始碼安裝

## 安裝方法

### 方法一：使用 brew install（推薦）

這是最簡單的安裝方式，會自動下載並編譯最新釋出版本：

```bash
brew tap voilelab/novelmaker https://github.com/voilelab/gonovelmaker
brew install voilelab/novelmaker/novelmaker-obs
```

確認安裝成功：

```bash
novelmaker-obs --version
```

### 方法二：從原始碼建置

如果您想要開發或自訂功能，可以從原始碼建置：

```bash
# 1. 複製儲存庫
git clone https://github.com/voilelab/gonovelmaker.git
cd gonovelmaker

# 2. 建置
go build ./cmd/novelmaker-obs
```

## 設定

### 建立設定檔

首次執行時，工具會自動在 `~/.novelmaker/` 建立範例設定檔：

```bash
novelmaker-obs --help
```

設定檔位置：`~/.novelmaker/config.toml`

### 設定方式

有兩種設定方式可供選擇：

=== "方式一：多後端（進階）"

    適合需要切換不同 LLM 服務的使用者：

    ```toml
    # 指定預設使用的後端
    user_llm_backend = "openai"

    # OpenAI 官方
    [llm_backend.openai]
    type = "openai"
    api_key = "sk-xxx"
    base_url = "https://api.openai.com/v1"
    model = "gpt-4o"
    image_model = "dall-e-3"

    # OpenRouter（多模型整合）
    [llm_backend.openrouter]
    type = "openai"
    api_key = "sk-or-xxx"
    base_url = "https://openrouter.ai/api/v1"
    model = "anthropic/claude-3.5-sonnet"
    image_model = "openai/dall-e-3"

    # 自訂端點
    [llm_backend.custom]
    type = "openai"
    api_key = "your-api-key"
    base_url = "https://your-endpoint.com/v1"
    model = "your-model"
    image_model = "your-image-model"
    ```

=== "方式二：使用 CLI 命令（推薦）"

    使用命令列工具管理後端配置，無需手動編輯設定檔：

    ```bash
    # 新增 OpenAI 後端
    novelmaker-obs backend add openai \
      --type openai \
      --api_key "sk-xxx" \
      --model "gpt-4o" \
      --image_model "dall-e-3"

    # 新增 OpenRouter 後端
    novelmaker-obs backend add openrouter \
      --type openrouter \
      --api_key "sk-or-xxx" \
      --base_url "https://openrouter.ai/api/v1" \
      --model "anthropic/claude-3.5-sonnet"

    # 設定預設後端
    novelmaker-obs backend use openai

    # 檢查後端連線
    novelmaker-obs backend check openai

    # 查看所有後端
    novelmaker-obs backend list
    ```

    !!! tip "推薦使用 CLI 命令"
        使用 CLI 命令管理後端更安全、更方便：
        
        - ✅ 自動驗證設定格式
        - ✅ 支援測試連線
        - ✅ 避免手動編輯錯誤
        - ✅ 便於腳本自動化

### 設定選項說明

#### 多後端設定

| 選項 | 說明 | 必需 |
|------|------|------|
| `user_llm_backend` | 預設後端名稱 | ✅ |
| `llm_backend.[名稱].type` | 後端類型 | ✅ |
| `llm_backend.[名稱].api_key` | API 金鑰 | ✅ |
| `llm_backend.[名稱].base_url` | API 端點 | ✅ |
| `llm_backend.[名稱].model` | 文字模型 | ✅ |
| `llm_backend.[名稱].image_model` | 圖片模型 | ✅ |

### 後端管理命令

使用 CLI 命令管理後端配置：

| 命令 | 功能 |
|------|------|
| `backend add <name>` | 新增或更新後端配置 |
| `backend list` | 列出所有後端 |
| `backend check <name>` | 測試後端連線 |
| `backend use <name>` | 設定預設後端 |
| `backend remove <name>` | 移除後端配置 |

詳細使用方式請參考 [CLI 命令參考](cli/commands.md#後端管理命令)。

## 取得 OpenAI API 金鑰

1. 前往 [OpenAI Platform](https://platform.openai.com/)
2. 註冊或登入帳號
3. 導航至 API Keys 頁面
4. 點擊 "Create new secret key"
5. 複製金鑰並儲存到設定檔

!!! danger "安全提醒"
    - 不要將 API 金鑰提交到版本控制系統
    - 定期輪換您的 API 金鑰
    - 監控 API 使用量以避免意外費用

## 驗證安裝

執行以下命令驗證設定：

```bash
novelmaker-obs config-check
```

如果設定正確，您應該會看到設定資訊的摘要。

## 下一步

- 📖 [CLI 概覽](cli/overview.md) - 了解可用的命令
- 🚀 [使用範例](cli/examples.md) - 查看實際使用案例
- 🏗️ [世界書結構](worldbook/schema.md) - 了解專案檔案結構

## 常見問題

### API 金鑰錯誤

檢查：

1. API 金鑰格式正確（以 `sk-` 開頭）
2. 金鑰尚未過期
3. OpenAI 帳號有足夠的額度

### 相容性問題

如果遇到相容性問題，確保：

- Go 版本 ≥ 1.25.4
- 使用最新版本的 gonovelmaker

## 更新

定期更新以取得最新功能和錯誤修復：

```bash
brew upgrade voilelab/novelmaker/novelmaker-obs
```
