package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func userIsEligible(email, password, age string) error {
	if strings.Trim(email, " ") == "" {
		return errors.New("email can't be empty")
	}

	if password == "" {
		return errors.New("password can't be empty")
	}

	if ageInput, err := strconv.Atoi(age); err != nil || ageInput < 18 {
		return fmt.Errorf("age must be at least %v years old", 18)
	}
	return nil

}
