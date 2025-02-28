package azure

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func validateResourceGroupName(fl validator.FieldLevel) bool {
	// More info: https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules#microsoftresources
	// Alphanumerics, underscores, parentheses, hyphens, periods, and unicode characters that match the regex documentation.
	// Can't end with period.

	name := fl.Field().String()

	// Alphanumerics, underscores, parentheses, hyphens, periods
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_\(\)-\.]+$`, name)
	if !matched {
		return false
	}

	// Unicode characters that match the regex documentation
	// More info about regex pattern here: https://docs.microsoft.com/en-us/rest/api/resources/resourcegroups/createorupdate
	matched, _ = regexp.MatchString(`^[-\w\._\(\)]+$`, name)
	if !matched {
		return false
	}

	// Can't end with period
	if string(name[len(name)-1]) == "." {
		return false
	}

	return true
}

func validateStorageAccountContainerName(fl validator.FieldLevel) bool {
	// More info: https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules#microsoftstorage
	// Lowercase letters, numbers, and hyphens.
	// Start with lowercase letter or number. Can't use consecutive hyphens.

	name := fl.Field().String()

	// Lowercase letters, numbers, and hyphens
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, name)
	if !matched {
		return false
	}

	// Can't use consecutive hyphens
	matched, _ = regexp.MatchString(`[-]{2}`, name)
	if matched {
		return false
	}

	// Start with lowercase letter or number
	matched, _ = regexp.MatchString(`^[a-z0-9]$`, string(name[0]))
	if !matched {
		return false
	}

	return true
}

func validateKeyVaultName(fl validator.FieldLevel) bool {
	// More info: https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules#microsoftkeyvault
	// Alphanumerics and hyphens.
	// Start with letter. End with letter or digit. Can't contain consecutive hyphens.

	name := fl.Field().String()

	// Alphanumerics and hyphens
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9-]+$`, name)
	if !matched {
		return false
	}

	// Can't contain consecutive hyphens
	matched, _ = regexp.MatchString(`[-]{2}`, name)
	if matched {
		return false
	}

	// Start with letter
	matched, _ = regexp.MatchString(`^[a-zA-Z]$`, string(name[0]))
	if !matched {
		return false
	}

	// End with letter or digit
	matched, _ = regexp.MatchString(`^[a-zA-Z0-9]$`, string(name[len(name)-1]))
	if !matched {
		return false
	}

	return true
}

func validateKeyVaultKeyName(fl validator.FieldLevel) bool {
	// More info: https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules#microsoftkeyvault
	// Alphanumerics and hyphens.

	name := fl.Field().String()

	// Alphanumerics and hyphens
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9-]+$`, name)
	if !matched {
		return false
	}

	return true
}
