package api

const (
	jsonUtfContent = "application/json; charset=utf-8"
	xmlContent     = "application/xml; charset=utf-8"
)

type ClientOption func(*Client) error

var JSONClient ClientOption = func(c *Client) error {
	c.accept = jsonUtfContent
	c.contentType = jsonUtfContent
	c.readResponse = readJSONResponse
	return nil
}

var XMLClient ClientOption = func(c *Client) error {
	c.accept = xmlContent
	c.contentType = "application/x-www-form-urlencoded"
	c.readResponse = readXMLResponse
	return nil
}
