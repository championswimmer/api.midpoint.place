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

var DefaultGroupPlaceTypes = []PlaceType{
	PlaceTypeRestaurant,
	PlaceTypeBar,
	PlaceTypeCafe,
	PlaceTypePark,
}

func IsSupportedPlaceType(placeType PlaceType) bool {
	switch placeType {
	case PlaceTypeRestaurant, PlaceTypeBar, PlaceTypeCafe, PlaceTypePark, PlaceTypeMuseum, PlaceTypeBookstore:
		return true
	default:
		return false
	}
}
