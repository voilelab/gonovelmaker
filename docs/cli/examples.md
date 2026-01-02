# 使用範例

本頁面提供實際的使用案例，展示如何在不同情境下使用 gonovelmaker。

## 後端管理案例

### 案例 1：設定多個 LLM 後端

配置不同的後端以便快速切換：

```bash
# 新增 OpenAI 後端（主要用於文字生成）
novelmaker-obs backend add openai \
  --type openai \
  --api_key "sk-xxx" \
  --model "gpt-4o" \
  --image_model "dall-e-3"

# 新增 OpenRouter 後端（用於測試 Claude）
novelmaker-obs backend add claude \
  --type openrouter \
  --api_key "sk-or-xxx" \
  --base_url "https://openrouter.ai/api/v1" \
  --model "anthropic/claude-3.5-sonnet"

# 查看所有後端
novelmaker-obs backend list
```

**輸出：**

```
Configured backends:

• openai (default)
  Type:        openai
  API Key:     ...xxx
  Model:       gpt-4o
  Image Model: dall-e-3

• claude
  Type:        openrouter
  API Key:     ...xxx
  Base URL:    https://openrouter.ai/api/v1
  Model:       anthropic/claude-3.5-sonnet

• local
  Type:        openai
  API Key:     ...test
  Base URL:    http://localhost:1234/v1
  Model:       llama-3.1-8b
```

---

### 案例 2：切換後端生成內容

為不同類型的內容使用不同的後端：

```bash
# 使用 GPT-4o 生成日常對話章節
novelmaker-obs backend use openai
novelmaker-obs gen-next --title "第三章：日常訓練"

# 切換到 Claude 生成需要深度思考的章節
novelmaker-obs backend use claude
novelmaker-obs gen-next --title "第四章：哲學辯論"

# 切換回 OpenAI 生成動作場景
novelmaker-obs backend use openai
novelmaker-obs gen-next --title "第五章：決戰時刻"
```

---

### 案例 3：測試後端連線

在開始生成前確認後端配置正確：

```bash
# 檢查 OpenAI 後端
novelmaker-obs backend check openai

# 檢查 Claude 後端
novelmaker-obs backend check claude

# 使用 JSON 格式檢查（便於腳本處理）
novelmaker-obs backend check openai --json
```

**成功輸出：**

```
Checking backend 'openai'...
  Type:     openai
  Model:    gpt-4o

Sending test request... ✅ SUCCESS

Response time: 892ms
Response: OK

Token usage:
  Input tokens:  15
  Output tokens: 2
  Total tokens:  17
```

---

### 案例 4：更新後端配置

修改現有後端的設定：

```bash
# 升級模型版本
novelmaker-obs backend add openai \
  --model "gpt-4o-mini"

# 調整超時時間
novelmaker-obs backend add claude \
  --timeout 180

# 更換 API 金鑰
novelmaker-obs backend add openai \
  --api_key "sk-new-key"

# 查看更新後的配置
novelmaker-obs backend list
```

---

### 案例 5：管理後端生命週期

```bash
# 新增測試後端
novelmaker-obs backend add test \
  --type openai \
  --api_key "test" \
  --model "test-model"

# 檢查是否正常
novelmaker-obs backend check test

# 測試完成後移除
novelmaker-obs backend remove test

# 確認已移除
novelmaker-obs backend list
```

---

### 案例 6：成本優化策略

使用不同後端節省成本：

```bash
# 配置便宜的模型用於草稿
novelmaker-obs backend add draft \
  --type openai \
  --api_key "$OPENAI_KEY" \
  --model "gpt-4o-mini"

# 配置高品質模型用於最終版本
novelmaker-obs backend add final \
  --type openai \
  --api_key "$OPENAI_KEY" \
  --model "gpt-4o"

# 先用草稿模型生成
novelmaker-obs backend use draft
for i in {1..10}; do
  novelmaker-obs gen-next --title "第${i}章（草稿）"
done

# 滿意後用高品質模型重新生成重要章節
novelmaker-obs backend use final
novelmaker-obs gen-curr --filepath "Story/001_ch1.md"
novelmaker-obs gen-curr --filepath "Story/010_ch10.md"
```

---

## 快速開始：建立第一部小說

### 步驟 1：初始化專案

```bash
# 建立專案目錄
mkdir my-fantasy-novel
cd my-fantasy-novel

# 初始化結構
novelmaker-obs init
```

