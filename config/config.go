package config

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

const (
	LOGOUT_FILE   = "file"
	LOGOUT_STDOUT = "stdout"

	LISTENER_TYPE_ETH_TREASURE = "ETH_TREASURE"
	LISTENER_TYPE_ZK_TREASURE  = "ZK_TREASURE"

	EVENT_TYPE_ERC20_DEPOSIT     = "ERC20_DEPOSIT"
	EVENT_TYPE_ERC20_WITHDRAW    = "ERC20_WITHDRAW"
	EVENT_TYPE_NFT_LOOT_WITHDRAW = "NFT_LOOT_WITHDRAW"
	EVENT_TYPE_NFT_GAME_MINT     = "NFT_GAME_MINT"
	EVENT_TYPE_NFT_WITHDRAW      = "NFT_WITHDRAW"
	EVENT_TYPE_NFT_DEPOSIT       = "NFT_DEPOSIT"
)

// one DB one instance
type RedisConfig struct {
	UseCA        bool
	IP           string
	ConnPort     int64
	SSHPort      int64
	SSHAccount   string
	SSHKey       string
	Name         string
	Password     string
	Host         string
	DB           int64
	MinIdleConns int64
}

// one database one instance
type MysqlConfig struct {
	UseCA        bool
	IP           string
	ConnPort     int64
	SSHPort      int64
	SSHAccount   string
	SSHKey       string
	Account      string
	Password     string
	SqlName      string
	MaxOpenConns int64
	MaxIdleConns int64
	MaxLifetime  int64
}
type ServerConfig struct {
	RunMode             string
	HttpPort            string
	ReadTimeout         int64
	WriteTimeout        int64
	TLSCAFile           string
	TLSCAKey            string
	JWTExpireTimeMinute int64
	JWTSecret           string
	LogOut              string
}

type MiddlewareLogConfig struct {
	VisitLogFile   string
	RecoverLogFile string
	SkipPath       []string
}

type EmailConfig struct {
	Account  string
	Password string
}

type OSSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	S3Region        string
	S3Bucket        string
	EndPoint        string
}

//for test

type ChainTokenInfo struct {
	TokenName string `toml:"TokenName"`
	Decimal   int64  `toml:"Decimal"`
}

type FilterConfig struct {
	NodeURL             []string
	ListenerType        string
	EfficientBlockNum   int64
	ParseLogsInterval   int64
	TxConfirmedInterval int64
	BlockHeightInterval int64

	ParseLogSwitch      bool
	InitFilterLogHeight int64
	SkipCheckOrphanNum  int64

	Events []EventConfig
}

type EventConfig struct {
	Type  string //ERC20_DEPOSIT,ERC20_WITHDRAW,NFT_LOOT_WITHDRAW,NFT_GAME_MINT,NFT_WITHDRAW,NFT_DEPOSIT
	Topic string
}

type SignMachineConfig struct {
	SignaturePublicFile string
	SignServerURL       string
	SignResultServerURL string
	SignRequestNum      int64
}

// skywalking config
type SkyWalkingConfig struct {
	SkyRPCNode        string
	ServerName        string
	KeyWithdraw       string
	ValueWithdraw     string
	KeyGameWithdraw   string
	ValueGameWithdraw string
	KeyRiskReview     string
	ValueRiskReview   string
}

// AssetSortConfig config
type AssetSortConfig struct {
	ArchLootPart    string
	BOSSChest       string
	CommonChest     string
	RepairingPotion string
}

// PassAssetConfig config
type PassAssetConfig struct {
	AdventurerPass string
	GoblinPass     string
	CollectorPass  string
}

type GameConfig struct {
	GameServerMock bool
	FilterConfig   FilterConfig
}

// struct decode must has tag
type Config struct {
	ServerConf        ServerConfig        `toml:"ServerConfig" mapstructure:"ServerConfig"`
	MysqlConf         MysqlConfig         `toml:"MysqlConfig" mapstructure:"MysqlConfig"`
	RedisConf         RedisConfig         `toml:"RedisConfig" mapstructure:"RedisConfig"`
	MiddlewareLogConf MiddlewareLogConfig `toml:"MiddlewareLogConfig" mapstructure:"MiddlewareLogConfig"`
	GameConf          map[int]GameConfig  `toml:"GameConfig" mapstructure:"GameConfig"`
	SignConf          SignMachineConfig   `toml:"SignConfig" mapstructure:"SignConfig"`
	SkyWalkingConf    SkyWalkingConfig    `toml:"SkyWalkingConfig" mapstructure:"SkyWalkingConfig"`

	//vipper unmarshal is case insensitive.keys will be lowercase.  https://github.com/spf13/viper/issues/1014
	AssetSortConf map[string]int    `toml:"AssetSortConfig"`
	PassAssetConf map[string]string `toml:"PassAssetConfig"`

	NetworkChain map[string]int64
	ChainNetwork map[int64]string
	ChainToken   map[string]ChainTokenInfo `toml:"ChainToken" mapstructure:"ChainToken"`

	NodeURL   []string
	ZKNodeURL []string

	GameEquipmentName string

	GraphClientPoolSize int64
	HttpPoolSize        int64 //http pool size

	KeyExpireTime                  int64
	KeyNoExpireTime                int64
	FTDepositWithdrawRedisLockTime int64
	APIDelayTime                   int64

	MarketplaceGraph       string
	MarketplaceContract    string
	GameLootEquipmentGraph string
	BlindBoxImage          string

	AdminNodeURL string
}

