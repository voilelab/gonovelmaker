# gonovelmaker 文檔

歡迎使用 **gonovelmaker**！這是一個專為管理 Obsidian vault 中小說專案設計的命令列工具，整合了 OpenAI API 功能，協助您更高效地創作小說。

![gonovelmaker](image.png)

## ✨ 主要功能

- 📁 **專案管理**：初始化結構化的 Obsidian vault 小說專案
- 🔍 **專案掃描**：分析現有專案結構，支援 JSON 格式輸出
- 📖 **章節生成**：使用 OpenAI GPT 模型自動生成章節內容
- 👤 **角色創建**：基於專案背景生成詳細的角色檔案
- 🎨 **圖片生成**：使用 DALL-E 為角色生成視覺圖像
- 📤 **內容匯出**：將完整小說匯出為文字檔案
- 🔧 **自訂範本**：支援自訂提示詞模板系統
- 🔌 **Obsidian 整合**：內建 Novel Maker 外掛

## 🚀 快速開始

### 安裝

```bash
brew tap voilelab/novelmaker https://github.com/voilelab/gonovelmaker
brew install voilelab/novelmaker/novelmaker-obs
```

### 初始化專案

```bash
novelmaker-obs init
```

### 生成第一個章節

```bash
novelmaker-obs gen-next --title "第一章：開始"
```

## 📚 文檔導覽

<div class="grid cards" markdown>

- 📥 **安裝指南**

    ---

    了解如何安裝和配置 gonovelmaker
    
    [前往安裝指南](install.md)

- 💻 **CLI 命令列工具**

    ---

    探索所有可用的命令和選項
    
    [查看命令參考](cli/commands.md)

- 🧩 **Obsidian 外掛**

    ---

    在 Obsidian 中使用 Novel Maker 外掛
    
    [外掛使用說明](plugin/obsidian.md)

- 📖 **世界書結構**

    ---

    了解專案檔案結構和 frontmatter 格式
    
    [檔案結構說明](worldbook/schema.md)

</div>

## 🎯 使用場景

### 小說作家
使用 AI 輔助生成章節大綱、角色設定，加速創作流程。

### 世界觀設計師
結構化管理複雜的世界觀設定、角色關係和劇情線。

### Obsidian 愛好者
在熟悉的 Obsidian 環境中進行小說創作，享受雙向連結和圖形化視圖。

## 🛠️ 技術特色

- **Go 語言開發**：高效能、跨平台支援
- **OpenAI 整合**：支援 GPT-4、DALL-E 等最新模型
- **多後端支援**：可配置多個 LLM 後端（OpenAI、OpenRouter 等）
- **YAML Frontmatter**：結構化元數據管理
- **範本系統**：靈活的提示詞自訂能力

## 📖 專案狀態

[![Go Report Card](https://goreportcard.com/badge/github.com/voilelab/gonovelmaker)](https://goreportcard.com/report/github.com/voilelab/gonovelmaker)

當前版本專注於生成中文小說內容，預設模板和範例均為中文。

## 🤝 貢獻與授權

本專案核心程式碼以 BSD 3-Clause License 釋出，歡迎貢獻和回饋。

[GitHub Repository](https://github.com/voilelab/gonovelmaker){ .md-button .md-button--primary }

---

!!! tip "需要幫助？"
    如果您在使用過程中遇到問題，請查看各章節的詳細說明，或在 GitHub 上提出 Issue。
