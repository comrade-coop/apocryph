import { type Transport } from '@connectrpc/connect'
import { type UniversalClientRequest, type UniversalClientResponse } from '@connectrpc/connect/protocol'
import { createTransport } from '@connectrpc/connect/protocol-connect'
import { Uint8ArrayList } from 'uint8arraylist'
import { type Stream } from '@libp2p/interface/connection'

const encoder = new TextEncoder()
const decoder = new TextDecoder()
const eol = encoder.encode('\r\n')

export interface Libp2pTransportOptions {
  dialStream: () => Promise<Stream>
  interceptors: []
  readMaxBytes: 10000
  useBinaryFormat: true
  writeMaxBytes: 10000
}

export function createLibp2pConnectTransport (options: Libp2pTransportOptions): Transport {
  return createTransport({
    async httpClient (req: UniversalClientRequest): Promise<UniversalClientResponse> {
      const stream = await options.dialStream() // NOTE: keepalive could be nice here?

      let requestIsChunked = false
      if (!req.header.has('Content-Length') && req.body !== undefined) {
        requestIsChunked = true
        req.header.append('Transfer-Encoding', 'chunked')
      }
      req.header.append('Host', '127.0.0.1')

      const requestHeadersBuffer = new Uint8ArrayList()
      requestHeadersBuffer.append(encoder.encode(`${req.method} ${req.url} HTTP/1.2`), eol)
      req.header.forEach((value: string, key: string): void => {
        requestHeadersBuffer.append(encoder.encode(`${key}: ${value}`), eol)
      })
      requestHeadersBuffer.append(eol)

      let signalEnd: () => Promise<void>
      const bodyPromise = stream.sink(writeBody(
        new Uint8ArrayList(requestHeadersBuffer),
        req.body,
        requestIsChunked,
        new Promise((resolve) => {
          signalEnd = async () => {
            resolve()
            await bodyPromise
          }
        })
      ))

      let isStatusLine = true
      let isBody = false
      let responseStatus = -1
      const responseHeader = new Headers()

      const buffer = new Uint8ArrayList()

      while (!isBody) {
        try {
          const res = await stream.source.next()
          if (res.done ?? false) {
            throw new Error('Invalid HTTP response (ended too early)')
          }
          buffer.append(res.value)
        } catch (e) {
          console.log(e)
          throw e
        }

        let eolIndex: number
        while ((eolIndex = buffer.indexOf(eol)) !== -1) {
          const line = decoder.decode(buffer.subarray(0, eolIndex))
          buffer.consume(eolIndex + eol.byteLength)

          if (isStatusLine) {
            const match = line.match(/^HTTP\/[0-9]\.[0-9] ([0-9]{3}) .+$/)
            if (match === null) {
              throw new Error('Invalid HTTP response (status line)')
            }
            responseStatus = parseInt(match[1])
            isStatusLine = false
          } else {
            if (line === '') {
              isBody = true
              break
            }
            const match = line.match(/^([^ :]+) *: *(.+)$/)
            if (match === null) {
              throw new Error('Invalid HTTP response (header line)')
            }
            responseHeader.append(match[1], match[2])
          }
        }
      }

      let responseContentLength: number
      const transferEncoding = responseHeader.get('Transfer-Encoding')
      if ((transferEncoding?.indexOf('chunked') ?? -1) !== -1) {
        responseContentLength = -1
      } else {
        responseContentLength = parseInt(responseHeader.get('Content-Length') ?? '-1')
        if (responseContentLength < 0) {
          throw new Error('Invalid HTTP response (content length line)')
        }
      }
      const responseTrailer = new Headers()

      return {
        status: responseStatus,
        header: responseHeader,
        trailer: responseTrailer,
        body: readBody(buffer, stream.source, responseContentLength, responseTrailer, signalEnd)
      }
    },
    baseUrl: '',
    acceptCompression: [],
    compressMinBytes: 10,
    interceptors: [],
    readMaxBytes: options.readMaxBytes,
    sendCompression: null,
    useBinaryFormat: options.useBinaryFormat,
    writeMaxBytes: options.writeMaxBytes
  })
}

