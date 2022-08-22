# File Age Exporter

Prometheus exporter for metrics about file modification times.

# Usage

```
fileage-exporter -file "/backup/*.backup" -file "/foo/bar.json"
```

Example output:

```
curl http://localhost:9123/metrics

# HELP file_age_age_seconds File age in seconds
# TYPE file_age_age_seconds counter
file_age_age_seconds{path="/backup/foobar.backup"} 123.456
file_age_age_seconds{path="/foo/bar.json"} 78.910
# HELP file_age_num_files_matching Number of files matching glob patterns
# TYPE file_age_num_files_matching gauge
file_age_num_files_matching 2
# HELP file_age_size_bytes File size in bytes
# TYPE file_age_size_bytes gauge
file_age_size_bytes{path="/backup/foobar.backup"} 4321
file_age_size_bytes{path="/foo/bar.json"} 1234
```
