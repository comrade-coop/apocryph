using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;

namespace Apocryph.Runtime.FunctionApp.Consensus
{
    public class SnowballLoop<T> where T : class, IEquatable<T>
    {
        private readonly Node _node;
        private readonly int _k;
        private readonly double _alpha;
        private readonly int _beta;

        private readonly Func<IEnumerable<QueryMessage<T>>, CancellationToken, Task<IAsyncEnumerable<AnswerMessage<T>>>> _query;

        private readonly Func<QueryMessage<T>, T?, T?, AnswerMessage<T>> _answer;
        private readonly T? _opinion;

        private T? _value;
        private readonly TaskCompletionSource<T> _initialValueTask;

        public SnowballLoop(Node node,
            Func<IEnumerable<QueryMessage<T>>, CancellationToken, Task<IAsyncEnumerable<AnswerMessage<T>>>> query,
            Func<QueryMessage<T>, T?, T?, AnswerMessage<T>> answer,
            T? opinion = null):this(node, 100, 0.6, 10, query, answer, opinion)
        {
        }

        public SnowballLoop(Node node,
            int k, double alpha, int beta,
            Func<IEnumerable<QueryMessage<T>>, CancellationToken, Task<IAsyncEnumerable<AnswerMessage<T>>>> query,
            Func<QueryMessage<T>, T?, T?, AnswerMessage<T>> answer,
            T? opinion = null)
        {
            _node = node;
            _k = k;
            _alpha = alpha;
            _beta = beta;
            _query = query;
            _answer = answer;
            _opinion = opinion;

            _initialValueTask = new TaskCompletionSource<T>();
        }

        public AnswerMessage<T> Query(QueryMessage<T> message)
        {
            if (_opinion is null)
            {
                _initialValueTask.TrySetResult(message.Value);
            }

            return _answer(message, _value, _opinion);
        }

        public async Task<T> Run(Node[] receivers,
            CancellationToken cancellationToken)
        {
            _value = _opinion ?? await _initialValueTask.Task;
            var lastValue = _value;
            var d = new Dictionary<T, int> {[_value] = 0};
            var count = 0;
            while (!cancellationToken.IsCancellationRequested)
            {
                var subset = receivers.OrderBy(_ => RandomNumberGenerator.GetInt32(receivers.Length)).Take(_k);
                var messages = subset.Select(receiver => new QueryMessage<T>(_value, _node, receiver));
                var result = await _query(messages, cancellationToken);
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