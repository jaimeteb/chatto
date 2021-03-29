FROM alpine
RUN apk add --update --no-cache ca-certificates
RUN mkdir /data

FROM scratch
COPY chatto /chatto
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 /data /data
VOLUME /data
EXPOSE 4770/tcp
ENTRYPOINT ["/chatto"]
CMD ["--path", "/data"]
