package image

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"net/http"
	"path/filepath"
	"strings"

	"io"
	"os"
	"strconv"
	"sync"

	"github.com/ephex2/go-gpt-cli/api"
	"github.com/gabriel-vasile/mimetype"
	"github.com/pborman/uuid"
)

const BaseImageRoute string = "/v1/images"
const createImageRoute string = "/generations"
const editImageRoute string = "/edits"
const createVariationRoute string = "/variations"

func CreateImage(folderPath string, prompt []string) (ImagePaths []string, revisedPrompt string, err error) {
	err = testFolderExists(folderPath)
	if err != nil {
		return
	}

	msg := formatChat(prompt)

	imageProfile, err := getDefaultProfile()
	if err != nil {
		return
	}

	imageProfile.CreateImageBody.Prompt = msg
	reqBody, err := json.Marshal(imageProfile.CreateImageBody)
	if err != nil {
		return
	}

	route := BaseImageRoute + createImageRoute
	buf, err := api.GenericRequest(nil, reqBody, route, "POST", imageProfile.OverrideUrl())
	if err != nil {
		return
	}

	var imageResponse CreateImageResponse
	err = json.Unmarshal(buf, &imageResponse)
	if err != nil {
		return
	}

	if imageProfile.CreateImageBody.ResponseFormat == nil {
		ImagePaths, err = getImages(folderPath, imageResponse, "url")
	} else {
		ImagePaths, err = getImages(folderPath, imageResponse, *imageProfile.CreateImageBody.ResponseFormat)
	}

	return
}

func CreateDalle3Image(folderPath string, prompt []string) (ImagePaths []string, revisedPrompt string, err error) {
	err = testFolderExists(folderPath)
	if err != nil {
		return
	}

	msg := formatChat(prompt)

	imageProfile, err := getDefaultProfile()
	if err != nil {
		return
	}

	if *imageProfile.CreateDalle3ImageBody.N != 1 {
		err = errors.New("Dalle3 image profile specifies an n > 1 , which is not supported. N must be 1")
		return
	}

	imageProfile.CreateDalle3ImageBody.Prompt = msg
	reqBody, err := json.Marshal(imageProfile.CreateDalle3ImageBody)
	if err != nil {
		return
	}

	route := BaseImageRoute + createImageRoute
	buf, err := api.GenericRequest(nil, reqBody, route, "POST", imageProfile.OverrideUrl())
	if err != nil {
		return
	}

	var imageResponse CreateImageResponse
	err = json.Unmarshal(buf, &imageResponse)
	if err != nil {
		return
	}

	if imageProfile.CreateDalle3ImageBody.ResponseFormat == nil {
		ImagePaths, err = getImages(folderPath, imageResponse, "url")
	} else {
		ImagePaths, err = getImages(folderPath, imageResponse, *imageProfile.CreateDalle3ImageBody.ResponseFormat)
	}

	return
}

