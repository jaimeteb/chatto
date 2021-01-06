package fsm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/rpc"

	"github.com/spf13/viper"
)

// BotConfig struct models the bot.yml configuration file
type BotConfig struct {
	Extensions ExtensionsConfig `mapstructure:"extensions"`
}

// ExtensionsConfig struct models the extensions object in BotConfig
type ExtensionsConfig struct {
	Type string `mapstructure:"type"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	URL  string `mapstructure:"url"`
}

// ExtensionRPC is an RPC Client for extension functions
type ExtensionRPC struct {
	Client *rpc.Client
}

// ExtensionREST is a REST API URL for extension functions
type ExtensionREST struct {
	URL string
}

// Extension interface models an extension that can be either RPC or REST
type Extension interface {
	GetAllFuncs() []string
	RunExtFunc(extName, text string, dom Domain, m *FSM) interface{}
}

// RunExtFunc runs an extension function over RPC
func (e *ExtensionRPC) RunExtFunc(extName, text string, dom Domain, m *FSM) interface{} {
	req := Request{
		FSM: m,
		Req: extName,
		Txt: text,
		Dom: dom.NoFuncs(),
	}

	res := Response{}
	err := (*e).Client.Call("ListenerRPC.GetFunc", &req, &res)
	if err != nil {
		log.Println(err)
		return ""
	}

	*m = *res.FSM
	return res.Res
}

// GetAllFuncs retrieves all functions in extension
func (e *ExtensionRPC) GetAllFuncs() []string {
	res := new(GetAllFuncsResponse)
	if err := e.Client.Call("ListenerRPC.GetAllFuncs", new(Request), &res); err != nil {
		log.Println(err)
		return make([]string, 0)
	}
	return res.Res
}

// RunExtFunc runs an extension function over REST
func (e *ExtensionREST) RunExtFunc(extName, text string, dom Domain, m *FSM) interface{} {
	req := Request{
		FSM: m,
		Req: extName,
		Txt: text,
		Dom: dom.NoFuncs(),
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return ""
	}

	// TODO: if fail -> don't change states
	// send error msg from dom

	resp, err := http.Post(fmt.Sprintf("%v/ext/get_func", e.URL), "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Println(err)
		return ""
	}

	defer resp.Body.Close()
	res := Response{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println(err)
		return ""
	}

	*m = *res.FSM
	return res.Res
}

// GetAllFuncs retrieves all functions in extension
func (e *ExtensionREST) GetAllFuncs() []string {
	resp, err := http.Get(fmt.Sprintf("%v/ext/get_all_funcs", e.URL))
	if err != nil {
		log.Println(err)
		return make([]string, 0)
	}

	defer resp.Body.Close()
	var res []string
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println(err)
		return make([]string, 0)
	}

	return res
}

// LoadExtensions loads the extensions configuration and connects to the server
func LoadExtensions(path *string) Extension {
	config := viper.New()
	config.SetConfigName("bot")
	config.AddConfigPath(*path)

	if err := config.ReadInConfig(); err != nil {
		log.Println(err)
		return nil
	}

	var botCfg BotConfig
	if err := config.Unmarshal(&botCfg); err != nil {
		log.Println(err)
		return nil
	}

	switch botCfg.Extensions.Type {
	case "RPC":
		client, err := rpc.Dial("tcp", fmt.Sprintf("%v:%v", botCfg.Extensions.Host, botCfg.Extensions.Port))
		if err != nil {
			return nil
		}
		ext := ExtensionRPC{client}
		log.Println("Loaded extensions:")
		for i, fun := range ext.GetAllFuncs() {
			log.Printf("%v\t%v\n", i, fun)
		}
		return &ext
	case "REST":
		ext := ExtensionREST{botCfg.Extensions.URL}
		log.Println("Loaded extensions:")
		for i, fun := range ext.GetAllFuncs() {
			log.Printf("%v\t%v\n", i, fun)
		}
		return &ext
	default:
		return nil
	}
}
