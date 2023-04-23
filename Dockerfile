# 运行阶段
FROM alpine:latest
COPY sso /usr/bin/sso
COPY sso.yaml /usr/bin/sso.yaml
WORKDIR /usr/bin/
RUN  apk add ca-certificates
RUN chmod +x /usr/bin/sso
CMD ["sso","start"]