func CreateEdit(filePath string, mask *os.File, folderPath string, prompt []string) (imagePaths []string, err error) {
	err = testFolderExists(folderPath)
	if err != nil {
		return
	}

	imageProfile, err := getDefaultProfile()
	if err != nil {
		return
	}

	msg := formatChat(prompt)
	if msg == "" {
		err = errors.New("prompt provided is empty, please provide a non-empty prompt to the function CreateEdit")
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	size, err := getOpenAiSize(filePath)
	if err != nil {
		return
	}

	fieldMap := imageProfile.CreateEditBody
	fieldMap["prompt"] = msg
	fieldMap["size"] = size

	details := []api.FileUploadDetails{
		{
			File:                f,
			UploadFormFieldName: "image",
		},
	}

	if mask != nil {
		maskDetails := api.FileUploadDetails{
			File:                mask,
			UploadFormFieldName: "mask",
		}

		details = append(details, maskDetails)
	}

	route := BaseImageRoute + editImageRoute
	buf, err := api.MultiPartFormRequest(details, fieldMap, route, "POST", imageProfile.OverrideUrl())
	if err != nil {
		return
	}

	var imageResponse CreateImageResponse
	err = json.Unmarshal(buf, &imageResponse)
	if err != nil {
		return
	}

	format, ok := fieldMap["response_format"]
	if !ok {
		format = "url"
	}

	imagePaths, err = getImages(folderPath, imageResponse, format)

	return
}

// Creates an image variation from the image provided at path filePath
// Returns a list of the new image paths for the new image variations created.
func CreateVariation(filePath string, folderPath string) (imagePaths []string, err error) {
	err = testFolderExists(folderPath)
	if err != nil {
		return
	}

	imageProfile, err := getDefaultProfile()
	if err != nil {
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	size, err := getOpenAiSize(filePath)
	if err != nil {
		return
	}

	fieldMap := imageProfile.CreateVariationBody
	fieldMap["size"] = size

	route := BaseImageRoute + createVariationRoute

	details := []api.FileUploadDetails{
		{
			File:                f,
			UploadFormFieldName: "image",
		},
	}

	buf, err := api.MultiPartFormRequest(details, fieldMap, route, "POST", imageProfile.OverrideUrl())
	if err != nil {
		return
	}

	var imageResponse CreateImageResponse
	err = json.Unmarshal(buf, &imageResponse)
	if err != nil {
		return
	}

	format, ok := fieldMap["response_format"]
	if !ok {
		format = "url"
	}

	imagePaths, err = getImages(folderPath, imageResponse, format)
	return
}

// HELPER FUNCTIONS
// TODO: Consider moving this to a utils lib ? It's been re-declared a few times now.
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

func getImages(folderPath string, imageResponse CreateImageResponse, format string) (ImagePaths []string, err error) {
	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	for _, data := range imageResponse.Data {
		wg.Add(1)
		switch format {
		case "url":
			go func(url string) {
				defer wg.Done()

				uuidHyphen := uuid.NewRandom()
				uuid := strings.Replace(uuidHyphen.String(), "-", "", -1)

				// handle http:// and https:// -> remove them, leaving the 's' prepended on the string in the https case
				// --> We will only be using these strings to get the extension of a file in a url
				noProtocol := strings.Replace(url, "http", "", -1)
				noSeparator := strings.Replace(noProtocol, "://", "", -1)
				noQueryParameters := strings.Split(noSeparator, "?")

				ext := filepath.Ext(noQueryParameters[0])

				actualPath := folderPath + "/" + uuid + ext

				e := downloadImage(url, actualPath)
				if e != nil {
					err = e
					return
				}

				mu.Lock()
				ImagePaths = append(ImagePaths, actualPath)
				mu.Unlock()
			}(data.Url)
		case "b64_json":
			go func(b64json string) {
				defer wg.Done()

				buf, err := base64.StdEncoding.DecodeString(b64json)
				if err != nil {
					return
				}

				uuidHyphen := uuid.NewRandom()
				uuid := strings.Replace(uuidHyphen.String(), "-", "", -1)

				jsonBuf, err := json.MarshalIndent(buf, "", "    ")
				if err != nil {
					return
				}

				m := mimetype.Detect(jsonBuf)
				if m.String() != "" {
					err = errors.New("mimetype not retrieved from image")
					return
				}

				actualPath := folderPath + "/" + uuid + m.Extension()

				file, err := os.Create(actualPath)
				if err != nil {
					return
				}
				defer file.Close()

				_, err = io.Copy(file, bytes.NewReader(jsonBuf))
				if err != nil {
					return
				}
			}(data.B64Json)
		default:
			err = errors.New("response_format not supported. The response format provided is: " + format)
			return
		}

		wg.Wait()
	}

	return
}

func downloadImage(url string, filePath string) (err error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	response, err := http.DefaultClient.Do(r)
	if err != nil {
		return
	}
	defer response.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return
	}

	return
}

func getOpenAiSize(filePath string) (size string, err error) {
	// Adding register format here since it is only used here, may move it to an init() or Init() function
	// At a later time. For now... KISS
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	img, _, err := image.DecodeConfig(f)
	if err != nil {
		return
	}

	size = strconv.Itoa(img.Width) + "x" + strconv.Itoa(img.Height)
	return
}

func GetB64Encoding(imagePath string) (b64 string, err error) {
	f, err := os.Open(imagePath)
	if err != nil {
		return
	}
	defer f.Close()

	m, err := mimetype.DetectFile(imagePath)
	if err != nil {
		return
	}

	switch m.String() {
	case "image/jpeg":
		b64 += "data:image/jpeg;base64,"
	case "image/png":
		b64 += "data:image/png;base64,"
	case "image/webp":
		b64 += "data:image/webp;base64,"
	case "image/gif":
		b64 += "data:image/gif;base64,"
	default:
		err = errors.New("mimetype of image not supported, it is: " + m.String())
	}

	buf, err := io.ReadAll(f)
	if err != nil {
		return
	}

	b64 += base64.StdEncoding.EncodeToString(buf)
	return
}

// Test folder exists + is folder
func testFolderExists(folderPath string) (err error) {
	stat, err := os.Stat(folderPath)
	if err != nil {
		return
	}

	if !stat.IsDir() {
		err = errors.New("The path provided is not a dir. path: " + folderPath)
		return
	}

	return
}
