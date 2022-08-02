package zerologwriter

import "github.com/rs/zerolog"

type ZerologWriter struct {
	Zerolog zerolog.Logger
}

func (zlw ZerologWriter) Write(p []byte) (n int, err error) {
	zlw.Zerolog.Error().Msg(string(p))

	return len(p), nil
}
