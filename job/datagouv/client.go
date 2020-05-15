const (
)

var (
	BaseURL = url
)

type Client struct {
	url        *url.URL
	httpClient *http.Client

	apiToken string
}

func NewClient() *Client {
	url, _ := url.Parse(BaseURL)
	c := &Client{
		httpClient: http.DefaultClient,
		url:        url,
	}

	//TODO set env variable
	// c.apiToken = os.Getenv("THETREEP_DATAGOUV_TOKEN")
	c.apiToken = ""

	return c
}

func (c *Client) NewRequest(method, path string, body interface{}, params ...map[string][]string) (*http.Request, error) {
	u := c.url
	u.Path = path
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
		if err := enc.Encode(body); err != nil {
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

	return req, nil
}

func (c *Client) Do(req *http.Request, output interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 || resp.StatusCode < 199 {
		//TODO read error properly
		return resp, err
	}

	if output != nil {
		if w, ok := output.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(output)
			if err == io.EOF { // ignore EOF errors caused by empty response body
				return resp, nil
			}
			if err != nil {
				return resp, fmt.Errorf("cannot decode body as json: %s", err.Error())
			}
		}
	}
	return resp, err
}
