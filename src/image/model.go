package image

// Note that the Create images endpoint seems to differ in its parameters for dall-e-2 and dall-e-3. Creating a generic implementation that can plug into localAi implementations.

// TODO -- test | for negative prompts

var AllowedQuality = struct {
	Standard string
	Hd       string
}{
	Standard: "standard",
	Hd:       "hd",
}

var AllowedFormat = struct {
	Url     string
	B64Json string
}{
	Url:     "url",
	B64Json: "b64_json",
}

var AllowedStyle = struct {
	Vivid    string
	Standard string
}{
	Vivid:    "vivid",
	Standard: "natural",
}

var AllowedSizeDalle2 = struct {
	S256x256   string
	S512x512   string
	S1024x1024 string
}{
	S256x256:   "256x256",
	S512x512:   "512x512",
	S1024x1024: "1024x1024",
}

var AllowedSizeDalle3 = struct {
	S1024x1024 string
	S1024x1792 string
	S1792x1024 string
}{
	S1024x1024: "1024x1024",
	S1024x1792: "1024x1792",
	S1792x1024: "1792x1024",
}

var CreateImageEditBody = map[string]string{
	"model": "dall-e-2",
	"n":     "1",
	"size":  "Calculated at runtime",
	"user":  "go-gpt-cli",
}

var CreateVariationBody = map[string]string{
	"model": "dall-e-2",
	"n":     "1",
	"size":  "Calculated at runtime",
	"user":  "go-gpt-cli",
}

type CreateImageBody struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model"`
	N              int    `json:"n"`
	ResponseFormat string `json:"response_format,omitempty"`
	Size           string `json:"size"`
	User           string `json:"user"`
}

type CreateDalle3ImageBody struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model"`
	N              int    `json:"n"`
	Quality        string `json:"quality"`
	ResponseFormat string `json:"response_format,omitempty"`
	Size           string `json:"size"`
	Style          string `json:"style"`
	User           string `json:"user"`
}

type CreateImageResponse struct {
	Created int             `json:"created"`
	Data    []ImageResponse `json:"data"`
}

type ImageResponse struct {
	Url           string `json:"url,omitempty"`
	B64Json       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

type ResponseBase64 struct {
	B64Json       string `json:"b64_json"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

func GetDefaultCreateImageBody() CreateImageBody {
	return CreateImageBody{
		Prompt:         "",
		Model:          "dall-e-2",
		N:              1,
		ResponseFormat: AllowedFormat.Url,
		Size:           AllowedSizeDalle2.S256x256,
		User:           "go-gpt-cli",
	}
}

func GetDefaultDalle3Body() CreateDalle3ImageBody {
	return CreateDalle3ImageBody{
		Prompt:         "",
		Model:          "dall-e-3",
		N:              1,
		Quality:        AllowedQuality.Hd,
		ResponseFormat: AllowedFormat.Url,
		Size:           AllowedSizeDalle3.S1792x1024,
		Style:          AllowedStyle.Vivid,
		User:           "go-gpt-cli",
	}
}
