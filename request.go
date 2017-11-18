package main

type Request struct {
	Id       string `json:"id"`
	Hostname string `json:"hostname"`
	Cluster  string `json:"cluster"`
	Template string `json:"template"`
	Ipv4     struct {
		Address string `json:"address"`
		Netmask string `json:"netmask"`
		Gateway string `json:"gateway,omitempty"`
	} `json:"ipv4,omitempty"`
	Ipv6 struct {
		Address string `json:"address"`
		Netmask string `json:"netmask"`
		Gateway string `json:"gateway,omitempty"`
	} `json:"ipv6,omitempty"`
}
