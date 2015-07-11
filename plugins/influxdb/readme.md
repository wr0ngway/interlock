# InfluxDB
This plugin sends stats to InfluxDB as reported from the Docker stats API.

# Configuration
The following configuration is available through environment variables:

- `INFLUXDB_ADDR`: InfluxDB address (i.e. `http://1.2.3.4:8086`)
- `INFLUXDB_NAME`: InfluxDB DB name (i.e. `foo`)
- `INFLUXDB_USER`: InfluxDB username
- `INFLUXDB_PASS`: InfluxDB password
- `INFLUXDB_MEASUREMENT`: InfluxDB measurement for stats (default: `stats`)
- `INFLUXDB_IMAGE_NAME_FILTER`: Regex to match against container image name to gather stats (default: `.*` - all containers)
- `INFLUXDB_INTERVAL`: Interval (in seconds) to send stats to InfluxDB (default: `10`)
