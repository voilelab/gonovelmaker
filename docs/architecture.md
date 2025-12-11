# 架構設計

本頁面深入探討 gonovelmaker 的技術架構、程式碼組織和設計決策。

## 專案概覽

gonovelmaker 是一個用 Go 語言開發的命令列工具,旨在透過 AI 協助管理和生成小說內容。

### 核心理念

1. **簡潔性**：CLI 優先,易於使用和整合
2. **模組化**：清晰的職責分離
3. **可擴展性**：支援多種 LLM 後端
4. **可靠性**：強大的錯誤處理和資料驗證

## 技術棧

### 主要依賴

```go
require (
    github.com/openai/openai-go/v3      // OpenAI API 客戶端
    github.com/spf13/cobra              // CLI 框架
    github.com/pelletier/go-toml/v2     // TOML 配置解析
    gopkg.in/yaml.v3                    // YAML frontmatter 解析
)
```

### Go 版本

- **最低要求**：Go 1.25.4
- **推薦版本**：Go 1.25.x 最新版

## 專案結構

```
gonovelmaker/
├── cmd/novelmaker-obs/              # 主程式入口
│   ├── main.go                      # CLI 應用程式
│   ├── cmd_*.go                     # 各個命令的實作
│   ├── llmbackend_maker.go          # LLM 後端工廠
│   └── obsidian-novelmaker/         # Obsidian 外掛
│       ├── main.js
│       ├── manifest.json
│       └── styles.css
├── internal/                        # 內部套件
│   ├── config/                      # 配置管理
│   │   ├── config.go                # 配置載入與驗證
│   │   ├── example_config.toml      # 範例配置
│   │   └── templates/               # 內建提示詞範本
│   ├── llmbackend/                  # LLM 後端抽象
│   │   ├── interface.go             # 後端介面定義
│   │   ├── openai.go                # OpenAI 實作
│   │   └── dummy.go                 # 測試用假後端
│   ├── nmutil/                      # 工具函式庫
│   │   ├── frontmatter.go           # YAML frontmatter 處理
│   │   ├── slugify.go               # 檔案名生成
│   │   └── copyfs.go                # 檔案系統操作
│   └── obsidian/                    # Obsidian vault 處理
│       ├── vault.go                 # Vault 資料結構
│       └── init_template/           # 專案初始化範本
├── novelmaker/                      # 核心生成邏輯
│   ├── struct.go                    # 資料結構定義
│   └── render.go                    # AI 生成與渲染
├── go.mod                           # Go 模組定義
├── go.sum                           # 依賴鎖定檔
└── README.md                        # 專案說明
```

## 核心模組

### 1. CMD 層（cmd/novelmaker-obs/）

CLI 應用程式的入口點,使用 Cobra 框架組織命令。

#### 主要檔案

**main.go**
```go
// CLI 應用程式的根命令
var rootCmd = &cobra.Command{
    Use:   "novelmaker-obs",
    Short: "管理 Obsidian 小說專案的 CLI 工具",
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

**命令實作模式**

每個命令都有專門的檔案：

- `cmd_init.go` - 初始化專案
- `cmd_scan.go` - 掃描專案
- `cmd_gen_next.go` - 生成下一章節
- `cmd_gen_char.go` - 生成角色
- `cmd_export.go` - 匯出小說

**命令結構範例：**

```go
var genNextCmd = &cobra.Command{
    Use:   "gen-next",
    Short: "生成下一個章節",
    RunE: func(cmd *cobra.Command, args []string) error {
        // 1. 載入配置
        cfg, err := config.Load()
        if err != nil {
            return err
        }
        
        // 2. 建立 LLM 後端
        backend := makeLLMBackend(cfg)
        
        // 3. 載入 vault
        vault, err := obsidian.Load(".")
        if err != nil {
            return err
        }
        
        // 4. 生成章節
        chapter, err := novelmaker.GenerateChapter(backend, vault, title, prompt)
        if err != nil {
            return err
        }
        
        // 5. 儲存結果
        return vault.SaveChapter(chapter)
    },
}
```

### 2. Config 模組（internal/config/）

處理應用程式配置,支援兩種模式：

#### 傳統單一後端模式

```go
type Config struct {
    OpenAIAPIKey   string `toml:"openai_api_key"`
    Model          string `toml:"model"`
    ImageModel     string `toml:"image_model"`
    BaseURL        string `toml:"base_url"`
    Timeout        int    `toml:"timeout"`
}
```

#### 多後端模式

```go
type Config struct {
    UserLLMBackend string                     `toml:"user_llm_backend"`
    LLMBackend     map[string]BackendConfig   `toml:"llm_backend"`
}

