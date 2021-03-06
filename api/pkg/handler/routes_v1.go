package handler

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/kubermatic/kubermatic/api/pkg/handler/middleware"
	"github.com/kubermatic/kubermatic/api/pkg/handler/v1"
	"github.com/kubermatic/kubermatic/api/pkg/handler/v1/addon"
	"github.com/kubermatic/kubermatic/api/pkg/handler/v1/cluster"
	"github.com/kubermatic/kubermatic/api/pkg/handler/v1/common"
	"github.com/kubermatic/kubermatic/api/pkg/handler/v1/dc"
	kubernetesdashboard "github.com/kubermatic/kubermatic/api/pkg/handler/v1/kubernetes-dashboard"
	"github.com/kubermatic/kubermatic/api/pkg/handler/v1/label"
	"github.com/kubermatic/kubermatic/api/pkg/handler/v1/node"
	"github.com/kubermatic/kubermatic/api/pkg/handler/v1/openshift"
	"github.com/kubermatic/kubermatic/api/pkg/handler/v1/presets"
	"github.com/kubermatic/kubermatic/api/pkg/handler/v1/project"
	"github.com/kubermatic/kubermatic/api/pkg/handler/v1/provider"
	"github.com/kubermatic/kubermatic/api/pkg/handler/v1/serviceaccount"
	"github.com/kubermatic/kubermatic/api/pkg/handler/v1/ssh"
	"github.com/kubermatic/kubermatic/api/pkg/handler/v1/user"
)

