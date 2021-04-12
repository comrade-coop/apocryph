using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Threading;
using Apocryph.Ipfs.Fake;
using Apocryph.Ipfs.Impl;
using Ipfs.Http;
using Xunit;

namespace Apocryph.Ipfs.Test
{
    public class IpfsFixture : IDisposable
    {
        #region PublicApi

        public static string[] HashResolverImplementations =>
            hasIpfs.Value ?
                new string[] { "FakeHashResolver", "IpfsHashResolver" } :
                new string[] { "FakeHashResolver" };
        public static string[] PeerConnectorImplementations =>
            hasIpfs.Value ?
                new string[] { "FakePeerConnector", "IpfsPeerConnector" } :
                new string[] { "FakePeerConnector" };

        public IHashResolver GetHashResolver(string implementation, int instanceId)
        {
            if (implementation == "FakeHashResolver")
                return fakeHashResolver;
            else if (implementation == "IpfsHashResolver")
                return new IpfsHashResolver(GetIpfsClient(instanceId));
            else
                throw new Exception($"Unknown HashResolver: {implementation}");
        }

        public IPeerConnector GetPeerConnector(string implementation, int instanceId)
        {
            if (implementation == "FakePeerConnector")
                return fakePeerConnectorProvider.GetConnector(new Peer(Hash.From(instanceId).Bytes));
            else if (implementation == "IpfsPeerConnector")
                return new IpfsPeerConnector(GetIpfsClient(instanceId));
            else
                throw new Exception($"Unknown PeerConnector: {implementation}");
        }

        #endregion PublicApi

        private FakeHashResolver fakeHashResolver = new FakeHashResolver();
        private FakePeerConnectorProvider fakePeerConnectorProvider = new FakePeerConnectorProvider();

        #region Ipfs

        private Dictionary<int, IpfsClient> _ipfsClients = new Dictionary<int, IpfsClient>();

        private int _nextIpfsApiPort = 5001;
        private int _nextIpfsSwarmPort = 10001;
        private string? _lastPeerAddr = null;
        private List<Process> processesToCleanup = new List<Process>();

        private static Process RunIpfsCommand(string ipfsDirectory, params string[] arguments)
        {
            var info = new ProcessStartInfo("ipfs");

            foreach (var arg in arguments)
                info.ArgumentList.Add(arg);

            info.EnvironmentVariables["IPFS_PATH"] = ipfsDirectory;
            info.UseShellExecute = false;
            info.RedirectStandardOutput = true;
            info.RedirectStandardError = true;

            return Process.Start(info);
        }

        public static Lazy<bool> hasIpfs = new Lazy<bool>(() =>
        {
#if SLOWTESTS
            try
            {
                var testProcess = RunIpfsCommand("", "help");
                testProcess.WaitForExit();
                return testProcess.ExitCode == 0;
            }
            catch(Exception) { }
#endif
            return false;
        });

        private IpfsClient GetIpfsClient(int instanceId)
        {
            if (!_ipfsClients.ContainsKey(instanceId))
            {
                _ipfsClients[instanceId] = InstanceIpfsClient();
            }
            return _ipfsClients[instanceId];
        }

        private IpfsClient InstanceIpfsClient()
        {
            var ipfsApiPort = Interlocked.Increment(ref _nextIpfsApiPort);
            var ipfsSwarmPort = Interlocked.Increment(ref _nextIpfsSwarmPort);

            var ipfsDirectory = Path.Combine(Path.GetTempPath(), $"test-ipfs-instance-{ipfsApiPort}");
            if (Directory.Exists(ipfsDirectory))
            {
                // Console.WriteLine("Clearing old test instance path: {0}", ipfsDirectory);
                Directory.Delete(ipfsDirectory, true);
            }

            RunIpfsCommand(ipfsDirectory, "init").WaitForExit();
            RunIpfsCommand(ipfsDirectory, "config", "profile", "apply", "test").WaitForExit();
            RunIpfsCommand(ipfsDirectory, "config", "--json", "Experimental.Libp2pStreamMounting", "true").WaitForExit();
            RunIpfsCommand(ipfsDirectory, "config", "Addresses.API", $"/ip4/127.0.0.1/tcp/{ipfsApiPort}").WaitForExit();
            RunIpfsCommand(ipfsDirectory, "config", "--json", "Addresses.Swarm", $"[\"/ip4/127.0.0.1/tcp/{ipfsSwarmPort}\"]").WaitForExit();
            if (_lastPeerAddr != null)
            {
                RunIpfsCommand(ipfsDirectory, "config", "--json", "Bootstrap", $"[\"{_lastPeerAddr}\"]").WaitForExit();
            }

            var peerIdProcess = RunIpfsCommand(ipfsDirectory, "config", "Identity.PeerID");
            var peerId = peerIdProcess.StandardOutput.ReadToEnd().Trim();
            _lastPeerAddr = $"/ip4/127.0.0.1/tcp/{ipfsSwarmPort}/p2p/{peerId}";
            peerIdProcess.WaitForExit();

            var ipfsProcess = RunIpfsCommand(ipfsDirectory, "daemon", "--enable-pubsub-experiment");
            processesToCleanup.Add(ipfsProcess);
            while (!ipfsProcess.HasExited)
            {
                var line = ipfsProcess.StandardOutput.ReadLine();
                // Console.WriteLine("{0} | {1}", ipfsDirectory, line);
                if (line != null && line.Contains("Daemon is ready")) break;
            }

            return new IpfsClient($"http://localhost:{ipfsApiPort}");
        }

        void IDisposable.Dispose()
        {
            foreach (var process in processesToCleanup)
            {
                process.Kill();
                process.WaitForExit();
            }
        }

        #endregion Ipfs
    }

    [CollectionDefinition("Ipfs Collection")]
    public class IpfsCollection : ICollectionFixture<IpfsFixture>
    {
    }

}