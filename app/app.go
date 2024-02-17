package app

//layer to provide business logic to http requests

type App struct {
	Srv *Server
}

func New() *App {
	return &App{}
}
