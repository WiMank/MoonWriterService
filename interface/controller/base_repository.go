package controller

type BaseRepository interface {
	OpenDataBase()
	CloseDataBase()
}