### 步驟 2：編輯專案設定

編輯 `Config/project.md`：

```yaml
---
id: fantasy-001
name: 魔法學院編年史
world: 艾瑟利亞大陸
system_prompt: 這是一個充滿魔法的世界，魔法學院培養著未來的魔法師...
---
```

### 步驟 3：建立世界觀

編輯 `World/001_world_sample.md`：

```markdown
---
tags: [world, magic-system, geography]
---

# 艾瑟利亞大陸

艾瑟利亞是一個魔法與科技交織的世界...

## 魔法系統

魔法分為五大元素：火、水、風、地、光...
```

### 步驟 4：生成主要角色

```bash
# 生成主角
novelmaker-obs gen-char \
  --name "艾莉絲" \
  --prompt "16歲的魔法學院新生，擁有罕見的雙元素天賦（火與風），性格活潑好奇，來自鄉村，對魔法充滿熱情"

# 生成導師角色
novelmaker-obs gen-char \
  --name "梅林教授" \
  --prompt "50歲的資深魔法教授，專精空間魔法，外表嚴肅但內心溫暖，曾是傳奇魔法師"

# 生成好友
novelmaker-obs gen-char \
  --name "托比" \
  --prompt "同班同學，擅長水系魔法，性格溫和，家族世代都是治療師"
```

### 步驟 5：生成章節

```bash
# 第一章
novelmaker-obs gen-next \
  --title "第一章：入學典禮" \
  --prompt "描寫艾莉絲第一次踏入魔法學院的場景，展現學院的宏偉和魔法的奇妙"

# 第二章
novelmaker-obs gen-next \
  --title "第二章：元素測試" \
  --prompt "新生進行元素親和力測試，艾莉絲展現出雙元素天賦，引起轟動"

# 第三章
novelmaker-obs gen-next \
  --title "第三章：初次課堂" \
  --prompt "梅林教授的第一堂課，艾莉絲認識托比，兩人成為朋友"
```

### 步驟 6：檢視進度

```bash
novelmaker-obs scan
```

輸出：

```
📚 專案：魔法學院編年史
🌍 世界觀條目：1
👥 角色：3
  - 艾莉絲 (main)
  - 梅林教授
  - 托比
📖 章節：3
  - [001] 第一章：入學典禮
  - [002] 第二章：元素測試
  - [003] 第三章：初次課堂
```

---

## 進階案例：重新生成與修正

### 案例：修改章節內容

如果對生成的章節不滿意：

```bash
# 方法一：編輯 frontmatter 中的 prompt，然後重新生成
# 編輯 Story/001_ch1.md，修改 prompt 欄位
novelmaker-obs gen-curr --filepath "Story/001_ch1.md"

# 方法二：建立新版本
novelmaker-obs gen-next \
  --title "第一章：入學典禮（修訂版）" \
  --prompt "更詳細地描寫學院建築，增加神秘氣氛，減少對話比例"
```

### 案例：調整生成風格

```bash
# 使用不同模型生成
novelmaker-obs gen-next \
  --title "第四章：秘密圖書館" \
  --model "gpt-4o" \
  --prompt "較為嚴肅的語氣，增加懸疑感"

# 或使用不同後端
# 先在 config.toml 中配置 Claude
novelmaker-obs gen-next \
  --title "第五章：禁忌魔法" \
  --base-url "https://openrouter.ai/api/v1" \
  --model "anthropic/claude-3.5-sonnet"
```

### 案例：使用更多上下文

```bash
# 生成後期章節時使用更多前文
novelmaker-obs gen-curr \
  --filepath "Story/010_ch10.md" \
  --prev-chapters 8
```

---

## 角色管理案例

### 案例：生成角色圖片

```bash
# 基本生成
novelmaker-obs gen-char-img --name "艾莉絲"

# 自訂風格
novelmaker-obs gen-char-img \
  --name "艾莉絲" \
  --prompt "Anime style portrait of a 16-year-old girl with long red hair, bright green eyes, wearing a blue magic academy uniform with a phoenix emblem, confident smile, magical aura, detailed illustration"

# 生成所有主要角色的圖片
for char in "艾莉絲" "梅林教授" "托比"; do
  novelmaker-obs gen-char-img --name "$char"
  sleep 10  # DALL-E 有速率限制
done
```

### 案例：批次生成角色

建立角色清單檔案 `characters.txt`：

