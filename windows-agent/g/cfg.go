package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"os"
	"sync"
	"time"
)

type PluginConfig struct {
	Enabled bool   `json:"enabled"`
	Dir     string `json:"dir"`
	Git     string `json:"git"`
	Rsync   string `json:"rsync"`
	LogDir  string `json:"logs"`
	Args    string `json:"args"`
}

type HeartbeatConfig struct {
	Enabled  bool   `json:"enabled"`
	Addr     string `json:"addr"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type TransferConfig struct {
	Enabled  bool   `json:"enabled"`
	Addrs    []string `json:"addrs"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type GlobalConfig struct {
	Debug         bool             `json:"debug"`
	Hostname      string           `json:"hostname"`
	IP            string           `json:"ip"`
	SmartAPI      string           `json:"smartapi"`
	Heartbeat     *HeartbeatConfig `json:"heartbeat"`
	Transfer      *TransferConfig  `json:"transfer"`
	IgnoreMetrics map[string]bool  `json:"ignore"`
	Port          []string           `json:"port"`
	DialTimeout   time.Duration  `json:"dialTimeOut"`
	Process       map[string]bool `json:"process"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func Hostname() (string, error) {
	hostname := Config().Hostname
	if hostname != "" {
		return hostname, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: os.Hostname() fail", err)
	}
	return hostname, err
}

func IP() string {
	ip := Config().IP
	if ip != "" {
		// use ip in configuration
		return ip
	}

	if len(LocalIps) > 0 {
		ip = LocalIps[0]
	}

	return ip
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}
