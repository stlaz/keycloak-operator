package main

import (
	"context"
	"runtime"
	"time"

	"github.com/openshift/keycloak-operator/pkg/keycloak"
	"github.com/openshift/keycloak-operator/pkg/stub"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	sdkVersion "github.com/operator-framework/operator-sdk/version"

	"github.com/sirupsen/logrus"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	printVersion()

	h := stub.NewHandler()

	resource := "keycloak.config.openshift.io/v1alpha"
	kind := "KeycloakOperator"
	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		logrus.Fatalf("failed to get watch namespace: %v", err)
	}
	resyncPeriod := time.Duration(10) * time.Minute
	logrus.Infof("Watching %s, %s, %s, %d", resource, kind, namespace, resyncPeriod)

	config := keycloak.KeycloakConfig{
		appName:                "keycloak",
		adminUsername:          "admin",
		adminPassword:          "leCheval123",
		proxyAddressForwarding: true,
		dbVendor:               keycloak.DBH2,
		loglevel:               keycloak.LOGDEBUG,
	}
	keycloak.CreateNewDeployment(config)

	sdk.Watch(resource, kind, namespace, resyncPeriod)
	sdk.Handle(h)
	sdk.Run(context.TODO())
}
