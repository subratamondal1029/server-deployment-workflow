FROM scratch

WORKDIR /

COPY ./app ./app
COPY ./.env ./.env

EXPOSE 8000

ENTRYPOINT [ "./app" ]
