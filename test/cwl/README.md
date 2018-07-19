# Goal

```sh
% hotsub run \
  --cwl ./hello.cwl \
  --cwl-param ./job.yml \
  --verbose
```

# Control

```sh
% docker exec -it c4cwl \
  cwltool /tmp/work/hello.cwl /tmp/work/job.yml

/usr/local/bin/cwltool 1.0.20180622214234
Resolved '/tmp/work/hello.cwl' to 'file:///tmp/work/hello.cwl'
[job hello.cwl] /tmp/tmp09Rfd5$ echo \
    'Hello, World!!'
Hello, World!!
[job hello.cwl] completed success
{}
Final process status is success
```