KUBE_CONFIG := --kubeconfig=$(DEPLOY)/kube_config.yaml
NAMESPACE := notice

.PHONY: k8s.create.ns
k8s.create.ns:
	@kubectl $(KUBE_CONFIG) create namespace $(NAMESPACE)

.PHONY: k8s.create.secret
k8s.create.secret:
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) create secret generic notice-config --from-file=$(DEPLOY)/config.yaml

.PHONY: k8s.apply.secret
k8s.apply.secret:
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) scale deployment notice-dep --replicas=0
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) delete secrets notice-config
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) create secret generic notice-config --from-file=$(DEPLOY)/config.yaml
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) scale deployment notice-dep --replicas=1

.PHONY: k8s.scale.%
k8s.scale.%:
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) scale deployment notice-dep --replicas=$*

.PHONY: k8s.create.deploy
k8s.create.deploy:
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) create -f $(DEPLOY)/deployment.yaml

.PHONY: k8s.set.image
k8s.set.image:
	@kubectl $(KUBE_CONFIG) -n $(NAMESPACE) set image deployment notice-dep notice-c=$(DOCKER_TAG)