package arcapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func requestAPI(endpoint string, cookie string) (*WebAPIResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://webapi.lowiro.com/webapi/%s", endpoint), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", fmt.Sprintf("sid=%s", cookie))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data WebAPIResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
