package fraudlabspro

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// The Payment struct is the main object used to query the Payment Feedback API.
type Payment struct {
	configuration *Configuration
	baseUrl       string
}

// OpenPayment initializes with the Configuration object
func OpenPayment(config *Configuration) (*Payment, error) {
	var pay = &Payment{}
	pay.configuration = config
	pay.baseUrl = "https://api.fraudlabspro.com/v2"
	return pay, nil
}

// Report the final payment status back to the system, helping improve fraud detection and risk assessment
func (a *Payment) Feedback(params map[string]string) (map[string]interface{}, error) {
	var res map[string]interface{}
	var ex ErrorObject
	myUrl := a.baseUrl + "/payment/feedback"

	params["key"] = a.configuration.apiKey
	params["source"] = a.configuration.source
	params["source_version"] = a.configuration.sourceVersion
	params["format"] = "json"

	allowedKeys := []string{"key", "source", "source_version", "format", "email", "status", "message", "fraudlabspro_id"} // whitelist
	values := url.Values{}
	for _, key := range allowedKeys {
		if val, ok := params[key]; ok {
			values.Set(key, val)
		}
	}

	resp, err := http.PostForm(myUrl, values)

	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			return res, err
		}

		err = json.Unmarshal(bodyBytes, &res)

		if err != nil {
			return res, err
		}

		return res, nil
	} else if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized {
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			return res, err
		}

		bodyStr := string(bodyBytes[:])
		if strings.Contains(bodyStr, "error_message") {
			err = json.Unmarshal(bodyBytes, &ex)

			if err != nil {
				return res, err
			}
			return res, errors.New("Error: " + ex.Error.ErrorMessage)
		}
	}

	return res, errors.New("Error HTTP " + strconv.Itoa(int(resp.StatusCode)))
}
