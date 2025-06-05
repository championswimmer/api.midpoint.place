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

type GroupUserRole string

const (
	GroupUserAdmin  GroupUserRole = "admin"
	GroupUserMember GroupUserRole = "member"
)

type PlaceType string

const (
	PlaceTypeRestaurant PlaceType = "restaurant"
	PlaceTypeBar        PlaceType = "bar"
	PlaceTypeCafe       PlaceType = "cafe"
	PlaceTypePark       PlaceType = "park"
	PlaceTypeMuseum     PlaceType = "museum"
	PlaceTypeBookstore  PlaceType = "bookstore"
)

var AllPlaceTypes = []PlaceType{
	PlaceTypeRestaurant,
	PlaceTypeBar,
	PlaceTypeCafe,
	PlaceTypePark,
	PlaceTypeMuseum,
	PlaceTypeBookstore,
}
