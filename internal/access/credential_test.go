package access

import (
	"testing"
	"unicode"

	"gotest.tools/v3/assert"

	"github.com/infrahq/infra/internal/server/data"
	"github.com/infrahq/infra/internal/server/models"
)

func TestSettingsPasswordRequirements(t *testing.T) {
	c, db, _ := setupAccessTestContext(t)

	username := "bruce@example.com"
	user := &models.Identity{Name: username}
	err := data.CreateIdentity(db, user)
	assert.NilError(t, err)

	_, err = CreateCredential(c, *user)
	assert.NilError(t, err)

	err = data.SaveSettings(db, &models.Settings{
		LengthMin: 8,
	})
	assert.NilError(t, err)
	t.Run("Update user credentials fails if less than min length", func(t *testing.T) {
		err := UpdateCredential(c, user, "short")
		assert.ErrorContains(t, err, "does not pass requirements")
		assert.ErrorContains(t, err, "needs minimum length of 8")
	})

	// Test min length success
	settings, err := data.GetSettings(db)
	assert.NilError(t, err)
	settings.LengthMin = 5
	err = data.SaveSettings(db, settings)
	assert.NilError(t, err)
	t.Run("Update user credentials passes if equal than min length", func(t *testing.T) {
		err := UpdateCredential(c, user, "short")
		assert.NilError(t, err)
	})
	t.Run("Update user credentials passes if equal than min length", func(t *testing.T) {
		err := UpdateCredential(c, user, "longer")
		assert.NilError(t, err)
	})

	// Test multiple failures
	settings.LengthMin = 10
	settings.SymbolMin = 1
	err = data.SaveSettings(db, settings)
	assert.NilError(t, err)
	t.Run("Update user credentials fails with multiple requirement failures", func(t *testing.T) {
		err := UpdateCredential(c, user, "badpw")
		assert.ErrorContains(t, err, "does not pass requirements")
		assert.ErrorContains(t, err, "needs minimum 1 symbols")
		assert.ErrorContains(t, err, "needs minimum length of 10")
	})
}

func TestCreateCredential(t *testing.T) {
	c, db, _ := setupAccessTestContext(t)

	username := "bruce@example.com"
	user := &models.Identity{Name: username}
	err := data.CreateIdentity(db, user)
	assert.NilError(t, err)

	oneTimePassword, err := CreateCredential(c, *user)
	assert.NilError(t, err)
	assert.Assert(t, oneTimePassword != "")

	_, err = data.GetCredential(db, data.ByIdentityID(user.ID))
	assert.NilError(t, err)
}

func TestUpdateCredentials(t *testing.T) {
	c, db, _ := setupAccessTestContext(t)

	username := "bruce@example.com"
	user := &models.Identity{Name: username}
	err := data.CreateIdentity(db, user)
	assert.NilError(t, err)

	_, err = CreateCredential(c, *user)
	assert.NilError(t, err)

	t.Run("Update user credentials IS single use password", func(t *testing.T) {
		err := UpdateCredential(c, user, "newPassword")
		assert.NilError(t, err)

		creds, err := data.GetCredential(db, data.ByIdentityID(user.ID))
		assert.NilError(t, err)
		assert.Equal(t, creds.OneTimePassword, true)
	})

	t.Run("Update own credentials is NOT single use password", func(t *testing.T) {
		c.Set("identity", user)

		err := UpdateCredential(c, user, "newPassword")
		assert.NilError(t, err)

		creds, err := data.GetCredential(db, data.ByIdentityID(user.ID))
		assert.NilError(t, err)
		assert.Equal(t, creds.OneTimePassword, false)
	})
}

func TestLowercaseRequirements(t *testing.T) {
	result := hasMinimumCount(2, "a", unicode.IsLower)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "A", unicode.IsLower)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "ab", unicode.IsLower)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "AB", unicode.IsLower)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "Ab", unicode.IsLower)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "abc", unicode.IsLower)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "abC", unicode.IsLower)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "AbC", unicode.IsLower)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "aBc", unicode.IsLower)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "", unicode.IsLower)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "!$!@#23", unicode.IsLower)
	assert.Equal(t, result, false)
}

func TestUppercaseRequirements(t *testing.T) {
	result := hasMinimumCount(2, "a", unicode.IsUpper)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "A", unicode.IsUpper)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "ab", unicode.IsUpper)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "AB", unicode.IsUpper)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "Ab", unicode.IsUpper)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "abc", unicode.IsUpper)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "abC", unicode.IsUpper)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "AbC", unicode.IsUpper)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "aBc", unicode.IsUpper)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "", unicode.IsUpper)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "!$!@#23", unicode.IsUpper)
	assert.Equal(t, result, false)
}

func TestNumberRequirements(t *testing.T) {
	result := hasMinimumCount(2, "abc", unicode.IsNumber)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "aBc", unicode.IsNumber)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "", unicode.IsNumber)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "!$!@#", unicode.IsNumber)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "!$!@#23", unicode.IsNumber)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "!$!@#23123", unicode.IsNumber)
	assert.Equal(t, result, true)
}

func TestSymbolRequirements(t *testing.T) {
	result := hasMinimumCount(2, "", isValidSymbol)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "abAB", isValidSymbol)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "abc!", isValidSymbol)
	assert.Equal(t, result, false)

	result = hasMinimumCount(2, "  ", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "!!", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, `""`, isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, `##`, isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, `$$`, isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, `%%`, isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "&&", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "''", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "((", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "))", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "**", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "++", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, ",,", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "--", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "..", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "))", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "//", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "::", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, ";;", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "<<", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "==", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, ">>", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "??", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "@@", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "^^", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "__", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "{{", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "}}", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "||", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "~~", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, "~~", isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, `//`, isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, `\\`, isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, `[[`, isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, `]]`, isValidSymbol)
	assert.Equal(t, result, true)

	result = hasMinimumCount(2, `@$%@#ss`, isValidSymbol)
	assert.Equal(t, result, true)
}
