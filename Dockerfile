# Start by building the application.
FROM  quay.io/samsung_cnct/golang-container:latest as build

WORKDIR /go/src/github.com/samsung-cnct/gitlab-operator
COPY . .
RUN go install

# kubectl v1.9.0
ARG     K8S_VERSION=v1.9.0
ARG     K8S_SHA256=9150691c3c9d0c3d6c0c570a81221f476e107994b35e33c193b1b90b7b7c0cb5

RUN     wget -q https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/linux/amd64/kubectl && \
        echo "${K8S_SHA256}  kubectl" | sha256sum -c - && \
        chmod a+x kubectl && \
        mv kubectl /usr/local/bin

# Now copy it into our base image.
FROM gcr.io/distroless/base

COPY --from=build /go/bin/gitlab-operator /
COPY --from=build /usr/local/bin/kubectl /usr/local/bin

ENTRYPOINT ["/gitlab-operator"]