type NacosConfig struct {
	ServerConfigs []constant.ServerConfig
	ClientConfig  constant.ClientConfig
	Group         string
}

var (
	configMutex     = sync.RWMutex{}
	configPath      = ""
	configFileAbs   = ""
	config          Config
	nacos           NacosConfig
	configClient    config_client.IConfigClient
	configViper     *viper.Viper
	configFlyChange []chan bool
)

func RegistConfChange(c chan bool) {
	configFlyChange = append(configFlyChange, c)
}

func notifyConfChange() {
	for i := 0; i < len(configFlyChange); i++ {
		configFlyChange[i] <- true
	}
}

func watchConfig(c *viper.Viper) error {
	c.WatchConfig()
	c.OnConfigChange(func(e fsnotify.Event) {
		logger.Logrus.WithFields(logrus.Fields{"change": e.String()}).Info("config change and reload it")
		reloadConfig(c)
		notifyConfChange()
	})
	return nil
}

func LoadConf(configFilePath string) error {
	config = Config{}
	configMutex.Lock()
	defer configMutex.Unlock()

	configViper = viper.New()
	configViper.SetConfigName("config")
	configViper.AddConfigPath(configFilePath) //endwith "/"
	configViper.SetConfigType("yaml")

	if err := configViper.ReadInConfig(); err != nil {
		return err
	}
	if err := configViper.Unmarshal(&config); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(configViper.GetString("PassAssetConfig")), &config.PassAssetConf); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(configViper.GetString("AssetSortConfig")), &config.AssetSortConf); err != nil {
		return err
	}

	chainNetwork := map[int64]string{}
	for k, v := range config.NetworkChain {
		chainNetwork[v] = k
	}
	config.ChainNetwork = chainNetwork

	s, _ := json.MarshalIndent(config, "", "\t")
	fmt.Printf("Load config: %s", s)

	if err := watchConfig(configViper); err != nil {
		return err
	}
	return nil
}

func LoadFromNacos(c *viper.Viper) error {
	config = Config{}
	configMutex.Lock()
	defer configMutex.Unlock()

	if err := c.Unmarshal(&config); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(c.GetString("PassAssetConfig")), &config.PassAssetConf); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(c.GetString("AssetSortConfig")), &config.AssetSortConf); err != nil {
		return err
	}

	chainNetwork := map[int64]string{}
	for k, v := range config.NetworkChain {
		chainNetwork[v] = k
	}
	config.ChainNetwork = chainNetwork

	s, _ := json.MarshalIndent(config, "", "\t")
	fmt.Printf("Load config: %s", s)
	return nil
}

func reloadConfig(c *viper.Viper) {
	configMutex.Lock()
	defer configMutex.Unlock()

	if err := c.ReadInConfig(); err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("config ReLoad failed")
	}

	if err := configViper.Unmarshal(&config); err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("unmarshal config failed")
	}
	if err := json.Unmarshal([]byte(configViper.GetString("PassAssetConfig")), &config.PassAssetConf); err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("unmarshal PassAssetConfig failed")
	}

	if err := json.Unmarshal([]byte(configViper.GetString("AssetSortConfig")), &config.AssetSortConf); err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("unmarshal AssetSortConfig failed")
	}
	chainNetwork := map[int64]string{}
	for k, v := range config.NetworkChain {
		chainNetwork[v] = k
	}
	config.ChainNetwork = chainNetwork

	logger.Logrus.WithFields(logrus.Fields{"config": config}).Info("Config ReLoad Success")
}

