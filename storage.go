package scrapit


type Storage interface {
	Process()
	Initialize() (chan *Scrapit, error)
	Uname() (string)
}
