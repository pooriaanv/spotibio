FROM ubuntu:20.04

RUN apt-get update && apt-get install cron -y && apt-get install ca-certificates -y

COPY cronfile /etc/cron.d/cronfile

RUN chmod 0644 /etc/cron.d/cronfile

RUN crontab /etc/cron.d/cronfile

CMD ["cron", "-f"]