// RegisterV1 declares all router paths for v1
func (r Routing) RegisterV1(mux *mux.Router, metrics common.ServerMetrics) {
	//
	// no-op endpoint that always returns HTTP 200
	mux.Methods(http.MethodGet).
		Path("/healthz").
		HandlerFunc(statusOK)
	//
	// Defines endpoints for managing data centers
	mux.Methods(http.MethodGet).
		Path("/dc").
		Handler(r.datacentersHandler())

	mux.Methods(http.MethodGet).
		Path("/dc/{dc}").
		Handler(r.datacenterHandler())

	//
	// Defines a set of HTTP endpoint for interacting with
	// various cloud providers
	mux.Methods(http.MethodGet).
		Path("/providers/aws/sizes").
		Handler(r.listAWSSizes())

	mux.Methods(http.MethodGet).
		Path("/providers/aws/{dc}/subnets").
		Handler(r.listAWSSubnets())

	mux.Methods(http.MethodGet).
		Path("/providers/aws/{dc}/vpcs").
		Handler(r.listAWSVPCS())

	mux.Methods(http.MethodGet).
		Path("/providers/gcp/disktypes").
		Handler(r.listGCPDiskTypes())

	mux.Methods(http.MethodGet).
		Path("/providers/gcp/sizes").
		Handler(r.listGCPSizes())

	mux.Methods(http.MethodGet).
		Path("/providers/gcp/{dc}/zones").
		Handler(r.listGCPZones())

	mux.Methods(http.MethodGet).
		Path("/providers/digitalocean/sizes").
		Handler(r.listDigitaloceanSizes())

	mux.Methods(http.MethodGet).
		Path("/providers/azure/sizes").
		Handler(r.listAzureSizes())

	mux.Methods(http.MethodGet).
		Path("/providers/openstack/sizes").
		Handler(r.listOpenstackSizes())

	mux.Methods(http.MethodGet).
		Path("/providers/openstack/tenants").
		Handler(r.listOpenstackTenants())

	mux.Methods(http.MethodGet).
		Path("/providers/openstack/networks").
		Handler(r.listOpenstackNetworks())

	mux.Methods(http.MethodGet).
		Path("/providers/openstack/securitygroups").
		Handler(r.listOpenstackSecurityGroups())

	mux.Methods(http.MethodGet).
		Path("/providers/openstack/subnets").
		Handler(r.listOpenstackSubnets())

	mux.Methods(http.MethodGet).
		Path("/version").
		Handler(r.getKubermaticVersion())

	mux.Methods(http.MethodGet).
		Path("/providers/vsphere/networks").
		Handler(r.listVSphereNetworks())

	mux.Methods(http.MethodGet).
		Path("/providers/vsphere/folders").
		Handler(r.listVSphereFolders())

	mux.Methods(http.MethodGet).
		Path("/providers/packet/sizes").
		Handler(r.listPacketSizes())

	mux.Methods(http.MethodGet).
		Path("/providers/hetzner/sizes").
		Handler(r.listHetznerSizes())

	mux.Methods(http.MethodGet).
		Path("/providers/{provider_name}/presets/credentials").
		Handler(r.listCredentials())

	//
	// Defines a set of HTTP endpoints for project resource
	mux.Methods(http.MethodGet).
		Path("/projects").
		Handler(r.listProjects())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}").
		Handler(r.getProject())

	mux.Methods(http.MethodPost).
		Path("/projects").
		Handler(r.createProject())

	mux.Methods(http.MethodPut).
		Path("/projects/{project_id}").
		Handler(r.updateProject())

	mux.Methods(http.MethodDelete).
		Path("/projects/{project_id}").
		Handler(r.deleteProject())

	//
	// Defines a set of HTTP endpoints for SSH Keys that belong to a project
	mux.Methods(http.MethodPost).
		Path("/projects/{project_id}/sshkeys").
		Handler(r.createSSHKey())

	mux.Methods(http.MethodDelete).
		Path("/projects/{project_id}/sshkeys/{key_id}").
		Handler(r.deleteSSHKey())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/sshkeys").
		Handler(r.listSSHKeys())

	//
	// Defines a set of HTTP endpoints for cluster that belong to a project.
	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/clusters").
		Handler(r.listClustersForProject())

	mux.Methods(http.MethodPost).
		Path("/projects/{project_id}/dc/{dc}/clusters").
		Handler(r.createCluster(metrics.InitNodeDeploymentFailures))

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters").
		Handler(r.listClusters())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}").
		Handler(r.getCluster())

	mux.Methods(http.MethodPatch).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}").
		Handler(r.patchCluster())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/events").
		Handler(r.getClusterEvents())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/kubeconfig").
		Handler(r.getClusterKubeconfig())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/oidckubeconfig").
		Handler(r.getOidcClusterKubeconfig())

	mux.Methods(http.MethodDelete).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}").
		Handler(r.deleteCluster())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/health").
		Handler(r.getClusterHealth())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/upgrades").
		Handler(r.getClusterUpgrades())

	mux.Methods(http.MethodPut).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodes/upgrades").
		Handler(r.upgradeClusterNodeDeployments())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/metrics").
		Handler(r.getClusterMetrics())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/namespaces").
		Handler(r.listNamespace())

	mux.Methods(http.MethodPost).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles").
		Handler(r.createClusterRole())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles").
		Handler(r.listClusterRole())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id}").
		Handler(r.getClusterRole())

	mux.Methods(http.MethodDelete).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id}").
		Handler(r.deleteClusterRole())

	mux.Methods(http.MethodPatch).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id}").
		Handler(r.patchClusterRole())

	mux.Methods(http.MethodPost).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id}/clusterbindings").
		Handler(r.createClusterRoleBinding())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id}/clusterbindings").
		Handler(r.listClusterRoleBinding())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id}/clusterbindings/{binding_id}").
		Handler(r.getClusterRoleBinding())

	mux.Methods(http.MethodDelete).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id}/clusterbindings/{binding_id}").
		Handler(r.deleteClusterRoleBinding())

	mux.Methods(http.MethodPatch).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id}/clusterbindings/{binding_id}").
		Handler(r.patchClusterRoleBinding())

	mux.Methods(http.MethodPost).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles").
		Handler(r.createRole())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles").
		Handler(r.listRole())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}").
		Handler(r.getRole())

	mux.Methods(http.MethodDelete).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}").
		Handler(r.deleteRole())

	mux.Methods(http.MethodPatch).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}").
		Handler(r.patchRole())

	mux.Methods(http.MethodPost).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings").
		Handler(r.createRoleBinding())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings").
		Handler(r.listRoleBinding())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings/{binding_id}").
		Handler(r.getRoleBinding())

	mux.Methods(http.MethodDelete).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings/{binding_id}").
		Handler(r.deleteRoleBinding())

	mux.Methods(http.MethodPatch).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings/{binding_id}").
		Handler(r.patchRoleBinding())

	//
	// Defines set of HTTP endpoints for SSH Keys that belong to a cluster
	mux.Methods(http.MethodPut).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/sshkeys/{key_id}").
		Handler(r.assignSSHKeyToCluster())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/sshkeys").
		Handler(r.listSSHKeysAssignedToCluster())

	mux.Methods(http.MethodDelete).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/sshkeys/{key_id}").
		Handler(r.detachSSHKeyFromCluster())

	mux.Methods(http.MethodPut).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/token").
		Handler(r.revokeClusterAdminToken())

	mux.Methods(http.MethodPut).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/viewertoken").
		Handler(r.revokeClusterViewerToken())

	//
	// Defines a set of HTTP endpoint for node deployments that belong to a cluster
	mux.Methods(http.MethodPost).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments").
		Handler(r.createNodeDeployment())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments").
		Handler(r.listNodeDeployments())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments/{nodedeployment_id}").
		Handler(r.getNodeDeployment())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments/{nodedeployment_id}/nodes").
		Handler(r.listNodeDeploymentNodes())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments/{nodedeployment_id}/nodes/metrics").
		Handler(r.listNodeDeploymentMetrics())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments/{nodedeployment_id}/nodes/events").
		Handler(r.listNodeDeploymentNodesEvents())

	mux.Methods(http.MethodPatch).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments/{nodedeployment_id}").
		Handler(r.patchNodeDeployment())

	mux.Methods(http.MethodDelete).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments/{nodedeployment_id}").
		Handler(r.deleteNodeDeployment())

	//
	// Defines a set of HTTP endpoints for managing addons

	mux.Methods(http.MethodGet).
		Path("/addons").
		Handler(r.listAccessibleAddons())

	mux.Methods(http.MethodPost).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/addons").
		Handler(r.createAddon())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/addons").
		Handler(r.listAddons())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/addons/{addon_id}").
		Handler(r.getAddon())

	mux.Methods(http.MethodPatch).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/addons/{addon_id}").
		Handler(r.patchAddon())

	mux.Methods(http.MethodDelete).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/addons/{addon_id}").
		Handler(r.deleteAddon())

	//
	// Defines a set of HTTP endpoints for various cloud providers
	// Note that these endpoints don't require credentials as opposed to the ones defined under /providers/*
	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/aws/sizes").
		Handler(r.listAWSSizesNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/aws/subnets").
		Handler(r.listAWSSubnetsNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/gcp/disktypes").
		Handler(r.listGCPDiskTypesNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/gcp/sizes").
		Handler(r.listGCPSizesNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/gcp/zones").
		Handler(r.listGCPZonesNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/hetzner/sizes").
		Handler(r.listHetznerSizesNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/digitalocean/sizes").
		Handler(r.listDigitaloceanSizesNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/azure/sizes").
		Handler(r.listAzureSizesNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/sizes").
		Handler(r.listOpenstackSizesNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/tenants").
		Handler(r.listOpenstackTenantsNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/networks").
		Handler(r.listOpenstackNetworksNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/securitygroups").
		Handler(r.listOpenstackSecurityGroupsNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/subnets").
		Handler(r.listOpenstackSubnetsNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/vsphere/networks").
		Handler(r.listVSphereNetworksNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/vsphere/folders").
		Handler(r.listVSphereFoldersNoCredentials())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/packet/sizes").
		Handler(r.listPacketSizesNoCredentials())

	//
	// Defines a set of openshift-specific endpoints
	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/openshift/console/login").
		Handler(r.openshiftConsoleLogin())
	mux.PathPrefix("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/openshift/console/proxy").
		Handler(r.openshiftConsoleProxy())

	//
	// Defines a set of kubernetes-dashboard-specific endpoints
	mux.PathPrefix("/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/dashboard/proxy").
		Handler(r.kubernetesDashboardProxy())

	//
	// Defines set of HTTP endpoints for Users of the given project
	mux.Methods(http.MethodPost).
		Path("/projects/{project_id}/users").
		Handler(r.addUserToProject())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/users").
		Handler(r.getUsersForProject())

	mux.Methods(http.MethodPut).
		Path("/projects/{project_id}/users/{user_id}").
		Handler(r.editUserInProject())

	mux.Methods(http.MethodDelete).
		Path("/projects/{project_id}/users/{user_id}").
		Handler(r.deleteUserFromProject())

	//
	// Defines set of HTTP endpoints for ServiceAccounts of the given project
	mux.Methods(http.MethodPost).
		Path("/projects/{project_id}/serviceaccounts").
		Handler(r.addServiceAccountToProject())

	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/serviceaccounts").
		Handler(r.listServiceAccounts())

	mux.Methods(http.MethodPut).
		Path("/projects/{project_id}/serviceaccounts/{serviceaccount_id}").
		Handler(r.updateServiceAccount())

	mux.Methods(http.MethodDelete).
		Path("/projects/{project_id}/serviceaccounts/{serviceaccount_id}").
		Handler(r.deleteServiceAccount())

	//
	// Defines set of HTTP endpoints for tokens of the given service account
	mux.Methods(http.MethodPost).
		Path("/projects/{project_id}/serviceaccounts/{serviceaccount_id}/tokens").
		Handler(r.addTokenToServiceAccount())
	mux.Methods(http.MethodGet).
		Path("/projects/{project_id}/serviceaccounts/{serviceaccount_id}/tokens").
		Handler(r.listServiceAccountTokens())
	mux.Methods(http.MethodPut).
		Path("/projects/{project_id}/serviceaccounts/{serviceaccount_id}/tokens/{token_id}").
		Handler(r.updateServiceAccountToken())
	mux.Methods(http.MethodPatch).
		Path("/projects/{project_id}/serviceaccounts/{serviceaccount_id}/tokens/{token_id}").
		Handler(r.patchServiceAccountToken())
	mux.Methods(http.MethodDelete).
		Path("/projects/{project_id}/serviceaccounts/{serviceaccount_id}/tokens/{token_id}").
		Handler(r.deleteServiceAccountToken())

	//
	// Defines set of HTTP endpoints for control plane and kubelet versions
	mux.Methods(http.MethodGet).
		Path("/upgrades/cluster").
		Handler(r.getMasterVersions())

	mux.Methods(http.MethodGet).
		Path("/upgrades/node").
		Handler(r.getNodeUpgrades())

	//
	// Defines an endpoint to retrieve information about the current token owner
	mux.Methods(http.MethodGet).
		Path("/me").
		Handler(r.getCurrentUser())

	mux.Methods(http.MethodGet).
		Path("/labels/system").
		Handler(r.listSystemLabels())
}

// swagger:route GET /api/v1/projects/{project_id}/sshkeys project listSSHKeys
//
//     Lists SSH Keys that belong to the given project.
//     The returned collection is sorted by creation timestamp.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []SSHKey
//       401: empty
//       403: empty
func (r Routing) listSSHKeys() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(ssh.ListEndpoint(r.sshKeyProvider, r.projectProvider)),
		ssh.DecodeListReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route POST /api/v1/projects/{project_id}/sshkeys project createSSHKey
//
//    Adds the given SSH key to the specified project.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: SSHKey
//       401: empty
//       403: empty
func (r Routing) createSSHKey() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(ssh.CreateEndpoint(r.sshKeyProvider, r.projectProvider)),
		ssh.DecodeCreateReq,
		setStatusCreatedHeader(encodeJSON),
		r.defaultServerOptions()...,
	)
}

// swagger:route DELETE /api/v1/projects/{project_id}/sshkeys/{key_id} project deleteSSHKey
//
//     Removes the given SSH Key from the system.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) deleteSSHKey() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(ssh.DeleteEndpoint(r.sshKeyProvider, r.projectProvider)),
		ssh.DecodeDeleteReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/{provider_name}/presets/credentials credentials listCredentials
//
// Lists credential names for the provider
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: CredentialList
func (r Routing) listCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(presets.CredentialEndpoint(r.presetsManager)),
		presets.DecodeProviderReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/aws/sizes aws listAWSSizes
//
// Lists available AWS sizes.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: AWSSizeList
func (r Routing) listAWSSizes() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
		)(provider.AWSSizeEndpoint()),
		provider.DecodeAWSSizesReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/aws/{dc}/subnets aws listAWSSubnets
//
// Lists available AWS subnets
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: AWSSubnetList
func (r Routing) listAWSSubnets() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.AWSSubnetEndpoint(r.presetsManager, r.seedsGetter)),
		provider.DecodeAWSSubnetReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/aws/{dc}/vpcs aws listAWSVPCS
//
// Lists available AWS vpc's
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: AWSVPCList
func (r Routing) listAWSVPCS() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.AWSVPCEndpoint(r.presetsManager, r.seedsGetter)),
		provider.DecodeAWSVPCReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/gcp/disktypes gcp listGCPDiskTypes
//
// Lists disk types from GCP
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: GCPDiskTypeList
func (r Routing) listGCPDiskTypes() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.GCPDiskTypesEndpoint(r.presetsManager)),
		provider.DecodeGCPTypesReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/gcp/sizes gcp listGCPSizes
//
// Lists machine sizes from GCP
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: GCPMachineSizeList
func (r Routing) listGCPSizes() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.GCPSizeEndpoint(r.presetsManager)),
		provider.DecodeGCPTypesReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/gcp/{dc}/zones gcp listGCPZones
//
// Lists available GCP zones
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: GCPZoneList
func (r Routing) listGCPZones() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.GCPZoneEndpoint(r.presetsManager, r.seedsGetter)),
		provider.DecodeGCPZoneReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/digitalocean/sizes digitalocean listDigitaloceanSizes
//
// Lists sizes from digitalocean
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: DigitaloceanSizeList
func (r Routing) listDigitaloceanSizes() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.DigitaloceanSizeEndpoint(r.presetsManager)),
		provider.DecodeDoSizesReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/azure/sizes azure listAzureSizes
//
// Lists available VM sizes in an Azure region
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: AzureSizeList
func (r Routing) listAzureSizes() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.AzureSizeEndpoint(r.presetsManager)),
		provider.DecodeAzureSizesReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/openstack/sizes openstack listOpenstackSizes
//
// Lists sizes from openstack
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []OpenstackSize
func (r Routing) listOpenstackSizes() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.OpenstackSizeEndpoint(r.seedsGetter, r.presetsManager)),
		provider.DecodeOpenstackReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/vsphere/networks vsphere listVSphereNetworks
//
// Lists networks from vsphere datacenter
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []VSphereNetwork
func (r Routing) listVSphereNetworks() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.VsphereNetworksEndpoint(r.seedsGetter, r.presetsManager)),
		provider.DecodeVSphereNetworksReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/vsphere/folders vsphere listVSphereFolders
//
// Lists folders from vsphere datacenter
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []VSphereFolder
func (r Routing) listVSphereFolders() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.VsphereFoldersEndpoint(r.seedsGetter, r.presetsManager)),
		provider.DecodeVSphereFoldersReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/packet/sizes packet listPacketSizes
