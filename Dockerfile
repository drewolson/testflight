FROM 864879987165.dkr.ecr.us-east-1.amazonaws.com/calm/docker-base-go:0.0.12 as base
ARG GITHUB_ACCESS_TOKEN
COPY go.* ./
RUN git config --global url."https://${GITHUB_ACCESS_TOKEN}:@github.com/".insteadOf "https://github.com/"  \
    && go mod download \
    && rm -rf /root/.gitconfig
COPY . .
RUN go build ./...
