import { describe, test, expect } from 'vitest'
import { projectRow, type RawProtocol } from '../lib/snapshot'

describe('projectRow', () => {
  const baseProtocol: RawProtocol = {
    slug: 'aave',
    name: 'Aave',
    category: 'Lending',
    chains: ['Ethereum', 'Polygon'],
    dataFile: 'aave.json',
  }

  test('writes 1 for present metrics and 0 for absent ones', () => {
    const present = new Set(['tvl', 'dailyFees'])
    const row = projectRow(baseProtocol, present)

    expect(row.cells.tvl).toBe(1)
    expect(row.cells.dailyFees).toBe(1)
    expect(row.cells.dailyRevenue).toBe(0)
    expect(row.cells.dailyVolume).toBe(0)
  })

  test('coverage equals the count of present metrics', () => {
    const present = new Set(['tvl', 'dailyFees', 'dailyRevenue'])
    const row = projectRow(baseProtocol, present)
    expect(row.coverage).toBe(3)
  })

  test('coverage is 0 when presence set is undefined', () => {
    const row = projectRow(baseProtocol, undefined)
    expect(row.coverage).toBe(0)
    expect(row.cells.tvl).toBe(0)
  })

  test('empty-string category collapses to undefined', () => {
    const row = projectRow({ ...baseProtocol, category: '' }, new Set())
    expect(row.category).toBeUndefined()
  })

  test('preserves a non-empty category verbatim', () => {
    const row = projectRow(baseProtocol, new Set())
    expect(row.category).toBe('Lending')
  })

  test('lowercases chain tokens', () => {
    const row = projectRow(baseProtocol, new Set())
    expect(row.chains).toEqual(['ethereum', 'polygon'])
  })

  test('carries slug and name through unchanged', () => {
    const row = projectRow(baseProtocol, new Set(['tvl']))
    expect(row.slug).toBe('aave')
    expect(row.name).toBe('Aave')
  })
})
