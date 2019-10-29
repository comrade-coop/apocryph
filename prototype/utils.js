class Stream {
  constructor() {
    this.values = []
    this.callbacks = []
    
    this.public = {
      receive: this.receive.bind(this),
      iterator: this.iterator.bind(this),
      foreach: this.foreach.bind(this),
    }
  }
  
  send(value) {
    this.values.push(value)
    if (this.callbacks[this.values.length - 1]) {
      for (let callback of this.callbacks[this.values.length - 1]) {
        callback(value)
      }
    }
    delete this.callbacks[this.values.length - 1]
  }
  
  sendGenerator(asyncGenerator) {
    let runner = async () => {
      for await (let x of asyncGenerator) {
        this.send(x)
      }
    }
    
    runner()
  }
  
  sendStream(stream, from=0) {
    this.sendGenerator(stream.iterator(from))
  }
  
  receive(i=-1) {
    if (i < 0) {
      i = this.values.length;
    }
    
    if (i < this.values.length) {
      return Promise.resolve(this.values[i])
    } else {
      return new Promise((resolve) => {
        this.callbacks[i] = this.callbacks[i] || []
        this.callbacks[i].push(resolve)
      })
    }
  }
  
  async* iterator(from=0) {
    if (from < 0) {
      from = this.values.length;
    }
    for (let i = from;; i++) {
      yield await this.receive(i)
    }
  }
  
  foreach(callback, from=0) {
    let runner = async () => {
      for await (let x of this.iterator(from)) {
        await callback(x)
      }
    }
    
    runner()
  }
}

Stream.race = function (...streams) {
  var result = new Stream()
  for (let stream of streams) {
    result.sendStream(stream)
  }
  return result
}

Stream.filter = function (inputStream, filter, from=0) {
  let result = new Stream()
  
  async function runner() {
    for await (let x of inputStream.iterator(from)) {
      if (filter(x)) {
        result.send(x)
      }
    }
  }
  
  runner()
  
  return result
}

Stream.delay = function (stream, {reorder = false, minAmount = 0, maxAmount = 100}, from=0) {
  var result = new Stream()
  async function runner() {
    for await (let item of stream.iterator(from)) {
      var time = minAmount + (maxAmount - minAmount) * Math.random()
      
      if (reorder) {
        setTimeout(function () {
          result.send(item)
        }, time)
      } else {
        await new Promise(function (resolve) {
          setTimeout(resolve, time)
        })
        result.send(item)
      }
    }
  }
  
  runner()
  
  return result
}

HashRegistry = {
  _indices: {},
  _items: [],

  add(data) {
    var serialized = JSON.stringify(data)
    if (this._indices[serialized] == undefined) {
      this._items.push(serialized)
      this._indices[serialized] = this._items.length - 1
    }
    return this._indices[serialized].toString(16)
  },

  get(hash) {
    return JSON.parse(this._items[parseInt(hash, 16)])
  },
}

class Keypair {
  constructor() {
    this.uniqueData = Math.random().toString(16).slice(2)
    this.hash = HashRegistry.add(this.uniqueData)
    
    this.public = {
      validate: this.validate.bind(this)
    }
  }
  
  validate(value, signature) {
    return signature.data == HashRegistry.add(value) + '|' + this.uniqueData && signature.signer == this.hash
  }
  
  sign(value) {
    return {data: HashRegistry.add(value) + '|' + this.uniqueData, signer: this.hash}
  }
}

module.exports = {
  Stream,
  HashRegistry,
  Keypair
}