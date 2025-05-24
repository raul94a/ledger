

# Troubleshooting Docker Compose

## Windows

Some problems have been found while following the [Elastic Stack Docker Compose guide](https://www.elastic.co/blog/getting-started-with-the-elastic-stack-and-docker-compose#conclusion):

- *Metricbeat*: the binding volume for the [metric beat yml file](../integrations/metricbeat.yml) had a read-only tag. It's been removed in the present project. It was giving an error that shut the container down.

- *Filebeat*: watch out with the identation. When copying the [Filebeat yml file](../integrations/filebeat.yml) from the official guide the identation is appearly not good.