```
艾莉絲|主角，魔法學院新生，雙元素天賦
托比|主角好友，水系治療師
梅林教授|空間魔法教授，嚴肅但溫暖
莉莉安|學院院長，強大的光系魔法師
達克斯|反派，禁忌魔法研究者
```

使用腳本批次生成：

```bash
#!/bin/bash
while IFS='|' read -r name prompt; do
  echo "生成角色：$name"
  novelmaker-obs gen-char --name "$name" --prompt "$prompt"
  sleep 2
done < characters.txt
```

---

## 專案管理案例

### 案例：多專案管理

使用不同目錄管理多個小說專案：

```bash
# 專案結構
~/novels/
  ├── fantasy-novel/
  ├── sci-fi-novel/
  └── romance-novel/

# 切換專案
cd ~/novels/fantasy-novel
novelmaker-obs scan

cd ~/novels/sci-fi-novel
novelmaker-obs scan
```

### 案例：版本控制

使用 Git 管理專案版本：

```bash
# 初始化 Git
git init
git add .
git commit -m "Initial project structure"

# 生成新章節後
novelmaker-obs gen-next --title "新章節"
git add Story/
git commit -m "Add new chapter"

# 查看變更
git diff HEAD~1

# 如果不滿意，回退
git checkout HEAD~1 Story/005_ch5.md
```

### 案例：備份與匯出

```bash
#!/bin/bash
# backup.sh - 自動備份腳本

DATE=$(date +%Y%m%d)
BACKUP_DIR="backups/$DATE"

# 建立備份目錄
mkdir -p "$BACKUP_DIR"

# 複製重要檔案
cp -r Config/ "$BACKUP_DIR/"
cp -r Character/ "$BACKUP_DIR/"
cp -r Story/ "$BACKUP_DIR/"
cp -r World/ "$BACKUP_DIR/"

# 匯出文字版本
novelmaker-obs export \
  --output "$BACKUP_DIR/novel-$DATE.txt" \
  --type txt

echo "✅ 備份完成：$BACKUP_DIR"
```

---

## 自動化工作流程

### 案例：每日創作腳本

```bash
#!/bin/bash
# daily-writing.sh

# 設定
CHAPTER_NUM=$(ls Story/ | wc -l)
NEXT_NUM=$((CHAPTER_NUM + 1))

# 提示使用者
read -p "今天要寫的章節標題：" TITLE
read -p "章節提示（可選）：" PROMPT

# 建立空白章節
novelmaker-obs gen-next-empty \
  --title "第${NEXT_NUM}章：${TITLE}" \
  --prompt "$PROMPT"

# 在 Obsidian 中開啟
CHAPTER_FILE="Story/$(printf '%03d' $NEXT_NUM)_*.md"
open -a "Obsidian" "$CHAPTER_FILE"

echo "✅ 章節已建立，開始創作吧！"
```

### 案例：AI 輔助創作流程

```bash
#!/bin/bash
# ai-assist.sh - 半自動創作流程

# 1. 生成章節大綱
novelmaker-obs gen-next \
  --title "第六章：試煉" \
  --prompt "只生成章節大綱，包含3-5個場景"

# 2. 人工審閱和編輯大綱

# 3. 根據修改後的大綱生成完整內容
novelmaker-obs gen-curr \
  --filepath "Story/006_ch6.md"

# 4. 自動生成摘要
novelmaker-obs scan --json | \
  jq -r '.chapters[-1]' > chapter-summary.json
```

### 案例：批次處理章節

```bash
#!/bin/bash
# batch-process.sh

# 為所有空白章節生成內容
for file in Story/*.md; do
  # 檢查是否為空白章節
  if ! grep -q "^[^#]" "$file"; then
    echo "處理：$file"
    novelmaker-obs gen-curr --filepath "${file#./}"
    sleep 5  # 避免 API 速率限制
  fi
done
```

---

## 整合與匯出案例

### 案例：匯出不同格式

```bash
# 匯出純文字
novelmaker-obs export \
  --output "novel.txt" \
  --type txt

# 使用 pandoc 轉換為其他格式
pandoc novel.txt -o novel.epub
pandoc novel.txt -o novel.pdf
pandoc novel.txt -o novel.docx
```

### 案例：生成網頁版本

