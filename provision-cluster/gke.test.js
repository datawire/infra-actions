const gke = require('./gke.js')

test('fibonacciDelaySequence unlimited', () => {
  let seq = gke.fibonacciDelaySequence(1)
  let prev = 0;
  let cur = 1;
  for (let i = 0; i < 100; i++) {
    let next = cur + prev
    prev = cur
    cur = next
    expect(seq()).toBe(cur)
  }
})

test('fibonacciDelaySequence limited', () => {
  let limit = 10
  let seq = gke.fibonacciDelaySequence(1, limit)
  let prev = 0;
  let cur = 1;
  for (let i = 0; i < 100; i++) {
    let next = cur + prev
    prev = cur
    cur = next
    if (cur > limit) {
      expect(seq()).toBe(limit)
    } else {
      expect(seq()).toBe(cur)
    }
  }
})

test('uid', () => {
  let ids = new Set()
  for (let i = 0; i < 100; i++) {
    let id = gke.uid()
    expect(ids.has(id)).toBeFalsy()
    ids.add(id)
  }
})

const mock = require('./mock.js')
let MOCK = mock.MOCK

test('gke mock e2e', async ()=> {
  let client = mock.Client()
  let cluster = await client.allocateCluster()
  let kubeconfig = await client.makeKubeconfig(cluster)
  expect(kubeconfig).toEqual(
    {
      "apiVersion": "v1",
      "kind":"Config",
      "clusters": [{
        "cluster": {
          "certificate-authority-data": cluster.masterAuth.clusterCaCertificate,
          "server":"https://34.172.65.239"
        },
        "name":"gke-cluster"
      }],
      "users": [{
        "name":"gke-user",
        "user":{
          "token":MOCK.ACCESS_TOKEN
        }}],
      "contexts": [{
        "context":{"cluster":"gke-cluster","namespace":"default","user":"gke-user"},
        "name":"gke-context"}],
      "current-context":"gke-context"
    })

  let op = await client.deleteCluster(cluster)
  expect(op.done).toBeTruthy()
})
