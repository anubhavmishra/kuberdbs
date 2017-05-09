FROM alpine:3.4
MAINTAINER Anubhav Mishra <anubhavmishra@me.com>

# copy binary
COPY kuberdbs /usr/local/bin/kuberdbs

ENTRYPOINT ["/usr/local/bin/kuberdbs"]
