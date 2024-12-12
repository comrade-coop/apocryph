//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// MockToken
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

export const mockTokenAbi = [
  { type: 'constructor', inputs: [], stateMutability: 'nonpayable' },
  {
    type: 'function',
    inputs: [
      { name: 'owner', internalType: 'address', type: 'address' },
      { name: 'spender', internalType: 'address', type: 'address' }
    ],
    name: 'allowance',
    outputs: [{ name: '', internalType: 'uint256', type: 'uint256' }],
    stateMutability: 'view'
  },
  {
    type: 'function',
    inputs: [
      { name: 'spender', internalType: 'address', type: 'address' },
      { name: 'value', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'approve',
    outputs: [{ name: '', internalType: 'bool', type: 'bool' }],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [{ name: 'account', internalType: 'address', type: 'address' }],
    name: 'balanceOf',
    outputs: [{ name: '', internalType: 'uint256', type: 'uint256' }],
    stateMutability: 'view'
  },
  {
    type: 'function',
    inputs: [],
    name: 'decimals',
    outputs: [{ name: '', internalType: 'uint8', type: 'uint8' }],
    stateMutability: 'view'
  },
  {
    type: 'function',
    inputs: [{ name: 'amount', internalType: 'uint256', type: 'uint256' }],
    name: 'mint',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [],
    name: 'name',
    outputs: [{ name: '', internalType: 'string', type: 'string' }],
    stateMutability: 'view'
  },
  {
    type: 'function',
    inputs: [],
    name: 'symbol',
    outputs: [{ name: '', internalType: 'string', type: 'string' }],
    stateMutability: 'view'
  },
  {
    type: 'function',
    inputs: [],
    name: 'totalSupply',
    outputs: [{ name: '', internalType: 'uint256', type: 'uint256' }],
    stateMutability: 'view'
  },
  {
    type: 'function',
    inputs: [
      { name: 'to', internalType: 'address', type: 'address' },
      { name: 'value', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'transfer',
    outputs: [{ name: '', internalType: 'bool', type: 'bool' }],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: 'from', internalType: 'address', type: 'address' },
      { name: 'to', internalType: 'address', type: 'address' },
      { name: 'value', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'transferFrom',
    outputs: [{ name: '', internalType: 'bool', type: 'bool' }],
    stateMutability: 'nonpayable'
  },
  {
    type: 'event',
    anonymous: false,
    inputs: [
      {
        name: 'owner',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'spender',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'value',
        internalType: 'uint256',
        type: 'uint256',
        indexed: false
      }
    ],
    name: 'Approval'
  },
  {
    type: 'event',
    anonymous: false,
    inputs: [
      { name: 'from', internalType: 'address', type: 'address', indexed: true },
      { name: 'to', internalType: 'address', type: 'address', indexed: true },
      {
        name: 'value',
        internalType: 'uint256',
        type: 'uint256',
        indexed: false
      }
    ],
    name: 'Transfer'
  },
  {
    type: 'error',
    inputs: [
      { name: 'spender', internalType: 'address', type: 'address' },
      { name: 'allowance', internalType: 'uint256', type: 'uint256' },
      { name: 'needed', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'ERC20InsufficientAllowance'
  },
  {
    type: 'error',
    inputs: [
      { name: 'sender', internalType: 'address', type: 'address' },
      { name: 'balance', internalType: 'uint256', type: 'uint256' },
      { name: 'needed', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'ERC20InsufficientBalance'
  },
  {
    type: 'error',
    inputs: [{ name: 'approver', internalType: 'address', type: 'address' }],
    name: 'ERC20InvalidApprover'
  },
  {
    type: 'error',
    inputs: [{ name: 'receiver', internalType: 'address', type: 'address' }],
    name: 'ERC20InvalidReceiver'
  },
  {
    type: 'error',
    inputs: [{ name: 'sender', internalType: 'address', type: 'address' }],
    name: 'ERC20InvalidSender'
  },
  {
    type: 'error',
    inputs: [{ name: 'spender', internalType: 'address', type: 'address' }],
    name: 'ERC20InvalidSpender'
  }
] as const

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Payment
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

export const paymentAbi = [
  {
    type: 'constructor',
    inputs: [
      { name: '_token', internalType: 'contract IERC20', type: 'address' }
    ],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: '_authorized', internalType: 'address', type: 'address' },
      { name: 'provider', internalType: 'address', type: 'address' },
      { name: 'podId', internalType: 'bytes32', type: 'bytes32' }
    ],
    name: 'authorize',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: 'publisher', internalType: 'address', type: 'address' },
      { name: 'provider', internalType: 'address', type: 'address' },
      { name: 'podId', internalType: 'bytes32', type: 'bytes32' }
    ],
    name: 'available',
    outputs: [{ name: '', internalType: 'uint256', type: 'uint256' }],
    stateMutability: 'view'
  },
  {
    type: 'function',
    inputs: [
      { name: '', internalType: 'address', type: 'address' },
      { name: '', internalType: 'address', type: 'address' },
      { name: '', internalType: 'bytes32', type: 'bytes32' }
    ],
    name: 'channels',
    outputs: [
      { name: 'investedByPublisher', internalType: 'uint256', type: 'uint256' },
      { name: 'withdrawnByProvider', internalType: 'uint256', type: 'uint256' },
      { name: 'unlockTime', internalType: 'uint256', type: 'uint256' },
      { name: 'unlockedAt', internalType: 'uint256', type: 'uint256' }
    ],
    stateMutability: 'view'
  },
  {
    type: 'function',
    inputs: [
      { name: 'publisher', internalType: 'address', type: 'address' },
      { name: 'provider', internalType: 'address', type: 'address' },
      { name: 'podId', internalType: 'bytes32', type: 'bytes32' }
    ],
    name: 'closeChannel',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: 'provider', internalType: 'address', type: 'address' },
      { name: 'podId', internalType: 'bytes32', type: 'bytes32' },
      { name: 'unlockTime', internalType: 'uint256', type: 'uint256' },
      { name: 'initialAmount', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'createChannel',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: 'publisher', internalType: 'address', type: 'address' },
      { name: 'provider', internalType: 'address', type: 'address' },
      { name: 'podId', internalType: 'bytes32', type: 'bytes32' },
      { name: 'newProvider', internalType: 'address', type: 'address' },
      { name: 'newPodId', internalType: 'bytes32', type: 'bytes32' },
      { name: 'amount', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'createSubChannel',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: 'provider', internalType: 'address', type: 'address' },
      { name: 'podId', internalType: 'bytes32', type: 'bytes32' },
      { name: 'amount', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'deposit',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: 'publisher', internalType: 'address', type: 'address' },
      { name: 'provider', internalType: 'address', type: 'address' },
      { name: 'podId', internalType: 'bytes32', type: 'bytes32' },
      { name: '_address', internalType: 'address', type: 'address' }
    ],
    name: 'isAuthorized',
    outputs: [{ name: '', internalType: 'bool', type: 'bool' }],
    stateMutability: 'view'
  },
  {
    type: 'function',
    inputs: [],
    name: 'token',
    outputs: [{ name: '', internalType: 'contract IERC20', type: 'address' }],
    stateMutability: 'view'
  },
  {
    type: 'function',
    inputs: [
      { name: 'publisher', internalType: 'address', type: 'address' },
      { name: 'provider', internalType: 'address', type: 'address' },
      { name: 'podId', internalType: 'bytes32', type: 'bytes32' }
    ],
    name: 'unlock',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: 'publisher', internalType: 'address', type: 'address' },
      { name: 'podId', internalType: 'bytes32', type: 'bytes32' },
      { name: 'amount', internalType: 'uint256', type: 'uint256' },
      { name: 'transferAddress', internalType: 'address', type: 'address' }
    ],
    name: 'withdraw',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: 'publisher', internalType: 'address', type: 'address' },
      { name: 'provider', internalType: 'address', type: 'address' },
      { name: 'podId', internalType: 'bytes32', type: 'bytes32' }
    ],
    name: 'withdrawUnlocked',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: 'publisher', internalType: 'address', type: 'address' },
      { name: 'podId', internalType: 'bytes32', type: 'bytes32' },
      { name: 'totalWithdrawAmount', internalType: 'uint256', type: 'uint256' },
      { name: 'transferAddress', internalType: 'address', type: 'address' }
    ],
    name: 'withdrawUpTo',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: 'publisher', internalType: 'address', type: 'address' },
      { name: 'provider', internalType: 'address', type: 'address' },
      { name: 'podId', internalType: 'bytes32', type: 'bytes32' }
    ],
    name: 'withdrawn',
    outputs: [{ name: '', internalType: 'uint256', type: 'uint256' }],
    stateMutability: 'view'
  },
  {
    type: 'event',
    anonymous: false,
    inputs: [
      {
        name: 'publisher',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'provider',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      { name: 'podId', internalType: 'bytes32', type: 'bytes32', indexed: true }
    ],
    name: 'ChannelClosed'
  },
  {
    type: 'event',
    anonymous: false,
    inputs: [
      {
        name: 'publisher',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'provider',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      { name: 'podId', internalType: 'bytes32', type: 'bytes32', indexed: true }
    ],
    name: 'ChannelCreated'
  },
  {
    type: 'event',
    anonymous: false,
    inputs: [
      {
        name: 'publisher',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'provider',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'podId',
        internalType: 'bytes32',
        type: 'bytes32',
        indexed: true
      },
      {
        name: 'depositAmount',
        internalType: 'uint256',
        type: 'uint256',
        indexed: false
      }
    ],
    name: 'Deposited'
  },
  {
    type: 'event',
    anonymous: false,
    inputs: [
      {
        name: 'publisher',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'provider',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'podId',
        internalType: 'bytes32',
        type: 'bytes32',
        indexed: true
      },
      {
        name: 'unlockedAt',
        internalType: 'uint256',
        type: 'uint256',
        indexed: false
      }
    ],
    name: 'UnlockTimerStarted'
  },
  {
    type: 'event',
    anonymous: false,
    inputs: [
      {
        name: 'publisher',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'provider',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'podId',
        internalType: 'bytes32',
        type: 'bytes32',
        indexed: true
      },
      {
        name: 'unlockedAmount',
        internalType: 'uint256',
        type: 'uint256',
        indexed: false
      }
    ],
    name: 'Unlocked'
  },
  {
    type: 'event',
    anonymous: false,
    inputs: [
      {
        name: 'publisher',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'provider',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'podId',
        internalType: 'bytes32',
        type: 'bytes32',
        indexed: true
      },
      {
        name: 'withdrawnAmount',
        internalType: 'uint256',
        type: 'uint256',
        indexed: false
      }
    ],
    name: 'Withdrawn'
  },
  {
    type: 'error',
    inputs: [{ name: 'target', internalType: 'address', type: 'address' }],
    name: 'AddressEmptyCode'
  },
  { type: 'error', inputs: [], name: 'AlreadyExists' },
  { type: 'error', inputs: [], name: 'AmountRequired' },
  { type: 'error', inputs: [], name: 'ChannelLocked' },
  { type: 'error', inputs: [], name: 'DoesNotExist' },
  { type: 'error', inputs: [], name: 'FailedCall' },
  {
    type: 'error',
    inputs: [
      { name: 'balance', internalType: 'uint256', type: 'uint256' },
      { name: 'needed', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'InsufficientBalance'
  },
  { type: 'error', inputs: [], name: 'InsufficientFunds' },
  { type: 'error', inputs: [], name: 'NotAuthorized' },
  {
    type: 'error',
    inputs: [{ name: 'token', internalType: 'address', type: 'address' }],
    name: 'SafeERC20FailedOperation'
  }
] as const

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// PaymentV2
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

export const paymentV2Abi = [
  {
    type: 'constructor',
    inputs: [
      { name: '_token', internalType: 'contract IERC20', type: 'address' },
      { name: '_unlockTime', internalType: 'uint256', type: 'uint256' }
    ],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: 'channelId', internalType: 'bytes32', type: 'bytes32' },
      { name: 'recipient', internalType: 'address', type: 'address' },
      { name: 'permissions', internalType: 'uint16', type: 'uint16' },
      { name: 'reservedAmount', internalType: 'uint256', type: 'uint256' },
      { name: 'extraLimit', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'authorize',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: '', internalType: 'bytes32', type: 'bytes32' },
      { name: '', internalType: 'address', type: 'address' }
    ],
    name: 'channelAuthorizations',
    outputs: [
      { name: 'reservation', internalType: 'uint128', type: 'uint128' },
      { name: 'limit', internalType: 'uint128', type: 'uint128' },
      { name: 'unlocksAt', internalType: 'uint240', type: 'uint240' },
      { name: 'permissions', internalType: 'uint16', type: 'uint16' }
    ],
    stateMutability: 'view'
  },
  {
    type: 'function',
    inputs: [{ name: '', internalType: 'bytes32', type: 'bytes32' }],
    name: 'channels',
    outputs: [
      { name: 'available', internalType: 'uint128', type: 'uint128' },
      { name: 'reserved', internalType: 'uint128', type: 'uint128' }
    ],
    stateMutability: 'view'
  },
  {
    type: 'function',
    inputs: [
      {
        name: 'channelDiscriminator',
        internalType: 'bytes32',
        type: 'bytes32'
      },
      { name: 'initialAmount', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'create',
    outputs: [{ name: 'channelId', internalType: 'bytes32', type: 'bytes32' }],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      {
        name: 'channelDiscriminator',
        internalType: 'bytes32',
        type: 'bytes32'
      },
      { name: 'initialAmount', internalType: 'uint256', type: 'uint256' },
      { name: 'recipient', internalType: 'address', type: 'address' },
      { name: 'permissions', internalType: 'uint16', type: 'uint16' },
      { name: 'reservedAmount', internalType: 'uint256', type: 'uint256' },
      { name: 'limit', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'createAndAuthorize',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: 'channelId', internalType: 'bytes32', type: 'bytes32' },
      { name: 'amount', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'deposit',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [
      { name: 'creator', internalType: 'address', type: 'address' },
      { name: 'channelDiscriminator', internalType: 'bytes32', type: 'bytes32' }
    ],
    name: 'getChannelId',
    outputs: [{ name: 'channelId', internalType: 'bytes32', type: 'bytes32' }],
    stateMutability: 'pure'
  },
  {
    type: 'function',
    inputs: [],
    name: 'token',
    outputs: [{ name: '', internalType: 'contract IERC20', type: 'address' }],
    stateMutability: 'view'
  },
  {
    type: 'function',
    inputs: [
      { name: 'channelId', internalType: 'bytes32', type: 'bytes32' },
      { name: 'recipient', internalType: 'address', type: 'address' }
    ],
    name: 'unlock',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'function',
    inputs: [],
    name: 'unlockTime',
    outputs: [{ name: '', internalType: 'uint256', type: 'uint256' }],
    stateMutability: 'view'
  },
  {
    type: 'function',
    inputs: [
      { name: 'channelId', internalType: 'bytes32', type: 'bytes32' },
      { name: 'transferAddress', internalType: 'address', type: 'address' },
      { name: 'transferAmount', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'withdraw',
    outputs: [],
    stateMutability: 'nonpayable'
  },
  {
    type: 'event',
    anonymous: false,
    inputs: [
      {
        name: 'channelId',
        internalType: 'bytes32',
        type: 'bytes32',
        indexed: true
      },
      {
        name: 'recipient',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'permissions',
        internalType: 'uint16',
        type: 'uint16',
        indexed: false
      },
      {
        name: 'reservedAmount',
        internalType: 'uint256',
        type: 'uint256',
        indexed: false
      },
      {
        name: 'limit',
        internalType: 'uint256',
        type: 'uint256',
        indexed: false
      }
    ],
    name: 'Authorize'
  },
  {
    type: 'event',
    anonymous: false,
    inputs: [
      {
        name: 'channelId',
        internalType: 'bytes32',
        type: 'bytes32',
        indexed: true
      },
      {
        name: 'payer',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'amount',
        internalType: 'uint256',
        type: 'uint256',
        indexed: false
      }
    ],
    name: 'Deposit'
  },
  {
    type: 'event',
    anonymous: false,
    inputs: [
      {
        name: 'channelId',
        internalType: 'bytes32',
        type: 'bytes32',
        indexed: true
      },
      {
        name: 'recipient',
        internalType: 'address',
        type: 'address',
        indexed: true
      }
    ],
    name: 'Unlock'
  },
  {
    type: 'event',
    anonymous: false,
    inputs: [
      {
        name: 'channelId',
        internalType: 'bytes32',
        type: 'bytes32',
        indexed: true
      },
      {
        name: 'recipient',
        internalType: 'address',
        type: 'address',
        indexed: true
      },
      {
        name: 'amount',
        internalType: 'uint256',
        type: 'uint256',
        indexed: false
      }
    ],
    name: 'Withdraw'
  },
  {
    type: 'error',
    inputs: [{ name: 'target', internalType: 'address', type: 'address' }],
    name: 'AddressEmptyCode'
  },
  { type: 'error', inputs: [], name: 'AlreadyInitialized' },
  { type: 'error', inputs: [], name: 'AuthorizationLocked' },
  { type: 'error', inputs: [], name: 'FailedCall' },
  {
    type: 'error',
    inputs: [
      { name: 'balance', internalType: 'uint256', type: 'uint256' },
      { name: 'needed', internalType: 'uint256', type: 'uint256' }
    ],
    name: 'InsufficientBalance'
  },
  { type: 'error', inputs: [], name: 'InsufficientFunds' },
  { type: 'error', inputs: [], name: 'NotAuthorized' },
  {
    type: 'error',
    inputs: [{ name: 'token', internalType: 'address', type: 'address' }],
    name: 'SafeERC20FailedOperation'
  }
] as const
