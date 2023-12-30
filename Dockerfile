FROM debian:stable-slim

WORKDIR app

COPY ./files .

RUN apt-get update && \
    apt-get install -y rtl-sdr rtl-433 mosquitto-clients jq

CMD ["./run.sh"]