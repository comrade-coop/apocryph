namespace Apocryph.Runtime.FunctionApp.Consensus.Core
{
    public class Message<T>
    {
        public T Value { get; }
        public MessageType Type { get; }

        public Message(T value, MessageType type)
        {
            Value = value;
            Type = type;
        }
    }

    public enum MessageType
    {
        Valid,
        Invalid
    }
}