type BackendConfig struct {
    Type       string `toml:"type"`
    APIKey     string `toml:"api_key"`
    BaseURL    string `toml:"base_url"`
    Model      string `toml:"model"`
    ImageModel string `toml:"image_model"`
}
```

#### 配置載入流程

```go
func Load() (*Config, error) {
    // 1. 確定配置檔案路徑
    configPath := filepath.Join(os.Getenv("HOME"), ".novelmaker", "config.toml")
    
    // 2. 如果不存在,建立範例配置
    if !fileExists(configPath) {
        if err := createExampleConfig(configPath); err != nil {
            return nil, err
        }
    }
    
    // 3. 讀取並解析 TOML
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }
    
    var cfg Config
    if err := toml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    
    // 4. 驗證配置
    if err := cfg.Validate(); err != nil {
        return nil, err
    }
    
    // 5. 載入環境變數（覆蓋配置）
    cfg.LoadFromEnv()
    
    return &cfg, nil
}
```

### 3. LLM Backend 模組（internal/llmbackend/）

提供統一的 LLM 介面,支援不同的後端實作。

#### 介面定義

```go
// interface.go
type LLMBackend interface {
    // 文字生成
    GenerateText(ctx context.Context, messages []Message) (string, error)
    
    // 圖片生成
    GenerateImage(ctx context.Context, prompt string) (imageURL string, err error)
    
    // 取得模型資訊
    GetModel() string
    GetImageModel() string
}

type Message struct {
    Role    string // "system", "user", "assistant"
    Content string
}
```

#### OpenAI 實作

```go
// openai.go
type OpenAIBackend struct {
    client     *openai.Client
    model      string
    imageModel string
}

func NewOpenAIBackend(apiKey, baseURL, model, imageModel string) *OpenAIBackend {
    client := openai.NewClient(
        option.WithAPIKey(apiKey),
        option.WithBaseURL(baseURL),
    )
    
    return &OpenAIBackend{
        client:     client,
        model:      model,
        imageModel: imageModel,
    }
}

func (b *OpenAIBackend) GenerateText(ctx context.Context, messages []Message) (string, error) {
    // 轉換訊息格式
    chatMessages := make([]openai.ChatCompletionMessageParamUnion, len(messages))
    for i, msg := range messages {
        chatMessages[i] = openai.UserMessage(msg.Content)
    }
    
    // 呼叫 API
    resp, err := b.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
        Model:    openai.F(b.model),
        Messages: openai.F(chatMessages),
    })
    if err != nil {
        return "", err
    }
    
    return resp.Choices[0].Message.Content, nil
}
```

#### 擴展性

新增其他 LLM 後端只需實作 `LLMBackend` 介面：

```go
// anthropic.go（範例）
type AnthropicBackend struct {
    apiKey string
    model  string
}

func (b *AnthropicBackend) GenerateText(ctx context.Context, messages []Message) (string, error) {
    // 實作 Claude API 呼叫
    // ...
}
```

### 4. Obsidian 模組（internal/obsidian/）

處理 Obsidian vault 的檔案結構和資料載入。

#### Vault 結構

```go
// vault.go
type Vault struct {
    RootPath  string
    Project   *Project
    World     []*Worldbook
    Character []*Character
    Story     []*Chapter
}

type Project struct {
    Name string `yaml:"name"`
    // 其他欄位...
}

