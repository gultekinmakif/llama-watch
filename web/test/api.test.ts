import { describe, test, expect } from 'vitest'
import { parseErrorEnvelope } from '../lib/api'

describe('parseErrorEnvelope', () => {
  test('extracts code and message from a well-formed envelope', async () => {
    const response = new Response(
      JSON.stringify({ error: { code: 'not_found', message: 'protocol missing' } }),
      { status: 404 },
    )
    const result = await parseErrorEnvelope(response)
    expect(result).toEqual({ code: 'not_found', message: 'protocol missing' })
  })

  test('falls back to unknown when body has no error key', async () => {
    const response = new Response(JSON.stringify({ status: 'bad' }), { status: 500 })
    const result = await parseErrorEnvelope(response)
    expect(result).toEqual({ code: 'unknown', message: 'http 500' })
  })

  test('falls back to unknown when body is non-JSON', async () => {
    const response = new Response('not json at all', { status: 502 })
    const result = await parseErrorEnvelope(response)
    expect(result).toEqual({ code: 'unknown', message: 'http 502' })
  })

  test('falls back when error.code is the wrong type', async () => {
    const response = new Response(
      JSON.stringify({ error: { code: 42, message: 'numeric code' } }),
      { status: 400 },
    )
    const result = await parseErrorEnvelope(response)
    expect(result).toEqual({ code: 'unknown', message: 'numeric code' })
  })

  test('falls back when error.message is the wrong type', async () => {
    const response = new Response(
      JSON.stringify({ error: { code: 'bad_request', message: null } }),
      { status: 400 },
    )
    const result = await parseErrorEnvelope(response)
    expect(result).toEqual({ code: 'bad_request', message: 'http 400' })
  })

  test('falls back when error is null', async () => {
    const response = new Response(JSON.stringify({ error: null }), { status: 500 })
    const result = await parseErrorEnvelope(response)
    expect(result).toEqual({ code: 'unknown', message: 'http 500' })
  })
})
