import { ReactElement, useEffect, useState } from 'react'
import { useAccount, useConnect, usePublicClient, useWalletClient } from 'wagmi'
import { injected } from 'wagmi/connectors'
import { formatUnits, parseUnits } from 'viem'

import { Light as SyntaxHighlighter } from 'react-syntax-highlighter'
import { tomorrowNight as syntaxStyle } from 'react-syntax-highlighter/dist/esm/styles/hljs'

import BlurUpdatedInput from './BlurUpdatedInput'
import ActionPopButton from './ActionPopButton'
import { watchAvailableFunds, depositFunds } from './contracts'
import apocryphLogo from '/apocryph.svg?url'
import metamaskLogo from '/metamask.svg?url'
import './App.css'
import { getSiweToken } from './signin'
import codeExamples, { envExample } from './codeExamples'
import { OpenExternalLink } from './icons'

const attestationLink: string | undefined = import.meta.env.VITE_PUBLIC_ATTESTATION_URL
const documentationLink = "https://comrade-coop.github.io/s3-aapp/"
const s3AappHost = (import.meta.env.VITE_GLOBAL_HOST || "s3-aapp.kubocloud.io").trim()
const s3consoleAappHost = (import.meta.env.VITE_GLOBAL_HOST_CONSOLE || "console-s3-aapp.kubocloud.io").trim()

