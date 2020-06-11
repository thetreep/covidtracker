package cds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/thetreep/covidtracker"
)

const (
	authenticatePath  = "Authenticate"
	hotelByPrefixPath = "Hotels"
)

var (
	testEndPoint = "https://bookings.cdsgroupe.com/cds-api-test/v1/"
	prodEndPoint = "https://bookings.cdsgroupe.com/cds-api/v1/"
	endpoint     = testEndPoint

	User          = os.Getenv("THETREEP_COVIDTRACKER_CDS_API_USER")
	Password      = os.Getenv("THETREEP_COVIDTRACKER_CDS_API_PASSWORD")
	AgentDutyCode = os.Getenv("THETREEP_COVIDTRACKER_CDS_API_DUTY_CODE")
)

var (
	Service cdsAPI
)

type cdsAPI interface {
	HotelsByPrefix(city string, prefix string) ([]*covidtracker.Hotel, error)
}

type tracedCDSService struct {
	service cdsAPI
}

func Init() {
	endpoint = testEndPoint
	Service = &tracedCDSService{service: newClient(nil)}
}

func (s *tracedCDSService) HotelsByPrefix(c string, p string) ([]*covidtracker.Hotel, error) {
	out, err := s.service.HotelsByPrefix(c, p)
	return out, err
}

type Client struct {
	authToken  string
	baseURL    *url.URL
	httpClient *http.Client
}

func newClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(endpoint)

	c := &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
	return c
}

func (c *Client) HotelsByPrefix(city string, prefix string) ([]*covidtracker.Hotel, error) {
	clientCode, user, pwd := AgentDutyCode, User, Password
	err := c.authenticate(user, pwd)
	if err != nil {
		return nil, fmt.Errorf("cannot login with CDS: %s", err)
	}
	params := map[string][]string{
		"agentDutyCode": []string{clientCode},
		"prefix":        []string{prefix},
	}

	if city != "" {
		params["city"] = []string{city}
		params["country"] = []string{"FR"}
	}

	req, err := c.NewRequest("GET", hotelByPrefixPath, nil, params)
	if err != nil {
		return nil, err
	}

	var resp *hotelResultResp
	_, err = c.Do(req, &resp)
	if err != nil {
		return nil, fmt.Errorf("cannot search hotel: %s", err)
	}
	var hotels []*covidtracker.Hotel
	for _, h := range resp.Hotels {
		hotel := h.ToHotel()
		hotels = append(hotels, hotel)
	}

	return hotels, nil
}

//authenticate client to get a authToken from username/password credentials
func (c *Client) authenticate(user, pwd string) error {
	req, err := c.NewRequest("POST", "Authenticate", AuthentificationReq{Username: user, Password: pwd})
	req.Header.Set("Content-Type", "application/json")
	var result AuthentificationResp
	_, err = c.Do(req, &result)
	if err != nil {
		return fmt.Errorf("cannot authenticate with CDS Rest API %v", err)
	}
	if result.Token == "" {
		return fmt.Errorf("cannot authenticate with CDS Rest API: got empty token %s et %s", user, pwd)
	}

	c.authToken = result.Token
	return nil
}

func (c *Client) NewRequest(method, urlStr string, body interface{}, params ...map[string][]string) (*http.Request, error) {
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if len(params) > 0 {
		v := url.Values{}
		for key, paramValues := range params[0] {
			for _, paramValue := range paramValues {
				v.Add(key, paramValue)
			}
		}
		u.RawQuery = v.Encode()
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err = enc.Encode(body); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, output interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 || resp.StatusCode < 199 {
		bodyJSON := make(map[string]interface{})
		if err = json.NewDecoder(resp.Body).Decode(&bodyJSON); err != nil {
			return resp, fmt.Errorf("error %d", resp.StatusCode)
		}
		errorMsg, ok := bodyJSON["Message"]
		if !ok {
			return resp, fmt.Errorf("error %d", resp.StatusCode)
		}
		return resp, fmt.Errorf("error %d: %v", resp.StatusCode, errorMsg)
	}

	if output != nil {
		if w, ok := output.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(&output)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}
	return resp, err
}
