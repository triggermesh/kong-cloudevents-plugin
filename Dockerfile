# Copyright 2021 TriggerMesh Inc.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.17 as builder

WORKDIR /go/src/app
COPY . .

ENV GOOS=linux GARCH=amd64 CGO_ENABLED=0
RUN go build -v -o ce-plugin ./plugin

FROM kong:2.5

USER root

RUN apk update && \
    apk add protobuf-dev

COPY --from=builder /go/src/app/ce-plugin /usr/local/bin/

USER kong

ENV KONG_PLUGINS=bundled,ce-plugin
ENV KONG_PLUGINSERVER_NAMES=ce-plugin
ENV KONG_PLUGINSERVER_CE_PLUGIN_QUERY_CMD="/usr/local/bin/ce-plugin -dump"

EXPOSE 8000 8443 8001 8444
STOPSIGNAL SIGQUIT
HEALTHCHECK --interval=10s --timeout=10s --retries=10 CMD kong health

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["kong", "docker-start"]
