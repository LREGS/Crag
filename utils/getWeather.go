package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	// headers "workspaces/github.com/lregs/Crag/headers"
	helpers "workspaces/github.com/lregs/Crag/helper"
	t "workspaces/github.com/lregs/Crag/types"
)

func GetForecast(url string, headers map[string]string, client *http.Client) (t.Forecast, error) {

	//eventually req functionality will be in a router, so when an end-point is hit, a request is made and sent to getForecast that returns a response

	//

	req, err := createRequest(url, headers)
	helpers.CheckError(err)

	res, err := sendRequest(req, client)
	helpers.CheckError(err)

	Forecast, err := parseResponse(res)
	helpers.CheckError(err)

	return Forecast, nil

}

func createRequest(apiUrl string, headers map[string]string) (*http.Request, error) {

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// jsonHeaders, err := json.Marshal(headers)
	// if err != nil {
	// 	return nil, err
	// }
	// print(jsonHeaders[0])
	// req.Header.Add("Headers", string(jsonHeaders))
	// fmt.Println(req.Header)
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return req, nil

}

func sendRequest(req *http.Request, client *http.Client) (*http.Response, error) {
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		log.Println(fmt.Sprintf("Status %v", res.Status))
	}

	return res, nil
}

func parseResponse(res *http.Response) (t.Forecast, error) {
	body, err := io.ReadAll(res.Body)
	helpers.CheckError(err)

	// var ResponseData = make(map[string]interface{})
	var Forecast t.Forecast

	err = json.Unmarshal(body, &Forecast)
	helpers.CheckError(err)

	return Forecast, nil
}
