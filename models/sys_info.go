package models

type SysInfo struct {
	Hostname string
	Platform string
	CPU      string
	RAM      uint64
	Disk     uint64
	Free     uint64
}
