# Quickstart

## Dependencies

This module requires API key to function. You may subscribe a free API key at https://www.fraudlabspro.com

## Installation

Install this module using the command below:

``` bash
go get github.com/fraudlabspro/fraudlabspro-go
```

## Sample Codes

### Validate Order

You can validate your order as below:

```go
package main

import (
	"github.com/fraudlabspro/fraudlabspro-go/fraudlabspro"
	"fmt"
)

func main() {
	apikey := "YOUR_API_KEY"

	config, err := fraudlabspro.OpenConfiguration(apikey)

	if err != nil {
		fmt.Print(err)
		return
	}

	flp, err := fraudlabspro.OpenOrder(config)

	if err != nil {
		fmt.Print(err)
		return
	}

	params := make(map[string]string)

	params["ip"] = "146.112.62.105"

	params["user_order_id"] = "67398"
	params["currency"] = "USD"
	params["amount"] = "79.89"
	params["quantity"] = "1"
	params["payment_gateway"] = "Gateway To Bliss"
	params["payment_mode"] = fraudlabspro.PaymentMethodCreditCard

	params["number"] = "4556553172971283"

	params["first_name"] = "Hector"
	params["last_name"] = "Henderson"
	params["email"] = "hh5566@gmail.com"
	params["user_phone"] = "561-628-8674"
	params["bill_addr"] = "1766 Powder House Road"
	params["bill_city"] = "West Palm Beach"
	params["bill_state"] = "FL"
	params["bill_zip_code"] = "33401"
	params["bill_country"] = "US"

	params["ship_first_name"] = "Hector"
	params["ship_last_name"] = "Henderson"
	params["ship__addr"] = "4469 Chestnut Street"
	params["ship_city"] = "Tampa"
	params["ship_state"] = "FL"
	params["ship_zip_code"] = "33602"
	params["ship_country"] = "US"

	res, err := flp.Validate(params)

	fmt.Printf("%+v\n", res)
}
```

### Get Transaction

You can get the details of a transaction as below:

```go
package main

import (
	"github.com/fraudlabspro/fraudlabspro-go/fraudlabspro"
	"fmt"
)

func main() {
	apikey := "YOUR_API_KEY"

	config, err := fraudlabspro.OpenConfiguration(apikey)

	if err != nil {
		fmt.Print(err)
		return
	}

	flp, err := fraudlabspro.OpenOrder(config)

	if err != nil {
		fmt.Print(err)
		return
	}

	res, err := flp.GetTransaction("20250827-ZZZZZZ") // replace with your actual FraudLabs Pro transaction ID

	
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%+v\n", res)
}
```

### Feedback

You can approve, reject or blacklist a transaction as below:

```go
package main

import (
	"github.com/fraudlabspro/fraudlabspro-go/fraudlabspro"
	"fmt"
)

func main() {
	apikey := "YOUR_API_KEY"

	config, err := fraudlabspro.OpenConfiguration(apikey)

	if err != nil {
		fmt.Print(err)
		return
	}

	flp, err := fraudlabspro.OpenOrder(config)

	if err != nil {
		fmt.Print(err)
		return
	}

	params := make(map[string]string)

	params["id"] = "20250827-ZZZZZZ" // replace with your actual FraudLabs Pro transaction ID
	params["action"] = fraudlabspro.OrderActionReject
	
	res, err := flp.Feedback(params)

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%+v\n", res)
}
```

### Send SMS Verification

You can send SMS verification for authentication purpose as below:

```go
package main

import (
	"github.com/fraudlabspro/fraudlabspro-go/fraudlabspro"
	"fmt"
)

func main() {
	apikey := "YOUR_API_KEY"

	config, err := fraudlabspro.OpenConfiguration(apikey)

	if err != nil {
		fmt.Print(err)
		return
	}

	sms, err := fraudlabspro.OpenSmsVerification(config)

	if err != nil {
		fmt.Print(err)
		return
	}

	params := make(map[string]string)

	params["tel"] = "+11111111111111" // replace with the actual recipient mobile phone number
	params["country_code"] = "XX" // replace with the actual recipient country
	params["mesg"] = "Hello, your OTP is <otp>." // custom message

	res, err := sms.SendSms(params)

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%+v\n", res)
}
```

### Get SMS Verification Result

You can verify the OTP sent by Fraudlabs Pro SMS verification API as below:

```go
package main

import (
	"github.com/fraudlabspro/fraudlabspro-go/fraudlabspro"
	"fmt"
)

func main() {
	apikey := "YOUR_API_KEY"

	config, err := fraudlabspro.OpenConfiguration(apikey)

	if err != nil {
		fmt.Print(err)
		return
	}

	sms, err := fraudlabspro.OpenSmsVerification(config)

	if err != nil {
		fmt.Print(err)
		return
	}

	params := make(map[string]string)

	params["tran_id"] = "XXXXXXXXXXXXXXXXXXXX" // replace with the actual transaction ID returned by the SendSms function
	params["otp"] = "000000" // replace with the actual OTP received by the recipient mobile phone

	res, err := sms.VerifySms(params)

	
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%+v\n", res)
}
```