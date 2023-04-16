## K8s quick start

```env example
        - name: OS_AUTH_URL
          value: "https://test
        - name: OS_IDENTITY_API_VERSION
          value: "3"
        - name: OS_INTERFACE
          value: "public"
        - name: OS_PASSWORD
          value: "test"
        - name: OS_PROJECT_ID
          value: "test"
        - name: OS_PROJECT_NAME
          value: "t-devops"
        - name: OS_REGION_NAME
          value: "RegionOne"
        - name: OS_USERNAME
          value: "sa-vvc-test"
        - name: OS_USER_DOMAIN_NAME
          value: "test"
```
## Usage

```console
Usage of ./openstack-bash-exporter:
  -debug
    	Debug log level
  -interval int
    	Interval for metrics collection in seconds (default 300)
  -web.listen-address string
    	Address on which to expose metrics (default ":9300")
```
