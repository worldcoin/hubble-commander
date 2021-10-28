FROM amazon/aws-cli

RUN yum install -y cronie

COPY docker_init.sh /docker_init.sh

RUN chmod +x /docker_init.sh

ENTRYPOINT ["/docker_init.sh"]
