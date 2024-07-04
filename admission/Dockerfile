FROM golang:1.20-alpine as build

WORKDIR /root

COPY . .

RUN go build -o admission-example .

FROM alpine

COPY --from=build /root/admission-example .

CMD [ "./admission-example" ]