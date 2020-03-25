using System;
using System.Security.Cryptography;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Model;
using Apocryph.FunctionApp.Service;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class RuntimeWorker
    {
        [FunctionName(nameof(RuntimeWorker))]
        [return: Perper("$return")]
        public static AgentContext<object> Run([PerperWorkerTrigger] object workerContext,
            [Perper("state")] object state,
            [Perper("sender")] AgentCapability sender,
            [Perper("message")] object message)
        {
            var context = new AgentContext<object>(state);
            if (message is InitMessage)
            {
                context.SampleStore("0", "1");
                context.SubscribeUserInput(new ValidatorKey
                {
                    Key = new ECParameters
                    {
                        Curve = ECCurve.NamedCurves.nistP521,
                        Q = new ECPoint
                        {
                            // Valid signed object: {"$type":"Apocryph.FunctionApp.Model.Signed`1[[Apocryph.FunctionApp.Model.Commit, Apocryph.FunctionApp]], Apocryph.FunctionApp","Value":{"For":{"Bytes":"AXESIAqgZa22xsCD2XY3SIz4RWk3M4kbUVw93vb+tik8+Wb2"}},"Signer":{"Key":{"X":"ADBRFw9bcnJGVJsrZmvG90YGv2iighBCuc6EtJzktGpfhqaGx6yPBUHgWiby/JwFsWMzSGUm39UxkARSNj/x0AJW","Y":"AbtoMEZ8kF4jc25HYb4VJHYlwU3Pv9BEzJyQMkngj+/z2WDcaVwwAo4zoUz0ls+wXeMoh3qWoM/d7jr61W3B1DjQ"}},"Signature":{"Bytes":"AY7MGR4IM+Qe6Z4wPXi3Ajna7u56B0iw77SZRrMkqPH5OfWAsT78nrWqpk1CbV63p9a+A0SKRs4/7RktyOrARRkDACBxFrlGbGtRrcH4TEvBPi6UR4L39TmBSelNKEQWA24ttGVwCH8NWYyTLFZmWlebx1owzboaXT0bp2ZVmyNb+Ch+"}}
                            X = Convert.FromBase64String("ADBRFw9bcnJGVJsrZmvG90YGv2iighBCuc6EtJzktGpfhqaGx6yPBUHgWiby/JwFsWMzSGUm39UxkARSNj/x0AJW"),
                            Y = Convert.FromBase64String("AbtoMEZ8kF4jc25HYb4VJHYlwU3Pv9BEzJyQMkngj+/z2WDcaVwwAo4zoUz0ls+wXeMoh3qWoM/d7jr61W3B1DjQ"),
                        }
                    }
                });
                context.AddReminder(TimeSpan.FromSeconds(5), "0");
            }
            else if (sender.AgentId == "Reminder" && message is string key)
            {
                context.SampleRestore(key);
            }
            else if (sender.AgentId == "Sample" && message is Tuple<string, object> data)
            {
                if (data.Item2 is string item)
                {
                    context.AddReminder(TimeSpan.FromSeconds(5), item);
                }
                else
                {
                    context.SampleStore(data.Item1, (int.Parse(data.Item1) + 1).ToString());
                    context.AddReminder(TimeSpan.FromSeconds(5), "0");
                }
            }
            else if (sender.AgentId == "IpfsInput")
            {
                Console.WriteLine(message);
            }
            return context;
        }
    }
}