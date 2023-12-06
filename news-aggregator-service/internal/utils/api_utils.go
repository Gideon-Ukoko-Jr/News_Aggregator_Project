// utils/api_utils.go

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"news-aggregator-service/internal/models"
	"time"
)

// Making a GET request to the News API and returning the response.
func GetNewsAPIResponse(apiURL string) (*models.NewsApiResponse, error) {
	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("Failed to fetch News API response: %v\n", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Non-OK status code from News API: %v\n", response.StatusCode)
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("Response Body: %s\n", body)
		return nil, fmt.Errorf("Non-OK status code from News API: %v", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errString := fmt.Sprintf("Failed to read News API response body: %v", err)
		fmt.Println(errString)
		return nil, errors.New(errString)
	}

	var newsResponse models.NewsApiResponse
	err = json.Unmarshal(body, &newsResponse)
	if err != nil {
		errString := fmt.Sprintf("Failed to unmarshal News API response: %v", err)
		fmt.Println(errString)
		return nil, errors.New(errString)
	}

	logString := fmt.Sprintf("News API Request - Good to go: %v", &newsResponse.Status)
	fmt.Println(logString)

	return &newsResponse, nil
}

// Making a GET request to the Guardian News API and returning the response.
func GetGuardianNewsAPIResponse(apiURL string) (*models.GuardianNewsApiResponse, error) {
	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("Failed to fetch Guardian News API response: %v\n", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Non-OK status code from Guardian News API: %v\n", response.StatusCode)
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("Response Body: %s\n", body)
		return nil, fmt.Errorf("Non-OK status code from Guardian News API: %v", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errString := fmt.Sprintf("Failed to read Guardian News API response body: %v", err)
		fmt.Println(errString)
		return nil, errors.New(errString)
	}

	var guardianResponse models.GuardianNewsApiResponse
	err = json.Unmarshal(body, &guardianResponse)
	if err != nil {
		errString := fmt.Sprintf("Failed to unmarshal Guardian News API response: %v", err)
		fmt.Println(errString)
		return nil, errors.New(errString)
	}

	logString := fmt.Sprintf("Guardian API Request - Good to go: %v", &guardianResponse.Response.Status)
	fmt.Println(logString)
	return &guardianResponse, nil
}

// Parsing the Guardian API date format to a time.Time object.
func ParseGuardianDate(dateString string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		// Returning current time as a fallback
		fmt.Printf("Error parsing Guardian API date: %v\n", err)
		return time.Now()
	}
	return parsedTime
}
