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

// The Order struct is the main object used to query the Fraudlabs Pro API.
type Order struct {
	configuration *Configuration
	baseUrl       string
}

// OpenOrder initializes with the Configuration object
func OpenOrder(config *Configuration) (*Order, error) {
	var flp = &Order{}
	flp.configuration = config
	flp.baseUrl = "https://api.fraudlabspro.com/v2"
	return flp, nil
}

// Validate will screen an order for fraud
func (a *Order) Validate(params map[string]string) (map[string]interface{}, error) {
	var res map[string]interface{}
	var ex ErrorObject

	myUrl := a.baseUrl + "/order/screen"

	params["key"] = a.configuration.apiKey
	params["source"] = a.configuration.source
	params["source_version"] = a.configuration.sourceVersion
	params["format"] = "json"

	if val, ok := params["email"]; ok && val != "" {
		parts := strings.Split(val, "@")
		if len(parts) != 2 {
			params["email_domain"] = ""
		} else {
			params["email_domain"] = parts[1]
			params["email_hash"] = doHash(val)
		}
	}

	if val, ok := params["user_phone"]; ok && val != "" {
		params["user_phone"] = stripNonDigits(val)
	}

	if val, ok := params["currency"]; !ok || val == "" {
		params["currency"] = "USD"
	}

	if val, ok := params["number"]; ok && val != "" {
		maxLen := 9
		if len(val) < maxLen {
			maxLen = len(val)
		}
		params["bin_no"] = val[:maxLen]
		params["card_hash"] = doHash(val)
	}

	if val, ok := params["amount"]; ok || val != "" {
		formatted, err := formatTo2Decimals(val)
		if err != nil {
			return res, errors.New("Error: " + err.Error())
		}
		params["amount"] = formatted
	}

	allowedKeys := []string{"ip", "key", "source", "source_version", "format", "last_name", "first_name", "bill_addr", "bill_city", "bill_state", "bill_country", "bill_zip_code", "ship_last_name", "ship_first_name", "ship_addr", "ship_city", "ship_state", "ship_country", "ship_zip_code", "user_phone", "email", "email_hash", "email_domain", "username", "bin_no", "card_hash", "avs_result", "cvv_result", "user_order_id", "amount", "quantity", "currency", "department", "payment_gateway", "payment_mode", "flp_checksum"} // whitelist
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

// Feedback will send a decision back to FraudLabs Pro
func (a *Order) Feedback(params map[string]string) (map[string]interface{}, error) {
	var res map[string]interface{}
	var ex ErrorObject

	myUrl := a.baseUrl + "/order/feedback"

	params["key"] = a.configuration.apiKey
	params["source"] = a.configuration.source
	params["source_version"] = a.configuration.sourceVersion
	params["format"] = "json"

	allowedActions := []string{OrderActionApprove, OrderActionReject, OrderActionRejectBlacklist} // whitelist
	allowedKeys := []string{"id", "key", "source", "source_version", "format", "action", "note"}  // whitelist
	values := url.Values{}
	for _, key := range allowedKeys {
		if val, ok := params[key]; ok {
			if key == "action" && !contains(allowedActions, val) {
				return res, errors.New("Error: Invalid order status provided.")
			}
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

// GetTransaction retrieves a validation result
func (a *Order) GetTransaction(transactionID string) (map[string]interface{}, error) {
	params := make(map[string]string)
	var res map[string]interface{}
	var ex ErrorObject

	myUrl := a.baseUrl + "/order/result"

	params["key"] = a.configuration.apiKey
	params["source"] = a.configuration.source
	params["source_version"] = a.configuration.sourceVersion
	params["format"] = "json"

	transactionID = strings.TrimSpace(transactionID)

	if transactionID == "" {
		return res, errors.New("Error: Invalid transaction ID.")
	}
	params["id"] = transactionID

	allowedKeys := []string{"id", "key", "source", "source_version", "format"} // whitelist
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
