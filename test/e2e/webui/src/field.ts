// SPDX-License-Identifier: GPL-3.0

import { bytesToHex, checksumAddress, hexToBytes, isHex } from 'viem'

// NOTE: A lot of this code can probably be converted to simple React/Vue/Svelte with components instead of fields, and it would probably work better too. It seemed simpler to do it this way when I was starting out, as it would also give the user full view of the configuration structure, but it's probably better to just hide it / make it available in some advanced view instead, as it is not particularly pretty to deal with.

/**
 * Represents a field that can be converted to a raw JSON value or to an HTML element. Compatible with the interface expected by JSON.stringify()
 */
export interface Field<T> {
  toElement: () => Element
  toJSON: () => T
}

/**
 * Represents either a raw value of type T or a field that produces a value of type T
 */
export type FieldOrRaw<T> = Field<T> | RecordField<T> | T

/**
 * Represents an object (or array) that has all its properties be either of the proper type or a field
 */
export type RecordField<T> = {
  [P in keyof T]: FieldOrRaw<T[P]>
}

/**
 * Type matcher for the Field<T> interface
 */
export function isField<T>(value: FieldOrRaw<T>): value is Field<T> {
  return (
    typeof value === 'object' &&
    value !== null &&
    'toElement' in value &&
    'toJSON' in value
  )
}

/**
 * Type matcher for the RecordField<T> interface. Note that due to how general that interface is, nearly every object will pass, other than class instances.
 */
export function isRecordField<T>(value: RecordField<T> | T): value is RecordField<T> {
  return (
    typeof value === 'object' &&
    value !== null &&
    (Object.getPrototypeOf(value) === Object.prototype || Array.isArray(value))
  )
}

/**
 * Convert a key to its textual representation (default implementation)
 */
function defaultKeyToText(key: string): string {
  return key
    .split(/(?<=[a-z])(?=[A-Z])|(?<=[A-Z])(?=[A-Z][a-z])/)
    .map((x) => x[0].toUpperCase() + x.slice(1).toLowerCase())
    .join(' ')
}

/**
 * Convert an object/field to an element that can be inserted in the DOM tree.
 *
 * @remarks
 * Note that some fields might cache the element, and thus calling toElement on the same object could result in errors.
 *
 * @param value the object to get the element for
 * @param opts options
 * @param opts.keyToText a function converting an object key to a label
 */
