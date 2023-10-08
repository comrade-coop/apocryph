from ipfs/kubo
RUN ipfs init
RUN ipfs config --json Experimental.Libp2pStreamMounting true
COPY . .
RUN  chmod +x start-client.sh
ENTRYPOINT ["./start-client.sh"]