function App() {
  const account = useAccount()
  const publicClient = usePublicClient()
  const walletClient = useWalletClient()
  const { connect } = useConnect()

  const oneGb = parseUnits('1', 6)
  const [ amountGb, setAmountGbRaw ] = useState<bigint>(100n * oneGb)
  const [ durationMultiplier, setDurationMultiplier ] = useState<number>(60 * 60 * 24 * 365)
  const currency = 'S3T'
  const decimals = 18
  const priceGbSec = parseUnits('0.000004', decimals)
  const [ funds, setFunds ] = useState<bigint>(() => BigInt(durationMultiplier) * amountGb * priceGbSec / oneGb)
  const [ existingDeposit, setExistingDeposit ] = useState<bigint | undefined>(undefined)
  const [ depositInProgress, setDepositInProgress ] = useState(false)
  const [ depositError, setDepositError ] = useState('')
  const [ siweToken, setSiweToken ] = useState<string>()
  const [ consoleAccessLink, setConsoleAccessLink ] = useState<string>()
  const [ profitText, setProfitText ] = useState<string>("Ok?")
  const [ codeExample, setCodeExample ] = useState<keyof typeof codeExamples>("Go")

  const bucketId = `${account.address?.slice(2)?.toLowerCase()}`
  const bucketLink = `https://${s3AappHost}` // ${bucketId}.${s3AappHost}
  const bucketLinkHref = `//${s3AappHost}`
  const consoleLink = `https://${s3consoleAappHost}` // `${bucketId}.${s3consoleAappHost}`
  const consoleLinkHref = `//${s3consoleAappHost}/browser/${bucketId}`
  const minDeposit = parseUnits('10', 18) // TODO

  const duration: number = Number(funds) / Number(amountGb) / Number(priceGbSec) * Number(oneGb)
  function setDuration(newDuration: number) {
    setFunds(BigInt(newDuration) * amountGb * priceGbSec / oneGb)
  }
  function setAmountGb(newAmountGb: bigint) {
    if (newAmountGb < 1n) newAmountGb = 1n

    setFunds(funds * newAmountGb / amountGb)
    setAmountGbRaw(newAmountGb)
  }

  useEffect(() => {
    if (publicClient && account?.address) {
      return watchAvailableFunds(publicClient, account.address, (availableFunds) => {
        setExistingDeposit(availableFunds)
      })
    }
  }, [publicClient, account])

  async function topUpDeposit() {
    if (existingDeposit !== undefined && publicClient && walletClient?.data) {
      setDepositInProgress(true)
      setDepositError('')
      try {
        if (funds <= 0n) {
          await depositFunds(publicClient, walletClient.data, 0n)
        } else {
          await depositFunds(publicClient, walletClient.data, funds + minDeposit)
        }
      } catch(err) {
        setDepositError(err + '')
      }
      setDepositInProgress(false)
    }
  }
  async function openConsole() {
    if (walletClient?.data) {
      if (consoleAccessLink) {
        open(consoleAccessLink)
      } else {
        const token = await getSiweToken(walletClient.data, 3600)

        const consoleAccessLink = `${consoleLink}/x/apocryphLogin#${encodeURIComponent(token)}#/browser/${bucketId}`
        setConsoleAccessLink(consoleAccessLink)
        setTimeout(() => {
          setConsoleAccessLink(undefined)
        }, 3600)
        open(consoleAccessLink)
      }
    }
  }
  async function getLonglivedToken() {
    if (walletClient?.data) {
      const token = await getSiweToken(walletClient!.data!, undefined)
      setSiweToken(token)
    }
  }

  let showNextStep = true
  function step(sectionElement: ReactElement, completionCondition: boolean): ReactElement {
    if (showNextStep) {
      showNextStep = completionCondition
      return sectionElement
    }
    return <></>
  }

  return (
    <>
      {depositError != '' ? <div className='error-toast' key={depositError} onClick={() => setDepositError('')}>{depositError}</div> : <></>}
      <img src={apocryphLogo} alt="Apocryph Logo" />
      <h1>Get your S3-compatible bucket!</h1>
      <section>
        <p className="hero">Hosting your S3-compatible data buckets in the Apocryph S3 network allows for the ultimate privacy peace of mind, through trasparent encryption at-rest and cryptocurrency-enabled payments.<br/><a href={documentationLink}>Read more.</a></p>
      </section>
      {step(<section>
        <h2>Step 1: Connect</h2>
        <div className="button-card">
          <button onClick={() => connect({connector: injected()})}>
            {
              account.isConnected ? <>Connected</> :
              account.isDisconnected ? <>Connect with MetaMask</> :
              account.isConnecting ? <>Connecting to MetaMask...</> :
              account.isReconnecting ? <>Reconnecting to MetaMask...</> :
              ''
            }
            {account.address ? <> ({account.address.slice(0, 7) + '...' + account.address.slice(-5)})</> : ''}
            <img src={metamaskLogo} alt="Metamask Logo" className='icon' />
          </button>
        </div>
      </section>, account.isConnected)}
      {step(<section>
        <h2>Step 2: Fund</h2>
        <label>
          <span>Data you want to store in S3</span>
          <BlurUpdatedInput
            value={amountGb}
            stringify={v => formatUnits(v, 6)}
            parse={v => parseUnits(v, 6)}
            onChange={setAmountGb}/>
          <span className="fake-field">GB</span>
        </label>
        <label>
          <span>Duration to store it for</span>
          <BlurUpdatedInput
            value={duration}
            stringify={v => (Number(v) / durationMultiplier).toFixed(durationMultiplier == 1 ? 0 :  2)}
            parse={v => (parseFloat(v) || 0) * durationMultiplier}
            onChange={setDuration}/>
          <select
            value={durationMultiplier}
            onChange={e => setDurationMultiplier(parseInt(e.target.value))}
          >
            <option value={1}>seconds</option>
            <option value={60}>minutes</option>
            <option value={60 * 60}>hours</option>
            <option value={60 * 60 * 24}>days</option>
            <option value={60 * 60 * 24 * 7}>weeks</option>
            <option value={60 * 60 * 24 * 30}>months</option>
            <option value={60 * 60 * 24 * 365}>years</option>
          </select>
        </label>
        <label>
          <span>Current storage price</span>
          <span className="fake-field">{formatUnits(priceGbSec, decimals)}</span>
          <span className="fake-field"> {currency}/GB/s</span>
        </label>
        <label>
          <span>Total required authorization</span>
          <BlurUpdatedInput
            value={funds}
            stringify={v => formatUnits(v, decimals)}
            parse={v => parseUnits(v, decimals)}
            onChange={setFunds}/>
          <span className="fake-field"> {currency}</span>
        </label>
          <label>
            <span>Existing authorization</span>
            <span className="fake-field">{existingDeposit === undefined ? 'Loading...' : formatUnits(existingDeposit, decimals)}</span>
            <span className="fake-field">{currency}</span>
          </label>
        <div className="button-card">
          <button onClick={() => topUpDeposit()}>
            {
              existingDeposit === undefined ? <>Loading...</> :
              depositInProgress ? <>Processing...</> :
              existingDeposit <= minDeposit ? <>Authorize! ({formatUnits(existingDeposit - minDeposit - funds, decimals)} {currency})</> :
              funds > existingDeposit - minDeposit ? <>Top-up authorization ({formatUnits(existingDeposit - minDeposit - funds, decimals)} {currency})</> :
              funds <= 0n ? <>Remove authorization (+{formatUnits(existingDeposit - minDeposit - funds, decimals)} {currency})</> :
              <>Reduce authorization (+{formatUnits(existingDeposit - minDeposit - funds, decimals)} {currency})</>
            }
          </button>
        </div>
      </section>, existingDeposit !== undefined && existingDeposit > 0n)}
      {step(<section>
        <h2>Step 3: Access</h2>
        <label>
          <span>Console </span>
          <a className="fake-field" href={consoleLinkHref}>{consoleLink}</a>
          <ActionPopButton onClick={() => navigator.clipboard.writeText(bucketLinkHref)}>Copy</ActionPopButton>
        </label>
        <label>
          <span>S3 endpoint URL </span>
          <a className="fake-field" href={bucketLinkHref}>{bucketLink}</a>
          <ActionPopButton onClick={() => navigator.clipboard.writeText(bucketLinkHref)}>Copy</ActionPopButton>
        </label>
        <label>
          <span>Bucket ID </span>
          <span className="fake-field">{bucketId}</span>
          <ActionPopButton onClick={() => navigator.clipboard.writeText(bucketId)}>Copy</ActionPopButton>
        </label>
        <div className="button-card">
          <button onClick={() => openConsole()}>
            Launch Console <OpenExternalLink/>
          </button>
          <button onClick={() => getLonglivedToken()}>
            Configure programmatic access...
          </button>
        </div>
      </section>, !!siweToken)}
      {step(<section className="two-columns">
        <h2>Step 4: Hack away! </h2>
        <div className="button-card">
          <select value={codeExample} onChange={e => setCodeExample(e.target.value as keyof typeof codeExamples)}>
            {Object.keys(codeExamples).map(x => <option key={x}>{x}</option>)}
          </select>
        </div>
        <div className="button-over-code">
          <ActionPopButton onClick={() => navigator.clipboard.writeText(envExample(bucketLink, siweToken!))}>Copy!</ActionPopButton>
        </div>
        <SyntaxHighlighter language={'bash'} style={syntaxStyle} className="code" wrapLines={true}>
          {envExample(bucketLink, siweToken!)}
        </SyntaxHighlighter>
        <div className="button-over-code">
          <ActionPopButton onClick={() => {
            navigator.clipboard.writeText(codeExamples[codeExample].code)
          }}>Copy!</ActionPopButton>
        </div>
        <SyntaxHighlighter language={codeExamples[codeExample]?.language} style={syntaxStyle} className="code">
          {codeExamples[codeExample]?.code}
        </SyntaxHighlighter>
      </section>, true)}
      {step(<section>
        <h2>Step 5: Profit</h2>
        <div className="button-card">
          <button onClick={() => {
            if (profitText.length < 67) {
              setProfitText(profitText + (Math.random() < 0.6 ? "?" : "!"))
            } else {
              setProfitText(["Yay", "Whee", "Wohoo", "Huzzah", "Wow"][Math.floor(Math.random() * 5)])
            }
          }}>{profitText}</button>
        </div>
      </section>, true)}
      <a href={documentationLink} className="read-the-docs" target="_blank">Documentation <OpenExternalLink/></a>
      { attestationLink ? <a href={attestationLink} className="read-the-docs" target="_blank">View Attestation <OpenExternalLink/></a> : '' }
    </>
  )
}

export default App
