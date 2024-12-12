package validator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	// "github.com/go-playground/locales/en"
	// ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	// en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/rs/zerolog/log"
)

type Validator struct {
	// trans     ut.Translator
	validator *validator.Validate
}

func NewValidator() *Validator {
	validatorCustom := &Validator{}

	// en := en.New()
	// uni := ut.New(en, en)
	// trans, _ := uni.GetTranslator("en")

	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		var name string

		name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("params"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("prop"), ",", 2)[0]
		}

		if name == "-" {
			return ""
		}

		return name
	})

	// en_translations.RegisterDefaultTranslations(v, trans)
	if err := v.RegisterValidation("email_blacklist", isEmailBlacklist); err != nil {
		log.Fatal().Err(err).Msg("Error while registering email_blacklist validator")
	}
	if err := v.RegisterValidation("strong_password", isStrongPassword); err != nil {
		log.Fatal().Err(err).Msg("Error while registering strong_password validator")
	}
	if err := v.RegisterValidation("unique_in_slice", isUniqueInSlice); err != nil {
		log.Fatal().Err(err).Msg("Error while registering unique validator")
	}
	if err := v.RegisterValidation("phone", phoneValidator); err != nil {
		log.Fatal().Err(err).Msg("Error while registering phone validator")
	}
	if err := v.RegisterValidation("otp_number", otpNumberValidation); err != nil {
		log.Fatal().Err(err).Msg("Error while registering otp_number validator")
	}
	if err := v.RegisterValidation("nik", validateNIK); err != nil {
		log.Fatal().Err(err).Msg("Error while registering nik validator")
	}
	if err := v.RegisterValidation("birth_date", validateBirthDate); err != nil {
		log.Fatal().Err(err).Msg("Error while registering birth_date validator")
	}
	if err := v.RegisterValidation("salary", validateSalary); err != nil {
		log.Fatal().Err(err).Msg("Error while registering salary validator")
	}
	if err := v.RegisterValidation("file_path", validateFilePath); err != nil {
		log.Fatal().Err(err).Msg("Error while registering file_path validator")
	}
	if err := v.RegisterValidation("valid_text", validateLongText); err != nil {
		log.Fatal().Err(err).Msg("Error while registering valid_text validator")
	}

	validatorCustom.validator = v
	// validatorCustom.trans = trans

	return validatorCustom
}

func (v *Validator) Validate(i any) error {
	return v.validator.Struct(i)
}

// blacklist email validator
func isEmailBlacklist(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	disallowedDomains := []string{"outlook", "hotmail", "aol", "live", "inbox", "icloud", "mail", "gmx", "yandex"}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := strings.Split(parts[1], ".")[0]

	for _, disallowed := range disallowedDomains {
		if domain == disallowed {
			return false
		}
	}

	return true
}

func isStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUppercase := false
	hasLowercase := false
	hasNumber := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUppercase = true
		case char >= 'a' && char <= 'z':
			hasLowercase = true
		case char >= '0' && char <= '9':
			hasNumber = true
		}
	}

	return hasUppercase && hasLowercase && hasNumber
}

func isUniqueInSlice(fl validator.FieldLevel) bool {
	// Get the slice from the FieldLevel interface
	val := fl.Field()

	// Ensure the field is a slice
	if val.Kind() != reflect.Slice {
		return false
	}

	// Use a map to check for duplicates
	elements := make(map[interface{}]bool)
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i).Interface()
		if _, found := elements[elem]; found {
			return false // Duplicate found
		}
		elements[elem] = true
	}
	return true
}

func phoneValidator(fl validator.FieldLevel) bool {
	// +62, 62, 0
	phoneRegex := `^(?:\+62|62|0)[2-9][0-9]{7,14}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(fl.Field().String())
}

func otpNumberValidation(fl validator.FieldLevel) bool {
	otp := fl.Field().String()
	// Check if the OTP matches exactly 6-digit numeric pattern
	matched, _ := regexp.MatchString(`^\d{6}$`, otp)
	return matched
}

func validateNIK(fl validator.FieldLevel) bool {
	// Regular expression to ensure NIK is exactly 16 numeric digits
	numericRegex := regexp.MustCompile(`^\d{16}$`)

	// Get the field value as a string
	nik := fl.Field().String()

	// Validate against the regex
	return numericRegex.MatchString(nik)
}

func validateBirthDate(fl validator.FieldLevel) bool {
	// Regular expression to match the date format YYYY-MM-DD
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	// Get the field value as a string
	birthDate := fl.Field().String()

	// Check if the format matches
	if !dateRegex.MatchString(birthDate) {
		return false
	}

	// Parse the date to ensure it's valid
	_, err := time.Parse("2006-01-02", birthDate)
	return err == nil
}

func validateSalary(fl validator.FieldLevel) bool {
	// Regular expression to match a valid numeric value (integer or decimal)
	numericRegex := regexp.MustCompile(`^\d+(\.\d{1,2})?$`)

	// Get the field value as a string
	salary := fl.Field().String()

	// Check if the format matches
	if !numericRegex.MatchString(salary) {
		return false
	}

	// Convert the string to a float to check the value
	value, err := strconv.ParseFloat(salary, 64)
	if err != nil || value < 0 {
		return false
	}

	return true
}

func validateFilePath(fl validator.FieldLevel) bool {
	// Regular expression to validate a valid file path
	filePathRegex := regexp.MustCompile(`^[a-zA-Z0-9_\-/]+(\.[a-zA-Z0-9]{2,5})$`)

	// Get the field value as a string
	filePath := fl.Field().String()

	// Return true if the field matches the regex
	return filePathRegex.MatchString(filePath)
}

func validateLongText(fl validator.FieldLevel) bool {
	// Regex pattern to detect SQL Injection or JavaScript Injection attempts
	sqlInjectionPattern := regexp.MustCompile(`(?i)(select|insert|delete|drop|update|--|;|union|alert|<script|onerror)`)

	// Get the field value as a string
	text := fl.Field().String()

	// Return true if there's no match (valid), false if there's a match (invalid)
	return !sqlInjectionPattern.MatchString(text)
}