//
// Lists sizes from packet
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []PacketSizeList
func (r Routing) listPacketSizes() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.PacketSizesEndpoint(r.presetsManager)),
		provider.DecodePacketSizesReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/packet/sizes packet listPacketSizesNoCredentials
//
// Lists sizes from packet
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []PacketSizeList
func (r Routing) listPacketSizesNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.PacketSizesWithClusterCredentialsEndpoint(r.projectProvider)),
		provider.DecodePacketSizesNoCredentialsReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/openstack/tenants openstack listOpenstackTenants
//
// Lists tenants from openstack
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []OpenstackTenant
func (r Routing) listOpenstackTenants() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.OpenstackTenantEndpoint(r.seedsGetter, r.presetsManager)),
		provider.DecodeOpenstackTenantReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/openstack/networks openstack listOpenstackNetworks
//
// Lists networks from openstack
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []OpenstackNetwork
func (r Routing) listOpenstackNetworks() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.OpenstackNetworkEndpoint(r.seedsGetter, r.presetsManager)),
		provider.DecodeOpenstackReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/openstack/subnets openstack listOpenstackSubnets
//
// Lists subnets from openstack
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []OpenstackSubnet
func (r Routing) listOpenstackSubnets() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.OpenstackSubnetsEndpoint(r.seedsGetter, r.presetsManager)),
		provider.DecodeOpenstackSubnetReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/openstack/securitygroups openstack listOpenstackSecurityGroups
