// Copyright (c) Comrade Coop. All rights reserved.

using System.Numerics;

namespace Apocryph.FunctionApp.Wetonomy.Messages
{
	public class MintMessage
	{
		public BigInteger Amount { get; set; }

		public string To { get; set; }
	}
}