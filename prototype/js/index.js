const logMessages = process.argv.indexOf('-q') == -1
const reorderMode = 
  process.argv.indexOf('--mix=all') != -1 ? 'locally unordered' :
  process.argv.indexOf('--mix=local') != -1 ? 'locally ordered' :
  process.argv.indexOf('--mix=global') != -1 ? 'globally unordered' :
  process.argv.indexOf('--mix=none') != -1 ? 'globally ordered' :
  'globally ordered'

const minReorderTime = 30
const maxReorderTime = 80
const validatorCount = 10

const {Stream, HashRegistry, Keypair} = require('./utils')

function debug(...args) {
  if (logMessages) {
    console.log(...args)
  }
}

function executorService(inputStream, outputStream) {
  let state = undefined
  inputStream.foreach(function (item) {
    if (item.type == 'state') {
      state = item.value
    } else if (item.type == 'message') {
      executeMessage(state, item.value, outputStream)
    }
  })
  function executeMessage(state, message, outputStream) {
    if (message.action == 'increment') {
      state[message.target] = 1 + (state[message.target] || 0)
    } else if (message.action == 'jump') {
      outputStream.send({
        type: 'message',
        value: {
          action: message.n > 1 ? 'jump' : 'increment',
          n: message.n - 1,
          target: state[message.target] || 0,
          nonce: message.nonce + 1,
        }
      })
    } else if (message.action == 'jump2') {
      outputStream.send({
        type: 'message',
        value: {
          action: message.n > 1 ? 'jump' : 'increment',
          n: message.n - 1,
          target: state[message.target] || 0,
          nonce: message.nonce + 1,
        }
      })
      outputStream.send({
        type: 'message',
        value: {
          action: message.m > 1 ? 'jump' : 'increment',
          n: message.m - 1,
          target: (state[message.target] || 0) + message.offset,
          nonce: message.nonce + 10,
        }
      })
    }
    outputStream.send({type: 'state', value: JSON.parse(JSON.stringify(state))})
  }
}

class VotingRound {
  constructor(validators) {
    this.votes = new Map()
    this.validators = validators
    this.stakePresent = 0
    let totalStake = Object.values(this.validators).reduce((acc, x) => acc + (x.stake || 1), 0)
    this.stakeNeeded = totalStake * 2 / 3
  }
  
  hasEnoughVotes() {
    return this.stakePresent >= this.stakeNeeded
  }
  
  addVote(from, data) {
    if (this.votes.has(from)) {
      return
    }
    
    this.votes.set(from, data)
    this.stakePresent += (this.validators[from].stake || 1)
  }
  
  collectBestExec() {
    let maxHash = undefined
    let hashToStake = new Map()
    for (let [validator, voteHash] of this.votes.entries()) {
      let vote = HashRegistry.get(voteHash)
      let validatedHash = vote.data.exec
      hashToStake.set(validatedHash, (hashToStake.get(validatedHash) || 0) + (this.validators[validator].stake || 1))
      if (maxHash == undefined) {
        maxHash = validatedHash
      } else {
        let iteratedMaxHash = maxHash
        let iteratedWantedHash = validatedHash
        while (iteratedMaxHash != validatedHash && iteratedMaxHash && iteratedWantedHash != maxHash) {
          iteratedMaxHash = HashRegistry.get(iteratedMaxHash).data.previous
          iteratedWantedHash = iteratedWantedHash && HashRegistry.get(iteratedWantedHash).data.previous
        }
        if (iteratedMaxHash == validatedHash) {
          maxHash = validatedHash
        }
      }
    }
    
    let resultHash = maxHash
    let collectedStake = (hashToStake.get(resultHash) || 0)
    
    while (collectedStake < this.neededStake) {
      resultHash = HashRegistry.get(resultHash).data.previous
      collectedStake += (hashToStake.get(resultHash) || 0)
    }
    
    return resultHash
  }
}

