

# Troubleshooting Docker Compose

## Windows

Some problems have been found while following the [Elastic Stack Docker Compose guide](https://www.elastic.co/blog/getting-started-with-the-elastic-stack-and-docker-compose#conclusion):

- **Metricbeat**: As the OS Host is Windows, a permission issue was encountered for [metribeat.yml](../integrations/metricbeat.yml) when running the Metricbeat container. Although a volume containing this file is bound to the container, when Windows is the host OS, the permissions inheritance appears to be problematic. Specifically, the metricbeat.yml file needs Read-Only permissions. Here are the steps to resolve the issue:

    1. Create a bash script for modifying the permissions once the container is running.

    ```sh
        #!/bin/bash
        # Set correct permissions for the config file
        chmod 600 /usr/share/metricbeat/metricbeat.yml
        # Execute the original Metricbeat entrypoint command
        exec /usr/local/bin/docker-entrypoint -e
    ```

    2. Add a volume in [docker-compose.yml](../integrations/docker-compose.yml) in the Metricbeat service to bind the recently created bash script file. Call the script by using the `entrypoint` property of the metricbeat service.

    ```yml
    volumes:
     ...
     ...
     - ./entrypoint_mf.sh:/entrypoint_mf.sh
   command: ["-e"]
   entrypoint: ["/entrypoint_mf.sh"] # Use your custom entrypoint
    ```

- **Filebeat**: The issue with the Filebeat container experiences the same problem as the Metricbeat one. The solution is the same: to create a bash script to modify the [filebeat.yml](../integrations/filebeat.yml) permissions, bind a volume with the script to the container and call the script with the `entrypoint` property in docker-compose.


