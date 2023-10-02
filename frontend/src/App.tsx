import { useEffect, useState } from 'react'
import { useAccount, useConnect } from 'wagmi'
import { injected } from 'wagmi/connectors'
import { formatUnits, parseUnits } from 'viem'
import { outdent } from 'outdent'

import { Light as SyntaxHighlighter } from 'react-syntax-highlighter'
import { tomorrowNight as syntaxStyle } from 'react-syntax-highlighter/dist/esm/styles/hljs'

import BlurUpdatedInput from './BlurUpdatedInput'
import ActionPopButton from './ActionPopButton'
import apocryphLogo from '../public/apocryph.svg'
import metamaskLogo from '../public/metamask.svg'

const documentationLink = "https://comrade-coop.github.io/apocryph/"

function App() {
  const account = useAccount()
  const { connect } = useConnect()

  const oneGb = parseUnits('1', 6)
  const [ amountGb, setAmountGbRaw ] = useState<bigint>(100n * oneGb)
  const [ durationMultiplier, setDurationMultiplier ] = useState<number>(60 * 60 * 24 * 365)
  const currency = 'S3T'
  const decimals = 18
  const priceGbSec = parseUnits('0.000004', decimals)
  const [ funds, setFunds ] = useState<bigint>(() => BigInt(durationMultiplier) * amountGb * priceGbSec / oneGb)
  const [existingDeposit, setExistingDeposit] = useState(0n) // parseUnits('1.32', decimals)
  useEffect(() => {
    setInterval(() => {
      setExistingDeposit(x => x > 400n ? x - 400n : x)
    }, 20000)
  }, [])


  const duration: number = Number(funds) / Number(amountGb) / Number(priceGbSec) * Number(oneGb)
  function setDuration(newDuration: number) {
    setFunds(BigInt(newDuration) * amountGb * priceGbSec / oneGb)
  }
  function setAmountGb(newAmountGb: bigint) {
    if (newAmountGb < 1n) newAmountGb = 1n

    setFunds(funds * newAmountGb / amountGb)
    setAmountGbRaw(newAmountGb)
  }

  const [ s3Token, setS3Token ] = useState<{accessKeyId: string, secretKeyId: string}>()
  const [ profitText, setProfitText ] = useState<string>("Ok?")

  function refreshApiTokens() {
    const r = () => Math.random().toString(36).slice(2)
    setS3Token({
      accessKeyId: r().toUpperCase(),
      secretKeyId: 'zuf+' + r() + r() + r()
    })
  }

  const bucketLink = `${account.address?.slice(2)?.toLowerCase()}.s3.apocryph.io`
  const bucketLinkHref = `https://${bucketLink}`
  const consoleLink = `console.${bucketLink}`
  const consoleLinkHref = `https://${consoleLink}`

  const codeExamples = {
    'JavaScript': {
      language: 'javascript',
      code: () => outdent`
        import * as Minio from 'minio'

        const minioClient = new Minio.Client({
          endPoint: '${bucketLinkHref}',
          port: 9000,
          useSSL: true,
          accessKey: '${s3Token ? s3Token.accessKeyId : '...'}',
          secretKey: '${s3Token ? s3Token.secretKeyId : '...'}',
        })
      `
    },
    'JavaScript (AWS)': {
      language: 'javascript',
      code: () => outdent`
        import { S3Client } from "@aws-sdk/client-s3"

        const client = new S3Client({
          endPoint: '${bucketLinkHref}',
          region: "us-east-1", // default for MinIO
          credentials: {
            accessKeyId: "${s3Token ? s3Token.accessKeyId : '...'}",
            secretAccessKey: "${s3Token ? s3Token.secretKeyId : '...'}",
          },
        })
      `
    }
  }
  const [ codeExample, setCodeExample ] = useState<keyof typeof codeExamples>("JavaScript")

  return (
    <>
      <img src={apocryphLogo} alt="Apocryph Logo" />
      <h1>Get your S3-compatible bucket!</h1>
      <section>
      <p className="hero">Hosting your S3-compatible data buckets in the Apocryph S3 network allows for the ultimate privacy peace of mind, through trasparent encryption at-rest and cryptocurrency-enabled payments.<br/><a href={documentationLink}>Read more.</a></p>
      </section>
      <section>
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
      </section>
      <section style={{display: account.isConnected ? '' : 'none'}}>
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
          <span>Total required deposit</span>
          <BlurUpdatedInput
            value={funds}
            stringify={v => formatUnits(v, decimals)}
            parse={v => parseUnits(v, decimals)}
            onChange={setFunds}/>
          <span className="fake-field"> {currency}</span>
        </label>
        { existingDeposit > 0n ?
          <label>
            <span>Existing deposit</span>
            <span className="fake-field">{formatUnits(existingDeposit, decimals)}</span>
            <span className="fake-field">{currency}</span>
          </label>
          :
          <></>
        }
        <div className="button-card">
          <button onClick={() => setExistingDeposit(funds)}>
            {
              existingDeposit <= 0n ? <>Make deposit! ({formatUnits(existingDeposit - funds, decimals)} {currency})</> :
              funds > existingDeposit ? <>Top-up deposit ({formatUnits(existingDeposit - funds, decimals)} {currency})</> :
              <>Withdraw deposit (+{formatUnits(existingDeposit - funds, decimals)} {currency})</>
            }
          </button>
        </div>
      </section>
      <section style={{display: account.isConnected && existingDeposit > 0n ? '' : 'none'}}>
        <h2>Step 3: Access</h2>
        <label>
          <span>Console </span>
          <a className="fake-field" href={consoleLinkHref}>{consoleLink}</a>
          <ActionPopButton onClick={() => navigator.clipboard.writeText(bucketLinkHref)}>Copy</ActionPopButton>
        </label>
        <label>
          <span>S3 endpoint URL </span>
          <a className="fake-field" href={bucketLinkHref}>{bucketLinkHref}</a>
          <ActionPopButton onClick={() => navigator.clipboard.writeText(bucketLinkHref)}>Copy</ActionPopButton>
        </label>
        {s3Token != undefined ?
          <>
          <label>
            <span>S3 Access Key </span>
            <span className="fake-field">{s3Token.accessKeyId}</span>
            <ActionPopButton disabled={s3Token == undefined} onClick={() => s3Token && navigator.clipboard.writeText(s3Token.accessKeyId)}>Copy</ActionPopButton>
          </label>
          <label>
            <span>S3 Secret Key </span>
            <span className="fake-field">{s3Token.secretKeyId}</span>
            <ActionPopButton disabled={s3Token == undefined} onClick={() => s3Token && navigator.clipboard.writeText(s3Token.secretKeyId)}>Copy</ActionPopButton>
          </label>
          </>
          : <>
          </>
        }
        <div className="button-card">
          <button onClick={refreshApiTokens}>
            {
              s3Token == undefined ? <>Get S3 access tokens</> :
              <>Refresh S3 access tokens</>
            }
          </button>
        </div>
      </section>
      <section style={{display: account.isConnected && existingDeposit > 0n && s3Token ? '' : 'none'}} className="two-columns">
        <h2>Step 4: Hack</h2>
        <div className="button-card">
          <ActionPopButton popText='Copied' onClick={() => {
            navigator.clipboard.writeText(codeExamples[codeExample].code())
          }}>Start hacking!</ActionPopButton>
          <select value={codeExample} onChange={e => setCodeExample(e.target.value as keyof typeof codeExamples)}>
            {Object.keys(codeExamples).map(x => <option key={x}>{x}</option>)}
          </select>
        </div>
        <SyntaxHighlighter language={codeExamples[codeExample]?.language} style={syntaxStyle} className="code">
          {codeExamples[codeExample]?.code()}
        </SyntaxHighlighter>
      </section>
      <section style={{display: account.isConnected && existingDeposit > 0n && s3Token ? '' : 'none'}}>
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
      </section>
      <a href={documentationLink} className="read-the-docs">Documentation</a>
    </>
  )
}

export default App
