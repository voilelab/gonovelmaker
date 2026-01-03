# Obsidian NovelMaker Plugin

Obsidian plugin for the NovelMaker project.

## 開發流程

### 前置需求

- Node.js (for esbuild)

### 安裝依賴

```bash
cd cmd/novelmaker-obs/obsidian-novelmaker
npm install
```

### 開發

1. 修改 `src/` 目錄下的源代碼
2. 重新打包: `npm run build` 或 `node build.js`
3. 重新編譯 Go 程式: `go build -o novelmaker-obs ./cmd/novelmaker-obs`
4. 更新 plugin: `./novelmaker-obs update-plugin`

### 打包

```bash
# 打包 JavaScript 代碼
npm run build

# 或使用 shell 腳本
./build.sh
```

打包後的文件會輸出到 `dist/` 目錄。

## 部署

Go 程式會自動將 `dist/` 目錄中的打包文件嵌入到可執行文件中，然後通過 `update-plugin` 命令複製到 vault 的 `.obsidian/plugins/obsidian-novelmaker/` 目錄。
