FROM alpine:latest
COPY civo /usr/local/bin/civo

RUN adduser -S civo
USER civo

ENTRYPOINT ["civo"]
CMD [ "version" ]
