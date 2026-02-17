package scaffold

import "fmt"

func (p *Project) dockerfileTemplate() string {
	bundleStep := ""
	if p.Mode.HasOpenAPI() {
		bundleStep = fmt.Sprintf("\nRUN go run tools/bundlespec/main.go api/openapi cmd/%s/docs\n", p.Name)
	}

	expose := "EXPOSE 2112 8085"
	if p.Mode.HasGRPC() {
		expose += " 60006 60011"
	}
	if p.Mode.HasOpenAPI() {
		expose += " 8080"
	}

	return fmt.Sprintf(`# Build stage
FROM git.web3gate.ru:5000/golang:1.25.5 AS builder

ARG VERSION=dev
ARG BUILD_TIME
ARG NETRC

RUN if [ -n "$NETRC" ]; then \
      echo "$NETRC" > ~/.netrc && \
      chmod 600 ~/.netrc; \
    fi

WORKDIR /app

ENV GOTOOLCHAIN=auto
ENV GOSUMDB=off
ENV GOPRIVATE="git.web3gate.ru"

COPY go.mod go.sum ./
RUN go mod download

COPY . .
%s
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}" \
    -o bin/%s ./cmd/%s

# Runtime stage
FROM git.web3gate.ru:5000/alpine:3.18.3

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder /app/bin/%s .
COPY --from=builder /app/config.yaml .

RUN chown -R appuser:appuser /app

USER appuser

%s

ENTRYPOINT ["/app/%s"]
`, bundleStep, p.Name, p.Name, p.Name, expose, p.Name)
}

func (p *Project) gitignoreTemplate() string {
	return `# Binaries
bin/
*.exe
*.dll
*.so
*.dylib

# Test
*.test
coverage.out
coverage.html

# IDE
.idea/
.vscode/
*.swp
*.swo

# Environment
.env
.env.local

# OS
.DS_Store
Thumbs.db

# Vendor (optional)
# vendor/

# Build
dist/

cmd/docs/*/openapi.json
example

.claude
`
}

func (p *Project) makefileTemplate() string {
	phony := ".PHONY: build run test lint clean docker proto fmt vet mod-download check"
	if p.Mode.HasOpenAPI() {
		phony += " bundle-spec generate-ogen"
	}

	return fmt.Sprintf(`BINARY_NAME := %s
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date -u +"%%Y-%%m-%%dT%%H:%%M:%%SZ")

%s

build:
	go build -ldflags="-s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)" \
		-o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

run: build
	STAGE=local	./bin/$(BINARY_NAME)

test:
	go test -race -cover -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

check : lint fmt vet

lint:
	golangci-lint run --timeout 5m

fmt:
	go fmt ./...

vet:
	go vet ./...

mod-download:
	go mod download
	go mod tidy

clean:
	rm -rf bin/ coverage.out

`, p.Name, phony)
}

func (p *Project) helmChartTemplate() string {
	return fmt.Sprintf(`apiVersion: v2
version: v3
name: %s
type: application
`, p.Name)
}

func (p *Project) helmValuesTemplate() string {
	return fmt.Sprintf(`app:
  project: aml
  app: %s
  repository: git.web3gate.ru:5000
  tag: latest

replicas: 1

k8s:
  namespace: aml-dev

resources:
  requests:
    cpu: 500m
    memory: 1024Mi
  limits:
    cpu: 750m
    memory: 1536Mi

ingress:
  class:
    name: nginx
  http:
    host: aml-%s.dev.web3gate.ru
    secret: aml-dev
    enabled: false
    port: 8080
  grpc:
    enabled: true
    port: 60006

srvmon:
  http:
    port: 8085
    enabled: true
  grpc:
    port: 60011
    enabled: true

metrics:
  http:
    port: 2112
    enabled: true

hpa:
  minReplicas: 1
  maxReplicas: 5

env:
  VAULT_SECRET_ID: ""
  VAULT_ROLE_ID: ""
  VAULT_ADDRESS: ""
  VAULT_SECRET_PATH: ""
  STAGE: ""
`, p.Name, p.Name)
}

func (p *Project) helmValuesProdTemplate() string {
	return fmt.Sprintf(`app:
  project: aml
  app: %s
  repository: git.web3gate.ru:5000
  tag: latest

replicas: 1

k8s:
  namespace: aml

resources:
  requests:
    cpu: 50m
    memory: 64Mi
  limits:
    cpu: 100m
    memory: 128Mi

ingress:
  class:
    name: nginx-infra
  http:
    host: aml-%s.web3gate.ru
    secret: aml
    enabled: false
    port: 8080
  grpc:
    enabled: true
    port: 60006

srvmon:
  http:
    port: 8085
    enabled: true
  grpc:
    port: 60011
    enabled: true

metrics:
  http:
    port: 2112
    enabled: true

hpa:
  minReplicas: 1
  maxReplicas: 5

env:
  VAULT_SECRET_ID: ""
  VAULT_ROLE_ID: ""
  VAULT_ADDRESS: ""
  VAULT_SECRET_PATH: ""
  STAGE: ""
`, p.Name, p.Name)
}

