# Obsidian Novel Maker 外掛

Novel Maker 是專為 gonovelmaker 設計的 Obsidian 外掛，提供圖形化介面來管理和創作小說專案。

## 功能概覽

- 🎨 **直覺介面**：在 Obsidian 中直接操作 gonovelmaker 功能
- 🔄 **即時同步**：與 CLI 工具完全相容
- 📝 **編輯增強**：提供寫作輔助功能
- 🎯 **快速操作**：命令面板整合
- 📊 **專案視圖**：視覺化專案結構

## 安裝外掛

### 方法一：使用 CLI 工具更新（推薦）

在您的 vault 根目錄執行：

```bash
cd /path/to/your/novel-vault
novelmaker-obs update-plugin
```

這會自動將最新的外掛檔案複製到正確位置：
- `.obsidian/plugins/obsidian-novelmaker/main.js`
- `.obsidian/plugins/obsidian-novelmaker/manifest.json`
- `.obsidian/plugins/obsidian-novelmaker/styles.css`

### 方法二：手動安裝

1. 在 vault 根目錄建立外掛目錄：
   ```bash
   mkdir -p .obsidian/plugins/obsidian-novelmaker
   ```

2. 從 gonovelmaker 儲存庫複製檔案：
   ```bash
   cp cmd/novelmaker-obs/obsidian-novelmaker/* \
      .obsidian/plugins/obsidian-novelmaker/
   ```

3. 重新啟動 Obsidian

4. 在 Obsidian 設定中啟用外掛：
   - Settings → Community plugins
   - 找到 "Novel Maker"
   - 開啟開關

## 啟用外掛

1. 開啟 Obsidian
2. 前往 **Settings**（設定）
3. 選擇 **Community plugins**（社群外掛）
4. 關閉 **Safe mode**（安全模式，首次使用時）
5. 點擊 **Browse**（瀏覽）
6. 在已安裝外掛列表中找到 **Novel Maker**
7. 點擊開關啟用

!!! warning "安全模式"
    首次使用社群外掛時需要關閉安全模式。Novel Maker 是本地外掛，不會連線到外部服務。

## 使用介面

### 側邊欄視圖

外掛提供專用的側邊欄視圖：

1. 點擊左側邊欄的書本圖示 📚
2. 或使用命令面板：`Ctrl/Cmd + P` → 輸入 "Novel Maker: Open view"

**視圖內容：**

```
📚 小說專案管理

專案資訊
  名稱：魔法學院編年史
  世界：艾瑟利亞大陸
  章節：15

快速操作
  [生成新章節]
  [生成新角色]
  [匯出小說]
  [專案掃描]

最近章節
  ✏️ 第15章：決戰前夕
  ✅ 第14章：秘密揭曉
  ✅ 第13章：陰謀浮現
```

### 命令面板

按 `Ctrl/Cmd + P` 開啟命令面板，輸入以下命令：

| 命令 | 功能 |
|------|------|
| `Novel Maker: Generate next chapter` | 生成下一章節 |
| `Novel Maker: Generate character` | 生成新角色 |
| `Novel Maker: Regenerate current` | 重新生成當前檔案 |
| `Novel Maker: Generate character image` | 生成角色圖片 |
| `Novel Maker: Export novel` | 匯出小說 |
| `Novel Maker: Scan project` | 掃描專案結構 |
| `Novel Maker: Open settings` | 開啟外掛設定 |

### 右鍵選單

在檔案瀏覽器中右鍵點擊檔案：

**Story/ 目錄中的章節檔案：**
- 🔄 重新生成此章節
- 📊 章節統計
- 🔗 顯示相關角色

**Character/ 目錄中的角色檔案：**
- 🔄 重新生成角色描述
- 🎨 生成角色圖片
- 📋 顯示出場章節

### 編輯器工具列

在編輯 Story/ 或 Character/ 中的檔案時，工具列會顯示額外按鈕：

```
[📝 編輯] [👁️ 預覽] [🔄 重新生成] [📊 統計]
```

## 外掛設定

在 Settings → Novel Maker 中配置：

### 基本設定

