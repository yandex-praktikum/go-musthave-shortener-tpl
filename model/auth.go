package model

type SignedUserID struct {
	ID   int
	HMAC []byte
}

type User struct {
	ID  int
	Key []byte
}
