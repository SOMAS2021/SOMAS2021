# syntax=docker/dockerfile:1

FROM node:latest AS FRONTEND

# make the 'app' folder the current working directory
RUN mkdir /app

# make the new app the working directory
WORKDIR /app

# add all the files from the existing
ADD . /app/

# make frontend folder
WORKDIR /app/frontend/

# install dependencies and build app
RUN npm ci
RUN npm run build


FROM golang:1.17-alpine AS APP
RUN mkdir /app
WORKDIR /app
RUN mkdir /build
RUN apk add git gcc libc-dev g++
ADD . /app/
COPY --from=FRONTEND /app/frontend/build /app/build
RUN go mod download
RUN go build -o main main.go
EXPOSE 9000
CMD ["/app/main"]