type Chapter struct {
    Index     int    `yaml:"index"`
    Title     string `yaml:"title"`
    Prompt    string `yaml:"prompt"`
    Content   string
    FilePath  string
}
```

#### 載入流程

```go
func Load(rootPath string) (*Vault, error) {
    vault := &Vault{RootPath: rootPath}
    
    // 1. 載入專案配置
    project, err := loadProject(filepath.Join(rootPath, "Config", "project.md"))
    if err != nil {
        return nil, err
    }
    vault.Project = project
    
    // 2. 掃描 World/ 目錄
    worldEntries, err := loadWorldEntries(filepath.Join(rootPath, "World"))
    if err != nil {
        return nil, err
    }
    vault.World = worldEntries
    
    // 3. 掃描 Character/ 目錄
    characters, err := loadCharacters(filepath.Join(rootPath, "Character"))
    if err != nil {
        return nil, err
    }
    vault.Character = characters
    
    // 4. 掃描 Story/ 目錄並排序
    chapters, err := loadChapters(filepath.Join(rootPath, "Story"))
    if err != nil {
        return nil, err
    }
    sort.Slice(chapters, func(i, j int) bool {
        return chapters[i].Index < chapters[j].Index
    })
    vault.Story = chapters
    
    return vault, nil
}
```

#### Frontmatter 解析

使用 `nmutil` 套件處理 YAML frontmatter：

```go
// internal/nmutil/frontmatter.go
func Parse(content []byte) (frontmatter map[string]interface{}, body string, err error) {
    // 檢查是否有 frontmatter
    if !bytes.HasPrefix(content, []byte("---\n")) {
        return nil, string(content), nil
    }
    
    // 分割 frontmatter 和正文
    parts := bytes.SplitN(content[4:], []byte("\n---\n"), 2)
    if len(parts) != 2 {
        return nil, "", fmt.Errorf("invalid frontmatter format")
    }
    
    // 解析 YAML
    var fm map[string]interface{}
    if err := yaml.Unmarshal(parts[0], &fm); err != nil {
        return nil, "", err
    }
    
    return fm, string(parts[1]), nil
}
```

### 5. Novelmaker 模組（novelmaker/）

核心 AI 生成邏輯。

#### 資料結構

```go
// struct.go
type ChapterPrompt struct {
    System            string
    AssistantTemplate *template.Template
}

