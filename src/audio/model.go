package audio

type verboseJson string

const (
	verboseJsonConst verboseJson = "verbose_json"
)

var AllowedSpeecModels = struct {
	TTS1   string
	TTS1HD string
}{
	TTS1:   "tts-1",
	TTS1HD: "tts-1-hd",
}

var AllowedVoices = struct {
	Alloy   string
	Echo    string
	Fable   string
	Onyx    string
	Nova    string
	Shimmer string
}{
	Alloy:   "alloy",
	Echo:    "echo",
	Fable:   "fable",
	Onyx:    "onyx",
	Nova:    "nova",
	Shimmer: "shimmer",
}

// Exclude response_json as it is reserved for VerboseTranscriptionRequestBody
var AllowedTranscriptionResponseFormats = struct {
	Json string
	Text string
	Srt  string
	Vtt  string
}{
	Json: "json",
	Text: "text",
	Srt:  "srt",
	Vtt:  "vtt",
}

var AllowedTranslationResponseFormats = struct {
	Json string
	Text string
	Srt  string
	Vtt  string
}{
	Json: "json",
	Text: "text",
	Srt:  "srt",
	Vtt:  "vtt",
}

var AllowedTimeStampGranularities = struct {
	Segment string
	Word    string
}{
	Segment: "segment",
	Word:    "word",
}

/*
In this map:
- `file` is a required field of type []byte representing the audio file object.
- `model` is a required field of type string representing the ID of the model to use.
- `language` is an optional field of type string representing the language of the input audio.
- `prompt` is an optional field of type string providing guidance for the model's style or continuing a previous audio segment.
- `response_format` is an optional field of type string representing the format of the transcript output.
- `temperature` is an optional field of type float64 representing the sampling temperature.
*/
var CreateTranscriptionBody = map[string]string{
	"model":           "whisper-1",
	"language":        AllowedIsoLanguages.English,
	"prompt":          "",
	"response_format": AllowedTranscriptionResponseFormats.Json,
	"temperature":     "0",
}

/*
In this map:
- `file` is a required field of type []byte representing the audio file object.
- `model` is a required field of type string representing the ID of the model to use.
- `language` is an optional field of type string representing the language of the input audio.
- `prompt` is an optional field of type string providing guidance for the model's style or continuing a previous audio segment.
- `response_format` is an optional field of type string representing the format of the transcript output.
- `temperature` is an optional field of type float64 representing the sampling temperature.
*/
var CreateTranslationBody = map[string]string{
	"model":           "whisper-1",
	"language":        AllowedIsoLanguages.English,
	"prompt":          "",
	"response_format": AllowedTranslationResponseFormats.Json,
	"temperature":     "0",
}

/*
In this map:
- `file` is a required field of type []byte representing the audio file object.
- `model` is a required field of type string representing the ID of the model to use.
- `language` is an optional field of type string representing the language of the input audio.
- `prompt` is an optional field of type string providing guidance for the model's style or continuing a previous audio segment.
- `response_format` is a mandatory field of type string representing the format of the transcript output. In this case, it must be 'verbose_json'.
- `temperature` is an optional field of type float64 representing the sampling temperature.
- `timestamp_granularities` can only be used when response_format is set to 'verbose_json'. 'timestamp_granularities' is an optional field of type array of strings representing the timestamp granularities for the transcription.
*/
var CreateVerboseTranscriptionBody = map[string]string{
	"model":                   "whisper-1",
	"language":                AllowedIsoLanguages.English,
	"prompt":                  "",
	"response_format":         string(verboseJsonConst),
	"temperature":             "0",
	"timestamp_granularities": AllowedTimeStampGranularities.Segment,
}

/*
In this map:
- `file` is a required field of type []byte representing the audio file object.
- `model` is a required field of type string representing the ID of the model to use.
- `language` is an optional field of type string representing the language of the input audio.
- `prompt` is an optional field of type string providing guidance for the model's style or continuing a previous audio segment.
- `response_format` is a mandatory field of type string representing the format of the transcript output. In this case, it must be 'verbose_json'.
- `temperature` is an optional field of type float64 representing the sampling temperature.
*/
var CreateVerboseTranslationBody = map[string]string{
	"model":           "whisper-1",
	"language":        AllowedIsoLanguages.English,
	"prompt":          "",
	"response_format": string(verboseJsonConst),
	"temperature":     "0",
}

type CreateSpeechBody struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice"`
	ResponseFormat *string  `json:"response_format,omitempty"`
	Speed          *float64 `json:"speed,omitempty"`
}

type CreateTranscriptionJsonResponse struct {
	Text string `json:"text"`
}

type CreateTranslationJsonResponse struct {
	Text string `json:"text"`
}

type CreateVerboseTranscriptionResponse struct {
	Task     string    `json:"task"`
	Language string    `json:"language"`
	Duration float64   `json:"duration"`
	Text     string    `json:"text"`
	Segments []Segment `json:"segments"`
}

type CreateVerboseTranslationResponse struct {
	Task     string    `json:"task"`
	Language string    `json:"language"`
	Duration float64   `json:"duration"`
	Text     string    `json:"text"`
	Segments []Segment `json:"segments"`
}

type Segment struct {
	ID               int     `json:"id"`
	Seek             int     `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []int   `json:"tokens"`
	Temperature      float64 `json:"temperature"`
	AvgLogprob       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
}

func DefaultCreateSpeechBody() CreateSpeechBody {
    format := new(string)
    *format = "mp3"

	return CreateSpeechBody{
		Model:          AllowedSpeecModels.TTS1,
		Input:          "",
		Voice:          AllowedVoices.Nova,
		ResponseFormat: format,
	}
}
