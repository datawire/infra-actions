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
    await this.createKluster(clusterName, version, lifespan);
    const kubeConfig = await this.getKlusterKubeconfig(clusterName);
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
      throw new Error("Kluster name is required");
    }

    if (!version) {
      throw Error("Kluster version is required");
    }

    if (
      typeof lifespan === typeof undefined ||
      lifespan === "" ||
      lifespan === 0
    ) {
      lifespan = defaultLifespan;
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
        case 202:
          return;
        case 425:
          // This will be deprecated in the future, pending rework of the API
          // 425 should be treated as any other 4xx error
          return;
        default:
          if (response.message.statusCode >= 400) {
            throw new utils.Transient(
              `Status code ${response.message.statusCode}`
            );
          } else {
            throw new Error(`Status code ${response.message.statusCode}`);
          }
      }
    });
  }

  async getKlusterKubeconfig(name) {
    if (!name) {
      throw new Error("Kluster name is required");
    }

    return utils.fibonacciRetry(async () => {
      const response = await this.client.get(
        `https://sw.bakerstreet.io/kubeception/api/klusters/${name}/kubeconfig`
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
      throw Error("Kluster name is required");
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

  const options = {
    keepAlive: false,
  };

  return new httpClient.HttpClient(userAgent, [credentialHandler], options);
}

module.exports = { Client, getHttpClient };
