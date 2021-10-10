package model

type SignedUserID struct {
	ID   int
	HMAC []byte
}

type UserToAdd struct {
	Key []byte
}

type User struct {
	ID  int
	Key []byte
}
