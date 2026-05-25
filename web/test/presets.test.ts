import { describe, test, expect } from 'vitest'
import { CATEGORIES, DIMTYPES, expectedColumnsFor, metricsFor } from '../lib/presets'

describe('CATEGORIES', () => {
  test('is a non-empty array of strings', () => {
    expect(CATEGORIES.length).toBeGreaterThan(0)
    for (const c of CATEGORIES) expect(typeof c).toBe('string')
  })

  test('includes known seeded categories', () => {
    expect(CATEGORIES).toContain('Lending')
    expect(CATEGORIES).toContain('Bridge')
    expect(CATEGORIES).toContain('Derivatives')
    expect(CATEGORIES).toContain('Options')
  })

  test('is sorted', () => {
    const sorted = [...CATEGORIES].sort()
    expect([...CATEGORIES]).toEqual(sorted)
  })
})

describe('expectedColumnsFor', () => {
  test('returns a non-empty array including dailyFees for Lending', () => {
    const cols = expectedColumnsFor('Lending')
    expect(cols.length).toBeGreaterThan(0)
    expect(cols).toContain('dailyFees')
  })

  test('returns [] for an unseeded category', () => {
    expect(expectedColumnsFor('UnseededCategory')).toEqual([])
  })

  test('returns a fresh array on each call', () => {
    const first = expectedColumnsFor('Lending')
    first.pop()
    const second = expectedColumnsFor('Lending')
    expect(second).toContain('dailyFees')
    expect(second.length).toBeGreaterThan(first.length)
  })
})

describe('DIMTYPES', () => {
  test('includes the core dimTypes', () => {
    expect(DIMTYPES).toContain('fees')
    expect(DIMTYPES).toContain('dexs')
    expect(DIMTYPES).toContain('options')
  })

  test('is sorted', () => {
    const sorted = [...DIMTYPES].sort()
    expect([...DIMTYPES]).toEqual(sorted)
  })
})

describe('metricsFor', () => {
  test('returns fees metrics including dailyFees and dailyRevenue', () => {
    const metrics = metricsFor('fees')
    expect(metrics).toContain('dailyFees')
    expect(metrics).toContain('dailyRevenue')
  })

  test('returns [] for an unknown dimType', () => {
    expect(metricsFor('unknown-dimtype')).toEqual([])
  })

  test('returns a fresh array on each call', () => {
    const first = metricsFor('fees')
    first.pop()
    const second = metricsFor('fees')
    expect(second.length).toBeGreaterThan(first.length)
    expect(second).toContain('dailyFees')
  })
})
