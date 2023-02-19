package httplib

import "net/url"

func NewURL(baseURL, path string, kv ...string) (*url.URL, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	uri := base.JoinPath(path)
	if err != nil {
		return nil, err
	}

	if len(kv)%2 != 0 {
		return uri, nil
	}

	q := uri.Query()
	for i := 0; i < len(kv); i += 2 {
		q.Add(kv[i], kv[i+1])
	}
	uri.RawQuery = q.Encode()

	return uri, nil
}
