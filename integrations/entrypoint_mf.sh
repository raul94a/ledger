#!/bin/bash
# Set correct permissions for the config file
chmod 600 /usr/share/metricbeat/metricbeat.yml
# Execute the original Metricbeat entrypoint command
exec /usr/local/bin/docker-entrypoint -e