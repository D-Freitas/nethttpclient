package nethttpclient

import "context"

type Http interface {
	Dispatch(
		ctx context.Context,
		response any,
		request any,
		customHeaders map[string]string,
	) (*Response, error)
}

type Request struct {
	Url    string
	Method string
	Body   string
}

type Response struct {
	Body   string
	Status int
}
