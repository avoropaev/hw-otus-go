package apitest

import (
	"encoding/json"
	"errors"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/gen/openapicli"
)

func (s *APISuiteActions) errBody(err error) string {
	var apiErr openapicli.GenericOpenAPIError

	if errors.As(err, &apiErr) {
		return string(apiErr.Body())
	}

	return ""
}

func (s *APISuiteActions) errMessage(err error) string {
	s.T().Helper()

	var apiErr openapicli.GenericOpenAPIError

	if !errors.As(err, &apiErr) {
		return ""
	}

	model := &struct {
		Code    int
		Message string
	}{}

	unmarshalErr := json.Unmarshal(apiErr.Body(), model)
	s.Require().NoError(unmarshalErr)

	return model.Message
}
