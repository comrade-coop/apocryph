import { ReactNode, useState } from 'react'
import './TooltipButton.css'

interface TooltipButton {
  className?: string,
  children: ReactNode,
  tooltip: ReactNode,
  onClick?: () => unknown,
}

function TooltipButton(props: TooltipButton) {
  const [ open, setOpen ] = useState<boolean>(false)

  return (
      <button
        className={'tooltip-button ' + (open ? 'tooltip-button-open ' : '') + (props.className || '')}
        onClick={() => {
          if (open && props.onClick) {
            props.onClick()
          }
          setOpen(!open)
        }}>
        <span className="tooltip-button-content">
          {props.children}
        </span>
        <div className="tooltip-button-tooltip">
          {props.tooltip}
        </div>
      </button>
  )
}

export default TooltipButton
