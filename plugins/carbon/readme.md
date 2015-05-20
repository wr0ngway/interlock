# Carbon
This plugin reports stats as reported from the Docker stats API.

# Configuration
The following configuration is available through environment variables:

- `CARBON_ADDRESS`: Carbon receiver address (i.e. `1.2.3.4:2003`)
- `CARBON_PREFIX`: Stat prefix (default: `docker.stats`)
- `CARBON_IMAGE_NAME_FILTER`: Regex to match against container image name to gather stats (default: `.*` - all containers)
- `CARBON_INTERVAL`: Interval (in seconds) to send stats to Carbon (default: `10`)
