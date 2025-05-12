package config

const (
	LOCALS_USER        = "user"
	GROUP_CODE_CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type GroupType string

const (
	GroupTypePublic    GroupType = "public"
	GroupTypeProtected GroupType = "protected"
	GroupTypePrivate   GroupType = "private"
)
