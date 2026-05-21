import { describe, test, expect } from 'vitest'
import { parseCsv, encodeCsv } from '../lib/url-state'

describe('parseCsv', () => {
  test('returns [] for null', () => {
    expect(parseCsv(null)).toEqual([])
  })

  test('returns [] for empty string', () => {
    expect(parseCsv('')).toEqual([])
  })

  test('splits a comma-separated value', () => {
    expect(parseCsv('ethereum,polygon,arbitrum')).toEqual(['ethereum', 'polygon', 'arbitrum'])
  })

  test('drops empty segments from trailing or doubled commas', () => {
    expect(parseCsv('ethereum,,polygon,')).toEqual(['ethereum', 'polygon'])
  })

  test('passes a single value through', () => {
    expect(parseCsv('tvl')).toEqual(['tvl'])
  })
})

describe('encodeCsv', () => {
  test('returns null for an empty array', () => {
    expect(encodeCsv([])).toBeNull()
  })

  test('returns null when every entry is falsy', () => {
    expect(encodeCsv(['', ''])).toBeNull()
  })

  test('joins multiple values with commas', () => {
    expect(encodeCsv(['ethereum', 'polygon'])).toBe('ethereum,polygon')
  })

  test('drops empty segments before joining', () => {
    expect(encodeCsv(['ethereum', '', 'polygon'])).toBe('ethereum,polygon')
  })
})

describe('parseCsv <-> encodeCsv round-trip', () => {
  const canonical: string[][] = [[], ['tvl'], ['ethereum', 'polygon'], ['a', 'b', 'c', 'd']]

  for (const input of canonical) {
    test(`parseCsv(encodeCsv(${JSON.stringify(input)})) preserves the array`, () => {
      expect(parseCsv(encodeCsv(input))).toEqual(input)
    })
  }
})
