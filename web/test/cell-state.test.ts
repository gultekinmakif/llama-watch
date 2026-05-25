import { describe, test, expect } from 'vitest'

import { classifyByCategory, classifyCell } from '../lib/cell-state'

describe('classifyCell (dimType bundle scoring)', () => {
  test('metric not in any bundle is unexpected when present, na when absent', () => {
    expect(classifyCell([], 'tvl', true)).toBe('unexpected')
    expect(classifyCell([], 'tvl', false)).toBe('na')
  })

  test('metric inside a registered dimType bundle is expected', () => {
    expect(classifyCell(['fees'], 'dailyFees', true)).toBe('present')
    expect(classifyCell(['fees'], 'dailyFees', false)).toBe('missing')
  })

  test('metric outside any bundle is na or unexpected', () => {
    expect(classifyCell(['fees'], 'dailyVolume', false)).toBe('na')
    expect(classifyCell(['fees'], 'dailyVolume', true)).toBe('unexpected')
  })
})

describe('classifyByCategory (category-driven scoring)', () => {
  test('Lending expects its full curated bundle', () => {
    expect(classifyByCategory('Lending', 'dailyFees', true)).toBe('present')
    expect(classifyByCategory('Lending', 'dailyFees', false)).toBe('missing')
    expect(classifyByCategory('Lending', 'dailyVolume', false)).toBe('missing')
  })

  test('metric outside a category expected set is unexpected when present', () => {
    expect(classifyByCategory('Lending', 'openInterestAtEnd', true)).toBe('unexpected')
    expect(classifyByCategory('Lending', 'openInterestAtEnd', false)).toBe('na')
  })

  test('uncategorized rows do not fire unexpected on emissions', () => {
    // No curated expected set means we cannot distinguish expected from unexpected;
    // treat emissions as present so the matrix does not light up yellow everywhere.
    expect(classifyByCategory(undefined, 'dailyFees', true)).toBe('present')
    expect(classifyByCategory(undefined, 'dailyFees', false)).toBe('na')
  })

  test('category absent from CATEGORIES_EXPECTED falls back to the safe path', () => {
    // 'Liquidations' has fewer than 3 protocols and was not auto-derived.
    expect(classifyByCategory('Liquidations', 'dailyFees', true)).toBe('present')
    expect(classifyByCategory('Liquidations', 'dailyFees', false)).toBe('na')
  })

  test('metric not in curated set is na when absent, present when present with undefined category', () => {
    expect(classifyByCategory(undefined, 'tvl', true)).toBe('present')
    expect(classifyByCategory(undefined, 'tvl', false)).toBe('na')
  })
})
