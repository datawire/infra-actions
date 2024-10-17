"use strict";

const core = require("@actions/core");
const httpClient = require("@actions/http-client");
const httpClientLib = require("@actions/http-client/lib/auth.js");
const utils = require("./utils.js");
const yaml = require("yaml");

const MAX_KLUSTER_NAME_LEN = 63;
const DEFAULT_ARCHETYPE = "small";
const DEFAULT_MODE = "active";
const DEFAULT_LIFESPAN = 60 * 60;
const KUBECEPTION_URL = "https://kubeception.datawire.io";

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
      throw new Error("Kluster version is required");
    }

    lifespan = lifespan || DEFAULT_LIFESPAN;

    let kluster = {
      name: name,
      version: version,
      archetype: DEFAULT_ARCHETYPE,
      mode: DEFAULT_MODE,
      timeoutSecs: lifespan,
    };

    return utils.fibonacciRetry(async () => {
      const response = await this.client.post(
        `${KUBECEPTION_URL}/api/klusters`,
        JSON.stringify(kluster),
        {
          "Content-Type": "application/json",
        }
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

    return utils.fibonacciRetry(
      async () => {
        const response = await this.client.get(
          `${KUBECEPTION_URL}/api/klusters/${name}/kubeconfig`
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
      },
      600000,
      1000,
      600000
    );
  }

  async deleteKluster(name) {
    if (!name) {
      throw new Error("Kluster name is required");
    }

    const response = await this.client.del(
      `${KUBECEPTION_URL}/api/klusters/${name}`
    );

    if (!response || !response.message) {
      throw new Error("Unknown error getting response");
    }

    if (response.message.statusCode === 200) {
      return {
        done: true,
        status: "deleted",
      };
    } else {
      throw new Error(
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

  return new httpClient.HttpClient(userAgent, [credentialHandler], {
    keepAlive: false,
  });
}

module.exports = { Client, getHttpClient };