```yaml
API Configuration
  ├─ OpenAI API Key: [your-api-key]
  ├─ Model: gpt-4o
  ├─ Image Model: dall-e-3
  └─ API Endpoint: https://api.openai.com/v1

Behavior
  ├─ Auto-save after generation: ✓
  ├─ Show notifications: ✓
  └─ Confirm before regenerate: ✓

Display
  ├─ Show sidebar by default: ✓
  ├─ Theme: Auto (follows Obsidian)
  └─ Font size: Medium
```

### 進階設定

```yaml
Generation
  ├─ Default prompt template: [選擇範本]
  ├─ Context chapters: 3
  ├─ Max tokens: 4000
  └─ Temperature: 0.7

File Management
  ├─ Auto-number chapters: ✓
  ├─ Chapter file pattern: {number}_{slug}.md
  └─ Character file pattern: {name}.md

Integration
  ├─ Enable CLI sync: ✓
  └─ Watch for external changes: ✓
```

## 工作流程範例

### 場景 1：開始新章節

1. 在側邊欄點擊 **[生成新章節]**
2. 輸入章節標題：「第16章：最終試煉」
3. （可選）輸入自訂提示
4. 點擊 **生成**
5. 等待 AI 生成內容
6. 在編輯器中檢視和修改

### 場景 2：修改現有章節

1. 開啟章節檔案：`Story/015_ch15.md`
2. 修改 frontmatter 中的 `prompt` 欄位
3. 點擊工具列的 **🔄 重新生成** 按鈕
4. 確認操作
5. 查看更新後的內容

### 場景 3：建立角色卡

1. 使用命令面板：`Novel Maker: Generate character`
2. 輸入角色名稱：「莉莉安院長」
3. 輸入角色描述提示
4. 生成完成後，點擊 **🎨 生成圖片**
5. 等待圖片生成並嵌入

### 場景 4：專案總覽

1. 開啟側邊欄視圖
2. 查看專案統計：
   - 總章節數
   - 總字數
   - 角色數量
   - 最近更新
3. 點擊章節快速導航

## 鍵盤快捷鍵

可在 Obsidian 的 Hotkeys 設定中自訂：

| 功能 | 建議快捷鍵 |
|------|------------|
| 生成下一章節 | `Ctrl/Cmd + Shift + N` |
| 重新生成當前 | `Ctrl/Cmd + Shift + R` |
| 開啟側邊欄 | `Ctrl/Cmd + Shift + M` |
| 專案掃描 | `Ctrl/Cmd + Shift + S` |
| 生成角色 | `Ctrl/Cmd + Shift + C` |

設定步驟：
1. Settings → Hotkeys
2. 搜尋 "Novel Maker"
3. 為各命令設定快捷鍵

## 範本系統

### 使用內建範本

外掛包含預設範本：

1. **章節範本**：`~/.novelmaker/templates/chapter_prompt.tmpl`
2. **角色範本**：`~/.novelmaker/templates/character_prompt.tmpl`

在外掛設定中選擇要使用的範本。

### 建立自訂範本

1. 在 vault 中建立 `Templates/` 目錄
2. 建立範本檔案，例如 `my-chapter-template.tmpl`
3. 在外掛設定中選擇自訂範本

**範本語法範例：**

```go
你是一位{{.Genre}}小說作家。

專案：{{.ProjectName}}
世界觀：{{.World}}

{{if .PrevChapter}}
前一章摘要：
{{.PrevChapter.Summary}}
{{end}}

請為以下章節生成內容：
標題：{{.Title}}

{{if .Prompt}}
特別要求：{{.Prompt}}
{{end}}
```

## 整合功能

### 與 CLI 同步

外掛與 CLI 工具共用相同的檔案結構：

- 在 CLI 中生成的內容會自動顯示在 Obsidian
- 在 Obsidian 中的修改會反映到 CLI 操作
- Frontmatter 格式完全相容

### Dataview 整合

使用 Dataview 外掛建立動態視圖：

**章節列表：**

```dataview
TABLE 
  title as "標題",
  index as "順序",
  word_count as "字數"
FROM "Story"
SORT index ASC
```

**角色列表：**

