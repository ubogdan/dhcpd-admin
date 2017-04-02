package admin

type HostConfig struct {
	Host string
	MAC  string
	IP   string
}

type Lease struct {
	Default int
	Max     int
}

type GlobalConfig struct {
	DNS1, DNS2 string
	//	Router        string
	Lease         Lease
	Authoritative string
}

type SubnetConfig struct {
	Subnet  string
	Netmask string
	Known   IPRange
	Unknown IPRange
}

type IPRange struct {
	Initial string
	Final   string
}
