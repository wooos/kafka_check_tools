# kafka_check_tools

Kafka check tools.

## Build

```bash
make
```

## Usage

```bash
kafka check tools

Usage:
  kafka_check [command]

Available Commands:
  consumer    kafka consumer check tool
  help        Help about any command
  producer    kafka producer check tools

Flags:
  -h, --help   help for kafka_check

Use "kafka_check [command] --help" for more information about a command.
```

### producer

```bash
kafka producer check tools

Usage:
  kafka_check producer [flags]

Flags:
  -b, --brokens string    The Kafka brokers to connect to, as a comma separated list
      --ca string         The SSL ca certificate, if ssl, this is required
      --cert string       The SSL client certificate
  -h, --help              help for producer
      --key string        The SSL client key
  -p, --password string   The SASL password, if sasl, this is required
      --sasl              Enable SASL
      --ssl               Enable SSL
  -t, --topic string      The Kafka topic to use
  -u, --username string   The SASL user, if sasl, this is required
```

### consumer

```bash
kafka consumer check tool

Usage:
  kafka_check consumer [flags]

Flags:
  -b, --brokens string    The Kafka brokers to connect to, as a comma separated list
      --ca string         The SSL ca certificate, if ssl, this is required
      --cert string       The SSL client certificate
  -h, --help              help for consumer
      --key string        The SSL client key
  -p, --password string   The SASL password, if sasl, this is required
      --sasl              Enable SASL
      --ssl               Enable SSL
  -t, --topic string      The Kafka topic to use
  -u, --username string   The SASL user, if sasl, this is required
```