package room

const(
	RoomEventJoin = iota
	RoomEventLeave
	RoomEventClientToDirectLine
)

type Event struct {
	Type   int
	Client *Client
}
