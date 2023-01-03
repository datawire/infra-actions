'use strict';

const utils = require('./utils.js')

test('fibonacciDelaySequence unlimited', () => {
  let seq = utils.fibonacciDelaySequence(1)
  let prev = 0
  let cur = 1
  for (let i = 0; i < 100; i++) {
    let next = cur + prev
    prev = cur
    cur = next
    expect(seq()).toBe(cur)
  }
})

test('fibonacciDelaySequence limited', () => {
  let limit = 10
  let seq = utils.fibonacciDelaySequence(1, limit)
  let prev = 0
  let cur = 1
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
    let id = utils.uid()
    expect(ids.has(id)).toBeFalsy()
    ids.add(id)
  }
})

test('fibonacciRetry success', async () => {
  let count = 0
  let result = await utils.fibonacciRetry(async ()=> {
    count += 1
    return count
  })
  expect(result).toBe(count)
  expect(count).toBe(1)
})

test('fibonacciRetry fail some', async () => {
  let count = 0
  let result = await utils.fibonacciRetry(async ()=> {
    count += 1
    if (count > 3) {
      return count
    }
    throw new utils.Transient(`${count} is not big enough`)
  }, 100, 1)
  expect(count).toBe(result)
  expect(count).toBe(4)
})

test('fibonacciRetry fail all', async () => {
  let count = 0
  let start = Date.now()
  let returned = false
  try {
    await utils.fibonacciRetry(async ()=> {
      count += 1
      throw new utils.Transient('never big enough')
    }, 100, 1)
    returned = true
  } catch (err) {
    let elapsed = Date.now() - start
    expect(err.message).toContain("Transient error")
    expect(err.message).toContain("never big enough")
    expect(err.message).toContain("failing after")
    expect(err.message).toContain("attempts over")
    expect(count > 0).toBeTruthy()
    expect(elapsed < 1000).toBeTruthy()
  }

  expect(returned).toBeFalsy()
})

test('fibonacciRetry max delay', async () => {
  let count = 0
  await utils.fibonacciRetry(async ()=> {
    count += 1
    if (count > 10) {
      return count
    }
    throw new utils.Transient('never big enough')
  }, 100, 1, 10)
})

test('fibonacciRetry error', async () => {
  let count = 0
  let returned = false
  try {
    await utils.fibonacciRetry(async ()=> {
      count += 1
      throw new Error('blah')
    }, 100, 1, 10)
    returned = true
  } catch (err) {
    expect(err.message).toEqual('blah')
    expect(count).toBe(1)
  }
  expect(returned).toBeFalsy()
})

test('getUniqueClusterName', () => {
  process.env["GITHUB_REPOSITORY"] = "repo"
  process.env["GITHUB_HEAD_REF"] = "head-ref"
  process.env["GITHUB_SHA"] = "1234"
  for (let i = 10; i < 1000; i++) {
    let name = utils.getUniqueClusterName(i)
    expect(name.length <= i).toBeTruthy()
    expect(name.endsWith("-")).toBeFalsy()
  }
})
