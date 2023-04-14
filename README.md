## Docker quick start

```console
$ docker run -v $PWD/examples:/scripts -p 9300:9300 greegorey/bash-exporter:2.0.0
```

```console
$ curl -s 127.1:9300/metrics | grep ^bash
bash{env="",hostname="node-1",job="job-2",verb="get"} 0.003
bash{env="",hostname="node-1",job="job-2",verb="put"} 0.13
bash{env="",hostname="node-1",job="job-2",verb="time"} 0.5
bash{env="dev",hostname="",job="job-1",verb="items"} 21
```

## Usage

```console
Usage of ./bash-exporter:
  -debug
    	Debug log level
  -interval int
    	Interval for metrics collection in seconds (default 300)
  -labels string
    	additioanal labels (default "hostname,env")
  -path string
    	path to directory with bash scripts (default "/scripts")
  -prefix string
    	Prefix for metrics (default "bash")
  -web.listen-address string
    	Address on which to expose metrics (default ":9300")
```

Just point `-path` flag to the directory with your bash scripts. Names of the files (`(.*).sh`) will be used as the `job` label. Bash scripts should return valid json.

Example output:
```console
# HELP bash bash exporter metrics
# TYPE bash gauge
bash{job="job-1",verb="items"} 21
bash{job="job-2",verb="get"} 0.003
bash{job="job-2",verb="put"} 0.13
bash{job="job-2",verb="time"} 0.5
...
```
