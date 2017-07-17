package cache

type Driver interface {
	Meta() Meta

	Put(name string, contents interface{}) error
	Get(name string) (interface{}, error)
	Exists(name string) bool
	Remove(name string) error
}

/*
	Not using a registry system to give the user of the package the ability configure settings more easily

	var (
		drivers = make(map[string]Driver)
	)

	func Register(name string, driver Driver) {
		if driver == nil {
			panic("Cache: Driver is nil")
		}
		if _, dup := drivers[name]; dup {
			panic("Cache: Driver already registered " + name)
		}
		drivers[name] = driver
	}

	func Drivers() (names []string) {
		for name := range drivers {
			names = append(names, name)
		}
		return
	}
*/
