# Cluster Provisioning

Use the `provision-cluster` action as described below:

```yaml
      - uses: ./provision-cluster
        with:
          # Tells the action what kind of cluster to create. One of: Kubeception, GKE, EKS, AKS, OpenShift
          distribution: ...
          # Tells the action what version of cluster to create.
          version: 1.23
          # Tells provision-cluster where to write the kubeconfig file.
          kubeconfig: path/to/kubeconfig.yaml

          ## For kubeception klusters

          # A kubeception secret token
          kubeceptionToken: ...

          ## For GKE clusters:

          # A json encoded string containing GKE credentials:
          gkeCredentials: ...
          # A json encoded string containing additional GKE cluster configuration. See GKE Cluster Config Options section for details.
          gkeConfig: ...
```

## GKE Cluster Config Options

The values included below are the defaults.

```json
{
  "resourceLabels": {
    "provisioned-category": "ephemeral",
    "provisioned-lifespan": "1800"
  },
  "description": "",
  "initialNodeCount": 1,
  "nodeConfig": {
    "oauthScopes": [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring"
    ],
    "tags": [],
    "accelerators": [],
    "taints": [],
    "metadata": {
      "disable-legacy-endpoints": "true"
    },
    "labels": {},
    "machineType": "e2-standard-2",
    "diskSizeGb": 100,
    "imageType": "COS_CONTAINERD",
    "localSsdCount": 0,
    "serviceAccount": "default",
    "preemptible": false,
    "diskType": "pd-standard",
    "minCpuPlatform": "",
    "workloadMetadataConfig": null,
    "sandboxConfig": null,
    "nodeGroup": "",
    "reservationAffinity": null,
    "shieldedInstanceConfig": {
      "enableSecureBoot": false,
      "enableIntegrityMonitoring": true
    },
    "linuxNodeConfig": null,
    "kubeletConfig": null,
    "bootDiskKmsKey": "",
    "gcfsConfig": null,
    "advancedMachineFeatures": null,
    "gvnic": null,
    "spot": false,
    "confidentialNodes": null,
    "loggingConfig": null
  },
  "loggingService": "logging.googleapis.com/kubernetes",
  "monitoringService": "monitoring.googleapis.com/kubernetes",
  "network": "default",
  "clusterIpv4Cidr": "10.0.0.0/14",
  "addonsConfig": {
    "httpLoadBalancing": null,
    "horizontalPodAutoscaling": null,
    "kubernetesDashboard": {
      "disabled": true
    },
    "networkPolicyConfig": {
      "disabled": true
    },
    "cloudRunConfig": null,
    "dnsCacheConfig": null,
    "configConnectorConfig": null,
    "gcePersistentDiskCsiDriverConfig": {
      "enabled": true
    },
    "gcpFilestoreCsiDriverConfig": null
  },
  "subnetwork": "default",
  "enableKubernetesAlpha": false,
  "labelFingerprint": "81c637a5",
  "legacyAbac": {
    "enabled": false
  },
  "networkPolicy": null,
  "ipAllocationPolicy": {
    "useIpAliases": true,
    "createSubnetwork": false,
    "subnetworkName": "",
    "clusterIpv4Cidr": "10.0.0.0/14",
    "nodeIpv4Cidr": "",
    "servicesIpv4Cidr": "10.124.16.0/20",
    "clusterSecondaryRangeName": "gke-test-3fceb6744f7639bc0d6e9b601a051071-pods-da48a07c",
    "servicesSecondaryRangeName": "gke-test-3fceb6744f7639bc0d6e9b601a051071-services-da48a07c",
    "clusterIpv4CidrBlock": "10.0.0.0/14",
    "nodeIpv4CidrBlock": "",
    "servicesIpv4CidrBlock": "10.124.16.0/20",
    "tpuIpv4CidrBlock": "",
    "useRoutes": false
  },
  "masterAuthorizedNetworksConfig": null,
  "maintenancePolicy": {
    "window": null,
    "resourceVersion": "e3b0c442"
  },
  "binaryAuthorization": null,
  "autoscaling": null,
  "networkConfig": {
    "network": "projects/datawireio/global/networks/default",
    "subnetwork": "projects/datawireio/regions/us-central1/subnetworks/default",
    "enableIntraNodeVisibility": false,
    "defaultSnatStatus": null,
    "enableL4ilbSubsetting": false,
    "datapathProvider": "DATAPATH_PROVIDER_UNSPECIFIED",
    "privateIpv6GoogleAccess": "PRIVATE_IPV6_GOOGLE_ACCESS_UNSPECIFIED",
    "dnsConfig": null,
    "serviceExternalIpsConfig": {
      "enabled": false
    }
  },
  "defaultMaxPodsConstraint": {
    "maxPodsPerNode": "110"
  },
  "resourceUsageExportConfig": null,
  "authenticatorGroupsConfig": null,
  "privateClusterConfig": null,
  "databaseEncryption": {
    "keyName": "",
    "state": "DECRYPTED"
  },
  "verticalPodAutoscaling": null,
  "shieldedNodes": {
    "enabled": true
  },
  "releaseChannel": {
    "channel": "REGULAR"
  },
  "workloadIdentityConfig": null,
  "notificationConfig": {
    "pubsub": {
      "enabled": false,
      "topic": "",
      "filter": null
    }
  },
  "confidentialNodes": null,
  "identityServiceConfig": null,
  "meshCertificates": null,
  "initialClusterVersion": "1.22.12-gke.2300",
  "nodeIpv4CidrSize": 0,
  "servicesIpv4Cidr": "10.124.16.0/20",
  "enableTpu": false,
  "tpuIpv4CidrBlock": "",
  "autopilot": null,
  "loggingConfig": {
    "componentConfig": {
      "enableComponents": [
        "SYSTEM_COMPONENTS",
        "WORKLOADS"
      ]
    }
  },
  "monitoringConfig": {
    "componentConfig": {
      "enableComponents": [
        "SYSTEM_COMPONENTS"
      ]
    },
    "managedPrometheusConfig": null
  },
  "nodePoolAutoConfig": {
    "networkTags": null
  },
  "nodePoolDefaults": {
    "nodeConfigDefaults": {
      "gcfsConfig": null,
      "loggingConfig": {
        "variantConfig": {
          "variant": "DEFAULT"
        }
      }
    }
  }
}
```
