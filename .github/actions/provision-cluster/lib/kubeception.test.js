"use strict";

const core = require("@actions/core");
const kubeception = require("./kubeception.js");
const common = require("./common_test.js");
const mock = require("./mock.js");
const MOCK = mock.MOCK;
const cluster = mock.cluster;

test("kubeception", async () => {
  let inputs = {
    kubeceptionToken: "mock-kube-token",
  };

  let count = 0;

  class MockHttpClient {
    async post() {
      return {
        message: {
          statusCode: 200,
        },
      };
    }
    async get() {
      let status = 200;
      if (count < 2) {
        status = 425;
        count = count + 1;
      } else {
        status = 200;
      }
      return {
        message: {
          statusCode: status,
        },
        readBody: () => {
          return JSON.stringify({
            apiVersion: "v1",
            kind: "Config",
            clusters: [
              {
                cluster: {
                  "certificate-authority-data":
                    cluster.masterAuth.clusterCaCertificate,
                  server: "https://34.172.65.239",
                },
                name: "gke-cluster",
              },
            ],
            users: [
              {
                name: "gke-user",
                user: {
                  token: MOCK.ACCESS_TOKEN,
                },
              },
            ],
            contexts: [
              {
                context: {
                  cluster: "gke-cluster",
                  namespace: "default",
                  user: "gke-user",
                },
                name: "gke-context",
              },
            ],
            "current-context": "gke-context",
          });
        },
      };
    }
    async del() {
      return {
        message: {
          statusCode: 200,
        },
      };
    }
  }

  process.env.GITHUB_REPOSITORY = "test-project-repo";
  process.env.GITHUB_HEAD_REF = "refs/pull/1234";
  process.env.GITHUB_SHA = "abc1234";
  core.getInput = (name) => {
    return inputs[name];
  };
  await common.lifecycle(new kubeception.Client(new MockHttpClient()));
});
