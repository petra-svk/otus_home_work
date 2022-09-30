package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
		test        string
	}{
		{
			in: User{
				ID:     "324324212312312432435435323132133122",
				Name:   "Kate",
				Age:    23,
				Email:  "qwer@ty.ru",
				Role:   "stuff",
				Phones: []string{"79981112233", "79973455677"},
				meta:   json.RawMessage(`{"is_good": true}`),
			},
			expectedErr: nil,
			test:        "all good",
		},
		{
			in:          "hello",
			expectedErr: ErrNotStruct,
			test:        "not a struct",
		},
		{
			in: User{
				ID:     "32432421231231243243543532313213312",
				Name:   "Kate",
				Age:    23,
				Email:  "qwer@ty.ru",
				Role:   "stuff",
				Phones: []string{"79981112233", "79973455677"},
				meta:   json.RawMessage(`{"is_good": true}`),
			},
			expectedErr: ErrLen,
			test:        "wrong len ID",
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
			test:        "good version",
		},
		{
			in: App{
				Version: "124324345",
			},
			expectedErr: ErrLen,
			test:        "bad version",
		},
		{
			in: Token{
				Header:    []byte("Your String"),
				Payload:   []byte("Your String"),
				Signature: []byte("Your String"),
			},
			expectedErr: nil,
			test:        "nothing check",
		},
		{
			in: Response{
				Code: 301,
				Body: "{ 'msg': 'Hello'}",
			},
			expectedErr: ErrInclude,
			test:        "bad include value",
		},
		{
			in: User{
				ID:     "3243242123123124324354353231321331",
				Name:   "Kate",
				Age:    15,
				Email:  "qwer@ty.ru",
				Role:   "stuff",
				Phones: []string{"79981112233", "79973455677"},
				meta:   json.RawMessage(`{"is_good": true}`),
			},
			expectedErr: ErrMin,
			test:        "error min",
		},
		{
			in: User{
				ID:     "3243242123123124324354353231321331",
				Name:   "Kate",
				Age:    15,
				Email:  "324234",
				Role:   "stuff",
				Phones: []string{"79981112233", "79973455677"},
				meta:   json.RawMessage(`{"is_good": true}`),
			},
			expectedErr: ErrRegexp,
			test:        "error email regexp",
		},
		{
			in: User{
				ID:     "3243242123123124324354353231321331",
				Name:   "Kate",
				Age:    15,
				Email:  "324234",
				Role:   "ewqeq",
				Phones: []string{"79981112233", "79973455677"},
				meta:   json.RawMessage(`{"is_good": true}`),
			},
			expectedErr: ErrRegexp,
			test:        "error role value",
		},
		{
			in: User{
				ID:     "3243242123123124324354353231321331",
				Name:   "Kate",
				Age:    15,
				Email:  "324234",
				Role:   "ewqeq",
				Phones: []string{"799811122", "79973455677"},
				meta:   json.RawMessage(`{"is_good": true}`),
			},
			expectedErr: ErrLen,
			test:        "error phone length in slice",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("case %s", tt.test), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			switch {
			case tt.expectedErr == nil:
				require.Empty(t, err)
			case errors.Is(err, tt.expectedErr):
				require.Error(t, err)
			default:
				var e *ValidationErrors
				if errors.As(err, &e) {
					foundError := false
					for _, ve := range *e {
						if errors.Is(&ve, tt.expectedErr) {
							foundError = true
						}
					}
					if !foundError {
						t.Errorf("AHTUNG: can't find expected error %v in slice %v\n", tt.expectedErr, err)
					}
				} else {
					t.Error("AHTUNG: got unexpected error")
				}
			}
		})
	}
}
