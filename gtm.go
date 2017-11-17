package bigip

import "encoding/json"
import "log"



type Datacenters struct {
	Datacenters []Datacenter `json:"items"`
}

type Datacenter struct {
	Name  string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Contact string `json:"contact,omitempty"`
	App_service  string `json:"appService,omitempty"`
	Disabled  bool `json:"disabled,omitempty"`
	Enabled  bool `json:"enabled,omitempty"`
	Prober_pool  string `json:"proberPool,omitempty"`
}


type Gtmmonitors struct {
	Gtmmonitors []Gtmmonitor `json:"items"`
}


type Gtmmonitor struct {
	Name  string `json:"name,omitempty"`
	Defaults_from string `json:"defaultsFrom,omitempty"`
	Interval int `json:"interval,omitempty"`
	Probe_timeout  int `json:"probeTimeout,omitempty"`
	Recv  string `json:"recv,omitempty"`
	Send  string `json:"send,omitempty"`
}

type Servers struct {
	Servers []Server `json:"items"`
}

type Server struct {
	Name      string
	Datacenter string
	Monitor  string
	Virtual_server_discovery bool
	Product  string
	Addresses     []ServerAddresses
	GTMVirtual_Server []VSrecord
}

type serverDTO struct {
	Name      string   `json:"name"`
	Datacenter string  `json:"datacenter,omitempty"`
	Monitor  string   `json:"monitor,omitempty"`
	Virtual_server_discovery bool `json:"virtual_server_discovery"`
	Product  string  `json:"product,omitempty"`
	Addresses     struct {
		Items []ServerAddresses`json:"items,omitempty"`
	} `json:"addressesReference,omitempty"`
	GTMVirtual_Server struct {
		Items []VSrecord `json:"items,omitempty"` } `json:"virtualServersReference,omitempty"`
	}

	func (p *Server) MarshalJSON() ([]byte, error) {
		return json.Marshal(serverDTO{
			Name:      p.Name,
			Datacenter: p.Datacenter,
			Monitor:  p.Monitor,
			Virtual_server_discovery:  p.Virtual_server_discovery,
			Product:  p.Product,
			Addresses: struct {
				Items []ServerAddresses `json:"items,omitempty"`
			}{Items: p.Addresses},
			GTMVirtual_Server: struct {
				Items []VSrecord `json:"items,omitempty"`
			}{Items: p.GTMVirtual_Server},
		})
	}

	func (p *Server) UnmarshalJSON(b []byte) error {
		var dto serverDTO
		err := json.Unmarshal(b, &dto)
		if err != nil {
			return err
		}

		p.Name = dto.Name
		p.Datacenter = dto.Datacenter
		p.Monitor = dto.Monitor
		p.Virtual_server_discovery = dto.Virtual_server_discovery
		p.Product = dto.Product
		p.Addresses = dto.Addresses.Items
		p.GTMVirtual_Server = dto.GTMVirtual_Server.Items
		return nil
	}


	type ServerAddressess struct {
		Items []ServerAddresses `json:"items,omitempty"`
	}



	type ServerAddresses struct {
		Name       string `json:"name"`
		Device_name   string `json:"deviceName,omitempty"`
		Translation  string `json:"translation,omitempty"`
	}



	type VSrecords struct {
		Items []VSrecord `json:"items,omitempty"`
	}



	type VSrecord struct {
		Name       string `json:"name"`
		Destination   string `json:"destination,omitempty"`
	}





const (
	uriGtm       = "gtm"
	uriServer    = "server"
	uriDatacenter = "datacenter"
	uriGtmmonitor    = "monitor"
	uriHttp       = "http"
)

func (b *BigIP) Datacenters() (*Datacenter, error) {
	var datacenter Datacenter
	err, _ := b.getForEntity(&datacenter, uriGtm, uriDatacenter)

	if err != nil {
		return nil, err
	}

	return &datacenter, nil
}

func (b *BigIP) CreateDatacenter(name, description,contact, app_service string, enabled, disabled bool, prober_pool string) error {
	config := &Datacenter{
		Name:    name,
		Description:  description,
		Contact:  contact,
		App_service: app_service,
		Enabled: enabled,
		Disabled: disabled,
		Prober_pool: prober_pool,
	}
	log.Printf("I am %#v\n  here %s   ", config)
	return b.post(config, uriGtm, uriDatacenter)
}

func (b *BigIP) ModifyDatacenter(*Datacenter) error {
	return b.patch(uriGtm, uriDatacenter)
}


func (b *BigIP) DeleteDatacenter(name string) error {
	return b.delete(uriGtm, uriDatacenter, name)
}



func (b *BigIP) Gtmmonitors() (*Gtmmonitor, error) {
	var gtmmonitor Gtmmonitor
	err, _ := b.getForEntity(&gtmmonitor, uriGtm, uriGtmmonitor, uriHttp)

	if err != nil {
		return nil, err
	}

	return &gtmmonitor, nil
}
func (b *BigIP) CreateGtmmonitor(name, defaults_from string, interval, probeTimeout int, recv, send string) error {
	config := &Gtmmonitor{
		Name:    name,
		Defaults_from:  defaults_from,
		Interval: interval,
		Probe_timeout: probeTimeout,
		Recv: recv,
		Send: send,
	}
	log.Printf("I am %#v\n  here %s   ", config)
	return b.post(config, uriGtm, uriGtmmonitor, uriHttp)
}

func (b *BigIP) ModifyGtmmonitor(*Gtmmonitor) error {
	return b.patch(uriGtm, uriGtmmonitor, uriHttp)
}


func (b *BigIP) DeleteGtmmonitor(name string) error {
	return b.delete(uriGtm, uriGtmmonitor, uriHttp, name)
}


func (b *BigIP) CreateGtmserver(p *Server) error {
	log.Println(" what is the complete payload    ", p)
	return b.post(p, uriGtm, uriServer)
}

//Update an existing policy.
func (b *BigIP) UpdateGtmserver(name string, p *Server) error {
	return b.put(p, uriGtm, uriServer, name)
}

//Delete a policy by name.
func (b *BigIP) DeleteGtmserver(name string) error {
	return b.delete(uriGtm, uriServer, name)
}


func (b *BigIP) GetGtmserver(name string) (*Server, error) {
 var p Server
 err, ok := b.getForEntity(&p, uriGtm, uriServer, name)
 if err != nil {
	 return nil, err
 }
 if !ok {
	 return nil, nil
 }

 return &p, nil
 }