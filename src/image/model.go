package image

// Note that the Create images endpoint seems to differ in its parameters for dall-e-2 and dall-e-3.
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
	Prompt         string  `json:"prompt"`
	Model          *string `json:"model,omitempty"`
	N              *int    `json:"n,omitempty"`
	ResponseFormat *string `json:"response_format,omitempty"`
	Size           *string `json:"size,omitempty"`
	User           string  `json:"user,omitempty"`
}

type CreateDalle3ImageBody struct {
	Prompt         string  `json:"prompt"`
	Model          *string `json:"model,omitempty"`
	N              *int    `json:"n,omitempty"`
	Quality        *string `json:"quality,omitempty"`
	ResponseFormat *string `json:"response_format,omitempty"`
	Size           *string `json:"size,omitempty"`
	Style          *string `json:"style,omitempty"`
	User           string  `json:"user,omitempty"`
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

func GetDefaultCreateImageBody() CreateImageBody {
	model := new(string)
	*model = "dall-e-2"

	n := new(int)
	*n = 1

	rFormat := new(string)
	*rFormat = "url"

	size := new(string)
	*size = "256x256"

	return CreateImageBody{
		Prompt:         "",
		Model:          model,
		N:              n,
		ResponseFormat: rFormat,
		Size:           size,
		User:           "go-gpt-cli",
	}
}

func GetDefaultDalle3Body() CreateDalle3ImageBody {
	model := new(string)
	*model = "dall-e-3"

	n := new(int)
	*n = 1

	quality := new(string)
	*quality = "hd"

	rFormat := new(string)
	*rFormat = "url"

	size := new(string)
	*size = "1792x1024"

	style := new(string)
	*style = "vivid"

	return CreateDalle3ImageBody{
		Prompt:         "",
		Model:          model,
		N:              n,
		Quality:        quality,
		ResponseFormat: rFormat,
		Size:           size,
		Style:          style,
		User:           "go-gpt-cli",
	}
}