```dataview
TABLE 
  name as "名稱",
  main as "主要角色",
  appearances as "出場次數"
FROM "Character"
WHERE name
```

**進度追蹤：**

```dataview
TABLE
  chapters as "已完成章節",
  target_chapters as "目標章節",
  round((chapters / target_chapters) * 100) + "%" as "進度"
FROM "Config"
```

### Graph View 視覺化

Obsidian 的圖形視圖可以顯示：

- 角色之間的關係（通過雙向連結）
- 章節之間的依賴
- 世界觀元素的關聯

在圖形視圖中：
- 🟦 藍色節點 = 章節
- 🟩 綠色節點 = 角色
- 🟨 黃色節點 = 世界觀

## 疑難排解

### 外掛無法載入

**症狀：**
- 外掛列表中沒有顯示
- 啟用後沒有反應

**解決方案：**

1. 檢查檔案結構：
   ```bash
   ls -la .obsidian/plugins/obsidian-novelmaker/
   # 應該看到：main.js, manifest.json, styles.css
   ```

2. 查看 Obsidian 控制台：
   - 按 `Ctrl/Cmd + Shift + I` 開啟開發者工具
   - 查看 Console 標籤中的錯誤訊息

3. 重新安裝外掛：
   ```bash
   rm -rf .obsidian/plugins/obsidian-novelmaker
   novelmaker-obs update-plugin
   ```

4. 重啟 Obsidian

### API 連線失敗

**症狀：**
- 生成內容時顯示錯誤
- 「API request failed」訊息

**解決方案：**

1. 檢查 API 金鑰設定：
   - Settings → Novel Maker → API Configuration
   - 確認 API Key 正確

2. 測試 CLI 工具：
   ```bash
   novelmaker-obs config-check
   ```

3. 檢查網路連線

4. 查看 API 用量限制（OpenAI Dashboard）

### 生成內容為空

**症狀：**
- 生成完成但檔案內容為空
- 只有 frontmatter 沒有正文

**解決方案：**

1. 檢查專案結構：
   ```bash
   novelmaker-obs scan
   ```

2. 確認 Config/project.md 存在且正確

3. 增加 API 超時時間：
   - Settings → Novel Maker → Advanced
   - Timeout: 300 seconds

4. 查看詳細錯誤日誌：
   - 開發者工具 → Console

### 外掛版本不相容

**症狀：**
- 外掛功能異常
- 與 CLI 工具版本不匹配

**解決方案：**

1. 更新 CLI 工具：
   ```bash
   go install github.com/voilelab/gonovelmaker/cmd/novelmaker-obs@latest
   ```

2. 更新外掛：
   ```bash
   cd /path/to/vault
   novelmaker-obs update-plugin
   ```

3. 重啟 Obsidian

## 最佳實踐

### 工作區配置

建議的工作區佈局：

```
┌─────────────────────────────────────────┐
│  [檔案]  [搜尋]  [📚]                    │
├──────────┬──────────────────────┬───────┤
│          │                      │       │
│  檔案    │    編輯器視窗        │  大綱 │
│  樹狀圖  │                      │       │
│          │                      │  標籤 │
│  Novel   │                      │       │
│  Maker   │                      │  統計 │
│  側邊欄  │                      │       │
│          │                      │       │
└──────────┴──────────────────────┴───────┘
```

### 寫作流程

1. **規劃階段**
   - 在 Graph View 中規劃章節結構
   - 使用空白章節建立大綱
   - 標註角色關係

2. **創作階段**
   - 使用側邊欄快速生成
   - 在編輯器中潤飾內容
   - 利用 Dataview 追蹤進度

3. **修改階段**
   - 使用重新生成功能優化章節
   - 檢視角色一致性
   - 匯出並審閱全文

### 效能優化

- 定期清理未使用的圖片
- 使用標籤組織內容
- 啟用自動儲存
- 定期備份 vault

## 下一步

- 🏗️ [世界書結構](../worldbook/schema.md) - 了解檔案組織
- 💡 [CLI 使用範例](../cli/examples.md) - 學習 CLI 整合
- 🏛️ [架構設計](../architecture.md) - 理解外掛實作
