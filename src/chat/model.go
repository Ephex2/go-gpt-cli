package chat

// For both system_prompt files and config files, I ran into default value issues:
// I wanted to use json files but this gives much easier control
var defaultSystemPrompt = Message{
	Role:    "system",
	Content: "You will be acting as an AI assistant. Your goal is to answer the user's requests in a detailed yet concise manner. Here are some important rules for the interaction:\n- Always respond in a neutral and professional tone.\n- If a coding question could use an example, it is ok to be more lenient on the need to be concise.",
}

// Ref: https://platform.openai.com/docs/api-reference/chat/create
type CreateCompletionBody struct {
	Messages         []Message       `json:"messages"`
	Model            string          `json:"model"`
	FrequencyPenalty *float64        `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int  `json:"logit_bias,omitempty"`
	Logprobs         *bool           `json:"logprobs,omitempty"`
	TopLogprobs      *int            `json:"top_logprobs,omitempty"`
	MaxTokens        *int            `json:"max_tokens,omitempty"`
	N                *int            `json:"n,omitempty"`
	PresencePenalty  *float64        `json:"presence_penalty,omitempty"`
	ResponseFormat   *ResponseFormat `json:"response_format,omitempty"`
	Seed             *int            `json:"seed,omitempty"`
	Stop             []string        `json:"stop,omitempty"`
	Stream           *bool           `json:"stream,omitempty"`
	Temperature      *float64        `json:"temperature,omitempty"`
	TopP             *float64        `json:"top_p,omitempty"`
	Tools            []Tool          `json:"tools,omitempty"`
	ToolChoice       *interface{}    `json:"tool_choice,omitempty"` // can be "auto", "none", nil, or a Tool{}
	User             string          `json:"user,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormat struct {
	Type string `json:"type"` // must be "text" or "json_object"
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Parameters  string `json:"parameters"` // json string of function parameters. Varies per function.
}

type CompletionResponse struct {
	Id                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

type Choice struct {
	Index         int     `json:"index"`
	Message       Message `json:"message"`
	Logprobs      bool    `json:"logprobs"`
	FinishSession string  `json:"finish_session"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}


type CreateVisionCompletionBody struct {
	Messages         []VisionMessage `json:"messages"`
	Model            string          `json:"model"`
	FrequencyPenalty *float64        `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int  `json:"logit_bias,omitempty"`
	Logprobs         *bool           `json:"logprobs,omitempty"`
	TopLogprobs      *int            `json:"top_logprobs,omitempty"`
	MaxTokens        *int            `json:"max_tokens,omitempty"`
	N                *int            `json:"n,omitempty"`
	PresencePenalty  *float64        `json:"presence_penalty,omitempty"`
	ResponseFormat   *ResponseFormat `json:"response_format,omitempty"`
	Seed             *int            `json:"seed,omitempty"`
	Stop             []string        `json:"stop,omitempty"`
	Stream           *bool           `json:"stream,omitempty"`
	Temperature      *float64        `json:"temperature,omitempty"`
	TopP             *float64        `json:"top_p,omitempty"`
	Tools            []Tool          `json:"tools,omitempty"`
	ToolChoice       *interface{}    `json:"tool_choice,omitempty"` // can be "auto", "none", nil, or a Tool{}
	User             string          `json:"user,omitempty"`
}

type VisionMessage struct {
	Role    string          `json:"role"`
	Content []VisionContent `json:"content"`
}

type VisionContent struct {
	Type     string   `json:"type"` // text or image_url. image_url can be a literal url to the image or "data:image/png;base64,<<base64 image data>>"
	Text     *string  `json:"text,omitempty"`
	ImageUrl ImageUrl `json:"image_url,omitempty"`
}

type ImageUrl struct {
	Url string `json:"url"`
}

// user default to pre-populate our client's desired default values
// commented out request properties I plan not to use and that cannot have a good default value
// might be a skill issue, but it is difficult in go for a property to not exist if not unmarshaled to.
func GetDefaultBody() CreateCompletionBody {
	n := new(int)
	*n = 1

	b := CreateCompletionBody{
		Model:    "gpt-3.5-turbo",
		Messages: []Message{},
		N:        n,
		ResponseFormat: &ResponseFormat{
			Type: "text",
		},
		User: "go-gpt-cli",
	}

	b.Messages = append(b.Messages, defaultSystemPrompt)
	return b
}

// Vision section
func GetDefaultVisionBody() CreateVisionCompletionBody {
	n := new(int)
	*n = 1

	max := new(int)
	*max = 300

	b := CreateVisionCompletionBody{
		Model:     "gpt-4-vision-preview",
		Messages:  []VisionMessage{},
		MaxTokens: max,
		N:         n,
		ResponseFormat: &ResponseFormat{
			Type: "text",
		},
		User: "go-gpt-cli",
	}

	return b
}