func (p *Project) helmDeploymentTemplate() string {
	return `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: "{{ .Release.Name }}"
  labels:
    k8s-app: "{{ .Release.Name }}"
    project: {{ .Values.app.project }}
  namespace: {{ .Values.k8s.namespace }}
spec:
  replicas: {{ .Values.replicas | default 1 }}
  selector:
    matchLabels:
      k8s-app: "{{ .Release.Name }}"
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 0
      maxSurge: 1
  template:
    metadata:
      labels:
        date: "{{ now | unixEpoch }}"
        k8s-app: "{{ .Release.Name }}"
    spec:
      terminationGracePeriodSeconds: 30
      restartPolicy: Always
      imagePullSecrets:
      - name: regcred
      containers:
      - name: "{{ .Release.Name }}"
        image: "{{ .Values.app.repository }}/{{ .Values.app.project }}/{{ .Values.app.app }}:{{ .Values.app.tag }}"
        imagePullPolicy: Always
        resources:
        {{- toYaml .Values.resources | nindent 10 }}
        {{- if .Values.ingress }}
        ports:
          {{- if and .Values.ingress.http .Values.ingress.http.enabled .Values.ingress.http.port }}
          - containerPort: {{ .Values.ingress.http.port }}
          {{- end }}
          {{- if and .Values.ingress.grpc .Values.ingress.grpc.enabled .Values.ingress.grpc.port }}
          - containerPort: {{ .Values.ingress.grpc.port }}
          {{- end }}
          {{- if and .Values.metrics.http .Values.metrics.http.enabled .Values.metrics.http.port }}
          - containerPort: {{ .Values.metrics.http.port }}
          {{- end }}
          {{- if and .Values.srvmon.http .Values.srvmon.http.enabled .Values.srvmon.http.port }}
          - containerPort: {{ .Values.srvmon.http.port }}
          {{- end }}
          {{- if and .Values.srvmon.grpc .Values.srvmon.grpc.enabled .Values.srvmon.grpc.port }}
          - containerPort: {{ .Values.srvmon.grpc.port }}
          {{- end }}
        {{- end }}
        env:
          - name: VAULT_SECRET_ID
            value: {{ .Values.env.VAULT_SECRET_ID }}
          - name: VAULT_ROLE_ID
            value: {{ .Values.env.VAULT_ROLE_ID }}
          - name: VAULT_ADDRESS
            value: {{ .Values.env.VAULT_ADDRESS }}
          - name: VAULT_SECRET_PATH
            value: {{ .Values.env.VAULT_SECRET_PATH }}
          - name: STAGE
            value: {{ .Values.env.STAGE }}
          - name: MY_NODE_NAME
            valueFrom:
              fieldRef:
                fieldPath: spec.nodeName
          - name: MY_POD_NAME
            valueFrom:
              fieldRef:
                fieldPath: metadata.name
          - name: MY_POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          - name: MY_POD_IP
            valueFrom:
              fieldRef:
                fieldPath: status.podIP
          - name: MY_POD_SERVICE_ACCOUNT
            valueFrom:
              fieldRef:
                fieldPath: spec.serviceAccountName
        readinessProbe:
          exec:
            command:
              - sh
              - -c
              - |
                wget -q -O- http://localhost:{{ .Values.srvmon.http.port }}/readyz | \
                grep -q '"ready":true'
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          successThreshold: 1
          failureThreshold: 3
        livenessProbe:
          exec:
            command:
              - sh
              - -c
              - |
                wget -q -O- http://localhost:{{ .Values.srvmon.http.port }}/healthz | \
                grep -qE '^\{"status":"(STATUS_UP|STATUS_DEGRADED)"'
          initialDelaySeconds: 15
          periodSeconds: 20
          timeoutSeconds: 3
          successThreshold: 1
          failureThreshold: 3
`
}

