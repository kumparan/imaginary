FROM 475170104714.dkr.ecr.ap-southeast-1.amazonaws.com/imaginary-service:builder-vips-8.12.2-1.22.2 as builder
ARG IMAGINARY_VERSION="dev"
WORKDIR ${GOPATH}/src/github.com/kumparan/imaginary
# Copy imaginary sources
COPY . .
# Making sure all dependencies are up-to-date
RUN go mod download
# Compile imaginary
RUN CGO_CFLAGS_ALLOW=-Xpreprocessor go test && go build -o ${GOPATH}/bin/imaginary -ldflags="-s -w -h -X main.Version=${IMAGINARY_VERSION}" github.com/kumparan/imaginary
FROM 475170104714.dkr.ecr.ap-southeast-1.amazonaws.com/imaginary-service:base-1713944528
ARG SPINNAKER_ID="dev"
LABEL maintainer="sre@kumparan.com" \
      org.label-schema.description="kumparan imaginary" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.url="https://github.com/kumparan/imaginary" \
      org.label-schema.vcs-url="https://github.com/kumparan/imaginary"
COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/bin/imaginary /usr/local/bin/imaginary
COPY --from=builder /go/src/github.com/kumparan/imaginary/config.yml.* /usr/local/bin/
RUN apt-get update && apt-get install wget -y
COPY --from=builder /go/src/github.com/kumparan/imaginary/Heebo-*.ttf /usr/share/fonts/googlefonts/
RUN ldconfig