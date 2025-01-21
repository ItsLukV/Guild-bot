package models

type User struct {
	ID                string `json:"id"`
	ActiveProfileUUID string `json:"active_profile_UUID"`
	FetchData         bool   `json:"fetch_data"`
}

type UsersResponse struct {
	Users []User `json:"users"`
}
