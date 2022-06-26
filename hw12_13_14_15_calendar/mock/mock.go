package mock

//nolint:lll
//go:generate go run -mod=mod github.com/golang/mock/mockgen -package=mock -destination=application.go -build_flags=-mod=mod github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/app Application
//go:generate go run -mod=mod github.com/golang/mock/mockgen -package=mock -destination=storage.go -build_flags=-mod=mod github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/internal/app Storage