```bash
#!/bin/bash
# generate-web.sh

# 匯出 JSON
novelmaker-obs scan --json > data.json

# 使用 jq 和範本生成 HTML
cat > index.html <<EOF
<!DOCTYPE html>
<html>
<head>
  <title>$(jq -r '.project.name' data.json)</title>
</head>
<body>
  <h1>$(jq -r '.project.name' data.json)</h1>
EOF

# 新增章節
jq -r '.chapters[] | "  <h2>\(.title)</h2>\n  <p>\(.content)</p>"' \
  data.json >> index.html

cat >> index.html <<EOF
</body>
</html>
EOF

echo "✅ 網頁版本已生成：index.html"
```

### 案例：統計分析

```bash
# 字數統計
novelmaker-obs scan --json | \
  jq -r '.chapters[].filepath' | \
  xargs wc -w

# 章節數量
novelmaker-obs scan --json | \
  jq '.chapters | length'

# 角色登場次數
for char in $(novelmaker-obs scan --json | jq -r '.characters[].name'); do
  echo "$char:"
  grep -r "$char" Story/ | wc -l
done
```

---

## Obsidian 整合案例

### 案例：更新外掛

```bash
cd /path/to/your/vault

# 更新外掛
novelmaker-obs update-plugin

# 重啟 Obsidian
osascript -e 'quit app "Obsidian"'
sleep 2
open -a Obsidian .
```

### 案例：自訂 Obsidian 工作區

在 `.obsidian/workspace.json` 中設定：

```json
{
  "main": {
    "id": "novel-workspace",
    "type": "split",
    "children": [
      {
        "type": "leaf",
        "state": {
          "type": "markdown",
          "file": "Config/project.md"
        }
      },
      {
        "type": "leaf",
        "state": {
          "type": "graph"
        }
      }
    ]
  }
}
```

---

## 疑難排解案例

### 問題：生成內容品質不佳

**解決方案 1：優化提示詞**

```bash
# 使用更詳細的提示
novelmaker-obs gen-next \
  --title "第七章" \
  --prompt "場景：魔法學院圖書館深夜。艾莉絲發現一本古老的魔法書。
  氛圍：神秘、緊張。
  重點：展現艾莉絲的好奇心和勇氣，為後續劇情埋下伏筆。
  避免：過度對話，冗長描述。
  風格：簡潔明快，注重動作和心理描寫。"
```

**解決方案 2：調整模型**

```bash
# 嘗試不同模型
novelmaker-obs gen-next --title "測試" --model "gpt-4o"
novelmaker-obs gen-next --title "測試" --model "gpt-4o-mini"
```

### 問題：API 超時

```bash
# 增加超時時間
novelmaker-obs gen-next \
  --title "長章節" \
  --timeout 300

# 或減少上下文
novelmaker-obs gen-curr \
  --filepath "Story/010_ch10.md" \
  --prev-chapters 2
```

### 問題：檔案編碼問題

```bash
# 檢查檔案編碼
file -I Story/*.md

# 轉換編碼（如果需要）
iconv -f ISO-8859-1 -t UTF-8 Story/001.md > Story/001_utf8.md
```

---

## 最佳實踐

### 提示詞撰寫技巧

1. **具體明確**
   ```bash
   # ❌ 不好
   --prompt "寫一個有趣的章節"
   
   # ✅ 好
   --prompt "在魔法競技場上，艾莉絲與托比對決，展現雙元素組合技"
   ```

2. **提供結構**
   ```bash
   --prompt "章節結構：
   1. 開場：競技場氣氛描寫
   2. 對決過程：展現魔法特效
   3. 結尾：意外的轉折
   風格：節奏明快，對話精簡"
   ```

3. **引用設定**
   ```bash
   --prompt "根據 World/001_magic_system.md 中的魔法規則，
   描寫艾莉絲學習新咒語的過程"
   ```

### 工作流程建議

1. **先建立框架**
   ```bash
   # 先生成所有空白章節
   for i in {1..20}; do
     novelmaker-obs gen-next-empty --title "第${i}章"
   done
   ```

2. **逐步細化**
   ```bash
   # 填寫每章的 prompt
   # 然後批次生成或逐章生成
   ```

3. **定期備份**
   ```bash
   # 使用 Git 或備份腳本
   git commit -am "Progress: Chapter 5 completed"
   ```

## 下一步

- 🏗️ [世界書結構](../worldbook/schema.md) - 了解檔案格式
- 🔌 [Obsidian 外掛](../plugin/obsidian.md) - 學習外掛使用
