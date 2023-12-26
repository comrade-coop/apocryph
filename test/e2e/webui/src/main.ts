// SPDX-License-Identifier: GPL-3.0

import './style.css'
import { template } from './template'
import { toElement, toJSON } from './field'
import {
  ProviderConfig,
  Pod,
  PaymentChannelConfig
} from 'apocryph-proto-ts'
import { fundPaymentChannel } from './fund'
import { provisionPod } from './provision'
import { bytesToHex, hexToBytes } from 'viem'

const appRoot = document.querySelector<HTMLDivElement>('#app')

const podTemplate = template()

const form = document.createElement('form')

form.appendChild(toElement(podTemplate, {}))

const error = document.createElement('div')
error.className = 'error-text'
form.appendChild(error)

const submit = document.createElement('button')
submit.type = 'submit'
submit.innerText = 'Send request'
form.appendChild(submit)

const results = document.createElement('div')
form.appendChild(results)

let submitPromise: Promise<any> | null = null

form.addEventListener('submit', (ev) => {
  ev.preventDefault()
  if (submitPromise === null) {
    submit.classList.add('loading')
    error.innerText = ''
    results.innerHTML = ''
    submitPromise = (async () => {
      const values = toJSON(podTemplate)
      const deployment = {
        pod: new Pod(values.pod),
        payment: new PaymentChannelConfig(values.payment),
        provider: new ProviderConfig(values.provider)
      }
      deployment.payment.podID = hexToBytes(
        `0x${crypto.randomUUID().replace(/-/g, '')}`,
        { size: 32 }
      )
      {
        const row = document.createElement('div')
        const name = document.createElement('span')
        name.className = 'key'
        name.innerText = 'Pod id: '
        row.appendChild(name)
        const id = document.createElement('span')
        id.innerText = bytesToHex(deployment.payment.podID)
        row.appendChild(id)
        results.appendChild(row)
      }

      await fundPaymentChannel(deployment, values.funds, {
        unlockTime: values.unlockTime,
        mintFunds: true
      })
      const response = await provisionPod(deployment)

      for (const address of response.addresses) {
        const row = document.createElement('div')
        const name = document.createElement('span')
        name.className = 'key'
        name.innerText = `Container ${address.containerName}: -> `
        row.appendChild(name)
        const addr = document.createElement('a')
        addr.href = `http://${address.multiaddr.split('/').slice(-1)[0]}:1234` // TODO: fix when we have actual addresses
        addr.innerText = address.multiaddr
        row.appendChild(addr)
        results.appendChild(row)
      }
    })()
      .catch((err) => {
        error.innerText = err.toString().split('\n')[0]
        console.error(err)
      })
      .then(() => {
        submitPromise = null
        submit.classList.remove('loading')
      })
  }
})

appRoot?.appendChild(form)
