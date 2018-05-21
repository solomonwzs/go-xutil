package ethernet

type NetworkData interface {
	Marshal() ([]byte, error)
}
