package embeddings

var AllowedEncodingModels = struct{
    TextEmbedding3Small string
    TextEmbedding3Large string
    TextEmbeddingAda002 string
}{
    TextEmbedding3Small: "text-embedding-3-small",
    TextEmbedding3Large: "text-embedding-3-large",
    TextEmbeddingAda002: "text-embedding-ada-002",
}

var AllowedEmbeddingEncodingFormats = struct{
    Float string
    Base64 string
}{
    Float: "float",
    Base64: "base64",
}

/* Note: the input can technically be a string, array of strings, integer, or array of integers. Simplified for the current case.

See the 'reducing embedding dimensions' section here for more information on Dimensions: https://platform.openai.com/docs/guides/embeddings/use-cases 
*/
type CreateEmbeddingBody struct {
    Input []string `json:"input"`
    Model string `json:"model"`
    EncodingFormat string `json:"encoding_format"`
    Dimensions int `json:"dimensions"`
    User string `json:"user"`
}

type CreateEmbeddingResponse struct {
    Object string `json:"object"`
    Data []Embedding `json:"data"`
    Model string `json:"model"`
    Usage Usage `json:"usage"`
}

type Embedding struct {
    Object string `json:"object"`
    Embedding []float64 `json:"embedding"`
    Index int `json:"index"`
}

type Usage struct {
    PromptTokens int `json:"prompt_tokens"`
    TotalTokens int `json:"total_tokens"`
}

func GetDefaultBody() CreateEmbeddingBody {
    return CreateEmbeddingBody{
        Input: []string{},
        Model: AllowedEncodingModels.TextEmbedding3Small,
        EncodingFormat: AllowedEmbeddingEncodingFormats.Float,
        Dimensions: 256, // Decided naively, depends on vector store used + accuracy desired
        User: "go-gpt-cli",
    }
}
