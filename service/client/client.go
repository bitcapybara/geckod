package client

type ClientManager interface {
	Add(name string) (uint64, error)
	Get(name string) (*Client, error)
	Del(name string) error
}

type Client struct {
}
