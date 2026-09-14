package std

type Realm interface {
	Address() Address
	PkgPath() string
	Previous() Realm
	IsCode() bool
	IsUser() bool
	IsUserCall() bool
	IsUserRun() bool
	IsEphemeral() bool
	IsCurrent() bool
	Sub(subpath string) Realm
	Subpath() string
	String() string
}

type Address string

func (a Address) String() string { return string(a) }
func (a Address) IsValid() bool  { return false }
