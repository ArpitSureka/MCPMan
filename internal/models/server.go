package models

var (
	serverDB = ServerDB {
		// Currently ServerDB can hold only 2O MCP Servers running at MAX. Later shift 20 value inside 
		// config also make an algo to stop running MCP Servers if running servers starts crossing this value.
		Servers: make(chan Server, 20),
	}
)


func GetServerDB() *ServerDB {
	return &serverDB;
}

// ServerDB contains a channel of Running MCP Servers.
type ServerDB struct {
	Servers chan Server
}

type Server struct {
	ServerJSON
	Status   ServerStatus
}

// ServerType represents the type of server, such as Stdio or SSE.
type ServerType int

const (
	Stdio ServerType = iota
	SSE
)

// ServerStatus represents the current status of a server.
type ServerStatus int

const (
	Loaded ServerStatus = iota
	Starting
	Running 
	Stopped
	Paused
	Restarting
	Updating
	Failed
)
