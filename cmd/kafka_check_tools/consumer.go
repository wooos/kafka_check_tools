package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

type consumerValue struct {
	brokens string
	topic string
	sasl bool
	username string
	password string
	ssl bool
	cert string
	key string
	ca string
}

func (v *consumerValue) runCommand(cmd *cobra.Command, args []string) {
	brokens := strings.Split(v.brokens, ",")

	conf := sarama.NewConfig()
	conf.Producer.Retry.Max = 1
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	conf.Metadata.Full = true
	conf.Version = sarama.V0_10_0_0
	conf.ClientID = "sasl_scram_client"
	conf.Metadata.Full = true
	conf.Net.ReadTimeout = time.Second * 5
	conf.Net.WriteTimeout = time.Second * 5

	if v.sasl {
		conf.Net.SASL.Enable = true
		conf.Net.SASL.User = v.username
		conf.Net.SASL.Password = v.password
		conf.Net.SASL.Handshake = true
	}

	if v.ssl {
		ca, err := ioutil.ReadFile(v.ca)
		if err != nil {
			log.Printf("failed open file, error: %s", err.Error())
			return
		}

		caCert := x509.NewCertPool()
		if ok := caCert.AppendCertsFromPEM(ca); !ok {
			log.Printf("kafka producer failed to parse root certificate")
		}

		if v.cert == "" || v.key == "" {
			conf.Net.TLS.Config = &tls.Config{
				RootCAs: caCert,
				InsecureSkipVerify: true,
			}
		} else {
			cert, err := tls.LoadX509KeyPair(v.cert, v.key)
			if err != nil {
				log.Printf("failed load key: %s", err.Error())
				return
			}

			conf.Net.TLS.Config = &tls.Config{
				Certificates:                []tls.Certificate{cert},
				RootCAs:                     caCert,
			}
		}

		conf.Net.TLS.Enable = true
	}

	consumer, err := sarama.NewConsumer(brokens, conf)
	if err != nil {
		log.Printf("failed create comsumer: %v", err.Error())
		return
	}

	log.Printf("consumer created")
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(v.topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d, Key: %s, Value: %s", msg.Offset, msg.Key, msg.Value)
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("Consumed: %d\n", consumed)

}

func newConsumerCommand() *cobra.Command {
	c := &consumerValue{}

	cmd := &cobra.Command{
		Use: "consumer",
		Short: "kafka consumer check tool",
		Run: c.runCommand,
	}

	flags := cmd.Flags()

	flags.StringVarP(&c.brokens, "brokens", "b", "", "The Kafka brokers to connect to, as a comma separated list")
	flags.StringVarP(&c.topic, "topic", "t", "", "The Kafka topic to use")
	flags.BoolVar(&c.sasl, "sasl", false, "Enable SASL")
	flags.StringVarP(&c.username, "username", "u", "", "The SASL user, if sasl, this is required")
	flags.StringVarP(&c.password, "password", "p", "", "The SASL password, if sasl, this is required")
	flags.BoolVar(&c.ssl, "ssl", false, "Enable SSL")
	flags.StringVar(&c.ca, "ca", "", "The SSL ca certificate, if ssl, this is required")
	flags.StringVar(&c.cert, "cert", "", "The SSL client certificate")
	flags.StringVar(&c.key, "key", "", "The SSL client key")

	_ = cmd.MarkFlagRequired("brokens")
	_ = cmd.MarkFlagRequired("topic")

	return cmd
}