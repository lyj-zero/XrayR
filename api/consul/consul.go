package consul

import (
	"encoding/json"
	b64 "encoding/base64"
	"log"
	"fmt"
	"net"
	"errors"
	"time"

	"github.com/XrayR-project/XrayR/api"
	//"github.com/XrayR-project/XrayR/service/controller"
	"github.com/go-resty/resty/v2"
)


type APIClient struct {
	client              *resty.Client
	Rule                string
	Src                 string
	Url                 string
	NodeID              int
	Tls                 bool
	Date                []DatePut
}

// New creat a api instance
//func New(apiConfig *api.Config, nodeid int, host string, port int) *APIClient {
func New(timeout int, consulhost string, nodeInfo *api.NodeInfo) *APIClient {
	client := resty.New()
	client.SetRetryCount(3)
	if timeout > 0 {
		client.SetTimeout(time.Duration(timeout) * time.Second)
	} else {
		client.SetTimeout(5 * time.Second)
	}
	client.OnError(func(req *resty.Request, err error) {
		if v, ok := err.(*resty.ResponseError); ok {
			// v.Response contains the last response from the server
			// v.Err contains the original error
			log.Print(v.Err)
		}
	})
	client.SetBaseURL(consulhost)

	// Create Key for each requests
	//client.SetQueryParam("key", apiConfig.Key)
	// Add support for muKey
	//client.SetQueryParam("muKey", apiConfig.Key)
	// Read local rule list
	
	//Get Node IP
	ip, err := getClientIp()
	if err != nil {
		log.Println(err)
	}

	return &APIClient{
		client:              client,
		Rule:                fmt.Sprintf("Host(`%s`)", nodeInfo.Host),
		Src:                 fmt.Sprintf("v2n%d-src", nodeInfo.NodeID),
		Url:                 fmt.Sprintf("http://%s:%d", ip, nodeInfo.Port),
		NodeID:              nodeInfo.NodeID,
		Tls:                 true,
		//Date:                []DatePut,
	}
}

func (c *APIClient) AddKV(key string, value string) {
	c.Date = append(c.Date, DatePut { KVPut {"set" , key, b64.StdEncoding.EncodeToString([]byte(value)), 0, 0, ""} })
}

func (c *APIClient) Post() {
	c.Update()
	date, err := json.Marshal(c.Date)
	if err != nil {
		log.Printf("Error : %s", err)
	}

	path := "v1/txn"

	res, err := c.client.R().
		SetHeader("Content-Type" , "application/json").
		SetResult(&Response{}).
		SetBody(date).
		ForceContentType("application/json").
		Put(path)

	if err != nil {
		log.Printf("Error : %s", err)
	}

	response := res.Result().(*Response)
	if response.Errors != nil {
		log.Printf("Error : %s", response.Errors)
	}

	log.Print("Info : consul update")
	//log.Println(res)
	//log.Println(string(date))
}

func (c *APIClient) Update() {

	c.AddKV(fmt.Sprintf("traefik/http/routers/v2n%d/rule", c.NodeID), c.Rule)
	c.AddKV(fmt.Sprintf("traefik/http/routers/v2n%d1/rule", c.NodeID), c.Rule)
	c.AddKV(fmt.Sprintf("traefik/http/routers/v2n%d/service", c.NodeID), c.Src)
	c.AddKV(fmt.Sprintf("traefik/http/routers/v2n%d1/service", c.NodeID), c.Src)
	c.AddKV(fmt.Sprintf("traefik/http/services/%s/loadbalancer/servers/%d/url", c.Src, 0), c.Url)
	c.AddKV(fmt.Sprintf("traefik/http/routers/v2n%d/entrypoints", c.NodeID), "websecure")
	c.AddKV(fmt.Sprintf("traefik/http/routers/v2n%d1/entrypoints", c.NodeID), "web")
	if c.Tls {
		c.AddKV(fmt.Sprintf("traefik/http/routers/v2n%d/tls", c.NodeID), "true")
		c.AddKV(fmt.Sprintf("traefik/http/routers/v2n%d/tls/certResolver", c.NodeID), "vsspro")
	}
	c.AddKV(fmt.Sprintf("traefik/http/routers/v2n%d1/middlewares", c.NodeID), "traefik-scheme@docker")

}

func getClientIp() (string ,error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), err
			}

		}
	}

	return "", errors.New("Can not find the client ip address!")

}