func (p *Project) helmServiceTemplate() string {
	return `{{- if .Values.ingress }}
---
apiVersion: v1
kind: Service
metadata:
  name: "{{ .Release.Name }}"
  namespace: {{ .Values.k8s.namespace }}
  labels:
    k8s-app: "{{ .Release.Name }}"
    project: {{ .Values.app.project }}
spec:
  ports:
  {{- if and .Values.ingress.http .Values.ingress.http.enabled .Values.ingress.http.port }}
  - protocol: TCP
    name: "tcp-{{ .Values.ingress.http.port }}"
    port: {{ .Values.ingress.http.port }}
    targetPort: {{ .Values.ingress.http.port }}
  {{- end }}
  {{- if and .Values.ingress.grpc .Values.ingress.grpc.enabled .Values.ingress.grpc.port }}
  - protocol: TCP
    name: "tcp-{{ .Values.ingress.grpc.port }}"
    port: {{ .Values.ingress.grpc.port }}
    targetPort: {{ .Values.ingress.grpc.port }}
  {{- end }}
  {{- if and .Values.metrics.http .Values.metrics.http.enabled .Values.metrics.http.port }}
  - protocol: TCP
    name: "tcp-{{ .Values.metrics.http.port }}"
    port: {{ .Values.metrics.http.port }}
    targetPort: {{ .Values.metrics.http.port }}
  {{- end }}
  {{- if and .Values.srvmon.http .Values.srvmon.http.enabled .Values.srvmon.http.port }}
  - protocol: TCP
    name: "tcp-{{ .Values.srvmon.http.port }}"
    port: {{ .Values.srvmon.http.port }}
    targetPort: {{ .Values.srvmon.http.port }}
  {{- end }}
  {{- if and .Values.srvmon.grpc .Values.srvmon.grpc.enabled .Values.srvmon.grpc.port }}
  - protocol: TCP
    name: "tcp-{{ .Values.srvmon.grpc.port }}"
    port: {{ .Values.srvmon.grpc.port }}
    targetPort: {{ .Values.srvmon.grpc.port }}
  {{- end }}

  selector:
    k8s-app: "{{ .Release.Name }}"
{{- end }}
`
}

func (p *Project) helmHPATemplate() string {
	return `---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: "{{ .Release.Name }}"
spec:
  behavior:
    scaleDown:
      policies:
      - periodSeconds: 60
        type: Pods
        value: 1
      selectPolicy: Max
      stabilizationWindowSeconds: 60
    scaleUp:
      policies:
      - periodSeconds: 15
        type: Pods
        value: 4
      - periodSeconds: 15
        type: Percent
        value: 100
      selectPolicy: Max
      stabilizationWindowSeconds: 0
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: "{{ .Release.Name }}"
  minReplicas: {{ .Values.hpa.minReplicas }}
  maxReplicas: {{ .Values.hpa.maxReplicas }}
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 60
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
`
}

func (p *Project) helmIngressTemplate() string {
	return `{{- if .Values.ingress }}
{{- if and .Values.ingress.http .Values.ingress.http.enabled .Values.ingress.http.port }}
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: "{{ .Release.Name }}-http"
  namespace: {{ .Values.k8s.namespace }}
  labels:
    project: {{ .Values.app.project }}
spec:
  ingressClassName: {{ .Values.ingress.class.name }}
  rules:
    - host: {{ .Values.ingress.http.host }}
      http:
        paths:
          - path: "/"
            pathType: Prefix
            backend:
              service:
                name: "{{ .Release.Name }}"
                port:
                  number: {{ .Values.ingress.http.port }}
  tls:
    - hosts:
        - {{ .Values.ingress.http.host }}
      secretName: {{ .Values.ingress.http.secret }}
{{- end }}
{{- end }}
`
}

func (p *Project) dockerComposeMonitoringTemplate() string {
	return `version: '3'
services:
  grafana:
    image: grafana/grafana
    ports:
      - 3000:3000
    networks:
      - monitoring
  prometheus:
    image: prom/prometheus
    ports:
      - 9090:9090
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - monitoring
networks:
  monitoring:
`
}

func (p *Project) prometheusConfigTemplate() string {
	return `global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['host.docker.internal:2112']
`
}

