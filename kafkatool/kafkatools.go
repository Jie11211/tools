package kafkatool

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	"github.com/Shopify/sarama"
)

func NewKafkaTool(addr []string) *Kafkatool {
	return &Kafkatool{
		Addr:   addr,
		Config: sarama.NewConfig(),
	}
}

func NewConfig() *sarama.Config {
	return sarama.NewConfig()
}

func NewDefaultConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Successes = true
	config.Version = sarama.V0_10_2_1
	config.Producer.Partitioner = sarama.NewHashPartitioner
	return config
}

// 啥都不带
func (k *Kafkatool) NewDefaultClient() (sarama.Client, error) {
	return sarama.NewClient(k.Addr, k.Config)
}

// 自带配置
func (k *Kafkatool) NewClientWithConfig(config *sarama.Config) (sarama.Client, error) {
	return sarama.NewClient(k.Addr, config)
}

// 配置携带密码
func (k *Kafkatool) AddPassAndUser(user, password string) {
	k.Config.Net.SASL.Enable = true
	k.Config.Net.SASL.User = user
	k.Config.Net.SASL.Password = password
}

// 添加ssl证书  clientcert:cert.pem   clientkey:key.pem  cacert:ca.pem
func (k *Kafkatool) AddTLS(clientcert, clientkey, cacert string) {
	tlsConfig, err := genTLSConfig(clientcert, clientkey, cacert)
	if err != nil {
		log.Fatal(err)
	}
	k.Config.Net.TLS.Enable = true
	k.Config.Net.TLS.Config = tlsConfig
}

func genTLSConfig(clientcertfile, clientkeyfile, cacertfile string) (*tls.Config, error) {
	// load client cert
	clientcert, err := tls.LoadX509KeyPair(clientcertfile, clientkeyfile)
	if err != nil {
		return nil, err
	}

	// load ca cert pool
	cacert, err := ioutil.ReadFile(cacertfile)
	if err != nil {
		return nil, err
	}
	cacertpool := x509.NewCertPool()
	cacertpool.AppendCertsFromPEM(cacert)

	// generate tlcconfig
	tlsConfig := tls.Config{}
	tlsConfig.RootCAs = cacertpool
	tlsConfig.Certificates = []tls.Certificate{clientcert}
	tlsConfig.BuildNameToCertificate()
	return &tlsConfig, err
}

// 新建生产者
func (k *Kafkatool) Newproducer(async bool) (*Producer, error) {
	if async {
		asyncProducer, err := sarama.NewAsyncProducer(k.Addr, k.Config)
		return &Producer{AsyncProducer: asyncProducer, SyncProducer: nil}, err
	}
	syncProducer, err := sarama.NewSyncProducer(k.Addr, k.Config)
	return &Producer{AsyncProducer: nil, SyncProducer: syncProducer}, err
}

// 携带配置文件新建生产者
func (k *Kafkatool) NewproducerWithConfig(async bool, config *sarama.Config) (*Producer, error) {
	if async {
		asyncProducer, err := sarama.NewAsyncProducer(k.Addr, config)
		return &Producer{AsyncProducer: asyncProducer, SyncProducer: nil}, err
	}
	syncProducer, err := sarama.NewSyncProducer(k.Addr, config)
	return &Producer{AsyncProducer: nil, SyncProducer: syncProducer}, err
}

// 新建消费者
func (k *Kafkatool) NewConsumer() (sarama.Consumer, error) {
	return sarama.NewConsumer(k.Addr, k.Config)
}

// 携带配置文件新建消费者
func (k *Kafkatool) NewConsumerWithConfig(config *sarama.Config) (sarama.Consumer, error) {
	return sarama.NewConsumer(k.Addr, config)
}

// 新建消费者组
func (k *Kafkatool) NewConsumerGroup(groupID string) (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroup(k.Addr, groupID, k.Config)
}
