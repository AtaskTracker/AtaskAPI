package googleCloudService

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const baseUrl = "https://google-cloud-task-processor.herokuapp.com"

type GoogleCloudService struct {
}

func (s *GoogleCloudService) UploadImage(name string, image string) (string, error) {
	url := baseUrl + "/storage/image"
	values := map[string]string{"name": name, "payload": image}
	jsonData, err := json.Marshal(values)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	var res map[string]string

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", err
	}

	return res["url"], nil
}
