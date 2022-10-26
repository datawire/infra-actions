const gke = require('./gke.js')

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
