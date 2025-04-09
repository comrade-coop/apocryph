import { ReactElement, useEffect, useState } from 'react'
import { useAccount, useConnect, usePublicClient, useWalletClient } from 'wagmi'
import { injected } from 'wagmi/connectors'
import { formatUnits, parseUnits } from 'viem'

import { Light as SyntaxHighlighter } from 'react-syntax-highlighter'
import { tomorrowNight as syntaxStyle } from 'react-syntax-highlighter/dist/esm/styles/hljs'

import BlurUpdatedInput from './BlurUpdatedInput'
import ActionPopButton from './ActionPopButton'
import { watchAvailableFunds, depositFunds, aappAddress, paymentAddress, debugMintFunds } from './contracts'
import apocryphLogo from '/apocryph.svg?url'
import metamaskLogo from '/metamask.svg?url'
import './App.css'
import { getSiweToken } from './signin'
import codeExamples, { envExample } from './codeExamples'
import { Error, InfoCircle, OpenExternalLink } from './icons'
import TooltipButton from './TooltipButton'

const attestationLink: string | undefined = import.meta.env.VITE_PUBLIC_ATTESTATION_URL
const documentationLink = "https://comrade-coop.github.io/s3-aapp/"
const s3AappHost = (import.meta.env.VITE_GLOBAL_HOST || "s3-aapp.kubocloud.io").trim()
const s3consoleAappHost = (import.meta.env.VITE_GLOBAL_HOST_CONSOLE || "console-s3-aapp.kubocloud.io").trim()

