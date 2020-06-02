using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.VirtualNodes;

namespace Apocryph.Core.Consensus
{
    public class Snowball<T> where T : class, IEquatable<T>
    {
        private readonly Node _node;
        private readonly int _k;
        private readonly double _alpha;
        private readonly int _beta;

        private readonly Func<Query<T>[], CancellationToken, IAsyncEnumerable<Query<T>>> _send;

        private readonly Func<Query<T>, T?, T?, Query<T>> _respond;
        private readonly T? _opinion;

        private T? _value;
        private readonly TaskCompletionSource<T> _initialValueTask;

        public Snowball(Node node,
            int k, double alpha, int beta,
            Func<Query<T>[], CancellationToken, IAsyncEnumerable<Query<T>>> send,
            Func<Query<T>, T?, T?, Query<T>> respond,
            T? opinion = null)
        {
            _node = node;
            _k = k;
            _alpha = alpha;
            _beta = beta;
            _send = send;
            _respond = respond;
            _opinion = opinion;

            _initialValueTask = new TaskCompletionSource<T>();
        }

        public Query<T> Query(Query<T> message)
        {
            if (_opinion is null)
            {
                _initialValueTask.TrySetResult(message.Value);
            }

            return _respond(message, _value, _opinion);
        }

        public async Task<T> Run(Node[] receivers,
            CancellationToken cancellationToken)
        {
            _value = _opinion ?? await _initialValueTask.Task;
            var lastValue = _value;
            var d = new Dictionary<T, int> { [_value] = 0 };
            var count = 0;
            while (!cancellationToken.IsCancellationRequested)
            {
                var subset = receivers.OrderBy(_ => RandomNumberGenerator.GetInt32(receivers.Length)).Take(_k);
                var messages = subset.Select(receiver => new Query<T>(_value, _node, receiver, QueryVerb.Request)).ToArray();
                var result = _send(messages, cancellationToken);
                var answers = await result.ToArrayAsync(cancellationToken);
                var answersValues =
                    from message in answers
                    group message by message.Value
                    into valueGroup
                    where valueGroup.Count() > _alpha * _k
                    select valueGroup.Key;
                foreach (var answerValue in answersValues)
                {
                    d[answerValue] = d.TryGetValue(answerValue, out var answerCount) ? answerCount + 1 : 1;
                    if (d[answerValue] > d[_value])
                    {
                        _value = answerValue;
                    }

                    if (answerValue != lastValue)
                    {
                        lastValue = answerValue;
                        count = 0;
                    }
                    else if (++count > _beta)
                    {
                        return _value;
                    }
                }
            }

            throw new OperationCanceledException();
        }
    }
}