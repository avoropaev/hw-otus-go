//go:build integration

package integration_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/apitest"
)

type IntegrationSuite struct {
	apitest.APISuite
}

func (s *IntegrationSuite) SetupTest() {
	godotenv.Load("../../.env")

	apiURL := "http://" + os.Getenv("CALENDAR_HTTP_HOST") + ":" + os.Getenv("CALENDAR_HTTP_PORT")

	s.Init(apiURL)
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
