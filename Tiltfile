# -*- mode: Python -*-
# SPDX-License-Identifier: GPL-3.0

# For more on Extensions, see: https://docs.tilt.dev/extensions.html
load("ext://restart_process", "docker_build_with_restart")
load("ext://namespace", "namespace_create")
load('ext://dotenv', 'dotenv')
dotenv()

local("which jq forge cast docker >/dev/null", echo_off=True)  # Check dependencies
root_dir = os.getcwd()


# Via https://github.com/comrade-coop/apocryph/blob/9cbf7e87200216a41619f7f40cda3bca56fe03e8/deploy/Tiltfile
def builder_resource(
    name,
    image,
    dir=".",
    write_dir="",
    builder_dir="/app",
    volumes=[],
    volumes_conf={},
    entrypoint=["sleep", "infinity"],
    env={},
    *args,
    **kwargs,
):
    services = {
        name: {
            "image": image,
            "container_name": name,
            "entrypoint": entrypoint,
            "working_dir": builder_dir,
            "network_mode": "host",  # TODO: Figure out a way to use docker's routing! (host likely only used for ipfs)
            "volumes": (
                [
                    "%s:%s" % (dir, builder_dir),
                ]
                if write_dir == ""
                else [
                    "%s:%s:ro" % (dir, builder_dir),
                    "%s/%s:%s/%s:rw" % (dir, write_dir, builder_dir, write_dir),
                ]
            )
            + volumes,
            "environment": env,
        },
    }
    docker_compose(encode_yaml({"services": services, "volumes": volumes_conf}))
    dc_resource(name, *args, **kwargs)


def local_resource_in_builder(name, cmd, builder, resource_deps=[], serve_cmd="", *args, **kwargs):
    local_resource(
        name,
        cmdline_in_builder(cmd, builder),
        resource_deps=resource_deps + [builder],
        serve_cmd=cmdline_in_builder(serve_cmd, builder),
        *args,
        **kwargs,
    )


def cmdline_in_builder(cmd, builder, *, interactive=False):
    if cmd == "":
        return ""
    flags = []
    if interactive:
        flags += ["-i"]
    if type(cmd) == "list":
        return ["docker", "exec"] + flags + [builder] + cmd
    else:
        return "docker exec %s %s %s" % (" ".join(flags), builder, cmd)


def s3_aapp_build_with_builder():
    docker_build(
        "comradecoop/s3-aapp/builder-go",
        root_dir,
        dockerfile=root_dir + "/Dockerfile",
        target="builder-go",
        only=[root_dir + "/Dockerfile"],
    )

    builder_resource(
        "apocryph-s3-go-builder",
        "comradecoop/s3-aapp/go-builder",
        dir=root_dir,
        write_dir="bin",
        volumes=["go-cache:/root/.cache/go-build", "go-mod-cache:/go/pkg/mod"],
        volumes_conf={"go-cache": {}, "go-mod-cache": {}},
    )

    local_resource_in_builder(
        "apocryph-s3-backend-go-compile",
        [
            "sh",
            "-c",
            'go build -v -buildvcs=false -ldflags="-s -w" -o /tmp ./backend/minio-manager ./backend/dns-build && cp /tmp/minio-manager bin/apocryph-s3-backend && cp /tmp/dns-build bin/apocryph-s3-dns',
        ],
        "apocryph-s3-go-builder",
        deps=[root_dir + "/backend"],
        allow_parallel=True,
    )
        
    docker_build(
        "comradecoop/s3-aapp/js-builder",
        root_dir,
        dockerfile=root_dir + "/Dockerfile",
        target="js-build-base",
        only=[root_dir + "/Dockerfile"],
    )

    builder_resource(
        "apocryph-s3-js-builder",
        "comradecoop/s3-aapp/js-builder",
        dir=root_dir,
        volumes=[
            "pnpm-cache:/pnpm/store",
            "../apocryph/pkg/abi-ts:/apocryph/pkg/abi-ts:ro",
        ],
        # TODO: ports=[5173],
        volumes_conf={"pnpm-cache": {}},
        env=get_build_args(local_eth) + {
            'COREPACK_INTEGRITY_KEYS': '0'
        }
    )

    local_resource_in_builder(
        "apocryph-s3-js-serve",
        ["sh", "-c", "cd frontend/ && pnpm install --frozen-lockfile && pnpm run build"],
        "apocryph-s3-js-builder",
        deps=[],
        allow_parallel=True,
    )

    docker_build(
        "comradecoop/s3-aapp",
        root_dir,
        dockerfile=root_dir + "/Dockerfile",
        target="run-all-singlenode",
        build_args=get_build_args(local_eth)
    )

