package interfaces

type Initialization interface {
	Initialized() bool
	Database() (bool, error)
	Seed() (bool, error)
}
