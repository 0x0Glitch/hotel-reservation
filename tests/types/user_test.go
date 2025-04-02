package types

import (
	"testing"

	"github.com/0x0Glitch/hotel-reservation/types"
)

// TestNewUserFromParams checks if user creation from params works correctly
// This validates that a new user is properly constructed with the right fields
func TestNewUserFromParams(t *testing.T) {
	// Set up test parameters
	params := types.CreateUserParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}

	// Create a new user from params
	user, err := types.NewUserFromParams(params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Validate all fields were set correctly
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstName %s, got %s", params.FirstName, user.FirstName)
	}

	if user.LastName != params.LastName {
		t.Errorf("expected lastName %s, got %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("expected email %s, got %s", params.Email, user.Email)
	}

	// Password should be encrypted, not stored as plain text
	if user.EncryptedPassword == params.Password {
		t.Errorf("password should be encrypted, not stored as plain text")
	}

	// Check that ID was generated
	if user.ID.IsZero() {
		t.Errorf("expected non-zero ID to be generated")
	}
}

// TestCreateUserParamsValidate checks all validation rules for user creation
func TestCreateUserParamsValidate(t *testing.T) {
	testCases := []struct {
		name          string
		params        types.CreateUserParams
		expectedError map[string]string
	}{
		{
			name: "valid params",
			params: types.CreateUserParams{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "password123",
			},
			expectedError: map[string]string{},
		},
		{
			name: "short firstName",
			params: types.CreateUserParams{
				FirstName: "J",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "password123",
			},
			expectedError: map[string]string{
				"firstName": "firstName length should be at least 2 characters",
			},
		},
		{
			name: "short lastName",
			params: types.CreateUserParams{
				FirstName: "John",
				LastName:  "D",
				Email:     "john@example.com",
				Password:  "password123",
			},
			expectedError: map[string]string{
				"lastName": "LastNamelength should be at least 2 characters",
			},
		},
		{
			name: "short password",
			params: types.CreateUserParams{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "pass",
			},
			expectedError: map[string]string{
				"password": "minimum password length should be at least 7 characters",
			},
		},
		{
			name: "invalid email",
			params: types.CreateUserParams{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "invalid-email",
				Password:  "password123",
			},
			expectedError: map[string]string{
				"email": "Email is invalid",
			},
		},
	}

	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errors := tc.params.Validate()
			
			if len(errors) != len(tc.expectedError) {
				t.Errorf("expected %d errors, got %d", len(tc.expectedError), len(errors))
			}

			// Check specific error messages
			for key, expectedMsg := range tc.expectedError {
				if actualMsg, ok := errors[key]; !ok || actualMsg != expectedMsg {
					t.Errorf("expected error for %s to be '%s', got '%s'", key, expectedMsg, actualMsg)
				}
			}
		})
	}
}

// TestIsEmailValid tests the email validation function with various email formats
func TestIsEmailValid(t *testing.T) {
	testCases := []struct {
		email    string
		expected bool
	}{
		{"john@example.com", true},           // Standard email
		{"john.doe@example.com", true},       // With dot in local part
		{"john+doe@example.com", true},       // With plus in local part
		{"john@example.co.uk", true},         // With multiple domain parts
		{"john@example", false},              // Missing TLD
		{"john@.com", false},                 // Missing domain
		{"john", false},                      // Missing @ and domain
		{"@example.com", false},              // Missing local part
	}

	for _, tc := range testCases {
		t.Run(tc.email, func(t *testing.T) {
			result := types.IsEmailValid(tc.email)
			if result != tc.expected {
				t.Errorf("isEmailValid(%s) = %v, expected %v", tc.email, result, tc.expected)
			}
		})
	}
}

// TestIsValidPassword verifies that password validation works correctly
func TestIsValidPassword(t *testing.T) {
	password := "password123"
	
	// Create a user with the test password
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  password,
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Test correct password
	if !types.IsValidPassword(user.EncryptedPassword, password) {
		t.Errorf("expected password validation to succeed for correct password")
	}

	// Test incorrect password
	if types.IsValidPassword(user.EncryptedPassword, "wrongpassword") {
		t.Errorf("expected password validation to fail for incorrect password")
	}
} 