//
// Lists security groups from openstack
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []OpenstackSecurityGroup
func (r Routing) listOpenstackSecurityGroups() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.OpenstackSecurityGroupEndpoint(r.seedsGetter, r.presetsManager)),
		provider.DecodeOpenstackReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/providers/hetzner/sizes hetzner listHetznerSizes
//
// Lists sizes from hetzner
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: HetznerSizeList
func (r Routing) listHetznerSizes() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.HetznerSizeEndpoint(r.presetsManager)),
		provider.DecodeHetznerSizesReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/dc datacenter listDatacenters
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: DatacenterList
func (r Routing) datacentersHandler() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(dc.ListEndpoint(r.seedsGetter)),
		dc.DecodeDatacentersReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// Get the datacenter
// swagger:route GET /api/v1/dc/{dc} datacenter getDatacenter
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: Datacenter
func (r Routing) datacenterHandler() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(dc.GetEndpoint(r.seedsGetter)),
		dc.DecodeLegacyDcReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/versions versions getMasterVersions
// swagger:route GET /api/v1/upgrades/cluster versions getMasterVersions
//
// Lists all versions which don't result in automatic updates
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []MasterVersion
func (r Routing) getMasterVersions() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
		)(cluster.GetMasterVersionsEndpoint(r.updateManager)),
		cluster.DecodeClusterTypeReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/version versions getKubermaticVersion
//
// Get versions of running Kubermatic components.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: KubermaticVersions
func (r Routing) getKubermaticVersion() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
		)(v1.GetKubermaticVersion()),
		decodeEmptyReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects project listProjects
//
//     Lists projects that an authenticated user is a member of.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []Project
//       401: empty
//       409: empty
func (r Routing) listProjects() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(project.ListEndpoint(r.projectProvider, r.privilegedProjectProvider, r.userProjectMapper, r.projectMemberProvider, r.userProvider)),
		decodeEmptyReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id} project getProject
//
//     Gets the project with the given ID
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: Project
//       401: empty
//       409: empty
func (r Routing) getProject() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(project.GetEndpoint(r.projectProvider, r.projectMemberProvider, r.userProvider)),
		common.DecodeGetProject,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route POST /api/v1/projects project createProject
//
//     Creates a brand new project.
//
//     Note that this endpoint can be consumed by every authenticated user.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       201: Project
//       401: empty
//       409: empty
func (r Routing) createProject() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(project.CreateEndpoint(r.projectProvider)),
		project.DecodeCreate,
		setStatusCreatedHeader(encodeJSON),
		r.defaultServerOptions()...,
	)
}

// swagger:route PUT /api/v1/projects/{project_id} project updateProject
//
//    Updates the given project
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: Project
//       400: empty
//       404: empty
//       500: empty
//       501: empty
func (r Routing) updateProject() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(project.UpdateEndpoint(r.projectProvider, r.projectMemberProvider, r.userProvider)),
		project.DecodeUpdateRq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route DELETE /api/v1/projects/{project_id} project deleteProject
//
//    Deletes the project with the given ID.
//
//     Produces:
//     - application/json
//
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) deleteProject() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(project.DeleteEndpoint(r.projectProvider)),
		project.DecodeDelete,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route POST /api/v1/projects/{project_id}/dc/{dc}/clusters project createCluster
//
//     Creates a cluster for the given project.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       201: Cluster
//       401: empty
//       403: empty
func (r Routing) createCluster(initNodeDeploymentFailures *prometheus.CounterVec) http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.SetPrivilegedClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.CreateEndpoint(r.sshKeyProvider, r.projectProvider, r.seedsGetter, initNodeDeploymentFailures, r.eventRecorderProvider, r.presetsManager, r.exposeStrategy)),
		cluster.DecodeCreateReq,
		setStatusCreatedHeader(encodeJSON),
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters project listClusters
//
//     Lists clusters for the specified project and data center.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: ClusterList
//       401: empty
//       403: empty
func (r Routing) listClusters() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.ListEndpoint(r.projectProvider)),
		cluster.DecodeListReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/clusters project listClustersForProject
//
//     Lists clusters for the specified project.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: ClusterList
//       401: empty
//       403: empty
func (r Routing) listClustersForProject() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.ListAllEndpoint(r.projectProvider, r.seedsGetter, r.clusterProviderGetter)),
		common.DecodeGetProject,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id} project getCluster
//
//     Gets the cluster with the given name
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: Cluster
//       401: empty
//       403: empty
func (r Routing) getCluster() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.SetPrivilegedClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.GetEndpoint(r.projectProvider)),
		common.DecodeGetClusterReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PATCH /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id} project patchCluster
//
//     Patches the given cluster using JSON Merge Patch method (https://tools.ietf.org/html/rfc7396).
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: Cluster
//       401: empty
//       403: empty
func (r Routing) patchCluster() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.PatchEndpoint(r.projectProvider, r.seedsGetter)),
		cluster.DecodePatchReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// getClusterEvents returns events related to the cluster.
// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/events project getClusterEvents
//
//     Gets the events related to the specified cluster.
//
//     Produces:
//     - application/yaml
//
//     Responses:
//       default: errorResponse
//       200: []Event
//       401: empty
//       403: empty
func (r Routing) getClusterEvents() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.SetPrivilegedClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.GetClusterEventsEndpoint()),
		cluster.DecodeGetClusterEvents,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// getClusterKubeconfig returns the kubeconfig for the cluster.
// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/kubeconfig project getClusterKubeconfig
//
//     Gets the kubeconfig for the specified cluster.
//
//     Produces:
//     - application/yaml
//
//     Responses:
//       default: errorResponse
//       200: Kubeconfig
//       401: empty
//       403: empty
func (r Routing) getClusterKubeconfig() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.GetAdminKubeconfigEndpoint(r.projectProvider)),
		cluster.DecodeGetAdminKubeconfig,
		cluster.EncodeKubeconfig,
		r.defaultServerOptions()...,
	)
}

// getOidcClusterKubeconfig returns the oidc kubeconfig for the cluster.
// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/oidckubeconfig project getOidcClusterKubeconfig
//
//     Gets the kubeconfig for the specified cluster with oidc authentication.
//
//     Produces:
//     - application/yaml
//
//     Responses:
//       default: errorResponse
//       200: Kubeconfig
//       401: empty
//       403: empty
func (r Routing) getOidcClusterKubeconfig() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.GetOidcKubeconfigEndpoint(r.projectProvider)),
		cluster.DecodeGetAdminKubeconfig,
		cluster.EncodeKubeconfig,
		r.defaultServerOptions()...,
	)
}

