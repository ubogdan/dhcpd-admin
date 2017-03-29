package admin

var knownHost = `
host {{.Host}}  {
  hardware ethernet {{.MAC}};
  fixed-address {{.IP}};
}
`
var globalConf = `
option domain-name-servers  {{.DNS1}} {{.DNS2}};
option routers              {{.Router}};
default-lease-time          {{.Lease.Default}};
max-lease-time              {{.Lease.Max}};
{{.Authoritative}};
`

var subnetConf = `
subnet  {{.Subnet}} netmask {{.Netmask}}  {
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
