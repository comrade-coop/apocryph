// Copyright (c) Comrade Coop. All rights reserved.

using System.Numerics;

namespace Apocryph.FunctionApp.Wetonomy.Messages
{
	public class TransferMessage
	{
		public BigInteger Amount { get; set; }

		public string From { get;set; }

		public string To { get;set; }
	}
}