export function toElement<T>(
  value: FieldOrRaw<T>,
  { keyToText = defaultKeyToText } = {}
): Element {
  if (isField(value)) {
    return value.toElement()
  } else if (isRecordField(value)) {
    value = value as RecordField<T> // (Huh?)
    const wrapper = document.createElement('div')
    wrapper.className = 'nested'
    for (const key in value) {
      if (Object.prototype.hasOwnProperty.call(value, key)) {
        const inner = toElement(value[key], { keyToText })
        const row = document.createElement(
          inner instanceof HTMLInputElement ? 'label' : 'div'
        )
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

/**
 * Convert an object/field to a JSON/object structure with all the regular parameters
 *
 * @param value the object to get the object for
 */
export function toJSON<T>(value: FieldOrRaw<T>): T {
  if (isField(value)) {
    return value.toJSON()
  } else if (isRecordField(value)) {
    value = value as RecordField<T> // (Huh?)
    const result: { [P in keyof T]?: T[P] } = Array.isArray(value)
      ? ([] as any)
      : {} /// Siigh
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

/**
 * Encoding for binary fields
 */
export interface Encoding {
  encode: (_: Uint8Array) => string
  decode: (_: string) => Uint8Array
}

// NOTE: typing is a bit odd here; ideally typescript would autodetect the list of encoding names without having to duplicate it, but eslint compalins about it.
/**
 * Recognized binary field encodings
 */
export const Encodings: { 'eth-address': Encoding; 'eth-hash': Encoding } = {
  'eth-address': {
    encode(value: Uint8Array) {
      return checksumAddress(bytesToHex(value, { size: 20 }))
    },
    decode(value: string) {
      if (isHex(value)) {
        return hexToBytes(value)
      } else {
        return hexToBytes('0x{value}')
      }
    }
  },
  'eth-hash': {
    encode(value: Uint8Array) {
      return checksumAddress(bytesToHex(value, { size: 32 }))
    },
    decode(value: string) {
      if (isHex(value)) {
        return hexToBytes(value)
      } else {
        return hexToBytes('0x{value}')
      }
    }
  }
}

/**
 * Options for encoded binary fields
 */
export interface EncodingOptions {
  encoding: keyof typeof Encodings
}

/**
 * Options for spinner (input type=number) fields
 */
export interface SpinnerOptions<T> {
  min?: T
  max?: T
  step?: T
}

/**
 * Type matcher for SpinnerOptions<T>
 */
export function isSpinnerOptions<T>(
  options?: object
): options is SpinnerOptions<T> {
  return (
    options !== undefined &&
    ('min' in options || 'max' in options || 'step' in options)
  )
}

/**
 * Create a string field with a given initial value
 */
export function field(initialValue: string, options?: object): Field<string>
/**
 * Create a number field with the given spinner options
 */
export function field(
  initialValue: number,
  options: SpinnerOptions<number>
): Field<number>
/**
 * Create a bigint field with the given spinner options
 */
export function field(
  initialValue: bigint,
  options: SpinnerOptions<bigint>
): Field<bigint>
/**
 * Create a binary field with the given encoding options
 */
export function field(
  initialValue: Uint8Array,
  options: EncodingOptions
): Field<Uint8Array>
/**
 * Create a field with the given initial value and options
 *
 * @param initialValue the default value of the field
 * @param options additional configuration options for the field
 */
export function field(
  initialValue: string | number | bigint | Uint8Array,
  options?: object
): Field<string> | Field<number> | Field<bigint> | Field<Uint8Array> {
  if (typeof initialValue === 'string') {
    return new StringField(initialValue)
  }
  if (typeof initialValue === 'number') {
    return new SpinnerNumberField(
      initialValue,
      options as SpinnerOptions<number>
    )
  }
  if (typeof initialValue === 'bigint') {
    return new SpinnerBigintField(
      initialValue,
      options as SpinnerOptions<bigint>
    )
  }
  if (initialValue instanceof Uint8Array) {
    return new EncodedField(initialValue, options as EncodingOptions)
  }
  throw new Error('Invalid options for field()')
}

/**
 * Class wrapping an input field as a Field<T>
 */
abstract class InputFieldBase<T> implements Field<T> {
  element: HTMLInputElement
  abstract value: T
  constructor() {
    this.element = document.createElement('input')
  }

  toElement(): Element {
    return this.element
  }

  toJSON(): T {
    return this.value
  }
}

/**
 * Class capturing common functionallity of spinner fields
 */
abstract class SpinnerFieldBase<T> extends InputFieldBase<T> {
  constructor(options: SpinnerOptions<T>) {
    super()
    this.element.type = 'number'
    if (options.min != null) this.element.min = options.min.toString()
    if (options.max != null) this.element.max = options.max.toString()
    if (options.step != null) this.element.step = options.step.toString()
  }
}

/**
 * Class representing a string field
 */
class StringField extends InputFieldBase<string> {
  constructor(value: string) {
    super()
    this.element.type = 'text'
    this.value = value
  }

  get value(): string {
    return this.element.value
  }

  set value(newValue: string) {
    this.element.value = newValue
  }
}

/**
 * Class representing a number field
 */
class SpinnerNumberField extends SpinnerFieldBase<number> {
  constructor(value: number, options: SpinnerOptions<number>) {
    super(options)
    this.value = value
  }

  get value(): number {
    return this.element.valueAsNumber
  }

  set value(newValue: number) {
    this.element.valueAsNumber = newValue
  }
}

/**
 * Class representing a bigint field
 */
class SpinnerBigintField extends SpinnerFieldBase<bigint> {
  constructor(value: bigint, options: SpinnerOptions<bigint>) {
    super(options)
    this.value = value
  }

  get value(): bigint {
    return BigInt(this.element.value)
  }

  set value(newValue: bigint) {
    this.element.value = newValue.toString()
  }
}

/**
 * Class representing a binary encoded field
 */
class EncodedField extends InputFieldBase<Uint8Array> {
  public encoding: Encoding
  constructor(value: Uint8Array, options: EncodingOptions) {
    super()
    this.encoding = Encodings[options.encoding]
    this.element.type = 'text'
    this.element.addEventListener('change', () => {
      // eslint-disable-next-line no-self-assign
      this.value = this.value // Re-encode
    })
    this.value = value
  }

  get value(): Uint8Array {
    return this.encoding.decode(this.element.value)
  }

  set value(newValue: Uint8Array) {
    this.element.value = this.encoding.encode(newValue)
  }
}
