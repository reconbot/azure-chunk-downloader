package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func request(token string, method string, url string) (*http.Response, error) {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("could not make request")
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("x-ms-version", "2020-06-12")
	request.Header.Add("x-ms-date", time.Now().UTC().Format(http.TimeFormat))

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal().Err(err).Msg("could not do request")
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().
			Any("statusCode", resp.StatusCode).
			Any("status", resp.Status).
			Any("responseHeader", resp.Header).
			// Any("requestHeader", request.Header).
			Msg("bad status code")
		return nil, fmt.Errorf("bad response")
	}
	return resp, nil
}