async function * writeBody (buffer: Uint8ArrayList, body?: AsyncIterable<Uint8Array>, isChunked: boolean = false, endPromise?: Promise<void>): AsyncGenerator<Uint8ArrayList, void, undefined> {
  yield buffer

  if (body !== undefined) {
    if (isChunked) {
      for await (const chunk of body) {
        if (chunk.length > 0) {
          yield new Uint8ArrayList(encoder.encode(chunk.byteLength.toString(16)), eol, chunk, eol)
        }
      }
      yield new Uint8ArrayList(encoder.encode('0\r\n\r\n\r\n'))
    } else {
      for await (const chunk of body) {
        yield new Uint8ArrayList(chunk)
      }
    }
  }

  if (endPromise !== undefined) {
    await endPromise
  }
}

async function * readBody (buffer: Uint8ArrayList, source: AsyncGenerator<Uint8ArrayList>, contentLength: number, trailers: Headers, signalEnd?: () => Promise<void>): AsyncGenerator<Uint8Array, void, undefined> {
  let remainingChunkBytes = 0
  let remainingChunkEolBytes = 0
  let isTrailers = false

  while (true) {
    if (contentLength === -1) {
      while (!isTrailers) { // Best of luck to whoever might end up having to debug this.. I am sorry :/
        if (remainingChunkEolBytes === 0) {
          if (remainingChunkBytes === 0) {
            const eolIndex = buffer.indexOf(eol)
            if (eolIndex !== -1) {
              const chunkLine = decoder.decode(buffer.subarray(0, eolIndex))
              buffer.consume(eolIndex + eol.byteLength)
              const chunkSize = parseInt(chunkLine, 16)
              console.log(chunkLine)
              if (chunkSize === 0) {
                isTrailers = true
                break
              }
              remainingChunkBytes = chunkSize + eol.byteLength
            } else {
              break
            }
          } else { // (remainingChunkBytes !== 0)
            if (buffer.byteLength >= remainingChunkBytes) {
              yield buffer.subarray(0, remainingChunkBytes)
              buffer.consume(remainingChunkBytes)
              remainingChunkBytes = 0
              remainingChunkEolBytes = eol.byteLength
            } else {
              if (buffer.byteLength > 0) {
                yield buffer.subarray()
                remainingChunkBytes -= buffer.byteLength
                buffer = new Uint8ArrayList()
              }
              break
            }
          }
        } else { // (remainingChunkEolBytes !== 0)
          if (buffer.byteLength < eol.byteLength) {
            break
          }
          const assertEol = buffer.subarray(0, eol.byteLength)
          for (let i = 0; i < eol.byteLength; i++) {
            if (assertEol[i] !== eol[i]) {
              throw new Error('Invalid HTTP response (chunk end)')
            }
          }
          buffer.consume(eol.byteLength)
        }
      }
      if (isTrailers) { // Lack of else is important (so we parse the trailers in the buffer right away)
        let eolIndex: number
        while ((eolIndex = buffer.indexOf(eol)) !== -1) {
          const line = decoder.decode(buffer.subarray(0, eolIndex))
          buffer.consume(eolIndex + eol.byteLength)
          if (line === '') {
            break
          }
          const match = line.match(/^([^ :]+) *: *(.+)$/)
          if (match === null) {
            throw new Error('Invalid HTTP response (header line)')
          }
          trailers.append(match[1], match[2])
        }
      }
    } else {
      if (buffer.byteLength >= contentLength) {
        yield buffer.subarray(0, contentLength)
        buffer.consume(contentLength)
        contentLength = 0
        break
      } else {
        if (buffer.byteLength > 0) {
          yield buffer.subarray()
          contentLength -= buffer.byteLength
          buffer = new Uint8ArrayList()
        }
      }
    }

    const res = await source.next()
    if (res.done ?? false) {
      break
    }
    buffer.append(res.value)
  }

  if (signalEnd !== undefined) {
    await signalEnd()
  }
}