// Delete the cluster
// swagger:route DELETE /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id} project deleteCluster
//
//     Deletes the specified cluster
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) deleteCluster() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.DeleteEndpoint(r.sshKeyProvider, r.projectProvider)),
		cluster.DecodeDeleteReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/health project getClusterHealth
//
//     Returns the cluster's component health status
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: ClusterHealth
//       401: empty
//       403: empty
func (r Routing) getClusterHealth() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.HealthEndpoint(r.projectProvider)),
		common.DecodeGetClusterReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PUT /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/sshkeys/{key_id} project assignSSHKeyToCluster
//
//     Assigns an existing ssh key to the given cluster
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       201: SSHKey
//       401: empty
//       403: empty
func (r Routing) assignSSHKeyToCluster() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.AssignSSHKeyEndpoint(r.sshKeyProvider, r.projectProvider)),
		cluster.DecodeAssignSSHKeyReq,
		setStatusCreatedHeader(encodeJSON),
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/sshkeys project listSSHKeysAssignedToCluster
//
//     Lists ssh keys that are assigned to the cluster
//     The returned collection is sorted by creation timestamp.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []SSHKey
//       401: empty
//       403: empty
func (r Routing) listSSHKeysAssignedToCluster() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.ListSSHKeysEndpoint(r.sshKeyProvider, r.projectProvider)),
		cluster.DecodeListSSHKeysReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route DELETE /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/sshkeys/{key_id} project detachSSHKeyFromCluster
//
//     Unassignes an ssh key from the given cluster
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) detachSSHKeyFromCluster() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.DetachSSHKeyEndpoint(r.sshKeyProvider, r.projectProvider)),
		cluster.DecodeDetachSSHKeysReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PUT /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/token project revokeClusterAdminToken
//
//     Revokes the current admin token
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) revokeClusterAdminToken() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.RevokeAdminTokenEndpoint(r.projectProvider)),
		cluster.DecodeAdminTokenReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PUT /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/viewertoken project revokeClusterViewerToken
//
//     Revokes the current viewer token
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) revokeClusterViewerToken() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.RevokeViewerTokenEndpoint(r.projectProvider)),
		cluster.DecodeAdminTokenReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/upgrades project getClusterUpgrades
//
//    Gets possible cluster upgrades
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []MasterVersion
//       401: empty
//       403: empty
func (r Routing) getClusterUpgrades() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.GetUpgradesEndpoint(r.updateManager, r.projectProvider)),
		common.DecodeGetClusterReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/upgrades/node versions getNodeUpgrades
//
//    Gets possible node upgrades for a specific control plane version
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []MasterVersion
//       401: empty
//       403: empty
func (r Routing) getNodeUpgrades() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
		)(cluster.GetNodeUpgrades(r.updateManager)),
		cluster.DecodeNodeUpgradesReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PUT /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodes/upgrades project upgradeClusterNodeDeployments
//
//    Upgrades node deployments in a cluster
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) upgradeClusterNodeDeployments() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.UpgradeNodeDeploymentsEndpoint(r.projectProvider)),
		cluster.DecodeUpgradeNodeDeploymentsReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route POST /api/v1/projects/{project_id}/users users addUserToProject
//
//     Adds the given user to the given project
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       201: User
//       401: empty
//       403: empty
func (r Routing) addUserToProject() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(user.AddEndpoint(r.projectProvider, r.userProvider, r.projectMemberProvider)),
		user.DecodeAddReq,
		setStatusCreatedHeader(encodeJSON),
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/users users getUsersForProject
//
//     Get list of users for the given project
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []User
//       401: empty
//       403: empty
func (r Routing) getUsersForProject() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(user.ListEndpoint(r.projectProvider, r.userProvider, r.projectMemberProvider)),
		common.DecodeGetProject,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PUT /api/v1/projects/{project_id}/users/{user_id} users editUserInProject
//
//     Changes membership of the given user for the given project
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: User
//       401: empty
//       403: empty
func (r Routing) editUserInProject() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(user.EditEndpoint(r.projectProvider, r.userProvider, r.projectMemberProvider)),
		user.DecodeEditReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route DELETE /api/v1/projects/{project_id}/users/{user_id} users deleteUserFromProject
//
//     Removes the given member from the project
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: User
//       401: empty
//       403: empty
func (r Routing) deleteUserFromProject() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(user.DeleteEndpoint(r.projectProvider, r.userProvider, r.projectMemberProvider)),
		user.DecodeDeleteReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/me users getCurrentUser
//
// Returns information about the current user.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: User
//       401: empty
func (r Routing) getCurrentUser() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
		)(user.GetEndpoint(r.userProjectMapper)),
		decodeEmptyReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route POST /api/v1/projects/{project_id}/serviceaccounts serviceaccounts addServiceAccountToProject
//
//     Adds the given service account to the given project
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       201: ServiceAccount
//       401: empty
//       403: empty
func (r Routing) addServiceAccountToProject() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(serviceaccount.CreateEndpoint(r.projectProvider, r.serviceAccountProvider)),
		serviceaccount.DecodeAddReq,
		setStatusCreatedHeader(encodeJSON),
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/serviceaccounts serviceaccounts listServiceAccounts
//
//     List Service Accounts for the given project
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []ServiceAccount
//       401: empty
//       403: empty
func (r Routing) listServiceAccounts() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(serviceaccount.ListEndpoint(r.projectProvider, r.serviceAccountProvider, r.userProjectMapper)),
		common.DecodeGetProject,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PUT /api/v1/projects/{project_id}/serviceaccounts/{serviceaccount_id} serviceaccounts updateServiceAccount
//
//     Updates service account for the given project
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: ServiceAccount
//       401: empty
//       403: empty
func (r Routing) updateServiceAccount() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(serviceaccount.UpdateEndpoint(r.projectProvider, r.serviceAccountProvider, r.userProjectMapper)),
		serviceaccount.DecodeUpdateReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route DELETE /api/v1/projects/{project_id}/serviceaccounts/{serviceaccount_id} serviceaccounts deleteServiceAccount
//
//     Deletes service account for the given project
//
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) deleteServiceAccount() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(serviceaccount.DeleteEndpoint(r.serviceAccountProvider, r.projectProvider)),
		serviceaccount.DecodeDeleteReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route POST /api/v1/projects/{project_id}/serviceaccounts/{serviceaccount_id}/tokens tokens addTokenToServiceAccount
