package main

/*
	A simple Continuum Echo Service Gateway, without any persistant storage.
	@author: Zia Syed <zia.syed@ericsson.com>

 	Requirements:
 	1 - A simple web app that sends back the request URL
 	  - Step 2 performed on that app
 		- App also listens at PORT 3000
 	2 - A consumer application to test service gateway
 	  - Simplest to get a capsule to run and check environment variables in it
 	  - Create capsule
 	  	- apc capsule create mycapp --image linux -ae
 	  	- Use mycapp in Step 5
 	  	- To verify,
 	  	 	- apc capsule connect mycapp
 	  	 	- shell> env
 	  	 	- look for environment variable name s-0 and curl on it

  Steps:
  1- Create service gateway
  	a) apc app create echo-sg --disable-routes
  	b) apc gateway promote echo-sg --type echosg
  2- Create service provider
  	a) apc app create echo-server --disable-routes
    b) apc app update echo-server --port-add 3000
    c) apc app start echo-server
  3- Register provider
    a) apc provider register echo --type echosg --job echo-server -port 3000 --url http://user:pass@example.com/ping
	4- Create Service
		a) apc service create s-0 --provider echo
	5- Bind Service
	  a) apc service bind s-0 -j mycapp
*/

import (
	"encoding/json"
	"github.com/drone/routes"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

//These are specific to the services we are exposing.
var apis = []byte(`{
  "/bindings": {
    "GET": [{}],
    "POST": {}
  },
  "/bindings/:id": {
    "GET": {},
    "DELETE": {}
  },
  "/providers": {
    "GET": [{}],
    "POST": {"params": {"url": "user:pass@http://example.com:3000/"}}
  },
  "/providers/:id": {
    "GET": {"params": {"url": "user:pass@http://example.com:3000/"}},
    "DELETE": {}
  },
  "/services": {
    "GET": [{}],
    "POST": {}
  },
  "/services/:id": {
    "GET": {"params": {"test_key": "test_val"}},
    "DELETE": {}
  }
}`)

type Provider struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Id          string            `json:"id"`
	Type        string            `json:"type"`
	Params      map[string]string `json:"params"`
}

type Service struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Id          string            `json:"id"`
	ProviderID  string            `json:"provider_id"`
	Params      map[string]string `json:"params"`
}

type Binding struct {
	Name        string            `json:"name"`
	Id          string            `json:"id"`
	ServiceID   string            `json:"service_id"`
	Url         string            `json:"url"`
	UrlTemplate string            `json:"url_template"`
	Protocol    map[string]string `json:"protocol"`
}

var Providers []Provider
var Bindings []Binding
var Services []Service

func init() {
}

//curl -v http://localhost:3000/
func rootHandler(rw http.ResponseWriter, req *http.Request) {
	data := make(map[string]interface{})
	json.Unmarshal(apis, &data)
	res, _ := json.Marshal(data)
	io.WriteString(rw, string(res))
}

//curl -v http://localhost:3000/bindings
func getAllBindingsHandler(rw http.ResponseWriter, req *http.Request) {
	response, _ := json.Marshal(Bindings)
	log.Println("Returning:", string(response))
	io.WriteString(rw, string(response))
}

