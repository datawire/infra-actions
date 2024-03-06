'use strict';

const common = require('./common_test.js')
const mock = require('./mock.js')

test('gke mock e2e', async ()=> {
  await common.lifecycle(mock.Client())
})