//
//     Generates a token for the given service account
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       201: ServiceAccountToken
//       401: empty
//       403: empty
func (r Routing) addTokenToServiceAccount() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(serviceaccount.CreateTokenEndpoint(r.projectProvider, r.serviceAccountProvider, r.serviceAccountTokenProvider, r.saTokenAuthenticator, r.saTokenGenerator)),
		serviceaccount.DecodeAddTokenReq,
		setStatusCreatedHeader(encodeJSON),
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/serviceaccounts/{serviceaccount_id}/tokens tokens listServiceAccountTokens
//
//     List tokens for the given service account
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []PublicServiceAccountToken
//       401: empty
//       403: empty
func (r Routing) listServiceAccountTokens() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(serviceaccount.ListTokenEndpoint(r.projectProvider, r.serviceAccountProvider, r.serviceAccountTokenProvider, r.saTokenAuthenticator)),
		serviceaccount.DecodeTokenReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PUT /api/v1/projects/{project_id}/serviceaccounts/{serviceaccount_id}/tokens/{token_id} tokens updateServiceAccountToken
//
//     Updates and regenerates the token
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: ServiceAccountToken
//       401: empty
//       403: empty
func (r Routing) updateServiceAccountToken() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(serviceaccount.UpdateTokenEndpoint(r.projectProvider, r.serviceAccountProvider, r.serviceAccountTokenProvider, r.saTokenAuthenticator, r.saTokenGenerator)),
		serviceaccount.DecodeUpdateTokenReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PATCH /api/v1/projects/{project_id}/serviceaccounts/{serviceaccount_id}/tokens/{token_id} tokens patchServiceAccountToken
//
//     Patches the token name
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: PublicServiceAccountToken
//       401: empty
//       403: empty
func (r Routing) patchServiceAccountToken() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(serviceaccount.PatchTokenEndpoint(r.projectProvider, r.serviceAccountProvider, r.serviceAccountTokenProvider, r.saTokenAuthenticator, r.saTokenGenerator)),
		serviceaccount.DecodePatchTokenReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route DELETE /api/v1/projects/{project_id}/serviceaccounts/{serviceaccount_id}/tokens/{token_id} tokens deleteServiceAccountToken
//
//     Deletes the token
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) deleteServiceAccountToken() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(serviceaccount.DeleteTokenEndpoint(r.projectProvider, r.serviceAccountProvider, r.serviceAccountTokenProvider)),
		serviceaccount.DecodeDeleteTokenReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/aws/sizes aws listAWSSizesNoCredentials
//
// Lists available AWS sizes
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: AWSSizeList
func (r Routing) listAWSSizesNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.AWSSizeNoCredentialsEndpoint(r.projectProvider, r.seedsGetter)),
		common.DecodeGetClusterReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/aws/subnets aws listAWSSubnetsNoCredentials
//
// Lists available AWS subnets
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: AWSSubnetList
func (r Routing) listAWSSubnetsNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.AWSSubnetWithClusterCredentialsEndpoint(r.projectProvider, r.seedsGetter)),
		common.DecodeGetClusterReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/gcp/sizes gcp listGCPSizesNoCredentials
//
// Lists machine sizes from GCP
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: GCPMachineSizeList
func (r Routing) listGCPSizesNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.GCPSizeWithClusterCredentialsEndpoint(r.projectProvider)),
		provider.DecodeGCPTypesNoCredentialReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/gcp/disktypes gcp listGCPDiskTypesNoCredentials
//
// Lists disk types from GCP
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: GCPDiskTypeList
func (r Routing) listGCPDiskTypesNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.GCPDiskTypesWithClusterCredentialsEndpoint(r.projectProvider)),
		provider.DecodeGCPTypesNoCredentialReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/gcp/zones gcp listGCPZonesNoCredentials
//
// Lists available GCP zones
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: GCPZoneList
func (r Routing) listGCPZonesNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.GCPZoneWithClusterCredentialsEndpoint(r.projectProvider, r.seedsGetter)),
		common.DecodeGetClusterReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/hetzner/sizes hetzner listHetznerSizesNoCredentials
//
// Lists sizes from hetzner
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: HetznerSizeList
func (r Routing) listHetznerSizesNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.HetznerSizeWithClusterCredentialsEndpoint(r.projectProvider)),
		provider.DecodeHetznerSizesNoCredentialsReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/digitalocean/sizes digitalocean listDigitaloceanSizesNoCredentials
//
// Lists sizes from digitalocean
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: DigitaloceanSizeList
func (r Routing) listDigitaloceanSizesNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.DigitaloceanSizeWithClusterCredentialsEndpoint(r.projectProvider)),
		provider.DecodeDoSizesNoCredentialsReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/azure/sizes azure listAzureSizesNoCredentials
//
// Lists available VM sizes in an Azure region
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: AzureSizeList
func (r Routing) listAzureSizesNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.AzureSizeWithClusterCredentialsEndpoint(r.projectProvider, r.seedsGetter)),
		provider.DecodeAzureSizesNoCredentialsReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/sizes openstack listOpenstackSizesNoCredentials
//
// Lists sizes from openstack
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []OpenstackSize
func (r Routing) listOpenstackSizesNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.OpenstackSizeWithClusterCredentialsEndpoint(r.projectProvider, r.seedsGetter)),
		provider.DecodeOpenstackNoCredentialsReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/tenants openstack listOpenstackTenantsNoCredentials
//
// Lists tenants from openstack
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []OpenstackTenant
func (r Routing) listOpenstackTenantsNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.OpenstackTenantWithClusterCredentialsEndpoint(r.projectProvider, r.seedsGetter)),
		provider.DecodeOpenstackNoCredentialsReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/networks openstack listOpenstackNetworksNoCredentials
//
// Lists networks from openstack
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []OpenstackNetwork
func (r Routing) listOpenstackNetworksNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.OpenstackNetworkWithClusterCredentialsEndpoint(r.projectProvider, r.seedsGetter)),
		provider.DecodeOpenstackNoCredentialsReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/securitygroups openstack listOpenstackSecurityGroupsNoCredentials
//
// Lists security groups from openstack
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []OpenstackSecurityGroup
func (r Routing) listOpenstackSecurityGroupsNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.OpenstackSecurityGroupWithClusterCredentialsEndpoint(r.projectProvider, r.seedsGetter)),
		provider.DecodeOpenstackNoCredentialsReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/subnets openstack listOpenstackSubnetsNoCredentials
//
// Lists subnets from openstack
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []OpenstackSubnet
func (r Routing) listOpenstackSubnetsNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.OpenstackSubnetsWithClusterCredentialsEndpoint(r.projectProvider, r.seedsGetter)),
		provider.DecodeOpenstackSubnetNoCredentialsReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/vsphere/networks vsphere listVSphereNetworksNoCredentials
//
// Lists networks from vsphere datacenter
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []VSphereNetwork
func (r Routing) listVSphereNetworksNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.VsphereNetworksWithClusterCredentialsEndpoint(r.projectProvider, r.seedsGetter)),
		provider.DecodeVSphereNetworksNoCredentialsReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/vsphere/folders vsphere listVSphereFoldersNoCredentials
