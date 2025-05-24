#!/bin/bash
# Set correct permissions for the config file
chmod 600 /usr/share/filebeat/filebeat.yml
# Execute the original Filebeat entrypoint command
exec /usr/local/bin/docker-entrypoint -e