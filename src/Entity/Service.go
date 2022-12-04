package Entity

type Service struct {
	Name         string
	Filename     string
	User         string
	Dll          string
	ServiceType  string
	StartType    string
	ErrorControl string
	Evidence     []string
}

var (
	ServiceTypeMap = map[int]string{4: "Adapter",
		1:   "KernelDriver",
		2:   "FileSystemDriver",
		8:   "RecognizerDriver",
		16:  "Win32OwnProcess",
		32:  "Win32ShareProcess",
		256: "InteractiveProcess",
	}
	StartTypeMap = map[int]string{
		0: "Boot",
		1: "System",
		2: "Automatic",
		4: "Disabled",
	}
	ErrorControlMap = map[int]string{
		0: "Ignore",
		1: "Normal",
		2: "Severe",
		3: "Critical",
	}
)

func containsService(services []Service, service Service) bool {
	for _, s := range services {
		if s.Name == service.Name && s.Filename == service.Filename {
			return true
		}
	}
	return false
}

func AddService(services []Service, service Service) []Service {
	if containsService(services, service) {
		services = append(services, service)
	}
	return services
}

func UnionServices(dest []Service, src []Service) []Service {
	for _, service := range src {
		dest = AddService(dest, service)
	}
	return dest
}

func NewService(pl PlasoLog) Service {
	var service Service

	service.Name = pl.Name
	service.Filename = pl.ImagePath
	service.ServiceType = ServiceTypeMap[pl.ServiceType]
	service.StartType = StartTypeMap[pl.StartType]
	service.ErrorControl = StartTypeMap[pl.ErrorControl]
	service.Dll = pl.ServiceDll
	service.Evidence = append(service.Evidence, pl.Message)
	service.User = pl.ObjectName

	return service
}
