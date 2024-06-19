package main

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}

func GetUserFromData(userData map[string]any) User {
	id := userData["id"].(float64)

	return User{
		Id:        int(id),
		Username:  userData["login"].(string),
		AvatarUrl: userData["avatar_url"].(string),
	}
}
