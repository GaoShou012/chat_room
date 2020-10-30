package models

type Rooms struct {
	Model
	Sort       *uint64 `json:"sort"`
	TenantCode *string `json:"tenantCode"`
	RoomCode   *string `json:"roomCode"`
	RoomName   *string `json:"roomName"`
}