//
// Lists folders from vsphere datacenter
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []VSphereFolder
func (r Routing) listVSphereFoldersNoCredentials() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(provider.VsphereFoldersWithClusterCredentialsEndpoint(r.projectProvider, r.seedsGetter)),
		provider.DecodeVSphereFoldersNoCredentialsReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route POST /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments project createNodeDeployment
//
//     Creates a node deployment that will belong to the given cluster
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       201: NodeDeployment
//       401: empty
//       403: empty
func (r Routing) createNodeDeployment() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(node.CreateNodeDeployment(r.sshKeyProvider, r.projectProvider, r.seedsGetter)),
		node.DecodeCreateNodeDeployment,
		setStatusCreatedHeader(encodeJSON),
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments project listNodeDeployments
//
//     Lists node deployments that belong to the given cluster
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []NodeDeployment
//       401: empty
//       403: empty
func (r Routing) listNodeDeployments() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(node.ListNodeDeployments(r.projectProvider)),
		node.DecodeListNodeDeployments,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments/{nodedeployment_id} project getNodeDeployment
//
//     Gets a node deployment that is assigned to the given cluster.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: NodeDeployment
//       401: empty
//       403: empty
func (r Routing) getNodeDeployment() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(node.GetNodeDeployment(r.projectProvider)),
		node.DecodeGetNodeDeployment,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments/{nodedeployment_id}/nodes project listNodeDeploymentNodes
//
//     Lists nodes that belong to the given node deployment.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []Node
//       401: empty
//       403: empty
func (r Routing) listNodeDeploymentNodes() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(node.ListNodeDeploymentNodes(r.projectProvider)),
		node.DecodeListNodeDeploymentNodes,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments/{nodedeployment_id}/nodes/metrics metric listNodeDeploymentMetrics
//
//     Lists metrics that belong to the given node deployment.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []NodeMetric
//       401: empty
//       403: empty
func (r Routing) listNodeDeploymentMetrics() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(node.ListNodeDeploymentMetrics(r.projectProvider)),
		node.DecodeListNodeDeploymentMetrics,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments/{nodedeployment_id}/nodes/events project listNodeDeploymentNodesEvents
//
//     Lists node deployment events. If query parameter `type` is set to `warning` then only warning events are retrieved.
//     If the value is 'normal' then normal events are returned. If the query parameter is missing method returns all events.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []Event
//       401: empty
//       403: empty
func (r Routing) listNodeDeploymentNodesEvents() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(node.ListNodeDeploymentNodesEvents()),
		node.DecodeListNodeDeploymentNodesEvents,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PATCH /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments/{nodedeployment_id} project patchNodeDeployment
//
//     Patches a node deployment that is assigned to the given cluster. Please note that at the moment only
//	   node deployment's spec can be updated by a patch, no other fields can be changed using this endpoint.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: NodeDeployment
//       401: empty
//       403: empty
func (r Routing) patchNodeDeployment() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(node.PatchNodeDeployment(r.sshKeyProvider, r.projectProvider, r.seedsGetter)),
		node.DecodePatchNodeDeployment,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route DELETE /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/nodedeployments/{nodedeployment_id} project deleteNodeDeployment
//
//    Deletes the given node deployment that belongs to the cluster.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) deleteNodeDeployment() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(node.DeleteNodeDeployment(r.projectProvider)),
		node.DecodeDeleteNodeDeployment,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route POST /api/v1/addons addon
//
//     Lists names of addons that can be configured inside the user clusters
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: AccessibleAddons
//       401: empty
//       403: empty
func (r Routing) listAccessibleAddons() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(addon.ListAccessibleAddons(r.accessibleAddons)),
		decodeEmptyReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route POST /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/addons addon createAddon
//
//     Creates an addon that will belong to the given cluster
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       201: Addon
//       401: empty
//       403: empty
func (r Routing) createAddon() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.Addons(r.addonProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(addon.CreateAddonEndpoint(r.projectProvider)),
		addon.DecodeCreateAddon,
		setStatusCreatedHeader(encodeJSON),
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/addons addon listAddons
//
//     Lists addons that belong to the given cluster
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []Addon
//       401: empty
//       403: empty
func (r Routing) listAddons() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.Addons(r.addonProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(addon.ListAddonEndpoint(r.projectProvider)),
		addon.DecodeListAddons,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/addons/{addon_id} addon getAddon
//
//     Gets an addon that is assigned to the given cluster.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: Addon
//       401: empty
//       403: empty
func (r Routing) getAddon() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.Addons(r.addonProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(addon.GetAddonEndpoint(r.projectProvider)),
		addon.DecodeGetAddon,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PATCH /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/addons/{addon_id} addon patchAddon
//
//     Patches an addon that is assigned to the given cluster.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: Addon
//       401: empty
//       403: empty
func (r Routing) patchAddon() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.Addons(r.addonProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(addon.PatchAddonEndpoint(r.projectProvider)),
		addon.DecodePatchAddon,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route DELETE /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/addons/{addon_id} addon deleteAddon
//
//    Deletes the given addon that belongs to the cluster.
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) deleteAddon() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.Addons(r.addonProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(addon.DeleteAddonEndpoint(r.projectProvider)),
		addon.DecodeGetAddon,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/metrics project getClusterMetrics
//
//    Gets cluster metrics
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: ClusterMetrics
//       401: empty
//       403: empty
func (r Routing) getClusterMetrics() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.SetPrivilegedClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.GetMetricsEndpoint(r.projectProvider)),
		common.DecodeGetClusterReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route POST /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles project createClusterRole
//
//    Creates cluster role
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       201: ClusterRole
//       401: empty
//       403: empty
func (r Routing) createClusterRole() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.CreateClusterRoleEndpoint()),
		cluster.DecodeCreateClusterRoleReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/openshift/console/login
//
//    Creates an oauth token for the user and redirects them to the Openshift Console
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       302: empty
//       401: empty
//       403: empty
func (r Routing) openshiftConsoleLogin() http.Handler {
	return openshift.ConsoleLoginEndpoint(
		r.log,
		middleware.TokenExtractor(r.tokenExtractors),
		r.projectProvider,
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			// TODO: Instead of using an admin client to talk to the seed, we should provide a seed
			// client that allows access to the cluster namespace only
			middleware.SetPrivilegedClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		),
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/openshift/console/proxy
//
//    Proxies the Openshift console. Requires a valid OIDC token. The token can be obtained
//    using the /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/openshift/console/login
//    endpoint.
//
//     Responses:
//       default: empty
func (r Routing) openshiftConsoleProxy() http.Handler {
	return openshift.ConsoleProxyEndpoint(
		r.log,
		middleware.TokenExtractor(r.tokenExtractors),
		r.projectProvider,
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			// TODO: Instead of using an admin client to talk to the seed, we should provide a seed
			// client that allows access to the cluster namespace only
			middleware.SetPrivilegedClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		),
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/openshift/console/proxy
//
//    Proxies the Kubernetes Dashboard. Requires a valid bearer token. The token can be obtained
//    using the /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/dashboard/login
//    endpoint.
//
//     Responses:
//       default: empty
func (r Routing) kubernetesDashboardProxy() http.Handler {
	return kubernetesdashboard.ProxyEndpoint(
		r.log,
		middleware.TokenExtractor(r.tokenExtractors),
		r.projectProvider,
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			// TODO: Instead of using an admin client to talk to the seed, we should provide a seed
			// client that allows access to the cluster namespace only
			middleware.SetPrivilegedClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		),
	)
}

// swagger:route POST /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles project createRole
//
//    Creates cluster role
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       201: Role
//       401: empty
//       403: empty
func (r Routing) createRole() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.CreateRoleEndpoint()),
		cluster.DecodeCreateRoleReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles project listClusterRole
//
//     Lists all ClusterRoles
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []ClusterRole
//       401: empty
//       403: empty
func (r Routing) listClusterRole() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.ListClusterRoleEndpoint()),
		cluster.DecodeListClusterRoleReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles project listRole
//
//     Lists all Roles
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []Role
//       401: empty
//       403: empty
func (r Routing) listRole() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.ListRoleEndpoint()),
		cluster.DecodeListClusterRoleReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{role_id} project getClusterRole
