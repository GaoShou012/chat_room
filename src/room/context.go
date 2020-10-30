package room

type ConnContext struct {
	TenantCode string
	RoomCode   string
	UserType   string
	UserId     uint64
	UserName   string
	UserThumb  string
	UserTags   string
}
