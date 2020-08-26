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
        private readonly int _round;

        private readonly Func<Query<T>, CancellationToken, Task<Query<T>>> _send;
        private readonly Func<Query<T>, T?, Query<T>> _respond;

        private T? _value;
        private readonly Dictionary<T, int> _d;
        private readonly TaskCompletionSource<T> _initialValueTask;

        public Snowball(Node node,
            int k, double alpha, int beta, int round,
            Func<Query<T>, CancellationToken, Task<Query<T>>> send,
            Func<Query<T>, T?, Query<T>> respond,
            T? initialValue = null)
        {
            _node = node;
            _k = k;
            _alpha = alpha;
            _beta = beta;
            _round = round;
            _send = send;
            _respond = respond;
            _value = initialValue;

            _initialValueTask = new TaskCompletionSource<T>();
            if (initialValue != null)
            {
                _initialValueTask.TrySetResult(initialValue);
            }
            _d = new Dictionary<T, int> { };
        }

        public Query<T> Query(Query<T> message)
        {
            var result = _respond(message, _value);

            if (!_d.ContainsKey(result.Value))
            {
                _d[result.Value] = 0;
            }

            if (_value is null)
            {
                _value = result.Value;
                _initialValueTask.TrySetResult(result.Value);
            }
            else
            {
                _value = result.Value;
            }

            return result;
        }

        public async Task<T> Run(Node?[] receivers,
            CancellationToken cancellationToken)
        {
            await _initialValueTask.Task;
            if (_value == null) throw new NullReferenceException();

            var lastValue = _value;
            var count = 0;
            while (!cancellationToken.IsCancellationRequested)
            {
                var subset = receivers.OrderBy(_ => RandomNumberGenerator.GetInt32(receivers.Length)).Take(_k);
                var messages = subset.Where(receiver => receiver != null).Select(receiver => new Query<T>(_value, _node, receiver!, _round, QueryVerb.Request)).ToArray();

                var timeoutTokenSource = CancellationTokenSource.CreateLinkedTokenSource(cancellationToken);
                timeoutTokenSource.CancelAfter(TimeSpan.FromMilliseconds(3000));

                var answers = await Task.WhenAll(messages.Select(message => _send(message, timeoutTokenSource.Token).ContinueWith(t => t.IsCanceled ? null : t.Result)).ToArray());
                var answersValues =
                    from message in answers
                    where message != null
                    group message by message.Value
                    into valueGroup
                    where valueGroup.Count() > _alpha * _k
                    select valueGroup.Key;
                foreach (var answerValue in answersValues)
                {
                    _d[answerValue] = _d.TryGetValue(answerValue, out var answerCount) ? answerCount + 1 : 1;
                    if (_d[answerValue] > _d[_value])
                    {
                        _value = answerValue;
                    }

                    if (!answerValue.Equals(lastValue))
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

        public IDictionary<T, int> GetConfirmedValues()
        {
            return _d;
        }
    }
}