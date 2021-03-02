package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

var defaultReadResponse responseReader = func(res *http.Response, resp interface{}) error {
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		resBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("%d Response: unknown error, invalid/empty response body", res.StatusCode)
		}

		return fmt.Errorf("%d Response: %s", res.StatusCode, resBytes)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("%d Response OK: Cannot read response body: %s", res.StatusCode, err)
	}

	resp = string(bytes)
	return nil

}

var readJSONResponse responseReader = func(res *http.Response, resp interface{}) error {
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return fmt.Errorf("%d Response: %s request %s failed: %v", res.StatusCode, res.Request.Method, res.Request.URL, errRes)
		}

		resBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("%d Response: unknown error, invalid/empty response body", res.StatusCode)
		}

		return fmt.Errorf("%d Response: %s", res.StatusCode, resBytes)
	}

	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		respBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("%d Response OK: Cannot read response body: %s", res.StatusCode, err)
		}
		return fmt.Errorf("%d Response OK: Invalid response body: %s", res.StatusCode, respBytes)
	}

	return nil

}

var readXMLResponse responseReader = func(res *http.Response, resp interface{}) error {
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		errBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("%d Response: unknown error, invalid/empty response body", res.StatusCode)
		}
		var errRes map[string]interface{}
		if err := xml.Unmarshal(errBytes, &errRes); err == nil {
			return fmt.Errorf("%d Response: %s request %s failed: %v", res.StatusCode, res.Request.Method, res.Request.URL, errRes)
		}

		return fmt.Errorf("%d Response: %s", res.StatusCode, errBytes)
	}

	respBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("%d Response OK: Cannot read response body: %s", res.StatusCode, err)
	}

	if err := xml.Unmarshal(respBytes, resp); err != nil {

		return fmt.Errorf("%d Response OK: Invalid XML response body: %s: body=%s", res.StatusCode, err, respBytes)
	}

	return nil

}
