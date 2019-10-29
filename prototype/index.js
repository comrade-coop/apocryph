const mode = 
  process.argv.indexOf('--mix=all') != -1 ? 'unordered local' :
  process.argv.indexOf('--mix=local') != -1 ? 'ordered local' :
  process.argv.indexOf('--mix=global') != -1 ? 'unordered global' :
  process.argv.indexOf('--mix=none') != -1 ? 'ordered global' :
  'ordered global'
const {Stream, HashRegistry, Keypair} = require('./utils')

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
        }
      })
    } else if (message.action == 'jump2') {
      outputStream.send({
        type: 'message',
        value: {
          action: message.n > 1 ? 'jump' : 'increment',
          n: message.n - 1,
          target: state[message.target] || 0,
        }
      })
      outputStream.send({
        type: 'message',
        value: {
          action: message.m > 1 ? 'jump' : 'increment',
          n: message.m - 1,
          target: (state[message.target] || 0) + message.offset,
        }
      })
    }
    outputStream.send({type: 'state', value: JSON.parse(JSON.stringify(state))})
  }
}


function coreService(keypair, validators, inputStream, outputStream) {
  function broadcastMessage(value) {
    let hash = HashRegistry.add({signature: keypair.sign(value), data: value})
    outputStream.send(hash)
    return hash
  }
  
  let finalizationVoteSent = false
  let finalizationState = 'voting' // voting, committing
  let finalizationStake = 0
  let finalizationVotes = new Map()
  
  function addFinalizationVote(from, data) {
    if (finalizationVotes.has(from)) {
      return
    }
    
    let stake = (validators[from].stake || 1)
    let neededStake = Object.values(validators).reduce((x, y) => (y.stake || 1) + x, 0) * 2 / 3
    
    finalizationStake += stake
    finalizationVotes.set(from, data)
      
    if (finalizationStake >= neededStake && finalizationStake - stake < neededStake) {
      return true
    } else {
      return false
    }
  }
  
  function clearFinalizationVotes() {
    let maxHash = undefined
    let hashToStake = new Map()
    for (let [validator, voteHash] of finalizationVotes.entries()) {
      let vote = HashRegistry.get(voteHash)
      let validatedHash = vote.data.exec
      hashToStake.set(validatedHash, (hashToStake.get(validatedHash) || 0) + (validators[validator].stake || 1))
      if (maxHash == undefined) {
        maxHash = validatedHash
      } else {
        let iteratedMaxHash = maxHash
        let iteratedWantedHash = validatedHash
        while (iteratedMaxHash != validatedHash && iteratedMaxHash && iteratedWantedHash && iteratedWantedHash != maxHash) {
          iteratedMaxHash = HashRegistry.get(iteratedMaxHash).data.previous
          iteratedWantedHash = HashRegistry.get(iteratedWantedHash).data.previous
        }
        if (iteratedMaxHash == validatedHash) {
          maxHash = validatedHash
        }
      }
    }
    
    let resultHash = maxHash
    let collectedStake = (hashToStake.get(resultHash) || 0)
    let neededStake = Object.values(validators).map(x => x.stake || 1) * 2 / 3
    while (collectedStake < neededStake) {
      resultHash = HashRegistry.get(resultHash).data.previous
      collectedStake += (hashToStake.get(resultHash) || 0)
    }
    
    finalizationVotes.clear()
    finalizationStake = 0
    
    return resultHash
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
  
  async function runProposer() {
    if (currentProposer != keypair.hash) {
      return
    }
    
    executorInput.send({type: 'state', value: HashRegistry.get(HashRegistry.get(lastExecHash).data.state)})
    
    while (messagePool.size != 0) {
      pendingMessagePool = new Set()
      for (let messageHash of messagePool) {
        console.log('[' + keypair.hash + ']', 'PROP', messageHash)
        messagePool.delete(messageHash) // HACK: might fail
        executorInput.send({type: 'message', value: HashRegistry.get(messageHash).data})
        let newState = (await executorStatesIterator.next()).value
        lastExecHash = broadcastMessage({type: 'exec', message: messageHash, state: HashRegistry.add(newState), previous: lastExecHash})
        if (currentProposer != keypair.hash) {
          return
        }
      }
      for (let messageHash of pendingMessagePool) {
        messagePool.add(messageHash)
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
    console.log('[' + keypair.hash + ']', 'COMMIT', lastExecHash, itemHash, HashRegistry.get(HashRegistry.get(lastExecHash).data.state));
    lastExecHash = itemHash
  }
  
  async function execTo(itemHash) {
    let iteratedItemHash = itemHash
    let iteratedCurrentHash = lastExecHash
    while (iteratedItemHash != lastExecHash && iteratedCurrentHash && iteratedItemHash && iteratedCurrentHash != itemHash) {
      iteratedItemHash = HashRegistry.get(iteratedItemHash).data.previous
      iteratedCurrentHash = HashRegistry.get(iteratedCurrentHash).data.previous
    }
    
    if (iteratedCurrentHash == itemHash) {
      // We are ahead, probably a reordered message
      return
    }
    
    let execsToDo = []
    let iteratedExecHash = itemHash
    while (iteratedExecHash != lastExecHash && iteratedExecHash) {
      execsToDo.push(iteratedExecHash)
      iteratedExecHash = HashRegistry.get(iteratedExecHash).data.previous
    }
    
    
    executorInput.send({type: 'state', value: HashRegistry.get(HashRegistry.get(lastExecHash).data.state)})
    
    execsToDo.reverse()
    console.log('[' + keypair.hash + ']', 'sync', lastExecHash, execsToDo)
    
    for (let execHash of execsToDo) {
      pendingMessagePool.clear()
      let exec = HashRegistry.get(execHash)
      let messageHash = exec.data.message
      let message = HashRegistry.get(messageHash)
      console.log('[' + keypair.hash + ']', 'exec', messageHash)
      
      if (!message.signature /* TODO */ && !messagePool.has(messageHash)) {
        if (!finalizationVoteSent) {
          console.log('[' + keypair.hash + ']', 'HALT', messagePool, messageHash)
          finalizationVoteSent = true
          broadcastMessage({type: 'finalizeVote', exec: execHash, proposer: currentProposer})
        }
        return
      }
      
      executorInput.send({type: 'message', value: message.data})
      let newState = (await executorStatesIterator.next()).value
      let newStateHash = HashRegistry.add(newState)
      
      if (exec.data.state != newStateHash) {
        if (!finalizationVoteSent) {
          finalizationVoteSent = true
          broadcastMessage({type: 'finalizeVote', exec: lastExecHash, proposer: currentProposer}) //  state: newStateHash
        }
        return
      }
      
      messagePool.delete(messageHash)
      for (let message of pendingMessagePool) {
        messagePool.add(message)
      }
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
    
    if (item.data.type == 'exec' && signer == currentProposer && signer != keypair.hash && !finalizationVoteSent) {
      await execTo(itemHash)
    }
    
    if (item.data.type == 'execFin' && signer == currentProposer && !finalizationVoteSent) {
      
      await execTo(item.data.previous)
      broadcastMessage({type: 'finalizeVote', exec: item.data.previous, proposer: currentProposer})
    }
    
    if (item.data.type == 'finalizeVote' && item.data.proposer == currentProposer && validSignature && finalizationState == 'voting') { // NOTE: Relies on the fact that we receive messages from ourselves
      if (addFinalizationVote(signer, itemHash)) {
        let votes = Array.from(finalizationVotes.values())
        let bestExec = clearFinalizationVotes()
        broadcastMessage({type: 'finalizeCommit', exec: bestExec, votes: votes, proposer: currentProposer})
        finalizationState = 'committing'
      }
    }
    
    if (item.data.type == 'finalizeCommit' && item.data.proposer == currentProposer && validSignature && finalizationState == 'voting') {
      for (let voteHash of item.data.votes) {
        let vote = HashRegistry.get(voteHash)
        // Whatever copypasta
        if (vote.data.type == 'finalizeVote' && vote.data.proposer == currentProposer && vote.signature && validators[vote.signature.signer] && validators[vote.signature.signer].validate(vote.data, vote.signature)) {
          if (addFinalizationVote(signer, itemHash)) {
            let votes = Array.from(finalizationVotes.values())
            let bestExec = clearFinalizationVotes()
            broadcastMessage({type: 'finalizeCommit', exec: bestExec, votes: votes, proposer: currentProposer})
            finalizationState = 'committing'
          }
        }
      }
    }
    
    if (item.data.type == 'finalizeCommit' && item.data.proposer == currentProposer && validSignature && finalizationState == 'committing') { // NOTE: Relies on the fact that we receive messages from ourselves
      if (addFinalizationVote(signer, itemHash)) {
        setTimeout(function () {
          let bestExec = clearFinalizationVotes()
          finalizationState = 'voting'
          
          resetTo(bestExec)
          finalizationVoteSent = false
          
          currentProposer = Object.keys(validators)[Object.keys(validators).indexOf(currentProposer) + 1]
          
          if (currentProposer == keypair.hash) {
            runProposer()
          }
        }, 1)
      }
    }
  })
  
  if (currentProposer == keypair.hash) {
    runProposer()
  }
}

let genesis = HashRegistry.add({
  data: {
    state: HashRegistry.add({}),
  }
})

{
  const validatorCount = 10
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
  for (let i = 0; i < validatorCount; i++) {
    if (mode.indexOf('local') != -1) {
      for (let j = 0; j < validatorCount; j++) {
        inputs[i].sendStream(Stream.delay(outputs[j], {reorder: mode.indexOf('order') == -1, maxTime: 30}))
      }
      broadcastStream.sendStream(outputs[i])
    } else if (mode.indexOf('global') != -1) {
      inputs[i].sendStream(broadcastStream)
      broadcastStream.sendStream(Stream.delay(outputs[i], {reorder: mode.indexOf('order') == -1, maxTime: 30}))
    }
  }
  
  let reachedEnd = false
  let messagesSent = 0
  broadcastStream.foreach(function (message) {
    messagesSent ++
    let messageValue = HashRegistry.get(message)
    console.log('!' + messageValue.signature.signer + '!', message, JSON.stringify(messageValue.data))
    if (messageValue.data.type == 'exec') {
      console.log(' Message:', messageValue.data.message, JSON.stringify(HashRegistry.get(messageValue.data.message)))
    }
    if (messageValue.data.type == 'exec' && HashRegistry.get(messageValue.data.message).data.D) {
      broadcastMessage({
        type: 'transaction',
        action: 'jump2',
        target: 0,
        offset: 0,
        n: 3,
        m: 0,
        D: HashRegistry.get(messageValue.data.message).data.D - 1
      })
    }
    if (messageValue.data.type == 'finalizeCommit' && messageValue.data.proposer == Object.keys(validators)[Object.keys(validators).length - 2]) {
      reachedEnd = true
    }
  })
  
  let masterKeypair = new Keypair()
  masterKeypair.stake = 0
  
  validators[masterKeypair.hash] = masterKeypair.public
  
  function broadcastMessage(data) {
    broadcastStream.send(HashRegistry.add({signature: masterKeypair.sign(data), data}))
  }
  
  // broadcastMessage({
  //   type: 'transaction',
  //   action: 'increment',
  //   target: 0
  // })
  // 
  // broadcastMessage({
  //   type: 'transaction',
  //   action: 'jump',
  //   target: 0,
  //   offset: 1,
  //   n: 1,
  //   m: 3,
  //   // D: 100
  // })
  
  let start = +(new Date())
  
  process.on('exit', function () {
    let end = +(new Date())
    console.log('Total messages:', messagesSent);
    console.log('Minimal transmissions:', messagesSent * (Object.keys(validators).length - 1));
    console.log('Simulation time:', (end - start) / 1000);
    console.log('Ran in ' + mode + ' mode');
    console.log('(try ./index.js --mix={all,local,global,none})');
    console.log(reachedEnd ? 'Finished executing all blocks for this prototype' : 'Failed to reach the end of execution, likely due to reordering');
    if (!reachedEnd) process.exit(1)
  })
}