def get_build_args(local_eth=False):
    return {
        'VITE_CHAIN_CONFIG': (str(encode_json({
            'id': 31337,
            'name': 'Localhost',
            'nativeCurrency': {'name': 'lETH', 'symbol': 'lETH', 'decimals': 18},
            'rpcUrls': {'default': {'http': ['http://anvil.local:8545'], 'websocket': ['ws://anvil.local:8545']}},
            'testnet': True,
        })) if local_eth else os.getenv('VITE_CHAIN_CONFIG')),
        'VITE_TOKEN': '0x5FbDB2315678afecb367f032d93F642f64180aa3' if local_eth else os.getenv('BACKEND_ETH_TOKEN'),
        'VITE_STORAGE_SYSTEM': (os.getenv('BACKEND_ETH_PRIVATE_KEY') and str(local('cast wallet a %s' % os.getenv('BACKEND_ETH_PRIVATE_KEY'), echo_off=True))),
        'VITE_GLOBAL_HOST': os.getenv('GLOBAL_HOST', 's3-aapp.local'),
        'VITE_GLOBAL_HOST_CONSOLE': os.getenv('GLOBAL_HOST_CONSOLE', 'console-s3-aapp.local'),
    }


def s3_aapp_build(local_eth=False):
    docker_build(
        "comradecoop/s3-aapp",
        root_dir,
        dockerfile=root_dir + "/Dockerfile",
        target="run-all-singlenode",
        build_args=get_build_args(local_eth)
    )


def s3_aapp_deploy(cluster_names=["zero"], local_eth=False, deploy_proxy=True, deployer_key="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"):
    compose_config = {
        "services": {},
        "volumes": {},
    }
    if deploy_proxy:
        compose_config["services"]["proxy"] = {
            "image": "nmaguiar/socksd",
            "ports": ["127.0.0.1:1080:1080"],
        }
        
    if local_eth:
        compose_config["services"]["anvil"] = {
            "container_name": "eth-anvil",
            "image": "ghcr.io/foundry-rs/foundry:nightly-25f24e677a6a32a62512ad4f561995589ac2c7dc",
            "entrypoint": ["anvil", "--host", "0.0.0.0"],
            "networks": {
                "default": {
                    "aliases": ["anvil.local"]
                }
            }
        }
        local_resource(  # TODO: Move to container!
            "anvil-deploy-contracts",
            dir=root_dir + "/contracts",
            cmd="forge script script/Deploy.s.sol --rpc-url http://$(docker inspect eth-anvil -f '{{range $x := .NetworkSettings.Networks}}{{$x.IPAddress}}{{end}}'                                                   ):8545 --private-key %s --broadcast || true"
            % (deployer_key,),
            resource_deps=["anvil"],
            deps=[
                root_dir + "/contracts/src",
                root_dir + "/contracts/script",
                root_dir + "/contracts/lib",
            ],
            allow_parallel=True,
        )
    replicate_sites = []
    for name in cluster_names:
        is_last = name == cluster_names[-1]
        # TODO: Fix networking with serf
        # TODO: Redo dns and firefox proxy
        compose_config["services"][name] = {
            "container_name": "s3-%s" % name,
            "image": "comradecoop/s3-aapp",
            "volumes": ["s3-%s-data:/data" % name],
            "environment": {
                "BACKEND_ETH_PRIVATE_KEY": os.getenv("BACKEND_ETH_PRIVATE_KEY", "4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356"),
                "BACKEND_ETH_RPC": ("ws://anvil:8545" if local_eth else os.getenv("BACKEND_ETH_RPC", "https://sepolia.base.org/")),
                "BACKEND_ETH_TOKEN": '0x5FbDB2315678afecb367f032d93F642f64180aa3' if local_eth else os.getenv('BACKEND_ETH_TOKEN'),
                "BACKEND_ETH_WITHDRAW": os.getenv("BACKEND_ETH_WITHDRAW", "0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f"),
                "BACKEND_EXTERNAL_URL": "http://%s.local" % name,
                "BACKEND_REPLICATE_SITES": ",".join(replicate_sites),
            },
            "networks": {
                "default": {
                    "aliases": ["%s.local" % name] + (["s3-aapp.local", "console-s3-aapp.local"] if is_last else [])
                }
            }
        }
        replicate_sites += ["http://%s.local" % name]
        compose_config["volumes"]["s3-%s-data" % name] = {}
    docker_compose(encode_yaml(compose_config))


config.define_string("scenario", args=True, usage="One of single-cluster, multi-cluster")
config.define_bool("external-eth", usage="Don't launch ethereum in the local cluster")
cfg = config.parse()
local_eth = not cfg.get("external-eth", False)
scenario = cfg.get("scenario", "single-cluster")

s3_aapp_build(local_eth=local_eth)
# s3_aapp_build_with_builder()

if scenario == "single-cluster" or scenario == "sc":
    s3_aapp_deploy(["zero"], local_eth=local_eth)
elif scenario == "multi-cluster" or scenario == "mc":
    s3_aapp_deploy(["zero", "one"], local_eth=local_eth)
elif scenario == "single-cluster-2" or scenario == "sc2":
    s3_aapp_deploy(["one"], local_eth=local_eth)
else:
    fail("Unexpected scenario value", scenario)

local_resource(
    "launch_firefox",
    cmd=[],
    trigger_mode=TRIGGER_MODE_MANUAL,
    auto_init=False,
    serve_cmd=["./launch-proxy-firefox.sh"])
