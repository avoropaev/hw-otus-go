package apitest

import (
	"net/http"
)

func (s *APISuite) TestNotFound() {
	tests := []struct {
		path   string
		method string
	}{
		{
			path:   "/unknown",
			method: http.MethodGet,
		},
		{
			path:   "/unknown",
			method: http.MethodPost,
		},
	}

	for _, test := range tests {
		test := test
		s.Run(test.method+test.path, func() {
			request, err := http.NewRequestWithContext(s.ctx, test.method, s.apiURL+test.path, nil)
			s.Require().NoError(err)

			resp, err := http.DefaultClient.Do(request)
			s.Require().NoError(err)

			err = resp.Body.Close()
			s.Require().NoError(err)

			s.Require().Equal(http.StatusNotFound, resp.StatusCode)
		})
	}
}
