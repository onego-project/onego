package services

import (
	"context"
	"fmt"

	"github.com/beevik/etree"
	"github.com/onego-project/xmlrpc"
)

// Service structure to make XML-RPC call
type Service struct {
	RPC *RPC
}

// RPC structure represents XML-RPC client and token string
type RPC struct {
	Client *xmlrpc.Client
	Token  string
}

// enum of result array index
const (
	successIndex = iota
	resultIndex
	errorCodeIndex
	idObjectCausedErrorIndex
)

// UpdateType for replace or merge resource template contents
type UpdateType int

const (
	// Replace the whole template
	Replace UpdateType = iota
	// Merge new template with the existing one
	Merge
)

func (s *Service) call(ctx context.Context, methodName string, args ...interface{}) ([]*xmlrpc.Result, error) {
	allArgs := append([]interface{}{s.RPC.Token}, args...)

	result, err := s.RPC.Client.Call(ctx, methodName, allArgs...)
	if err != nil {
		return nil, err
	}

	resArr := result.ResultArray()
	if !resArr[successIndex].ResultBoolean() {
		if len(resArr) == 4 {
			return nil, fmt.Errorf("%s, error code: %d, id of the object that caused the error %d",
				resArr[resultIndex].ResultString(), resArr[errorCodeIndex].ResultInt(),
				resArr[idObjectCausedErrorIndex].ResultInt())
		}
		return nil, fmt.Errorf("%s, code: %d", resArr[resultIndex].ResultString(),
			resArr[errorCodeIndex].ResultInt())
	}

	return resArr, nil
}

// retrieveInfo retrieves information for the given object
func (s *Service) retrieveInfo(ctx context.Context, methodName string, objectID int) (*etree.Document, error) {
	resArr, err := s.call(ctx, methodName, objectID)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	err = doc.ReadFromString(resArr[resultIndex].ResultString())
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *Service) list(ctx context.Context, methodName string) (*etree.Document, error) {
	resArr, err := s.call(ctx, methodName)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromString(resArr[resultIndex].ResultString()); err != nil {
		return nil, err
	}

	return doc, nil
}
