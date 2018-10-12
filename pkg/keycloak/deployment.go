package keycloak

import (
	"fmt"

	apps "github.com/openshift/api/apps/v1"
	routev1 "github.com/openshift/api/route/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
)

const (
	DBH2       = "H2"
	DBPOSTGRES = "POSTGRES"
	DBMYSQL    = "MYSQL"
	DBMARIADB  = "MARIADB"

	LOGDEBUG = "DEBUG"
	LOGINFO  = "INFO"
	LOGERROR = "ERROR"

	CABundlePath = "/etc/x509/client/tls.crt"
)

type KeycloakConfig struct {
	appName                string
	adminUsername          string
	adminPassword          string
	proxyAddressForwarding bool
	dbVendor               string
	loglevel               string
	hostnameHTTP           string
	hostnameHTTPS          string
}

func CreateNewDeployment(config KeycloakConfig) {) {
	// TODO:
	//     call adm ca create-server-cert
	//     separate only certs from the ^.crt file (possibly not necessary?)
	// 	   create server and client tls secrets

	//  create a keycloak deployment
	sdk.Create(newPersistenVolumeClaim(config.appName + "-data"))
	sdk.Create(newService(config.appName))
	skd.Create(newDeployment(config))
}

func newPersistenVolumeClaim(name string) *corev1.PersistentVolumeClaim {
	resourceMap := corev1.ResourceList{}
	resourceMap["storage"] = resource.MustParse("1Gi")
	return &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			Resources: corev1.ResourceRequirements{
				Requests: resourceMap,
			},
		},
	}
}

func newService(appName string) *corev1.Service {
	labels := map[string]string{}
	annotations := map[string]string{}
	labels["application"] = appName
	annotations["description"] = "The web server's https port."

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        appName,
			Annotations: annotations,
			Labels:      labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Port:       8443,
					TargetPort: intstr.FromInt(8443),
				},
			},
		},
	}
}

func newRoute(appName, httpsHost string) *routev1.Route {
	labels := map[string]string{}
	annotations := map[string]string{}
	labels["application"] = appName
	annotations["description"] = "Route for application's https service."

	return &routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: "v1",
			// Id:         appName,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        appName,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: routev1.RouteSpec{
			Host: httpsHost,
			To: routev1.RouteTargetReference{
				Name: appName,
			},
			TLS: &routev1.TLSConfig{
				Termination: "passthrough",
			},
		},
	}
}

func newDeployment(config KeycloakConfig) *apps.DeploymentConfig {
	labels := map[string]string{}
	selectors := map[string]string{}
	podLabels := map[string]string{}
	labels["application"] = config.appName
	selectors["deploymentConfig"] = config.appName
	podLabels["deploymentConfig"] = config.appName
	podLabels["application"] = config.appName

	return &apps.DeploymentConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DeploymentConfig",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   config.appName,
			Labels: labels,
		},
		Spec: apps.DeploymentConfigSpec{
			Strategy: apps.DeploymentStrategy{Type: "recreate"},
			Triggers: apps.DeploymentTriggerPolicies{
				{Type: "ConfigChange"},
			},
			Replicas: 1,
			Selector: selectors,
			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: config.appName,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name: config.appName,
							// TODO: have the image configurable
							Image:           "jboss/keycloak:nightly-openshift-integration2",
							ImagePullPolicy: "always",
							Args: []string{
								"-Djavax.net.debug=ssl:handshake:verbose",
								"-b",
								"0.0.0.0", // XXX: maybe reserved for bootstrap?
								"--debug",
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
									Protocol:      "TCP",
								},
								{
									Name:          "https",
									ContainerPort: 8443,
									Protocol:      "TCP",
								},
							},
							Env: []corev1.EnvVar{
								{Name: "KEYCLOAK_USER", Value: config.adminUsername},
								{Name: "KEYCLOAK_PASSWORD", Value: config.adminPassword},
								{Name: "DB_VENDOR", Value: config.dbVendor},
								{Name: "PROXY_ADDRESS_FORWARDING", Value: fmt.Sprintf("%v", config.proxyAddressForwarding)},
								{Name: "KEYCLOAK_LOGLEVEL", Value: config.loglevel},
								{Name: "X509_CA_BUNDLE", Value: CABundlePath},
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: boolAddr(false),
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "auth/realms/master",
										Port: intstr.FromInt(8080),
									},
								},
								InitialDelaySeconds: 10,
								PeriodSeconds:       10,
								TimeoutSeconds:      1,
								FailureThreshold:    5,
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      config.appName + "-data",
									MountPath: "/opt/jboss/keycloak/standalone/data",
									ReadOnly:  false,
								},
								{
									Name:      "keycloak-server-tls-volume",
									MountPath: "/etc/x509/https",
									ReadOnly:  true,
								},
								{
									Name:      "keycloak-client-tls-volume",
									MountPath: "/etc/x509/client",
									ReadOnly:  true,
								},
							},
						}, // Container
					}, // []Containers
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: "keycloak-server-tls-volume",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "keycloak-server-tls",
								},
							},
						},
						corev1.Volume{
							Name: "keycloak-client-tls-volume",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "keycloak-client-tls",
								},
							},
						},
						corev1.Volume{
							Name: "keycloak-data",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: config.appName + "-data",
								},
							},
						},
					},
				},
			},
		},
	}
}

func boolAddr(b bool) *bool {
	return &b
}
