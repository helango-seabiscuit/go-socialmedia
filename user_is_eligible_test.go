package main

import (
	"errors"
	"testing"
)

func TestUserIsEligible(tst *testing.T) {
	var tests = []struct {
		email       string
		password    string
		age         string
		expectedErr error
	}{
		{
			email:       "test@example.com",
			password:    "12345",
			age:         "18",
			expectedErr: nil,
		},
		{
			email:       "",
			password:    "12345",
			age:         "18",
			expectedErr: errors.New("email can't be empty"),
		},
		{
			email:       "test@example.com",
			password:    "",
			age:         "18",
			expectedErr: errors.New("password can't be empty"),
		},
		{
			email:       "test@example.com",
			password:    "12345",
			age:         "16",
			expectedErr: errors.New("age must be at least 18 years old"),
		},
	}

	for _, t := range tests {
		err := userIsEligible(t.email, t.password, t.age)
		errString := ""
		expectedErrString := ""
		if err != nil {
			errString = err.Error()
		}

		if t.expectedErr != nil {
			expectedErrString = t.expectedErr.Error()
		}
		if errString != expectedErrString {
			tst.Errorf("got %s, expected %s ", errString, expectedErrString)
		}
	}
}
