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

// The SmsVerification struct is the main object used to query the SMS Verification API.
type SmsVerification struct {
	configuration *Configuration
	baseUrl       string
}

// OpenSmsVerification initializes with the Configuration object
func OpenSmsVerification(config *Configuration) (*SmsVerification, error) {
	var sms = &SmsVerification{}
	sms.configuration = config
	sms.baseUrl = "https://api.fraudlabspro.com/v2"
	return sms, nil
}

// Sends SMS for verification purposes
func (a *SmsVerification) SendSms(params map[string]string) (map[string]interface{}, error) {
	var res map[string]interface{}
	var ex ErrorObject
	myUrl := a.baseUrl + "/verification/send"

	params["key"] = a.configuration.apiKey
	params["source"] = a.configuration.source
	params["source_version"] = a.configuration.sourceVersion
	params["format"] = "json"

	if val, ok := params["tel"]; ok && val != "" {
		if !strings.HasPrefix(val, "+") {
			params["tel"] = "+" + params["tel"]
		}
	}

	allowedKeys := []string{"key", "source", "source_version", "format", "tel", "country_code", "mesg", "otp_timeout"} // whitelist
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

// Verifies the validity of the OTP sent
func (a *SmsVerification) VerifySms(params map[string]string) (map[string]interface{}, error) {
	var res map[string]interface{}
	var ex ErrorObject
	myUrl := a.baseUrl + "/verification/result"

	params["key"] = a.configuration.apiKey
	params["source"] = a.configuration.source
	params["source_version"] = a.configuration.sourceVersion
	params["format"] = "json"

	allowedKeys := []string{"key", "source", "source_version", "format", "tran_id", "otp"} // whitelist
	values := url.Values{}
	for _, key := range allowedKeys {
		if val, ok := params[key]; ok {
			values.Set(key, val)
		}
	}

	myUrl = myUrl + "?" + values.Encode()
	resp, err := http.Get(myUrl)

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