func (p *Project) ciTemplate() string {
	return fmt.Sprintf(`variables:
  NAMESPACE: web3gate
  SERVICE_NAME: %s
  HELM_PATH: deploy/helm

stages:
  - lint
  - test
  - build
  - deploy

.go_setup:
  variables:
    GOPATH: $CI_PROJECT_DIR/.go
  cache:
    key: ${CI_COMMIT_REF_SLUG}-go
    paths:
      - .go/pkg/mod/
  before_script:
    - git config --global url."https://gitlab-ci-token:${CI_JOB_TOKEN}@git.web3gate.ru/".insteadOf "https://git.web3gate.ru/"
    - go env -w GOSUMDB=off
    - go env -w GOPRIVATE="git.web3gate.ru"

.code_changes: &code_changes
  - "cmd/**/*"
  - "internal/**/*"
  - "pkg/**/*"
  - "api/**/*"
  - "*.go"
  - "go.mod"
  - "go.sum"

.rules:code:
  rules:
    - if: $CI_COMMIT_BRANCH
      changes: *code_changes
    - if: $CI_COMMIT_BRANCH
      when: manual
      allow_failure: true

.rules:build:
  rules:
    - if: $CI_COMMIT_BRANCH == "main" || $CI_COMMIT_BRANCH == "develop"
      changes: *code_changes
    - if: $CI_COMMIT_BRANCH
      when: manual
      allow_failure: true

.rules:deploy-dev:
  rules:
    - if: $CI_COMMIT_BRANCH == "develop"
    - if: $CI_COMMIT_BRANCH
      when: manual
      allow_failure: true

.rules:deploy-prod:
  rules:
    - if: $CI_COMMIT_BRANCH == "main"
      when: manual
    - when: never

lint:
  stage: lint
  image: git.web3gate.ru:5000/golangci/golangci-lint:v2.9.0-alpine
  extends:
    - .go_setup
    - .rules:code
  script:
    - golangci-lint run --timeout 5m

test:
  stage: test
  image: git.web3gate.ru:5000/golang:1.25.5
  extends:
    - .go_setup
    - .rules:code
  script:
    - go test -race -cover -coverprofile=coverage.out ./...
    - go tool cover -func=coverage.out
  coverage: '/total:\s+\(statements\)\s+(\d+\.\d+)%%/'
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.out
    expire_in: 7 days

build:
  stage: build
  image: git.web3gate.ru:5000/docker:24.0.6-git
  extends: .rules:build
  before_script:
    - echo "$DOCKER_REGISTRY_PASS" | docker login -u "$DOCKER_REGISTRY_LOGIN" --password-stdin "$DOCKER_REGISTRY_FQN"
  script:
    - |
      TAG="${CI_COMMIT_REF_SLUG}.${CI_COMMIT_SHA:0:8}"
      IMAGE="${DOCKER_REGISTRY_FQN}/aml/${SERVICE_NAME}"

      docker build \
        --build-arg NETRC="$(cat $NETRC)" \
        --build-arg VERSION="${CI_COMMIT_SHA:0:8}" \
        --build-arg BUILD_TIME="$(date -u +%%Y-%%m-%%dT%%H:%%M:%%SZ)" \
        --cache-from "${IMAGE}:latest" \
        -t "${IMAGE}:${TAG}" \
        .

      docker push "${IMAGE}:${TAG}"

      if [ "$CI_COMMIT_BRANCH" = "main" ]; then
        docker tag "${IMAGE}:${TAG}" "${IMAGE}:production.latest"
        docker push "${IMAGE}:production.latest"
      elif [ "$CI_COMMIT_BRANCH" = "develop" ]; then
        docker tag "${IMAGE}:${TAG}" "${IMAGE}:development.latest"
        docker push "${IMAGE}:development.latest"
      fi

deploy:dev:
  stage: deploy
  image: git.web3gate.ru:5000/dtzar/helm-kubectl:3.13
  extends: .rules:deploy-dev
  needs: [build]
  environment: dev
  script:
    - kubectl config use-context aml/aml-infra:dev-agent
    - |
      helm upgrade --install ${SERVICE_NAME}-dev ${HELM_PATH} \
        --set env.VAULT_SECRET_ID=$VAULT_SECRET_ID \
        --set env.VAULT_ROLE_ID=$VAULT_ROLE_ID \
        --set env.VAULT_SECRET_PATH=$VAULT_SECRET_PATH \
        --set env.VAULT_ADDRESS=$VAULT_ADDRESS \
        --set env.STAGE=dev \
        --set app.tag=development.latest \
        --values ${HELM_PATH}/values.yaml \
        --namespace ${NAMESPACE}-dev \
        --atomic --timeout 120s

deploy:prod:
  stage: deploy
  image: git.web3gate.ru:5000/dtzar/helm-kubectl:3.13
  extends: .rules:deploy-prod
  needs: [build]
  environment: prod
  script:
    - kubectl config use-context aml/aml-infra:prod-agent
    - |
      helm upgrade --install ${SERVICE_NAME} ${HELM_PATH} \
        --set env.VAULT_SECRET_ID=$PROD_VAULT_SECRET_ID \
        --set env.VAULT_ROLE_ID=$PROD_VAULT_ROLE_ID \
        --set env.VAULT_SECRET_PATH=$PROD_VAULT_SECRET_PATH \
        --set env.VAULT_ADDRESS=$PROD_VAULT_ADDRESS \
        --set env.STAGE=prod \
        --set app.tag=production.latest \
        --values ${HELM_PATH}/values-prod.yaml \
        --namespace ${NAMESPACE} \
        --atomic --timeout 120s
`, p.Name)
}
