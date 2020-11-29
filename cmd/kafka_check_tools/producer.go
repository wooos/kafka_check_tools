package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"strings"
)

type producerValue struct {
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

func (v *producerValue) runCommand(cmd *cobra.Command, args []string) {
	brokens := strings.Split(v.brokens, ",")

	conf := sarama.NewConfig()
	conf.Producer.Retry.Max = 1
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	conf.Metadata.Full = true
	conf.Version = sarama.V0_10_0_0
	conf.ClientID = "sasl_scram_client"
	conf.Metadata.Full = true

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
			return
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

	if err := conf.Validate(); err != nil {
		log.Printf("config error: %v", err.Error())
		return
	}


	producer, err := sarama.NewSyncProducer(brokens, conf)
	if err != nil {
		log.Printf("failed create producer: %s", err.Error())
		return
	}

	patition, offset, err :=producer.SendMessage(&sarama.ProducerMessage{
		Topic:     v.topic,
		Key:       sarama.StringEncoder("test"),
		Value:     sarama.StringEncoder("This is a check message"),
	})

	if err != nil {
		log.Printf("failed send message to %s", v.topic)
		return
	}

	log.Printf("wrote message success, partition: %d, offset: %d", patition, offset)
}

func newProducerCommand() *cobra.Command {
	p := &producerValue{}

	cmd := &cobra.Command{
		Use: "producer",
		Short: "kafka producer check tools",
		Run: p.runCommand,
	}

	flags := cmd.Flags()
	flags.StringVarP(&p.brokens, "brokens", "b", "", "The Kafka brokers to connect to, as a comma separated list")
	flags.StringVarP(&p.topic, "topic", "t", "", "The Kafka topic to use")
	flags.BoolVar(&p.sasl, "sasl", false, "Enable SASL")
	flags.StringVarP(&p.username, "username", "u", "", "The SASL user, if sasl, this is required")
	flags.StringVarP(&p.password, "password", "p", "", "The SASL password, if sasl, this is required")
	flags.BoolVar(&p.ssl, "ssl", false, "Enable SSL")
	flags.StringVar(&p.ca, "ca", "", "The SSL ca certificate, if ssl, this is required")
	flags.StringVar(&p.cert, "cert", "", "The SSL client certificate")
	flags.StringVar(&p.key, "key", "", "The SSL client key")

	_ = cmd.MarkFlagRequired("brokens")
	_ = cmd.MarkFlagRequired("topic")

	return cmd
}
