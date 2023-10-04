package user

func GetUser(userId string) (User, error) {
	return GetUserById(userId)
}
