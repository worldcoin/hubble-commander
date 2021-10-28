FROM amazon/aws-cli

RUN yum install -y cronie rsync tar pigz

RUN amazon-linux-extras install -y postgresql13

COPY docker_init.sh /docker_init.sh

RUN chmod +x /docker_init.sh

ENTRYPOINT ["/docker_init.sh"]