//curl -i -X POST -H "Content-Type: application/json" -d '{"service_id":"7cf096d"}' "http://user:pass@localhost:3000/bindings"
func addBindingsHandler(rw http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	log.Println("addBindingsHandler created:", params.Get(":id"))
	defer req.Body.Close()
	postBody, _ := ioutil.ReadAll(req.Body)
	targetBinding := Binding{}
	json.Unmarshal(postBody, &targetBinding)
	log.Println("Received", targetBinding)
	found := false
	var service Service
	var provider Provider
	for _, v := range Services {
		if v.Id == targetBinding.ServiceID {
			service = v
			found = true
		}
	}
	log.Println("ServiceID Found?", found)

	if !found {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	found = false
	for _, v := range Providers {
		if v.Id == service.ProviderID {
			provider = v
			found = true
		}
	}
	log.Println("ProviderID Found?", found)

	if !found {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	binding := Binding{}
	binding.ServiceID = service.Id
	err := json.Unmarshal(postBody, &binding)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if binding.Id == "" {
		binding.Id = uuid.NewV4().String()
	}
	binding.Name = "SAPC binding to " + service.Name
	binding.Url = provider.Params["url"] + "/"
	binding.Protocol = map[string]string{"scheme": "http"}
	Bindings = append(Bindings, binding)
	response, _ := json.Marshal(binding)
	io.WriteString(rw, string(response))
}

//curl -v http://localhost:3000/bindings/12
func getBindingsHandler(rw http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	for _, v := range Bindings {
		if v.Id == params.Get(":id") {
			response, _ := json.Marshal(v)
			io.WriteString(rw, string(response))
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
}

//curl -v -X DELETE http://localhost:3000/bindings/12
func delBindingsHandler(rw http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	log.Println("delBindingsHandler deleted:", params.Get(":id"))
	for index, v := range Bindings {
		if v.Id == params.Get(":id") {
			Bindings = append(Bindings[:index], Bindings[index+1:]...)
			rw.WriteHeader(http.StatusOK)
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
}

//curl -v http://localhost:3000/services
func getAllServicesHandler(rw http.ResponseWriter, req *http.Request) {
	response, _ := json.Marshal(Services)
	log.Println("Returning:", string(response))
	io.WriteString(rw, string(response))
}

//curl -v http://localhost:3000/services/12
func getServicesHandler(rw http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	for _, v := range Services {
		if v.Id == params.Get(":id") {
			response, _ := json.Marshal(v)
			log.Println("Replying:", string(response))
			io.WriteString(rw, string(response))
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
}

//curl -i -X POST -H "Content-Type: application/json" -d '{"provider_id":"1d0223f","params":{"database":"an_example_db"}}' "http://localhost:3000/services"
func addServicesHandler(rw http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	log.Println("addServicesHandler created:", params.Get(":id"))
	defer req.Body.Close()
	postBody, _ := ioutil.ReadAll(req.Body)
	service := Service{}
	err := json.Unmarshal(postBody, &service)
	log.Println("Received", service)

	found := false
	for _, v := range Providers {
		if v.Id == service.ProviderID {
			found = true
		}
	}
	log.Println("ProviderID", service.ProviderID, "found?", found)

	if !found {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if service.Id == "" {
		service.Id = uuid.NewV4().String()
	}
	if len(service.Name) < 1 {
		service.Name = "echo-" + strconv.Itoa(len(Services))
	}
	service.Description = "Echo service at " + req.Host
	service.Params = map[string]string{"db": "0"}
	Services = append(Services, service)
	response, _ := json.Marshal(service)
	log.Println("Service create", string(response))
	rw.WriteHeader(http.StatusCreated)
	io.WriteString(rw, string(response))
}

//curl -v -X DELETE http://localhost:3000/services/12
func delServicesHandler(rw http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	log.Println("delServicesHandler deleted:", params.Get(":id"))
	for index, v := range Services {
		if v.Id == params.Get(":id") {
			Services = append(Services[:index], Services[index+1:]...)
			rw.WriteHeader(http.StatusOK)
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
}

//curl -v http://localhost:3000/providers
func getAllProvidersHandler(rw http.ResponseWriter, req *http.Request) {
	response, _ := json.Marshal(Providers)
	log.Println("Returning:", string(response))
	io.WriteString(rw, string(response))
}

//curl -v http://localhost:3000/providers/12
func getProvidersHandler(rw http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	for _, v := range Providers {
		if v.Id == params.Get(":id") {
			response, _ := json.Marshal(v)
			log.Println("Replying:", string(response))
			io.WriteString(rw, string(response))
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
}

//curl -i -X POST -H "Content-Type: application/json" -d '{"name":"provider_name","type":"echo","params": {"url": "mysql://mysql:mysql@localhost:3306/"}}' "http://localhost:3000/providers"
func addProvidersHandler(rw http.ResponseWriter, req *http.Request) {
	log.Println("addProvidersHandler received:")
	defer req.Body.Close()
	postBody, _ := ioutil.ReadAll(req.Body)
	provider := Provider{}
	err := json.Unmarshal(postBody, &provider)
	log.Println("Received", provider)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if provider.Id == "" {
		provider.Id = uuid.NewV4().String()
	}
	provider.Type = "hello"
	provider.Name = "hello-service-" + req.Host
	provider.Description = "hello service provider at " + req.Host
	Providers = append(Providers, provider)

	response, _ := json.Marshal(provider)
	log.Println("Provider created", string(response))

	io.WriteString(rw, string(response))
}

//curl -v -X DELETE http://localhost:3000/providers/12
func delProvidersHandler(rw http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	log.Println("delProvidersHandler deleted:", params.Get(":id"))
	for index, v := range Providers {
		if v.Id == params.Get(":id") {
			Providers = append(Providers[:index], Providers[index+1:]...)
			rw.WriteHeader(http.StatusOK)
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
}

func setContentHeader(rw http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, "] Request:", req.Host, req.URL.String(), req.URL.Query())
	rw.Header().Set("Content-Type", "application/json")
}
func main() {
	router := routes.New()
	router.Filter(setContentHeader)

	router.Get("/", rootHandler)
	router.Get("/bindings", getAllBindingsHandler)
	router.Get("/bindings/?", getAllBindingsHandler)

	router.Post("/bindings", addBindingsHandler)
	router.Get("/bindings/:id", getBindingsHandler)
	router.Del("/bindings/:id", delBindingsHandler)

	router.Get("/providers", getAllProvidersHandler)
	router.Get("/providers/?", getAllProvidersHandler)
	router.Post("/providers", addProvidersHandler)
	router.Get("/providers/:id", getProvidersHandler)
	router.Del("/providers/:id", delProvidersHandler)

	router.Get("/services", getAllServicesHandler)
	router.Get("/services/?", getAllServicesHandler)

	router.Post("/services", addServicesHandler)
	router.Get("/services/:id", getServicesHandler)
	router.Del("/services/:id", delServicesHandler)

	http.Handle("/", router)
	defer func() {
		log.Println("Service Gateway Terminating")
	}()
	log.Println("Service Gateway starting..")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
