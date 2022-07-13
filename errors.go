package cassinirpc

import (
	"errors"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func NewErrorWithDetails(code connect.Code, details string, extras ...proto.Message) error {
	err := connect.NewError(code, errors.New(details))
	for _, e := range extras {
		d, anyErr := anypb.New(e)
		if anyErr == nil {
			err.AddDetail(d)
		}
	}
	return err
}
