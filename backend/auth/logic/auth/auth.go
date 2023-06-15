package auth

//todo: this must reside somewhere else
//maybe logic/validation
//or logic/utilities

func validateUsername(username string) bool {
	//may be better to have a separate file for validation
	if username == "" {
		return false
	}
	//todo: other checks to prevent sql injection eg
	return true
}

func validatePassword(password string) bool {
	if password == "" {
		return false
	}

	return true
}
