# omada-to-ntfy

## Purpose

This is a small program written in Go which spawns a server that'll receive
webhook messages from a TP-Link Omada Network Controller, it converts them
into ntfy notifications and delivers them to ntfy.

Run it in Docker, in a LXC, or really anywhere you like (anywhere as long as
the Omada Network Controller can talk to it, and it can talk to your ntfy
server). I'm running it in a docker-compose stack with the stack managed
through Portainer. My compose file is further below.

ntfy is a simple HTTP-based pub-sub notification service that allows you to
send notifications to your phone or desktop via scripts from any computer,
entirely without signup, cost or setup. It can be self-hosted or you can use
the free public ntfy.sh server.

## Features

- **Priority Mapping**: Automatically maps Omada priorities (0-10) to ntfy priorities (1-5)
  - Omada 10 → ntfy 5 (Max/Urgent) 🚨
  - Omada 7 → ntfy 4 (High)
  - Omada 4 → ntfy 3 (Default)
  - Omada 0 → ntfy 2 (Low)
- **Emoji Tags**: Automatic emoji tags based on message type
  - Offline notifications: 🚨 (rotating_light)
  - Online notifications: ✅ (white_check_mark)
  - Test messages: 🧪 (test_tube)
  - Unrecognized messages: ⚠️ (warning)
- **Optional Authentication**: Supports Basic Auth for protected ntfy instances
- **Simple Setup**: No external dependencies beyond standard Go libraries

## Installation / Configuration

Environment variables are used for configuration. They are:

### Required environment variables

- `NTFY_URL` - The base URL of your ntfy server (e.g., `https://ntfy.sh` or `https://ntfy.example.com`)
- `NTFY_TOPIC` - The ntfy topic to publish to (e.g., `my_omada_alerts`)
- `OMADA_SHARED_SECRET` - The shared secret configured on the Omada Network Controller for this webhook

### Optional environment variables

- `NTFY_USER` - Username for ntfy authentication (if your ntfy instance requires auth)
- `NTFY_PASSWORD` - Password for ntfy authentication (if your ntfy instance requires auth)
- `PORT` - The port on which to run the server (default is `8080`)

## Usage

To use this project directly without Docker:

1. Configure the webhook in Omada using the "Omada format", match the server and port where you are running this program. For example: `http://192.168.12.34:8080/`.
2. Set the required environment variables, making sure to include the shared secret from Omada.
3. Launch the executable with those environment variables set.
4. Enable the events to monitor in both the global view and your sites.
5. Wait for a message to come through from your Omada Controller and see it appear in ntfy.

At the moment there are no delivery retries should delivery fail, but each time it fails to either parse or deliver it will log an error to the console and then try connecting to ntfy again on the next request. However, Omada itself allows you to set up retries and see information about both successful and failed webhook requests so that should be adequate.

### docker

A docker image can be built from this repository. Use the included Dockerfile to build your own image.

### docker-compose

Here's an example docker-compose.yml file for running omada-to-ntfy:

```yaml
services:

  omada-to-ntfy:
    image: omada-to-ntfy:latest
    environment:
      NTFY_URL: https://ntfy.sh
      NTFY_TOPIC: ${NTFY_TOPIC}
      NTFY_USER: ${NTFY_USER}       # Optional
      NTFY_PASSWORD: ${NTFY_PASSWORD} # Optional
      OMADA_SHARED_SECRET: ${OMADA_SHARED_SECRET}
    volumes:
      - /etc/timezone:/etc/timezone:ro
      - /etc/localtime:/etc/localtime:ro
    ports:
      - "8080:8080"
    restart: always
```

Mounting the timezone volumes like in this example will avoid the timestamps being reported in UTC (remove them if you *do* want UTC timestamps). You will need to set up the environment variables for your ntfy configuration and the Omada shared secret.

### Using with ntfy.sh (public server)

You can use the free public ntfy.sh server without authentication:

```bash
export NTFY_URL="https://ntfy.sh"
export NTFY_TOPIC="my_unique_omada_topic_12345"  # Choose something unique!
export OMADA_SHARED_SECRET="your-secret-here"
./omada-to-ntfy
```

**Note**: Since ntfy.sh is public, anyone who knows your topic name can subscribe to it. Choose a unique, hard-to-guess topic name!

### Using with self-hosted ntfy

If you're running your own ntfy server with authentication:

```bash
export NTFY_URL="https://ntfy.example.com"
export NTFY_TOPIC="omada_alerts"
export NTFY_USER="your-username"
export NTFY_PASSWORD="your-password"
export OMADA_SHARED_SECRET="your-secret-here"
./omada-to-ntfy
```

## Future

Possible additions to come (and feel free to contribute):

- Improving the instructions further, maybe also provide a basic LXC setup script.
- Specific support for more types of events from the Omada Controller, such as detecting more message patterns and providing appropriate priorities and emoji tags.
- MacOS support? I've got no way to test it works on MacOS, but I'll take pull requests for it if someone needs that. Then we'll blame you for any problems from then on. :wink:

## Migration from Gotify

This project was originally `omada-to-gotify` and has been migrated to use ntfy instead. If you're migrating from the old version:

1. Update your environment variables from `GOTIFY_*` to `NTFY_*`
2. Change `GOTIFY_URL` to `NTFY_URL` and `GOTIFY_APP_TOKEN` to `NTFY_TOPIC`
3. Optionally add `NTFY_USER` and `NTFY_PASSWORD` if your ntfy instance requires authentication
4. Update your docker image or rebuild from source

The webhook endpoint and Omada configuration remain the same.

## LICENSE

Copyright (c) 2025 Lianna Eeftinck <liannaee@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
