package model

type SignedUserID struct {
	ID   int64
	HMAC string
}

type UserToAdd struct {
	Key string
}

type User struct {
	ID  int64
	Key string
}
