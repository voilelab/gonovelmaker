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

### 方法一：使用 go install（推薦）

這是最簡單的安裝方式，會自動下載並編譯最新版本：

```bash
go install github.com/voilelab/gonovelmaker/cmd/novelmaker-obs@latest
```

確認安裝成功：

```bash
novelmaker-obs --version
```

!!! tip "PATH 設定"
    確保 `$GOPATH/bin` 或 `$HOME/go/bin` 已加入您的 PATH 環境變數。

### 方法二：從原始碼建置

如果您想要開發或自訂功能，可以從原始碼建置：

```bash
# 1. 複製儲存庫
git clone https://github.com/voilelab/gonovelmaker.git
cd gonovelmaker

# 2. 建置
go build ./cmd/novelmaker-obs

# 3. (可選) 安裝到系統
go install ./cmd/novelmaker-obs
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

=== "方式一：傳統單一後端（簡單）"

    適合只使用 OpenAI API 的使用者：

    ```toml
    openai_api_key = "sk-xxx"
    model = "gpt-4o"
    image_model = "dall-e-3"
    base_url = ""  # 可選：自訂 API 端點
    timeout = 0    # 可選：超時時間（秒），0 表示無限制
    ```

=== "方式二：多後端（進階）"

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

=== "方式三：使用 CLI 命令（推薦）"

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

#### 傳統設定

| 選項 | 說明 | 必需 | 預設值 |
|------|------|------|--------|
| `openai_api_key` | OpenAI API 金鑰 | ✅ | - |
| `model` | 文字生成模型 | ❌ | `gpt-4o` |
| `image_model` | 圖片生成模型 | ❌ | `dall-e-3` |
| `base_url` | API 端點 URL | ❌ | OpenAI 官方端點 |
| `timeout` | 請求超時（秒） | ❌ | `0` (無限制) |

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

### 使用環境變數（僅傳統模式）

您也可以透過環境變數設定（不支援多後端模式）：

```bash
export OPENAI_API_KEY=sk-xxx
export OPENAI_MODEL=gpt-4o              # 選用
export OPENAI_IMAGE_MODEL=dall-e-3      # 選用
export OPENAI_BASE_URL=https://api.example.com/v1  # 選用
```

!!! warning "優先順序"
    設定檔的設定會優先於環境變數。

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

### 找不到 novelmaker-obs 命令

確保 Go 的 bin 目錄在您的 PATH 中：

```bash
# 在 ~/.zshrc 或 ~/.bashrc 中加入：
export PATH=$PATH:$(go env GOPATH)/bin
```

然後重新載入 shell 設定：

```bash
source ~/.zshrc  # 或 source ~/.bashrc
```

### API 金鑰錯誤

檢查：
1. API 金鑰格式正確（以 `sk-` 開頭）
2. 金鑰尚未過期
3. OpenAI 帳號有足夠的額度

### 相容性問題

如果遇到相容性問題，確保：
- Go 版本 ≥ 1.25.4
- 使用最新版本的 gonovelmaker

更新到最新版：
```bash
go install github.com/voilelab/gonovelmaker/cmd/novelmaker-obs@latest
```

## 更新

定期更新以取得最新功能和錯誤修復：

```bash
go install github.com/voilelab/gonovelmaker/cmd/novelmaker-obs@latest
```

## 解除安裝

如需移除 gonovelmaker：

```bash
# 移除執行檔
rm $(which novelmaker-obs)

# 移除設定檔（可選）
rm -rf ~/.novelmaker
```
