package httplib

import "net/url"

type GetURLFunc func() (*url.URL, error)

type URLBuilder interface {
	Build() (*url.URL, error)
}

type URLBuilderOptions struct {
	BaseURL string
	Path    string
	Params  []string
}

func NewURLBuilder(opts URLBuilderOptions) URLBuilder {
	return &urlBuilder{
		baseURL: opts.BaseURL,
		path:    opts.Path,
		params:  opts.Params,
	}
}

type urlBuilder struct {
	baseURL string
	path    string
	params  []string
}

func (ub *urlBuilder) Build() (*url.URL, error) {
	base, err := url.Parse(ub.baseURL)
	if err != nil {
		return nil, err
	}

	uri := base.JoinPath(ub.path)
	if err != nil {
		return nil, err
	}

	if len(ub.params)%2 != 0 {
		return uri, nil
	}

	q := uri.Query()
	for i := 0; i < len(ub.params); i += 2 {
		q.Add(ub.params[i], ub.params[i+1])
	}
	uri.RawQuery = q.Encode()

	return uri, nil
}
