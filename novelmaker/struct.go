package novelmaker

import "time"

type Project struct {
	Name             string    `json:"name"`
	World            string    `json:"world"`
	SystemPrompt     string    `json:"system_prompt"`
	SystemPromptChar string    `json:"system_prompt_char"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Worldbook struct {
	Tags      []string  `json:"tags"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Character struct {
	Name      string    `json:"name"`
	Main      bool      `json:"main"`
	Prompt    string    `json:"prompt"`
	Profile   string    `json:"profile"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Chapter struct {
	Index     int       `json:"index"`
	Title     string    `json:"title"`
	Prompt    string    `json:"prompt"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}
