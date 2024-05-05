package services

import (
	"io"
	"log"
	"net/http"
)

func HttpGet(url string) (string, int, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("[services.Util] Error getting url: %s", err)
		return "", http.StatusInternalServerError, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[services.Util] Error reading body: %s", err)
		return "", res.StatusCode, err
	}

	if c.Config.Debug {
		log.Printf("[services.Util] HTTP-GET %d Response: %s", res.StatusCode, body)
	}

	return string(body), res.StatusCode, nil
}