function App() {
  const account = useAccount()
  const publicClient = usePublicClient()
  const walletClient = useWalletClient()
  const { connect } = useConnect()

  const oneGb = 1000*1000
  const [ amountGb, setAmountGbRaw ] = useState<bigint>(100n * 1000n*1000n)
  const [ durationMultiplier, setDurationMultiplier ] = useState<number>(12)
  const currency = 'USDC'
  const decimals = 6
  const priceGbMonth = parseUnits('0.025', decimals)
  const [ funds, setFunds ] = useState<bigint>(() => BigInt(Math.round(durationMultiplier * Number(amountGb * priceGbMonth) / oneGb)))
  const [ existingDeposit, setExistingDeposit ] = useState<bigint | undefined>(undefined)
  const [ depositInProgress, setDepositInProgress ] = useState(false)
  const [ balance, setBalance ] = useState<bigint | undefined>(undefined)
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
  const minDeposit = parseUnits('1', 6)

  const duration: number = Number(funds) / Number(amountGb) / Number(priceGbMonth) * oneGb
  function setDuration(newDuration: number) {
    setFunds(BigInt(Math.round(newDuration * Number(amountGb * priceGbMonth) / oneGb)))
  }
  function setAmountGb(newAmountGb: bigint) {
    if (newAmountGb < 1n) newAmountGb = 1n

    setFunds(funds * newAmountGb / amountGb)
    setAmountGbRaw(newAmountGb)
  }

  useEffect(() => {
    if (publicClient && account?.address) {
      return watchAvailableFunds(publicClient, account.address, (availableFunds, balance) => {
        setExistingDeposit(availableFunds)
        setBalance(balance)
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
        <p className="hero">Hosting your S3-compatible data buckets in the Apocryph S3 network allows for the ultimate privacy peace of mind, through trasparent encryption at-rest and cryptocurrency-enabled payments.<br/><a href={documentationLink}>Read more</a> <OpenExternalLink/></p>
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
            <option value={1 / 30 / 24}>hours</option>
            <option value={1 / 30}>days</option>
            <option value={1 / 30 * 7}>weeks</option>
            <option value={1}>months</option>
            <option value={12}>years</option>
          </select>
        </label>
        <label>
          <span>Current storage price</span>
          <span className="fake-field">{formatUnits(priceGbMonth, decimals)}</span>
          <span className="fake-field"> {currency}/GB/month</span>
        </label>
        <label>
          <span>
            Estimated funds needed
            <TooltipButton
              tooltip={<>
                The amount of {currency} that will be needed to store your data for the specified amount of time. <br/>
                <a href={documentationLink + "/PAYMENT.html#authorized-funds"} target="_blank">Read more</a> <OpenExternalLink/>
              </>}
            ><InfoCircle/></TooltipButton>
          </span>
          <BlurUpdatedInput
            value={funds}
            stringify={v => formatUnits(v, decimals)}
            parse={v => parseUnits(v, decimals)}
            onChange={setFunds}/>
          <span className="fake-field"> {currency}</span>
        </label>
        <label>
          <span>
            Minimal authorization 
            <TooltipButton
              tooltip={<>
                The minimal amount of spendable {currency} required by the Aapp before allowing you to log in. <br/>
                <a href={documentationLink + "/PAYMENT.html#minimal-required-authorization"} target="_blank">Read more</a> <OpenExternalLink/>
              </>}
            ><InfoCircle/></TooltipButton>
          </span>
          <span className="fake-field">{formatUnits(funds <= 0n ? 0n : minDeposit, decimals)}</span>
          <span className="fake-field"> {currency}</span>
        </label>
        <label>
          <span>Total required authorization</span>
          <span className="fake-field">{existingDeposit === undefined ? 'Loading...' : formatUnits(funds <= 0n ? 0n : funds + minDeposit, decimals)}</span>
          <span className="fake-field">{currency}</span>
        </label>
        <label>
          <span>Existing authorization</span>
          <span className="fake-field">{existingDeposit === undefined ? 'Loading...' : formatUnits(existingDeposit, decimals)}</span>
          <span className="fake-field">{currency}</span>
        </label>
        <div className="button-card">
          <button onClick={() => topUpDeposit()} disabled={balance === undefined || balance < minDeposit}>
            {
              existingDeposit === undefined ? <>Loading...</> :
              depositInProgress ? <>Processing...</> :
              existingDeposit <= minDeposit ? <>Authorize! ({formatUnits(existingDeposit - minDeposit - funds, decimals)} {currency})</> :
              funds > existingDeposit - minDeposit ? <>Top-up authorization ({formatUnits(existingDeposit - minDeposit - funds, decimals)} {currency})</> :
              funds <= 0n ? <>Remove authorization (+{formatUnits(existingDeposit, decimals)} {currency})</> :
              <>Reduce authorization (+{formatUnits(existingDeposit - minDeposit - funds, decimals)} {currency})</>
            }
          </button>
          {balance === undefined || balance < minDeposit + funds ?
          <TooltipButton
            tooltip={<>
              {balance === undefined || balance < minDeposit ? 
                <> Using the Aapp requires having a minimum of {formatUnits(minDeposit, decimals)} {currency} in your wallet. </> :
                <> You only have {formatUnits(balance, decimals)} {currency}. </>
              } <br/>
              Even if you increased the Aapp's spending cap, the effective authorization is limited by the {currency} backing the allowance. <br/>
              <a href={documentationLink + "/PAYMENT.html#authorized-funds"} target="_blank">Read more</a> <OpenExternalLink/>
            </>}
            onClick={publicClient?.chain?.id == 31337 ? () => walletClient?.data && debugMintFunds(publicClient, walletClient.data, funds + minDeposit) : undefined}
          ><Error/></TooltipButton> : <></>}
          {existingDeposit !== undefined && existingDeposit > minDeposit && funds <= 0n ?
          <TooltipButton
            tooltip={<>
              Removing the authorization will prevent you from logging into the Aapp until it's re-authorized it. <br/>
              <a href={documentationLink + "/PAYMENT.html#maximum-overdraft"} target="_blank">Read more</a> <OpenExternalLink/>
            </>}
          ><Error/></TooltipButton> : <></>}
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
      <p className="deployment-info">
        S3 aApp address: {aappAddress}<br/>
        Payment contract address: {paymentAddress}
      </p>
    </>
  )
}

export default App
