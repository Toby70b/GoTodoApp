package models

type Response struct {
	responseCode int
	body         any
}

func NewResponse(responseCode int, body any) Response {
	return Response{responseCode: responseCode, body: body}
}

func (r *Response) ResponseCode() int {
	return r.responseCode
}

func (r *Response) SetResponseCode(responseCode int) {
	r.responseCode = responseCode
}

func (r *Response) Body() any {
	return r.body
}

func (r *Response) SetBody(body any) {
	r.body = body
}
