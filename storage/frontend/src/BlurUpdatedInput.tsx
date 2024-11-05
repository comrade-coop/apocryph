import { useState } from 'react'

interface BlurUpdatedInputProps<T> {
  stringify: (value: T) => string,
  parse: (value: string) => T,
  value: T,
  onChange: (newValue: T) => void
}

function BlurUpdatedInput<T>(props: BlurUpdatedInputProps<T>) {
  const [ isFocused, setFocused ] = useState<boolean>(false)
  const [ valueString, setValueString ] = useState<string>(() => props.stringify(props.value))

  const valueDisplayed = isFocused ? valueString : props.stringify(props.value)

  return (
      <input
        value={valueDisplayed}
        onChange={e => {
          setValueString(e.target.value)
          props.onChange(props.parse(e.target.value))
        }}
        onFocus={_ => {
          setValueString(props.stringify(props.value))
          setFocused(true)
        }}
        onBlur={_ => {
          setFocused(false)
          setValueString("")
        }}/>
  )
}

export default BlurUpdatedInput
