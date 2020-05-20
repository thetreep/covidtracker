FROM golang:1.14


WORKDIR /app

# copy the linux executable
COPY covidtracker covidtracker

EXPOSE 3456

ENTRYPOINT ["/app/covidtracker"]
