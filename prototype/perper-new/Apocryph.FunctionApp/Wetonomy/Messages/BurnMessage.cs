// Copyright (c) Comrade Coop. All rights reserved.

using System.Numerics;
using Apocryph.FunctionApp.Model.Message;

namespace Apocryph.FunctionApp.Wetonomy.Messages
{
	public class BurnMessage : IMessage
	{
		public BigInteger Amount { get; set; }

		public string From { get; set;}
	}
}