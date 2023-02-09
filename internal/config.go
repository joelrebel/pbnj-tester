package internal

type Config struct {
	Tests   []string `json:"tests"`
	Servers []struct {
		Name     string `json:"name"`
		Vendor   string `json:"vendor"`
		Model    string `json:"model"`
		BmcHost  string `json:"bmcHost" yaml:"bmcHost"` // yaml struct tags defined because the library requires it
		BmcUser  string `json:"bmcUser" yaml:"bmcUser"`
		BmcPass  string `json:"bmcPass" yaml:"bmcPass"`
		IpmiPort string `json:"ipmiPort" yaml:"IpmiPort"`
	} `json:"servers"`
}
