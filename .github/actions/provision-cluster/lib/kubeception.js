"use strict";

const core = require("@actions/core");
const httpClient = require("@actions/http-client");
const httpClientLib = require("@actions/http-client/lib/auth.js");
const utils = require("./utils.js");
const yaml = require("yaml");

const MAX_KLUSTER_NAME_LEN = 63;
const defaultLifespan = 60 * 60; // One hour worth of seconds

class Client {
  constructor(client) {
    this.client = client;
  }

  async allocateCluster(version, lifespan) {
    const clusterName = utils.getUniqueClusterName(MAX_KLUSTER_NAME_LEN);
    const kubeConfig = await this.createKluster(clusterName, version, lifespan);
    return {
      name: clusterName,
      config: kubeConfig,
    };
  }

  async makeKubeconfig(cluster) {
    return yaml.parse(cluster.config);
  }

  async getCluster(clusterName) {
    return clusterName;
  }

  async deleteCluster(clusterName) {
    return this.deleteKluster(clusterName);
  }

  async expireClusters() {
    // Kubeception automatically expires klusters, no client side expiration is required.
    return [];
  }

  async createKluster(name, version, lifespan) {
    if (!name) {
      throw new Error("Function createKluster() needs a Kluster name");
    }

    if (!version) {
      throw Error("Function createKluster() needs a Kluster version");
    }

    if (
      typeof lifespan === typeof undefined ||
      lifespan === "" ||
      lifespan === 0
    ) {
      lifespan = defaultLifespan;
    }

    const kubeceptionToken = core.getInput("kubeceptionToken");
    if (!kubeceptionToken) {
      throw Error(
        `kubeceptionToken is missing. Make sure that input parameter kubeceptionToken was provided`
      );
    }

    let kubeceptionProfile = core.getInput("kubeceptionProfile");
    if (
      typeof kubeceptionProfile !== typeof "" ||
      kubeceptionProfile.trim() === ""
    ) {
      kubeceptionProfile = "default";
    }

    return utils.fibonacciRetry(async () => {
      const response = await this.client.put(
        `https://sw.bakerstreet.io/kubeception/api/klusters/${name}?version=${version}&profile=${kubeceptionProfile}&timeoutSecs=${lifespan}`
      );

      if (!response || !response.message) {
        throw new utils.Transient("Unknown error getting response");
      }

      switch (response.message.statusCode) {
        case 200:
        case 201:
          return await response.readBody();
        case 202:
          throw new utils.Retry("Request is still pending");
        default:
          if (response.message.statusCode >= 400) {
            throw new utils.Transient(
              `Status code ${response.message.statusCode}`
            );
          } else {
            let body = await response.readBody();
            throw new Error(
              `Status code ${response.message.statusCode}: ${body}`
            );
          }
      }
    });
  }

  async deleteKluster(name) {
    if (!name) {
      throw Error("Function deleteKluster() needs a Kluster name");
    }

    const response = await this.client.del(
      `https://sw.bakerstreet.io/kubeception/api/klusters/${name}`
    );
    if (!response || !response.message) {
      throw Error("Unknown error getting response");
    }

    if (response.message.statusCode == 200) {
      return {
        done: true,
        status: "deleted",
      };
    } else {
      throw Error(
        `Expected status code 200 but got ${response.message.statusCode}`
      );
    }
  }
}

function getHttpClient() {
  const userAgent = "datawire/provision-cluster";

  const kubeceptionToken = core.getInput("kubeceptionToken");
  if (!kubeceptionToken) {
    throw Error(
      `kubeceptionToken is missing. Make sure that input parameter kubeceptionToken was provided`
    );
  }

  const credentialHandler = new httpClientLib.BearerCredentialHandler(
    kubeceptionToken
  );
  return new httpClient.HttpClient(userAgent, [credentialHandler]);
}

module.exports = { Client, getHttpClient };
