package stub

import (
	"context"

	"github.com/openshift/keycloak-operator/pkg/keycloak"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	//switch o := event.Object.(type) {
	// case *v1alpha.KeycloakOperator:
	// 	err := sdk.Create(newbusyBoxPod(o))
	// 	if err != nil && !errors.IsAlreadyExists(err) {
	// 		logrus.Errorf("failed to create busybox pod : %v", err)
	// 		return err
	// 	}
	// }
	config := keycloak.KeycloakConfig{
		AppName:                "keycloak",
		AdminUsername:          "admin",
		AdminPassword:          "leCheval123",
		ProxyAddressForwarding: true,
		DbVendor:               keycloak.DBH2,
		Loglevel:               keycloak.LOGDEBUG,
	}
	keycloak.CreateNewDeployment(config)
	return nil
}
