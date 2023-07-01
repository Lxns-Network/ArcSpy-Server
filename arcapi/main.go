package arcapi

import (
	"fmt"
)

type WebAPIResponse struct {
	Success bool        `json:"success"`
	Value   interface{} `json:"value"`
}

func GetPlayerUserMe(cookie string) (*WebAPIResponse, error) {
	return requestAPI("user/me", cookie)
}

func GetPlayerBestScore(cookie string, userId int, songId string, difficulty int) (*WebAPIResponse, error) {
	data, err := requestAPI(fmt.Sprintf(
		"score/song/me?song_id=%s&difficulty=%d&start=0&limit=30", songId, difficulty), cookie)
	if err != nil {
		return &WebAPIResponse{
			Success: false,
			Value:   nil,
		}, err
	}

	if !data.Success {
		return data, nil
	}

	for _, v := range data.Value.([]interface{}) {
		score := v.(map[string]interface{})
		if int(score["user_id"].(float64)) == userId {
			return &WebAPIResponse{
				Success: true,
				Value:   score,
			}, nil
		}
	}

	return &WebAPIResponse{
		Success: false,
		Value:   nil,
	}, nil
}