type Project struct {
    Name      string
    World     string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Character struct {
    Name    string
    Main    bool
    Prompt  string
    Profile string
}

type Chapter struct {
    Index   int
    Title   string
    Prompt  string
    Content string
}
```

#### 生成流程

```go
// render.go
func GenerateChapter(
    backend llmbackend.LLMBackend,
    vault *obsidian.Vault,
    title string,
    prompt string,
) (*Chapter, error) {
    // 1. 準備上下文
    ctx := prepareContext(vault, title, prompt)
    
    // 2. 載入提示詞範本
    template, err := loadChapterTemplate()
    if err != nil {
        return nil, err
    }
    
    // 3. 渲染系統提示詞
    systemPrompt, err := renderTemplate(template, ctx)
    if err != nil {
        return nil, err
    }
    
    // 4. 呼叫 LLM
    messages := []llmbackend.Message{
        {Role: "system", Content: systemPrompt},
        {Role: "user", Content: prompt},
    }
    
    content, err := backend.GenerateText(context.Background(), messages)
    if err != nil {
        return nil, err
    }
    
    // 5. 建立章節物件
    chapter := &Chapter{
        Index:   vault.GetNextChapterIndex(),
        Title:   title,
        Prompt:  prompt,
        Content: content,
    }
    
    return chapter, nil
}
```

#### 提示詞範本

使用 Go template 引擎：

```go
// 範本檔案：~/.novelmaker/templates/chapter_prompt.tmpl
const chapterTemplate = `
你是一位{{.Genre}}小說作家。

專案：{{.ProjectName}}
世界觀：{{.World}}

{{if .Characters}}
主要角色：
{{range .Characters}}
- {{.Name}}: {{.Profile}}
{{end}}
{{end}}

{{if .PrevChapters}}
前面章節：
{{range .PrevChapters}}
## {{.Title}}
{{.Summary}}
{{end}}
{{end}}

請生成以下章節：
標題：{{.Title}}
{{if .Prompt}}
特別要求：{{.Prompt}}
{{end}}
`

func renderTemplate(tmpl *template.Template, data interface{}) (string, error) {
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }
    return buf.String(), nil
}
```

## 設計模式

### 1. 工廠模式

用於建立 LLM 後端：

```go
// llmbackend_maker.go
func makeLLMBackend(cfg *config.Config) llmbackend.LLMBackend {
    // 多後端模式
    if cfg.UserLLMBackend != "" {
        backendCfg := cfg.LLMBackend[cfg.UserLLMBackend]
        switch backendCfg.Type {
        case "openai":
            return llmbackend.NewOpenAIBackend(
                backendCfg.APIKey,
                backendCfg.BaseURL,
                backendCfg.Model,
                backendCfg.ImageModel,
            )
        case "anthropic":
            return llmbackend.NewAnthropicBackend(backendCfg)
        }
    }
    
    // 傳統模式
    return llmbackend.NewOpenAIBackend(
        cfg.OpenAIAPIKey,
        cfg.BaseURL,
        cfg.Model,
        cfg.ImageModel,
    )
}
```

### 2. 策略模式

LLM Backend 介面允許切換不同的 AI 服務：

```go
type Generator struct {
    backend llmbackend.LLMBackend
}

func (g *Generator) Generate(prompt string) (string, error) {
    // 不管底層是 OpenAI、Claude 或其他服務
    // 使用統一的介面
    return g.backend.GenerateText(context.Background(), messages)
}
```

### 3. 範本方法模式

命令執行的通用流程：

```go
func executeGenerationCommand(
    loadVault func() (*obsidian.Vault, error),
    generate func(*obsidian.Vault) (interface{}, error),
    save func(interface{}) error,
) error {
    // 1. 載入 vault
    vault, err := loadVault()
    if err != nil {
        return err
    }
    
    // 2. 生成內容
    result, err := generate(vault)
    if err != nil {
        return err
    }
    
    // 3. 儲存結果
    return save(result)
}
```

## 錯誤處理

### 錯誤類型

```go
var (
    ErrConfigNotFound    = errors.New("config file not found")
    ErrInvalidAPIKey     = errors.New("invalid API key")
    ErrVaultNotFound     = errors.New("vault not found")
    ErrInvalidFrontmatter = errors.New("invalid frontmatter format")
    ErrAPICallFailed     = errors.New("API call failed")
)
```

### 錯誤包裝

```go
if err := vault.Load(); err != nil {
    return fmt.Errorf("failed to load vault: %w", err)
}
```

### 使用者友善的錯誤訊息

```go
func handleError(err error) {
    switch {
    case errors.Is(err, ErrConfigNotFound):
        fmt.Println("❌ 找不到設定檔")
        fmt.Println("提示：執行 'novelmaker-obs config-check' 建立設定")
    case errors.Is(err, ErrVaultNotFound):
        fmt.Println("❌ 當前目錄不是有效的 vault")
        fmt.Println("提示：執行 'novelmaker-obs init' 初始化專案")
    default:
        fmt.Printf("❌ 錯誤：%v\n", err)
    }
}
```

## 測試策略

### 單元測試

```go
// internal/nmutil/frontmatter_test.go
func TestParseFrontmatter(t *testing.T) {
    content := []byte("---\ntitle: Test\n---\nBody content")
    
    fm, body, err := Parse(content)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if fm["title"] != "Test" {
        t.Errorf("expected title=Test, got %v", fm["title"])
    }
    
    if body != "Body content" {
        t.Errorf("expected body='Body content', got %q", body)
    }
}
```

### 整合測試

```go
// cmd/novelmaker-obs/cmd_init_test.go
func TestInitCommand(t *testing.T) {
    // 建立臨時目錄
    tmpDir := t.TempDir()
    
    // 執行 init 命令
    if err := initVault(tmpDir); err != nil {
        t.Fatalf("init failed: %v", err)
    }
    
    // 驗證檔案結構
    requiredFiles := []string{
        "Config/project.md",
        "World/001_world_sample.md",
        "Character/character_sample.md",
        "Story/001_prologue.md",
    }
    
    for _, file := range requiredFiles {
        path := filepath.Join(tmpDir, file)
        if _, err := os.Stat(path); os.IsNotExist(err) {
            t.Errorf("required file not created: %s", file)
        }
    }
}
```

### 測試 Dummy Backend

```go
// internal/llmbackend/dummy.go
type DummyBackend struct {
    responses []string
    callCount int
}

