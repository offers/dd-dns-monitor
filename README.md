# dd-dns-monitor
Query DNS servers every x seconds and report failures to Datadog

# Docker
To run in Docker you must set the following environment variables:
* DATADOG_HOST - your Datadog collector hostname
* DNS_NAME - known DNS name to lookup
* DNS_IP - IP address of the DNS name
* DNS_SERVERS - comma-separated list of DNS servers to monitor

You may also set the check interval with the DNS_INTERVAL env var.

# Datadog
The following stats are reported to Datadog:
* dd_dns_monitor.error - count of DNS errors, tagged with failed server IP
* dd_dns_monitor.time - response time of the DNS requests
