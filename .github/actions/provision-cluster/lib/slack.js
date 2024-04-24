"use strict";

const core = require("@actions/core");
const http = require("@actions/http-client");

async function notify(message) {
  const slackWebhook = core.getInput("slackWebhook");
  if (typeof slackWebhook === typeof undefined || slackWebhook === "") {
    return;
  }

  const channel = core.getInput("slackChannel");
  const runbook = core.getInput("slackRunbook");
  const username = core.getInput("slackUsername");

  const client = new http.HttpClient("datawire/provision-cluster", [], {
    keepAlive: false,
  });

  const body = {
    channel: channel,
    username: username,
    text: `${message}\n\nSee runbook: ${runbook}`,
    icon_emoji: `:kubernetes:`,
  };

  const result = await client.postJson(slackWebhook, body);
  if (result.statusCode != 200) {
    throw new Error(`Status ${result.statusCode} posting to slack: ${result}`);
  }
}

module.exports = {
  notify,
};
