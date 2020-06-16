FROM alpine:latest
COPY civo /usr/local/bin/civo

RUN adduser -S user
USER user

ENTRYPOINT ["civo"]
CMD [ "version" ]
