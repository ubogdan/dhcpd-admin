package admin

var KnownHostTmpl = `
host 	{{.Host}}	{
  hardware ethernet	{{.MAC}};
  fixed-address 	{{.IP}};
}
`
var GlobalTmpl = `
option	 domain-name-servers  	{{.DNS1}},{{.DNS2}};
default-lease-time          	{{.Lease.Default}};
max-lease-time              	{{.Lease.Max}};
{{.Authoritative}};
`

var SubnetTmpl = `
subnet  {{.Subnet}} netmask {{.Netmask}}  {
  deny duplicates;
  pool  {
    range {{.Unknown.Initial}}  {{.Unknown.Final}};
    allow unknown-clients;
  }
  pool  {
    range {{.Known.Initial}}  {{.Known.Final}};
    deny  unknown-clients;
  }
}
`
