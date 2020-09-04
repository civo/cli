FROM alpine:latest
COPY civo /usr/local/bin/civo

RUN apk add --update ca-certificates \
    && apk add --update -t deps curl \
	&& curl -L "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" -o /usr/local/bin/kubectl \
	&& chmod +x /usr/local/bin/kubectl

ENTRYPOINT ["civo", "--config", "/.civo.json"]
CMD [ "version" ]
