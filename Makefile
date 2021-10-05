init:
	terraform init

.PHONY: fmt
fmt:
	terraform fmt

.PHONY: validate
validate:
	terraform validate

apply: fmt validate
	terraform apply

destroy:
	terraform destroy
