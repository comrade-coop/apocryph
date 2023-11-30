import './style.css'
import { template } from './template'
import { toElement, toJSON } from './field'
import { ProviderConfig, Pod } from 'trusted-pods-proto-ts'

let appRoot = document.querySelector<HTMLDivElement>('#app')

let podTemplate = template()

let form = document.createElement('form')

form.appendChild(toElement(podTemplate, {}))

let submit = document.createElement('button')
submit.type = 'submit'
submit.innerText = 'Send request'
form.appendChild(submit)

form.addEventListener('submit', ev => {
  ev.preventDefault()
  let partialMessage = toJSON(podTemplate)
  console.log(partialMessage)
  console.log({
    pod: new Pod(partialMessage.pod),
    provider: new ProviderConfig(partialMessage.provider)
  })
})

appRoot?.appendChild(form)