func (b *DummyBackend) GenerateText(ctx context.Context, messages []Message) (string, error) {
    if b.callCount >= len(b.responses) {
        return "", errors.New("no more responses")
    }
    
    response := b.responses[b.callCount]
    b.callCount++
    return response, nil
}
```

## 效能考量

### 1. 延遲載入

只在需要時載入大型檔案：

```go
type Vault struct {
    rootPath string
    project  *Project
    chapters []*Chapter // 只儲存元數據
}

func (v *Vault) GetChapterContent(index int) (string, error) {
    // 延遲載入實際內容
    return loadChapterContent(v.chapters[index].FilePath)
}
```

### 2. 快取

快取常用資料：

```go
var templateCache = make(map[string]*template.Template)

func getTemplate(name string) (*template.Template, error) {
    if tmpl, ok := templateCache[name]; ok {
        return tmpl, nil
    }
    
    tmpl, err := loadTemplate(name)
    if err != nil {
        return nil, err
    }
    
    templateCache[name] = tmpl
    return tmpl, nil
}
```

### 3. 並行處理

批次操作時使用 goroutine：

```go
func loadChapters(dir string) ([]*Chapter, error) {
    files, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }
    
    chapters := make([]*Chapter, len(files))
    errors := make([]error, len(files))
    
    var wg sync.WaitGroup
    for i, file := range files {
        wg.Add(1)
        go func(i int, file os.DirEntry) {
            defer wg.Done()
            chapters[i], errors[i] = loadChapter(filepath.Join(dir, file.Name()))
        }(i, file)
    }
    
    wg.Wait()
    
    // 檢查錯誤
    for _, err := range errors {
        if err != nil {
            return nil, err
        }
    }
    
    return chapters, nil
}
```

## 未來擴展

### 計劃功能

1. **更多 LLM 後端**
   - Anthropic Claude
   - Google Gemini
   - Local LLMs (Ollama)

2. **進階生成選項**
   - Temperature 控制
   - Top-p sampling
   - 停止序列

3. **匯出格式**
   - EPUB
   - PDF
   - HTML

4. **協作功能**
   - Git 整合
   - 版本對比
   - 合併工具

### 貢獻指南

歡迎貢獻！請遵循：

1. Fork 專案
2. 建立功能分支
3. 撰寫測試
4. 確保通過 `go test ./...`
5. 提交 Pull Request

## 開發環境設定

```bash
# 1. Clone 儲存庫
git clone https://github.com/voilelab/gonovelmaker.git
cd gonovelmaker

# 2. 安裝依賴
go mod download

# 3. 執行測試
go test ./...

# 4. 建置
go build ./cmd/novelmaker-obs

# 5. 本地安裝
go install ./cmd/novelmaker-obs
```

## 相關資源

- [Cobra CLI 文檔](https://cobra.dev/)
- [OpenAI Go SDK](https://github.com/openai/openai-go)
- [Go Template 語法](https://pkg.go.dev/text/template)
- [TOML 規範](https://toml.io/)

---

若您對架構有任何疑問或建議,歡迎在 GitHub 上開 Issue 討論！
