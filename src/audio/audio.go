package audio

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ephex2/go-gpt-cli/api"
	"github.com/ephex2/go-gpt-cli/chat"
	"github.com/ephex2/go-gpt-cli/log"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

const BaseRoute = "/v1/audio"
const createSpeechRoute = BaseRoute + "/speech"
const createTranscriptionRoute = BaseRoute + "/transcriptions"
const createTranslationRoute = BaseRoute + "/translations"

// Creates a speech using TTS and saves it to speechPath.
func CreateSpeech(prompt []string) (speechPath string, err error) {
	msg := formatChat(prompt)

	// do audio stuff
	audioP, err := getDefaultProfile()
	if err != nil {
		return
	}

	audioP.CreateSpeechBody.Input = msg

	bodyBuf, err := json.Marshal(audioP.CreateSpeechBody)
	if err != nil {
		return
	}

	buf, err := api.GenericRequest(nil, bodyBuf, createSpeechRoute, "POST")
	if err != nil {
		return
	}

	// Setup audio file
	m := mimetype.Detect(buf)

	if m.Extension() == "" {
		err = errors.New("unable to determine mime type of returned audio file")
		return
	}

	speechPath, err = audioP.Save(buf, m.Extension())
	if err != nil {
		return
	}

	return
}

func CreateTranscription(filePath string, prompt []string) (s string, err error) {
	msg := formatChat(prompt)
	audioP, err := getDefaultProfile()
	if err != nil {
		return
	}

	fieldMap := audioP.CreateTranscriptionBody

	if msg != "" {
		fieldMap["prompt"] = msg
	}

	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	details := []api.FileUploadDetails{
		{
			File:                f,
			UploadFormFieldName: "file",
		},
	}

	buf, err := api.MultiPartFormRequest(details, fieldMap, createTranscriptionRoute, "POST")
	if err != nil {
		return
	}

	format := fieldMap["response_format"]

	if format != "json" {
		log.Debug("CreateTranscriptionBody response_format unsupported, returning raw string:\n")
		s = string(buf)
		return
	}

	var transcriptionResponse CreateTranscriptionJsonResponse
	err = json.Unmarshal(buf, &transcriptionResponse)
	if err != nil {
		return
	}

	s = transcriptionResponse.Text
	return
}

func CreateTranslation(filePath string, prompt []string) (s string, err error) {
	msg := formatChat(prompt)
	audioP, err := getDefaultProfile()
	if err != nil {
		return
	}

	fieldMap := audioP.CreateTranslationBody

	if msg != "" {
		fieldMap["prompt"] = msg
	}

	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	details := []api.FileUploadDetails{
		{
			File:                f,
			UploadFormFieldName: "file",
		},
	}

	buf, err := api.MultiPartFormRequest(details, fieldMap, createTranslationRoute, "POST")

	format := fieldMap["response_format"]

	if format != "json" {
		log.Debug("CreateTranslationBody response_format unsupported, returning raw string:\n")
		s = string(buf)
		return
	}

	var translationResponse CreateTranslationJsonResponse
	err = json.Unmarshal(buf, &translationResponse)
	if err != nil {
		return
	}

	s = translationResponse.Text
	return
}

func CreateVerboseTranscription(filePath string, prompt []string) (resp CreateVerboseTranscriptionResponse, err error) {
	msg := formatChat(prompt)
	audioP, err := getDefaultProfile()
	if err != nil {
		return
	}

	// Main difference from other function, consider modularizing fieldMap and passing it as arg
	fieldMap := audioP.CreateVerboseTranscriptionBody

	format := fieldMap["response_format"]

	if format != string(verboseJsonConst) {
		err = errors.New("the response_format in the map is not verbose_json")
		return
	}

	if msg != "" {
		fieldMap["prompt"] = msg
	}

	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	details := []api.FileUploadDetails{
		{
			File:                f,
			UploadFormFieldName: "file",
		},
	}

	buf, err := api.MultiPartFormRequest(details, fieldMap, createTranscriptionRoute, "POST")
	if err != nil {
		return
	}

	// response type is different from other function as well
	var transcriptionResponse CreateVerboseTranscriptionResponse
	err = json.Unmarshal(buf, &transcriptionResponse)
	if err != nil {
		return
	}

	resp = transcriptionResponse

	return
}

func CreateVerboseTranslation(filePath string, prompt []string) (resp CreateVerboseTranslationResponse, err error) {
	msg := formatChat(prompt)
	audioP, err := getDefaultProfile()
	if err != nil {
		return
	}

	// Main difference from other function, consider modularizing fieldMap and passing it as arg
	fieldMap := audioP.CreateVerboseTranslationBody

	format := fieldMap["response_format"]

	if format != string(verboseJsonConst) {
		err = errors.New("the response_format in the map is not verbose_json")
		return
	}

	if msg != "" {
		fieldMap["prompt"] = msg
	}

	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	details := []api.FileUploadDetails{
		{
			File:                f,
			UploadFormFieldName: "file",
		},
	}

	buf, err := api.MultiPartFormRequest(details, fieldMap, createTranscriptionRoute, "POST")
	if err != nil {
		return
	}

	// response type is different from other function as well
	var translationResponse CreateVerboseTranslationResponse
	err = json.Unmarshal(buf, &translationResponse)
	if err != nil {
		return
	}

	resp = translationResponse

	return
}

// Synchronously play audio over speakers
func PlayAudioFile(path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	// Do not use DetectReader; if you do, you will need to rewind the file's io.ReadCloser
	m, err := mimetype.DetectFile(path)
	if err != nil {
		return
	}

	var streamer beep.StreamCloser
	var format beep.Format

	switch m.String() {
	case "audio/mpeg":
		streamer, format, err = mp3.Decode(f)
	case "audio/wav":
		streamer, format, err = wav.Decode(f)
	case "audio/flac":
		streamer, format, err = flac.Decode(f)
	default:
		// Extend this switch statement as needed, adding appropriate play functions
		// Currently support all media types in "github.com/gopxl/beep"
		errorString := `
Below are some alternative ways to play the audio file from the command line -

For Linux, using ffplay should do the trick:
ffplay /path/to/audioFile.unsupported -autoexit

For Windows, consider playing the file from PowerShell:
# Play a single file - source: https://stackoverflow.com/questions/25895428/how-to-play-mp3-with-powershell-simple
Add-Type -AssemblyName presentationCore
$mediaPlayer = New-Object system.windows.media.mediaplayer
$mediaPlayer.open('C:\temp\audioFile.unsupported')
$mediaPlayer.Play()`

		err = errors.New("Mimetype of audio file not supported. It is: " + m.String() + "\n" + errorString)
		return
	}

	// Shared clean-up of switch statement
	if err != nil {
		return
	}
	defer streamer.Close()

	// Code resumes
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return
	}

	done := make(chan struct{})
	defer close(done)

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- struct{}{}
	})))

	<-done

	return
}

func ReadAudioPrompt(prompt []string) (err error) {
	msg, err := chat.CreateChatCompletion(prompt)
	if err != nil {
		return
	}

	fmt.Println(msg)
	path, err := CreateSpeech([]string{msg})
	if err != nil {
		return
	}

	err = PlayAudioFile(path)
	if err != nil {
		return
	}

	return
}

func formatChat(chat []string) string {
	var formattedChat string

	for i, word := range chat {
		if i+1 == len(chat) {
			formattedChat += word
		} else {
			formattedChat += word + " "
		}
	}

	return formattedChat
}
