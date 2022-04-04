
BIN_NAME = ./out/kubmetheus
VENDOR_FILES := $(shell find vendor/)
IMAGE_CHECK_FILE = out/build-img.check
K8S_DEPLOYMENT = kubmetheus-deployment
K8S_SERVICE = kubmetheus-svc

K8S_APPLY_FILE = kubectl apply -f

.PHONY: build
build: ${BIN_NAME}
	
run: ${BIN_NAME}
	$(shell $^)
clean:
	rm out/*

${BIN_NAME}: main.go
	go build  -gcflags=all="-N -l" -o $@ .

.PHONY: docker-image
docker-image: ${IMAGE_CHECK_FILE}

${IMAGE_CHECK_FILE}: ${BIN_NAME} Dockerfile 
	$(shell go mod vendor)
	docker build -t 10.100.124.45:5000/kubmetheus:0.1 .
	docker push 10.100.124.45:5000/kubmetheus:0.1
	touch $@

test:
	@echo $(DOCKER_HOST)

.PHONY: undeploy
undeploy: 
	kubectl delete deployments,pods,services,roles,rolebindings,serviceaccounts -l app=kubmetheus
	kubectl delete deployments,pods,services -l app=prometheus
	rm out/deploy.check out/service.check

.PHONY: deploy service restart
deploy: out/deploy.check
service: out/service.check
restart: ${IMAGE_CHECK_FILE}
	kubectl rollout restart	deployment ${K8S_DEPLOYMENT}


out/deploy.check: deploy/deployment.yml Dockerfile ${BIN_NAME} ${IMAGE_CHECK_FILE}
	$(K8S_APPLY_FILE) deploy/deployment.yml
	touch $@

out/service.check: deploy/service.yml 
	$(K8S_APPLY_FILE) deploy/service.yml
	touch $@