function coreService(keypair, validators, inputStream, outputStream) {
  function broadcastMessage(value) {
    let hash = HashRegistry.add({signature: keypair.sign(value), data: value})
    outputStream.send(hash)
    return hash
  }
  
  let currentRound = 0
  let votingRounds = {}
  getVoteRound(0).current = true
  
  function getVotingRound(id, type) {
    if (!votingRounds[id + type]) {
      votingRounds[id + type] = new VotingRound(validators)
    }
    return votingRounds[id + type]
  }
  
  function getVoteRound(id, type) {
    return getVotingRound(id, 'vote')
  }
  
  function tryCompleteVoteRound(id) {
    let round = getVoteRound(id)
    round.current = true
    if (round.hasEnoughVotes()) {
      round.current = false
      // let votes = Array.from(round.votes.values())
      let bestExec = round.collectBestExec()
      broadcastMessage({type: 'finalizeCommit', exec: bestExec, /*votes: votes,*/ proposer: currentProposer, round: id})
      tryCompleteCommitRound(id)
    }
  }
  
  function getCommitRound(id, type) {
    return getVotingRound(id, 'commit')
  }
  
  function tryCompleteCommitRound(id) {
    let round = getCommitRound(id)
    round.current = true
    if (round.hasEnoughVotes()) {
      round.current = false
      currentRound++
      resetTo(round.collectBestExec())
      
      currentProposer = Object.keys(validators)[Object.keys(validators).indexOf(currentProposer) + 1]

      debug('[' + keypair.hash + ']', 'proposer', currentProposer);
      
      let nextRound = getVoteRound(id + 1)
      if (!nextRound.hasEnoughVotes() && currentProposer == keypair.hash) {
        runProposer()
      }
      
      tryCompleteVoteRound(id + 1)
    }
  }
  
  let pendingMessagePool = new Set()
  let messagePool = new Set()
  let currentProposer = Object.keys(validators)[0]
  
  let executorInput = new Stream()
  let executorOutput = new Stream()
  executorService(executorInput, executorOutput)
  let executorStates = new Stream()
  
  executorOutput.foreach(function (item) {
    if (item.type == 'message') {
      pendingMessagePool.add(HashRegistry.add({data: item.value}))
    }
    if (item.type == 'state') {
      executorStates.send(item.value)
    }
  })
  
  let lastExecHash = genesis
  let executorStatesIterator = executorStates.iterator()
  
  async function executeMessage(state, messageHash, stateAlreadySet=false) {
    if (!stateAlreadySet) {
      executorInput.send({type: 'state', value: state.data})
    }
    
    if (messageHash != undefined) {
      pendingMessagePool.clear()
      executorInput.send({type: 'message', value: HashRegistry.get(messageHash).data})
      
      state.data = (await executorStatesIterator.next()).value
      
      delete state.pending[messageHash]
      for (let pendingMessageHash of pendingMessagePool) {
        state.pending[pendingMessageHash] = true
      }
    }
  }
  
  async function runProposer() {
    if (currentProposer != keypair.hash) {
      return
    }
    
    let state = HashRegistry.get(HashRegistry.get(lastExecHash).data.state)
    
    await executeMessage(state, undefined)
    
    async function propose(messageHash) {
      debug('[' + keypair.hash + ']', 'PROP', messageHash)
      
      await executeMessage(state, messageHash, true)
      
      lastExecHash = broadcastMessage({type: 'exec', message: messageHash, state: HashRegistry.add(state), previous: lastExecHash})
    }
    
    while (messagePool.size != 0 || Object.keys(state.pending).length != 0) {
      for (let messageHash of Object.keys(state.pending)) {
        await propose(messageHash)
        if (currentProposer != keypair.hash) {
          return
        }
      }
      
      for (let messageHash of Array.from(messagePool.values())) {
        await propose(messageHash)
        messagePool.delete(messageHash)
        if (currentProposer != keypair.hash) {
          return
        }
      }
      
      await new Promise(function (resolve) {
        setTimeout(resolve, 10)
      })
      
      if (currentProposer != keypair.hash) {
        return
      }
    }
    
    broadcastMessage({type: 'execFin', previous: lastExecHash})
  }
  
  function resetTo(itemHash) {
    debug('[' + keypair.hash + ']', 'COMMIT', lastExecHash, itemHash, HashRegistry.get(HashRegistry.get(itemHash).data.state));
    lastExecHash = itemHash
  }
  
  async function execTo(itemHash) {
    let iteratedItemHash = itemHash
    let iteratedCurrentHash = lastExecHash
    while (iteratedItemHash != lastExecHash && iteratedItemHash && iteratedCurrentHash != itemHash) {
      iteratedItemHash = HashRegistry.get(iteratedItemHash).data.previous
      iteratedCurrentHash = iteratedCurrentHash && HashRegistry.get(iteratedCurrentHash).data.previous
    }
    
    if (iteratedItemHash != lastExecHash) {
      // We are ahead, probably a reordered message
      return
    }
    
    let execsToDo = []
    let iteratedExecHash = itemHash
    while (iteratedExecHash != lastExecHash && iteratedExecHash) {
      execsToDo.push(iteratedExecHash)
      iteratedExecHash = HashRegistry.get(iteratedExecHash).data.previous
    }
    
    if (iteratedExecHash != lastExecHash) {
      if (!getVotingRound(currentRound).voted) {
        debug('[' + keypair.hash + ']', 'Fork?', itemHash, lastExecHash, iteratedItemHash)
        getVotingRound(currentRound).voted = true
        broadcastMessage({type: 'finalizeVote', exec: lastExecHash, proposer: currentProposer, round: currentRound})
      }
      return
    }
    
    let state = HashRegistry.get(HashRegistry.get(lastExecHash).data.state)
    
    await executeMessage(state, undefined)
    
    execsToDo.reverse()
    // debug('[' + keypair.hash + ']', 'sync', lastExecHash, execsToDo)
    
    for (let execHash of execsToDo) {
      let exec = HashRegistry.get(execHash)
      let messageHash = exec.data.message
      let message = HashRegistry.get(messageHash)
      debug('[' + keypair.hash + ']', 'exec', execHash, messageHash)
      
      if (!message.signature /* TODO */ && !state.pending[messageHash]) {
        if (!getVotingRound(currentRound).voted) {
          debug('[' + keypair.hash + ']', 'HALT (message)', state.pending, messageHash)
          getVotingRound(currentRound).voted = true
          broadcastMessage({type: 'finalizeVote', exec: lastExecHash, proposer: currentProposer, round: currentRound})
        }
        return
      }
      
      await executeMessage(state, messageHash, true)
      
      let stateHash = HashRegistry.add(state)
      
      if (exec.data.state != stateHash) {
        if (!getVotingRound(currentRound).voted) {
          debug('[' + keypair.hash + ']', 'HALT (state)', exec.data.state, stateHash)
          getVotingRound(currentRound).voted = true
          broadcastMessage({type: 'finalizeVote', exec: lastExecHash, proposer: currentProposer, round: currentRound}) //  state: newStateHash
        }
        return
      }
      
      messagePool.delete(messageHash)
      lastExecHash = itemHash
    }
  }

  inputStream.foreach(async function (itemHash) {
    let item = HashRegistry.get(itemHash)
    
    let validSignature = item.signature && validators[item.signature.signer] && validators[item.signature.signer].validate(item.data, item.signature)
    
    let signer = validSignature && item.signature && item.signature.signer
    
    if (item.data.type == 'transaction' && item.signature) { // TODO
      messagePool.add(itemHash)
    }
    
    if (item.data.type == 'exec' && signer == currentProposer && signer != keypair.hash && !getVotingRound(currentRound).voted) {
      await execTo(itemHash)
    }
    
    if (item.data.type == 'execFin' && signer == currentProposer && !getVotingRound(currentRound).voted) {
      await execTo(item.data.previous)
      
      broadcastMessage({type: 'finalizeVote', exec: item.data.previous, proposer: currentProposer, round: currentRound})
    }
    
    if (item.data.type == 'finalizeVote' && item.data.proposer == currentProposer && validSignature) { // NOTE: Relies on the fact that we receive messages from ourselves
      getVoteRound(item.data.round).addVote(signer, itemHash)
      if (getVoteRound(item.data.round).current) {
        tryCompleteVoteRound(item.data.round)
      }
    }
    
    if (item.data.type == 'finalizeCommit' && item.data.proposer == currentProposer && validSignature) { // NOTE: Relies on the fact that we receive messages from ourselves
      getCommitRound(item.data.round).addVote(signer, itemHash)
      if (getCommitRound(item.data.round).current) {
        tryCompleteCommitRound(item.data.round)
      }
    }
  })
  
  if (currentProposer == keypair.hash) {
    runProposer()
  }
}

