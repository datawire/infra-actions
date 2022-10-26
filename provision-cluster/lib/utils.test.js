const utils = require('./utils.js')

test('fibonacciDelaySequence unlimited', () => {
  let seq = utils.fibonacciDelaySequence(1)
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
  let seq = utils.fibonacciDelaySequence(1, limit)
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
    let id = utils.uid()
    expect(ids.has(id)).toBeFalsy()
    ids.add(id)
  }
})
