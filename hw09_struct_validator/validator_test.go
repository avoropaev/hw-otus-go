package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36|skip:value"`
		Name   string
		Age    int      `validate:"min:18|max:50|skip:value"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
		private string `validate:"len:5"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	SliceOfInts struct {
		Value []int `validate:"min:1|max:100"`
	}

	WrongValidate struct {
		Value string `validate:"len:"`
	}

	WrongValidate2 struct {
		Value string `validate:":"`
	}

	WrongValidate3 struct {
		Value string `validate:"len"`
	}

	WrongValidate4 struct {
		Value int `validate:"min:"`
	}

	WrongValidate5 struct {
		Value int `validate:":"`
	}

	WrongValidate6 struct {
		Value int `validate:"min"`
	}

	WithStruct struct {
		Code int `validate:"in:200,404,500"`
		App  App `validate:""`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in                 interface{}
		expectedErr        error
		expectedErrMessage string
	}{
		{
			in:                 123,
			expectedErr:        ErrUnsupportedType,
			expectedErrMessage: ErrUnsupportedType.Error(),
		},
		{
			in:                 WrongValidate{},
			expectedErr:        ErrConstraintIsInvalid,
			expectedErrMessage: ErrConstraintIsInvalid.Error(),
		},
		{
			in:                 WrongValidate2{},
			expectedErr:        ErrConstraintIsInvalid,
			expectedErrMessage: ErrConstraintIsInvalid.Error(),
		},
		{
			in:                 WrongValidate3{},
			expectedErr:        ErrConstraintIsInvalid,
			expectedErrMessage: ErrConstraintIsInvalid.Error(),
		},
		{
			in:                 WrongValidate4{},
			expectedErr:        ErrConstraintIsInvalid,
			expectedErrMessage: ErrConstraintIsInvalid.Error(),
		},
		{
			in:                 WrongValidate5{},
			expectedErr:        ErrConstraintIsInvalid,
			expectedErrMessage: ErrConstraintIsInvalid.Error(),
		},
		{
			in:                 WrongValidate6{},
			expectedErr:        ErrConstraintIsInvalid,
			expectedErrMessage: ErrConstraintIsInvalid.Error(),
		},
		{
			in: App{
				Version: "a",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Version",
					Err:   ErrStringLengthIsInvalid,
				},
			},
			expectedErrMessage: "validation errors:\nVersion: string length is invalid",
		},
		{
			in: User{
				ID:    "",
				Name:  "",
				Age:   0,
				Email: "",
				Role:  "",
				Phones: []string{
					"1",
				},
				meta: nil,
			},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Err:   ErrStringLengthIsInvalid,
				},
				{
					Field: "Age",
					Err:   ErrIntLessThen(18),
				},
				{
					Field: "Email",
					Err:   ErrStringDontMatchRegexp,
				},
				{
					Field: "Role",
					Err:   ErrValueNotIn{"admin", "stuff"},
				},
				{
					Field: "Phones[0]",
					Err:   ErrStringLengthIsInvalid,
				},
			},
			expectedErrMessage: "validation errors:\nID: string length is invalid\nAge: value less then 18\nEmail: string dont match" +
				" regexp\nRole: String not in: admin,stuff\nPhones[0]: string length is invalid",
		},
		{
			in: User{
				ID:    "363636363636363636363636363636363636",
				Name:  "",
				Age:   55,
				Email: "email@gmail.com",
				Role:  "admin",
				Phones: []string{
					"12345678901",
				},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Age",
					Err:   ErrIntGreaterThen(50),
				},
			},
			expectedErrMessage: "validation errors:\nAge: value greater then 50",
		},
		{
			in: Response{
				Code: 0,
				Body: "",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Code",
					Err:   ErrValueNotIn{"200", "404", "500"},
				},
			},
			expectedErrMessage: "validation errors:\nCode: String not in: 200,404,500",
		},
		{
			in: SliceOfInts{
				Value: []int{-10, 0, 150},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Value[0]",
					Err:   ErrIntLessThen(1),
				},
				{
					Field: "Value[1]",
					Err:   ErrIntLessThen(1),
				},
				{
					Field: "Value[2]",
					Err:   ErrIntGreaterThen(100),
				},
			},
			expectedErrMessage: "validation errors:\nValue[0]: value less then 1\nValue[1]: value less then 1\n" +
				"Value[2]: value greater then 100",
		},
		{
			in: WithStruct{
				Code: 0,
				App: App{
					Version: "1",
				},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Code",
					Err:   ErrValueNotIn{"200", "404", "500"},
				},
				{
					Field: "App.Version",
					Err:   ErrStringLengthIsInvalid,
				},
			},
			expectedErrMessage: "validation errors:\nCode: String not in: 200,404,500\nApp.Version: string length is invalid",
		},
		{
			in: App{
				Version: "aaaaa",
			},
		},
		{
			in: User{
				ID:    "363636363636363636363636363636363636",
				Name:  "",
				Age:   25,
				Email: "email@gmail.com",
				Role:  "admin",
				Phones: []string{
					"12345678901",
				},
			},
		},
		{
			in: Response{
				Code: 404,
				Body: "",
			},
		},
		{
			in: SliceOfInts{
				Value: []int{1, 50, 100},
			},
		},
		{
			in: WithStruct{
				Code: 200,
				App: App{
					Version: "abcde",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)

			if tt.expectedErr != nil {
				require.Equal(t, err.Error(), tt.expectedErrMessage)
			}
		})
	}
}

func TestValidationErrorsError(t *testing.T) {
	err := ValidationErrors{}
	require.Equal(t, "", err.Error())
}