//
//     Gets the cluster role with the given name
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: ClusterRole
//       401: empty
//       403: empty
func (r Routing) getClusterRole() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.GetClusterRoleEndpoint()),
		cluster.DecodeGetClusterRoleReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id} project getRole
//
//     Gets the role with the given name
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: Role
//       401: empty
//       403: empty
func (r Routing) getRole() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.GetRoleEndpoint()),
		cluster.DecodeGetRoleReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route DELETE /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id} project deleteClusterRole
//
//     Delete the cluster role with the given name
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) deleteClusterRole() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.DeleteClusterRoleEndpoint()),
		cluster.DecodeGetClusterRoleReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route DELETE /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id} project deleteRole
//
//     Delete the cluster role with the given name
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) deleteRole() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.DeleteRoleEndpoint()),
		cluster.DecodeGetRoleReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/namespaces project listNamespace
//
//     Lists all namespaces in the cluster
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []Namespace
//       401: empty
//       403: empty
func (r Routing) listNamespace() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.ListNamespaceEndpoint(r.projectProvider)),
		common.DecodeGetClusterReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PATCH /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id} project patchRole
//
//     Patch the role with the given name
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: Role
//       401: empty
//       403: empty
func (r Routing) patchRole() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.PatchRoleEndpoint()),
		cluster.DecodePatchRoleReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PATCH /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id} project patchClusterRole
//
//     Patch the cluster role with the given name
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: ClusterRole
//       401: empty
//       403: empty
func (r Routing) patchClusterRole() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.PatchClusterRoleEndpoint()),
		cluster.DecodePatchClusterRoleReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route POST /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings project createRoleBinding
//
//    Creates role binding
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       201: RoleBinding
//       401: empty
//       403: empty
func (r Routing) createRoleBinding() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.CreateRoleBindingEndpoint()),
		cluster.DecodeCreateRoleBindingReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings project listRoleBinding
//
//    List role binding
//
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []RoleBinding
//       401: empty
//       403: empty
func (r Routing) listRoleBinding() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.ListRoleBindingEndpoint()),
		cluster.DecodeListRoleBindingReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings/{binding_id} project getRoleBinding
//
//    Get role binding
//
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: RoleBinding
//       401: empty
//       403: empty
func (r Routing) getRoleBinding() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.GetRoleBindingEndpoint()),
		cluster.DecodeRoleBindingReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route DELETE /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings/{binding_id} project deleteRoleBinding
//
//    Delete role binding
//
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) deleteRoleBinding() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.DeleteRoleBindingEndpoint()),
		cluster.DecodeRoleBindingReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PATCH /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings/{binding_id} project patchRoleBinding
//
//    Update role binding
//
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: RoleBinding
//       401: empty
//       403: empty
func (r Routing) patchRoleBinding() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.PatchRoleBindingEndpoint()),
		cluster.DecodePatchRoleBindingReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route POST /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id}/clusterbindings project createClusterRoleBinding
//
//    Creates cluster role binding
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       201: ClusterRoleBinding
//       401: empty
//       403: empty
func (r Routing) createClusterRoleBinding() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.CreateClusterRoleBindingEndpoint()),
		cluster.DecodeCreateClusterRoleBindingReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id}/clusterbindings project listClusterRoleBinding
//
//    List cluster role binding
//
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: []ClusterRoleBinding
//       401: empty
//       403: empty
func (r Routing) listClusterRoleBinding() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.ListClusterRoleBindingEndpoint()),
		cluster.DecodeListClusterRoleBindingReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id}/clusterbindings/{binding_id} project getClusterRoleBinding
//
//    Get cluster role binding
//
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: ClusterRoleBinding
//       401: empty
//       403: empty
func (r Routing) getClusterRoleBinding() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.GetClusterRoleBindingEndpoint()),
		cluster.DecodeClusterRoleBindingReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route DELETE /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id}/clusterbindings/{binding_id} project deleteClusterRoleBinding
//
//    Delete cluster role binding
//
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: empty
//       401: empty
//       403: empty
func (r Routing) deleteClusterRoleBinding() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.DeleteClusterRoleBindingEndpoint()),
		cluster.DecodeClusterRoleBindingReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PATCH /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/clusterroles/{role_id}/clusterbindings/{binding_id} project patchClusterRoleBinding
//
//    Update cluster role binding
//
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: ClusterRoleBinding
//       401: empty
//       403: empty
func (r Routing) patchClusterRoleBinding() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
			middleware.UserSaver(r.userProvider),
			middleware.SetClusterProvider(r.clusterProviderGetter, r.seedsGetter),
			middleware.UserInfoExtractor(r.userProjectMapper),
		)(cluster.PatchClusterRoleBindingEndpoint()),
		cluster.DecodePatchClusterRoleBindingReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}

// swagger:route PATCH /api/v1/labels/system listSystemLabels
//
//    List restricted system labels
//
//
//     Produces:
//     - application/json
//
//     Responses:
//       default: errorResponse
//       200: ResourceLabelMap
//       401: empty
//       403: empty
func (r Routing) listSystemLabels() http.Handler {
	return httptransport.NewServer(
		endpoint.Chain(
			middleware.TokenVerifier(r.tokenVerifiers),
		)(label.ListSystemLabels()),
		decodeEmptyReq,
		encodeJSON,
		r.defaultServerOptions()...,
	)
}
