# b64-service

b64-service is a small backend that is responsible for serving files in base64 to be consumed through an http client

## Configuration
_for each file directory (service) you should add a settings array in the configuration file._

_config.yml_
```yml
settings:
    prefix: <prefix>
    querys: <query-params-array>
    path: <path>

```
