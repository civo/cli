FROM alpine:latest
COPY civo /usr/local/bin/civo

ENTRYPOINT ["civo", "--config", "/.civo.json"]
CMD [ "version" ]
