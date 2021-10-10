package model

type SignedUserID struct {
	ID   int
	HMAC string
}

type UserToAdd struct {
	Key string
}

type User struct {
	ID  int
	Key string
}
