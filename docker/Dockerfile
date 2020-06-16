FROM golang:1.14-alpine
COPY dest/civo_linux /usr/local/bin/civo

RUN adduser -S civo
USER civo

ENTRYPOINT ["civo"]
CMD [ "version" ]
