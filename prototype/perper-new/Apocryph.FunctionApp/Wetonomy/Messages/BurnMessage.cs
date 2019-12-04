// Copyright (c) Comrade Coop. All rights reserved.

using System.Numerics;

namespace Apocryph.FunctionApp.Wetonomy.Messages
{
	public class BurnMessage
	{
		public BigInteger Amount { get; set; }

		public string From { get; set;}
	}
}