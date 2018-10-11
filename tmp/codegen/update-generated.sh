#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

vendor/k8s.io/code-generator/generate-groups.sh \
deepcopy \
github.com/openshift/keycloak-operator/pkg/generated \
github.com/openshift/keycloak-operator/pkg/apis \
keycloak:v1alpha \
--go-header-file "./tmp/codegen/boilerplate.go.txt"