let genesis = HashRegistry.add({
  data: {
    state: HashRegistry.add({
      data: {},
      pending: {},
    }),
  }
})

{
  let validators = {}
  let outputs = []
  let inputs = []
  
  for (let i = 0; i < validatorCount; i++) {
    let input = new Stream()
    let output = new Stream()
    let keypair = new Keypair()
    
    validators[keypair.hash] = keypair.public
    
    coreService(keypair, validators, input, output)
    
    inputs.push(input)
    outputs.push(output)
  }
  
  let broadcastStream = new Stream()
  let transactionStream = new Stream()
  let delayOptions = {
    reorder: reorderMode.indexOf('unorder') != -1,
    minTime: reorderMode.indexOf('unorder') != -1 ? minReorderTime / 40 : minReorderTime,
    maxTime: reorderMode.indexOf('unorder') != -1 ? maxReorderTime / 40 : maxReorderTime,
  }
  if (reorderMode.indexOf('local') != -1) {
    for (let i = 0; i < validatorCount; i++) {
      for (let j = 0; j < validatorCount; j++) {
        if (i == j) {
          inputs[i].sendStream(outputs[j])
        } else {
          inputs[i].sendStream(Stream.delay(outputs[j], delayOptions))
        }
      }
      inputs[i].sendStream(Stream.delay(transactionStream, delayOptions))
      
      broadcastStream.sendStream(outputs[i])
    }
    broadcastStream.sendStream(transactionStream)
  } else if (reorderMode.indexOf('global') != -1) {
    for (let i = 0; i < validatorCount; i++) {
      inputs[i].sendStream(broadcastStream)
      broadcastStream.sendStream(Stream.delay(outputs[i], delayOptions))
    }
    broadcastStream.sendStream(Stream.delay(transactionStream, delayOptions))
  }
  
  let reachedBlock = 0
  let endHash = genesis
  let transactionsSent = 0
  let messagesSent = 0
  let bandwidth = 0
  
  debug(Object.keys(validators))
  
  broadcastStream.foreach(function (message) {
    let messageValue = HashRegistry.get(message)
    let stringifiedData = JSON.stringify(messageValue)
    
    messagesSent ++
    bandwidth += stringifiedData.length
    
    debug('!' + messageValue.signature.signer + '!', message, JSON.stringify(messageValue.data))
    if (messageValue.data.type == 'exec') {
      debug(' Message:', messageValue.data.message, JSON.stringify(HashRegistry.get(messageValue.data.message)))
    }
    
    if (messageValue.data.type == 'finalizeCommit') {
      endHash = messageValue.data.exec
    }
    
    if (messageValue.data.type == 'transaction') {
      transactionsSent += 1
    }
    
    if (messageValue.data.type == 'execFin') {
      reachedBlock ++
      for (var i = 0; i < reachedBlock; i++) {
        broadcastMessage({
          type: 'transaction',
          action: 'jump2',
          target: 0,
          offset: 1,
          n: 3,
          m: 2,
          nonce: Math.floor(1000 * Math.random()),
        })
      }
    }
  })
  
  let masterKeypair = new Keypair()
  masterKeypair.stake = 0
  
  validators[masterKeypair.hash] = masterKeypair.public
  
  function broadcastMessage(data) {
    transactionStream.send(HashRegistry.add({signature: masterKeypair.sign(data), data}))
  }
  
  let start = +(new Date())
  let printedAlready = false
  let peers = Object.keys(validators).length - 1
  
  console.log('Running in ' + reorderMode + ' mode');
  
  function printDetails() {
    if (printedAlready) return
    let end = +(new Date())
    let mb = bandwidth / 1024 / 1024
    
    let chainLength = 0
    let iteratedEndHash = endHash
    while (iteratedEndHash) {
      chainLength ++
      iteratedEndHash = HashRegistry.get(iteratedEndHash).data.previous
    }
    
    console.log('Simulation details:');
    console.log('  Validators:', peers);
    console.log('  Blocks:', reachedBlock);
    console.log('  Transactions (submitted):', transactionsSent);
    console.log('  Messages (sub-transactions):', chainLength);
    console.log('  Broadcasts:', messagesSent, '/', '~' + Math.round(mb * 100) / 100, 'MiB');
    console.log('  Minimal transmissions:', messagesSent * peers, '/', '~' + Math.round(mb * peers * 100) / 100, 'MiB');
    console.log('  Ping:', minReorderTime + '-' + maxReorderTime);
    console.log('  Time:', (end - start) / 1000);
  }
  
  process.on('exit', function () {
    if (!printedAlready) {
      printDetails()
      console.log('Ran in ' + reorderMode + ' mode');
      console.log('(try ./index.js --mix={all,global,local,none} {-q,})');
      console.log(reachedBlock == peers ? 'Finished executing all blocks for this prototype' : 'Failed to reach the end of execution, likely due to reordering');
    }
    printedAlready = true
    if (reachedBlock != peers) process.exit(1)
    process.exit(0)
  })
  
  let previousInterrupt = 0
  process.on('SIGINT', function () {
    let time = +(new Date())
    if (time - previousInterrupt < 250) {
      printedAlready = true
      console.log(' (try ./index.js --mix={all,global,local,none} {-q,})');
      process.exit(0)
    }
    previousInterrupt = time
    
    console.log(' Interrupted...');
    
    try {
      printDetails()
    } catch(e) {}
    
    console.log('(Interrupt twice in quick successsion to stop)');
  })
  // process.on('uncaughtException', printDetails)
}
