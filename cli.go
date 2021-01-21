package jcli

import "reflect"

type IClient interface {
	Call(operName string, operReq IRequest, resType reflect.Type) (operRes IResponse, err error)
}

type IRequest interface {
}

type IRequestValidator interface {
	IRequest
	Validate() error
}

type IResponse interface{}
