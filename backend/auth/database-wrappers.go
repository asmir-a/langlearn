package auth

//low level drivers

func insertUser(username string, hash string, salt string) (err error) {
	return
}

func insertSession(username string) {
}

func deleteSession(username string) {
}

// todo: the following two functions should be grouped into one, which returns struct credentials {...}, which in turn should be in struct userInfo{...}
func getUsername(username string) string {
	return ""
}

func getUserSalt(username string) string {
	return ""
}

func getUserHash(username string) string {
	return ""
}
