using Perper.WebJobs.Extensions.Fake;
using Xunit;

namespace Apocryph.ServiceRegistry.Test
{
    public class HandlerTests
    {
        [Fact]
        public async void Invoke_WhenUsed_CallsMethodWithParameters()
        {
            var called = 0;
            object? lastParameters = null;

            var testAgent = new FakeAgent();
            testAgent.RegisterFunction("testMethod", (object? parameters) =>
            {
                called++;
                lastParameters = parameters;
            });

            var handler = new Handler(testAgent, "testMethod");

            await handler.InvokeAsync("Sample string");
            Assert.Equal(1, called);
            Assert.Equal("Sample string", lastParameters);

            await handler.InvokeAsync(null);
            Assert.Equal(2, called);
            Assert.Null(lastParameters);
        }
    }
}