func NewNacosConfig(cctx *cli.Context) (*viper.Viper, error) {
	configViper = viper.New()
	configViper.SetConfigType("yaml")

	log.Println("read config from nacos")

	group := cctx.String("Group")
	DataIds := cctx.String("DataIds")
	NacosAddrs := cctx.String("NacosAddrs")
	NacosAddrsList := strings.Split(NacosAddrs, ",")
	ServerConfigs := make([]constant.ServerConfig, len(NacosAddrsList))
	nacos.ServerConfigs = ServerConfigs
	nacos.Group = group
	for i, addr := range NacosAddrsList {
		addrArray := strings.Split(addr, ":")
		intNum, _ := strconv.Atoi(addrArray[1])
		nacos.ServerConfigs[i] = constant.ServerConfig{
			Scheme:      "http",
			ContextPath: "/nacos",
			IpAddr:      addrArray[0],
			Port:        uint64(intNum),
		}
	}
	NamespaceId := cctx.String("NamespaceId")
	NacosLogLevel := cctx.String("NacosLogLevel")
	nacos.ClientConfig = constant.ClientConfig{
		NamespaceId:         NamespaceId,
		LogLevel:            NacosLogLevel,
		NotLoadCacheAtStart: true,
		TimeoutMs:           30000,
	}
	for _, DataId := range strings.Split(DataIds, ",") {
		content, err := GetConfigByDataId(DataId)
		if err != nil {
			log.Fatal("Read remote config error:", err)
		}
		err = configViper.MergeConfig(strings.NewReader(content))
		if err != nil {
			log.Fatal("Read remote config error:", err)
		}
		err = configClient.ListenConfig(vo.ConfigParam{
			DataId: DataId,
			Group:  group,
			OnChange: func(namespace, group, dataId, data string) {
				//configViper = viper.New()
				//configViper.SetConfigType("yaml")
				configViper.MergeConfig(strings.NewReader(data))
				LoadFromNacos(configViper)
				notifyConfChange()
			},
		})
		if err != nil {
			log.Fatal("listen remote config error:", err)
		}
	}

	return configViper, nil
}

func GetConfigByDataId(DataId string) (string, error) {
	if configClient == nil {
		cli, err := clients.CreateConfigClient(
			map[string]interface{}{
				"clientConfig":  nacos.ClientConfig,
				"serverConfigs": nacos.ServerConfigs,
			},
		)
		if err != nil {
			return "", err
		}
		configClient = cli
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: DataId,
		Group:  nacos.Group})
	return content, err
}

func GetServerConfig() ServerConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.ServerConf
}

func GetMysqlConfig() MysqlConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.MysqlConf
}

func GetMiddlewareLogConfig() MiddlewareLogConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.MiddlewareLogConf
}

func GetRedisConfig() RedisConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.RedisConf
}

func GetSkyWalkingConfig() SkyWalkingConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.SkyWalkingConf
}

func GetAssetSortConfig() map[string]int {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.AssetSortConf
}

func GetPassAssetConfig() map[string]string {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.PassAssetConf
}

func GetSignConfig() SignMachineConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.SignConf
}

func GetChainNetwork() map[int64]string {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.ChainNetwork
}

func GetNetworkChain() map[string]int64 {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.NetworkChain
}

func GetChainToken() map[string]ChainTokenInfo {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.ChainToken
}

func GetNodeURL() []string {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.NodeURL
}

func GetZKNodeURL() []string {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.ZKNodeURL
}

func GetGameEquipmentName() string {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.GameEquipmentName
}

func GetGraphClientPoolSize() int64 {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.GraphClientPoolSize
}

func GetHttpPoolSize() int64 {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.HttpPoolSize
}

func GetKeyExpireTime() int64 {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.KeyExpireTime
}

func GetKeyNoExpireTime() int64 {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.KeyNoExpireTime
}

func GetFTDepositWithdrawRedisLockTime() int64 {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.FTDepositWithdrawRedisLockTime
}

func GetAPIDelayTime() int64 {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.APIDelayTime
}

func GetMarketplaceGraph() string {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.MarketplaceGraph
}

func GetMarketplaceContract() string {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.MarketplaceContract
}

func GetGameLootEquipmentGraph() string {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.GameLootEquipmentGraph
}

func GetBlindBoxImage() string {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.BlindBoxImage
}

func GetAdminNodeURL() string {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.AdminNodeURL
}

func GetGameConfig() map[int]GameConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.GameConf
}

// check if logout equal file
func (c ServerConfig) LogOutFile() bool {
	return c.LogOut == LOGOUT_FILE
}

// check if logout equal stdout
func (c ServerConfig) LogOutStdout() bool {
	return c.LogOut == LOGOUT_STDOUT
}

func (c FilterConfig) IsEthTreasure() bool {
	return c.ListenerType == LISTENER_TYPE_ETH_TREASURE
}

func (c FilterConfig) IsZkTreasure() bool {
	return c.ListenerType == LISTENER_TYPE_ZK_TREASURE
}

func (c EventConfig) IsErc20Deposit() bool {
	return c.Type == EVENT_TYPE_ERC20_DEPOSIT
}

func (c EventConfig) IsErc20Withdraw() bool {
	return c.Type == EVENT_TYPE_ERC20_WITHDRAW
}
func (c EventConfig) IsNftGameMint() bool {
	return c.Type == EVENT_TYPE_NFT_GAME_MINT
}
func (c EventConfig) IsNftLootWithdraw() bool {
	return c.Type == EVENT_TYPE_NFT_LOOT_WITHDRAW
}
func (c EventConfig) IsNftWithdraw() bool {
	return c.Type == EVENT_TYPE_NFT_WITHDRAW
}

func (c EventConfig) IsNftDeposit() bool {
	return c.Type == EVENT_TYPE_NFT_DEPOSIT
}
