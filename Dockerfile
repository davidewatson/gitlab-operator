# Start by building the application.
FROM  quay.io/samsung_cnct/golang-container:latest as build

WORKDIR /go/src/github.com/samsung-cnct/gitlab-operator
COPY . .
RUN go install

# Now copy it into our base image.
FROM gcr.io/distroless/base

COPY --from=build /go/bin/gitlab-operator /
ENTRYPOINT ["/gitlab-operator"]
