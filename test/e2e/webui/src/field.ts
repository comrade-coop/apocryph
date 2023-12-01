import { bytesToHex, checksumAddress, hexToBytes, isHex } from 'viem'

export interface Field<T> {
  toElement: () => Element
  toJSON: () => T
}
export type FieldOrRaw<T> = Field<T> | RecordField<T> | T
export type RecordField<T> = {
  [P in keyof T]: FieldOrRaw<T[P]>;
}
export function isField<T> (value: FieldOrRaw<T>): value is Field<T> {
  return typeof value === 'object' && value !== null && 'toElement' in value && 'toJSON' in value
}
function isRecordField<T> (value: RecordField<T> | T): value is RecordField<T> {
  return typeof value === 'object' && value !== null && (Object.getPrototypeOf(value) === Object.prototype || Array.isArray(value))
}

function defaultKeyToText (key: string): string {
  return key.split(/(?<=[a-z])(?=[A-Z])|(?<=[A-Z])(?=[A-Z][a-z])/).map(x => x[0].toUpperCase() + x.slice(1).toLowerCase()).join(' ')
}

export function toElement<T> (value: FieldOrRaw<T>, { keyToText = defaultKeyToText } = {}): Element {
  if (isField(value)) {
    return value.toElement()
  } else if (isRecordField(value)) {
    value = value as RecordField<T> // (Huh?)
    const wrapper = document.createElement('div')
    wrapper.className = 'nested'
    for (const key in value) {
      if (Object.prototype.hasOwnProperty.call(value, key)) {
        const inner = toElement(value[key], { keyToText })
        const row = document.createElement((inner instanceof HTMLInputElement) ? 'label' : 'div')
        const span = document.createElement('span')
        span.className = 'key'
        span.innerText = keyToText(key) + ': '
        row.appendChild(span)
        row.appendChild(inner)
        wrapper.appendChild(row)
      }
    }
    return wrapper
  } else {
    value = value as T
    const element = document.createElement('span')
    if (value instanceof Uint8Array) {
      element.innerText = bytesToHex(value)
    } else {
      element.innerText = '' + (value as any)
    }
    return element
  }
}

export function toJSON<T> (value: FieldOrRaw<T>): T {
  if (isField(value)) {
    return value.toJSON()
  } else if (isRecordField(value)) {
    value = value as RecordField<T> // (Huh?)
    const result: { [P in keyof T]?: T[P] } = Array.isArray(value) ? [] as any : {} /// Siigh
    for (const key in value) {
      if (Object.prototype.hasOwnProperty.call(value, key)) {
        result[key] = toJSON<T[Extract<keyof T, string>]>(value[key])
      }
    }
    return result as T
  } else {
    return value
  }
}

export interface Encoding {
  encode: (_: Uint8Array) => string
  decode: (_: string) => Uint8Array
}
export const Encodings: { 'eth-address': Encoding, 'eth-hash': Encoding } = {
  'eth-address': {
    encode (value: Uint8Array) {
      return checksumAddress(bytesToHex(value, { size: 20 }))
    },
    decode (value: string) {
      if (isHex(value)) {
        return hexToBytes(value)
      } else {
        return hexToBytes('0x{value}')
      }
    }
  },
  'eth-hash': {
    encode (value: Uint8Array) {
      return checksumAddress(bytesToHex(value, { size: 32 }))
    },
    decode (value: string) {
      if (isHex(value)) {
        return hexToBytes(value)
      } else {
        return hexToBytes('0x{value}')
      }
    }
  }
}
export interface EncodingOptions { encoding: keyof typeof Encodings }
export interface SpinnerOptions<T> { min?: T, max?: T, step?: T }
export function isSpinnerOptions<T> (options?: object): options is SpinnerOptions<T> {
  return options !== undefined && ('min' in options || 'max' in options || 'step' in options)
}
export function field (initialValue: string, options?: object): Field<string>
export function field (initialValue: number, options: SpinnerOptions<number>): Field<number>
export function field (initialValue: bigint, options: SpinnerOptions<bigint>): Field<bigint>
export function field (initialValue: Uint8Array, options: EncodingOptions): Field<Uint8Array>
export function field (initialValue: string | number | bigint | Uint8Array, options?: object): Field<string> | Field<number> | Field<bigint> | Field<Uint8Array> {
  if (typeof initialValue === 'string') {
    return new StringField(initialValue)
  }
  if (typeof initialValue === 'number') {
    return new SpinnerNumberField(initialValue, options as SpinnerOptions<number>)
  }
  if (typeof initialValue === 'bigint') {
    return new SpinnerBigintField(initialValue, options as SpinnerOptions<bigint>)
  }
  if (initialValue instanceof Uint8Array) {
    return new EncodedField(initialValue, options as EncodingOptions)
  }
  throw new Error('Invalid options for field()')
}

abstract class InputFieldBase<T> implements Field<T> {
  element: HTMLInputElement
  abstract value: T
  constructor () {
    this.element = document.createElement('input')
  }

  toElement (): Element {
    return this.element
  }

  toJSON (): T {
    return this.value
  }
}

abstract class SpinnerFieldBase<T> extends InputFieldBase<T> {
  constructor (options: SpinnerOptions<T>) {
    super()
    this.element.type = 'number'
    if (options.min != null) this.element.min = options.min.toString()
    if (options.max != null) this.element.max = options.max.toString()
    if (options.step != null) this.element.step = options.step.toString()
  }
}

class StringField extends InputFieldBase<string> {
  constructor (value: string) {
    super()
    this.element.type = 'text'
    this.value = value
  }

  get value (): string {
    return this.element.value
  }

  set value (newValue: string) {
    this.element.value = newValue
  }
}

class SpinnerNumberField extends SpinnerFieldBase<number> {
  constructor (value: number, options: SpinnerOptions<number>) {
    super(options)
    this.value = value
  }

  get value (): number {
    return this.element.valueAsNumber
  }

  set value (newValue: number) {
    this.element.valueAsNumber = newValue
  }
}

class SpinnerBigintField extends SpinnerFieldBase<bigint> {
  constructor (value: bigint, options: SpinnerOptions<bigint>) {
    super(options)
    this.value = value
  }

  get value (): bigint {
    return BigInt(this.element.value)
  }

  set value (newValue: bigint) {
    this.element.value = newValue.toString()
  }
}

class EncodedField extends InputFieldBase<Uint8Array> {
  public encoding: Encoding
  constructor (value: Uint8Array, options: EncodingOptions) {
    super()
    this.encoding = Encodings[options.encoding]
    this.element.type = 'text'
    this.element.addEventListener('change', () => {
      // eslint-disable-next-line no-self-assign
      this.value = this.value // Re-encode
    })
    this.value = value
  }

  get value (): Uint8Array {
    return this.encoding.decode(this.element.value)
  }

  set value (newValue: Uint8Array) {
    this.element.value = this.encoding.encode(newValue)
  }
}
