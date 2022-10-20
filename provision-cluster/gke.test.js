const gke = require('./gke.js')
const fs = require('fs')

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

test('gke kubeconfig', async ()=> {
  let client = mock.Client()
  let cluster = JSON.parse(fs.readFileSync('cluster.json'))
  let kubeconfig = await client.makeKubeconfig(cluster)
  expect(kubeconfig).toEqual(
    {
      "apiVersion": "v1",
      "kind":"Config",
      "clusters": [{
        "cluster": {
          "certificate-authority-data": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVMVENDQXBXZ0F3SUJBZ0lSQUpXa0xadDFDVS9RNGpHcDkwTjZiWkl3RFFZSktvWklodmNOQVFFTEJRQXcKTHpFdE1Dc0dBMVVFQXhNa09HUmhNamhqTTJRdE9ETXlOUzAwWTJGakxXSm1NRGN0WkRKbE9EbGxOakF5TXpRdwpNQ0FYRFRJeU1UQXlNREUzTWpZME4xb1lEekl3TlRJeE1ERXlNVGd5TmpRM1dqQXZNUzB3S3dZRFZRUURFeVE0ClpHRXlPR016WkMwNE16STFMVFJqWVdNdFltWXdOeTFrTW1VNE9XVTJNREl6TkRBd2dnR2lNQTBHQ1NxR1NJYjMKRFFFQkFRVUFBNElCandBd2dnR0tBb0lCZ1FEVU1yWHdwMGdIVkt4cDVubTdJMlZDRGNtYldBS3l6eHowYmtsSAp0Z0tLMUk3MWt0Uyt4MGRSZFFvWEFvY3Q4aTBmVTlMT3V0MExCaERna3lVczhRamZ6V2NULytZL3ZpclhwZGJUCjF5R2R1MjlEVTJ6emZtSTFNOU9vNlRRK21mZmVSZjBpY3N2eE5JRjI2RE5Jdk0xcG1WUm4yZDRwUW03SFdMb04KSGpJWlJlS2I0RmtjSkpSWW9aVTNSUWg3blhsUWREcXhRaUxGZnBtL1dBTUVPVkdlUmVTb0prODdTSDBXU0IydQpmQTk0YzZRekdVQmwwRUtnQzJnU0ZlQXNpZmZNVUdTTjJHWHdTTXRlTTJtL3FKa2I5ZVV5bmVleEVwM3czRmlCClZ5Nmo1bXVCakdweStNWmhrN3VqY1ROS1lJekpMcEhMZ2d5cGxkVm5abEhxMnB5c09Zazc0VVlZTzFXY012ZkYKTDM0K1NYV0kzc3REWTFVdXFGNklNT3FSRVB5UGlGdmRic3VTMmdFVGNLV0U1Vk4rOHRoS3prQlA4VEduNlFaMgowTEpDSGJ3eVRxQ1E4Q2lJR2sxdTBUZEVzMjdJVXIvcmlQL1R6aWdDRHo3cGhhK1liNUo2RnJvVWFmMEhxL0hpCnRDUDVhOTBEV29lMzA4aEdieFhaMldVUWdwOENBd0VBQWFOQ01FQXdEZ1lEVlIwUEFRSC9CQVFEQWdJRU1BOEcKQTFVZEV3RUIvd1FGTUFNQkFmOHdIUVlEVlIwT0JCWUVGT2JhL1lReWVMbElnWVU1RjZCMXg2VjVkYTlITUEwRwpDU3FHU0liM0RRRUJDd1VBQTRJQmdRQXZwTFVBTUR1NXprQ0tjNE45YlJiOGhyRHBTOHkzYXBxcEZwY1dHYytRCnRmMmZ0ZkpYc2lxQW1Tems1Uy9VNlJYYjBSc0NIRkFlK1dGQWpTbkFYRFE4eEU0MUJWODNOYW9wNlZxQ3pXSFMKS2JUZ2ZBVmZ4NFBCQjRqMmhBZDN1YWpqZGE3Z3BpakZmSkdHN3NuNVR1d3hvQ0tKQW16aW9WUkIzT2psQTB3RQoyZVBhQ1RGVTEzM0EweU5PNUc1bkhrU3hGUFdWUUs4T05WeWJZNGFZRThjUkRuVDFSUTZSVHorMkdodmlRN2tNClhaRVY4Q2FCSzlKeW1UTVB4T3Mvc0hsYUlaMHdURDEwYndKSFZaU3krRHVjbGd6NjM5R3FlY1V0bXBzdE1QcG8KN2l6ejlmejhicVY4citZMmE4Y0FBd1dWOHVEMU1HZDBqZFB1TDNsOTI2U09ENUJrUTM2ZHh6aHB1cjFJekVibwpnM1ltaDNaVkhzVXRBQ2Z0aytsdUxEckpDNHNJai9sRjNBQXh0MDIyY0VHMnpOSzN5b2VuV0tacWVseC9xdUluCm01TWJQMndmcm12RUpxcUo1UkpiNm9jZlhyS3NUQW0rajRTSzJPM1JBNVloa2tjVTFGTzZUcnZCUnFtNTBBSVIKN216YWs5M1dFdWw1NnFTVTNsUTZYaE09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",
          "server":"https://34.172.65.239"
        },
        "name":"gke-cluster"
      }],
      "users": [{
        "name":"gke-user",
        "user":{
          "token":"mock-access-token"
        }}],
      "contexts": [{
        "context":{"cluster":"gke-cluster","namespace":"default","user":"gke-user"},
        "name":"gke-context"}],
      "current-context":"gke-context"
    })
})
