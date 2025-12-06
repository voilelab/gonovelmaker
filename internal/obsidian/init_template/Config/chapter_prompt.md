## 專案資訊
專案名稱：{{.ProjectName}}
世界設定：{{.World}}

## 世界觀資訊
{{range .Worldbook -}}
- 【{{join .Tags "、"}}】{{.Content}}
{{end}}

## 角色設定
{{range .Characters -}}
{{if .Main}}### [主角] {{.Name}}{{else}}### {{.Name}}{{end}}
{{.Profile}}

{{end}}

## 前面的章節
{{range .PreChapters -}}
### 第 {{.Index}} 章：{{.Title}}
{{.Content}}

{{end -}}

## 寫作任務
請根據以上資訊撰寫下一章節。

目標章節標題：{{.Title}}
{{if .Prompt}}
額外指示：{{.Prompt}}
{{end}}

## 寫作要求
1. 保持角色性格的一致性
2. 延續前面章節的情節發展
3. 注意世界觀設定的細節
4. 使用生動的描寫和對話
5. 字數約 2000-3000 字

