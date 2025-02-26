package password

// 检查密码长度是否在8-16位之间
func CheckPasswordLength(password string) bool {
	const minPasswordLength = 8
	const maxPasswordLength = 16
	return len(password) >= minPasswordLength && len(password) <= maxPasswordLength
}

// 检查密码复杂度是否满足要求
func CheckPasswordComplexity(password string) bool {
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char == '@' || char == '$' || char == '!' || char == '%' || char == '*' || char == '?' || char == '&':
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSpecial
}

// 检查邮箱格式是否正确
func IsValidEmail(email string) bool {
	// 简单的邮箱格式检查
	hasAt := false
	hasDot := false
	atIndex := -1

	for i, char := range email {
		if char == '@' {
			if hasAt {
				return false // 不能有多个@
			}
			hasAt = true
			atIndex = i
		} else if char == '.' {
			hasDot = true
		}
	}

	return hasAt && hasDot && atIndex > 0 && atIndex < len(email)-1
}
