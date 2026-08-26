FROM scratch

WORKDIR /

COPY ./app ./app

EXPOSE 8000

ENTRYPOINT [ "/app" ]