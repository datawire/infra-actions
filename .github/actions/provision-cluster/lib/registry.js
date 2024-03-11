"use strict";

const gke = require("./gke.js");
const kubeception = require("./kubeception.js");

const CLUSTER_NAME = "CLUSTER_NAME";
const clusterZone = "us-central1-b";

const DistributionType = {
  GKE: "gke",
  KUBECEPTION: "kubeception",
};

function getProvider(distribution) {
  const lowerCaseDistribution = distribution.toLowerCase();

  switch (lowerCaseDistribution) {
    case DistributionType.GKE:
      return new gke.Client(clusterZone);
    case DistributionType.KUBECEPTION:
      return new kubeception.Client(kubeception.getHttpClient());
    default:
      throw new Error(`unknown distribution: ${distribution}`);
  }
}

module.exports = {
  getProvider,
  CLUSTER_NAME,
};
