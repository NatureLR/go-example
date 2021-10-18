FROM golang:1.16.3-alpine3.13 as build

WORKDIR /root

COPY . .

RUN go build .

FROM alpine

COPY --from=build /root/admission-example .

CMD [ "./admission-example" ]