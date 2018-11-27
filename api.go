package api

import (
	"errors"
)

var (
	mainAPIURL string
)

// Init sets API url and authorization parameters
func Init(url string, username string, password string) {
	mainAPIURL = url
	SetBasicAuth(username, password)
}

// GetConfigInfo returns current configuration
func GetConfigInfo() (ConfigInfo, error) {
	var result ConfigInfo
	err := Get(mainAPIURL+"/v1/config", &result)
	return result, err
}

// GetSlaveList returns list of defined slave probes
func GetSlaveList() ([]string, error) {
	var result map[string]string
	err := Get(mainAPIURL+"/v1/slaves", &result)
	list := make([]string, len(result))
	for slave := range result {
		list = append(list, slave)
	}
	return list, err
}

// AddIP is simple interface for single IP adding
func AddIP(ip string, slaves []string, description string, groups []string, favourite bool) error {
	var r result

	err := Send("PUT", mainAPIURL+"/v1/config/ping/"+ip, TestDesc{
		Description: ip + " " + description,
		Favourite:   favourite,
		Groups:      groups,
		Slaves:      slaves,
	}, &r)
	if nil != err {
		return err
	}

	if r.Result != "OK" && "" != r.Error {
		return errors.New(r.Error)
	}

	if r.Result != "OK" {
		return errors.New("unknown error")
	}

	return nil
}

// AddIPs function adds multiply ips using only one API call
func AddIPs(ips []string, slaves []string, description string, groups []string, favourite bool) error {
	payload := make(map[string]TestDesc, len(ips))

	for _, ip := range ips {
		payload[ip] = TestDesc{
			Description: ip + " " + description,
			Favourite:   favourite,
			Groups:      groups,
			Slaves:      slaves,
		}
	}

	var r result

	err := Send("PUT", mainAPIURL+"/v1/mconfig/add", map[string]interface{}{
		"ips": payload,
	}, &r)
	if nil != err {
		return err
	}

	if r.Result != "OK" && "" != r.Error {
		return errors.New(r.Error)
	}

	if r.Result != "OK" {
		return errors.New("unknown error")
	}

	return nil

}

// DeleteIP removes one IP from cocopacket instance
func DeleteIP(ip string) error {
	var r result

	err := Send("DELETE", mainAPIURL+"/v1/config/ping/"+ip, nil, &r)
	if nil != err {
		return err
	}

	if r.Result != "OK" && "" != r.Error {
		return errors.New(r.Error)
	}

	if r.Result != "OK" {
		return errors.New("unknown error")
	}

	return nil
}