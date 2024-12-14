import { ReactNode, useState } from 'react'
import './ActionPopButton.css'

interface ActionPopButton {
  className?: string,
  children: ReactNode,
  onClick: () => unknown,
  popText?: string,
  disabled?: boolean,
}

function ActionPopButton(props: ActionPopButton) {
  const [ animating, setAnimating ] = useState<boolean>(false)

  return (
      <button
        className={"action-pop-button " + (animating ? 'action-pop-animate ' : '') + (props.className || '')}
        onClick={() => {
          setAnimating(true)
          props.onClick()
        }}
        onAnimationEnd={() => {
          // https://stackoverflow.com/a/34700273
          setAnimating(false)
        }}
        data-pop-text={props.popText}
        disabled={props.disabled}>
        <span className="action-pop-content">
          {props.children}
        </span>
      </button>
  )
}

export default ActionPopButton
