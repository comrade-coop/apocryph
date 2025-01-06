# -*- mode: Python -*-
# SPDX-License-Identifier: GPL-3.0

# For more on Extensions, see: https://docs.tilt.dev/extensions.html
load("ext://restart_process", "docker_build_with_restart")
load("ext://namespace", "namespace_create")
load("ext://helm_resource", "helm_resource", "helm_repo")

local(
    "which jq forge cast helm kubectl docker >/dev/null", echo_off=True
)  # Check dependencies
cluster_ip = local(
    "kubectl get no -o jsonpath --template '{$.items[0].status.addresses[?(.type==\"InternalIP\")].address}'"
)
root_dir = os.getcwd()
apocryph_dir = root_dir + "/.." + "/apocryph"


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
        },
    }
    docker_compose(encode_yaml({"services": services, "volumes": volumes_conf}))
    dc_resource(name, *args, **kwargs)


def local_resource_in_builder(name, cmd, builder, resource_deps=[], serve_cmd='', *args, **kwargs):
    local_resource(
        name,
        cmdline_in_builder(cmd, builder),
        resource_deps=resource_deps + [builder],
        serve_cmd=cmdline_in_builder(serve_cmd, builder),
        *args,
        **kwargs,
    )


def cmdline_in_builder(cmd, builder, *, interactive=False):
    if cmd == '':
        return ''
    flags = []
    if interactive:
        flags += ["-i"]
    if type(cmd) == "list":
        return ["docker", "exec"] + flags + [builder] + cmd
    else:
        return "docker exec %s %s %s" % (" ".join(flags), builder, cmd)


def s3_aapp_build_with_builder():
    docker_build(
        "comradecoop/s3-aapp/go-builder",
        root_dir,
        dockerfile=root_dir + "/Dockerfile",
        target="go-build-base",
        only=[root_dir + "/Dockerfile"],
    )

    builder_resource(
        "apocryph-s3-go-builder",
        "comradecoop/s3-aapp/go-builder",
        dir=root_dir,
        write_dir="bin",
        volumes=["go-cache:/root/.cache/go-build", "go-mod-cache:/go/pkg/mod"],
        volumes_conf={"go-cache": {}, "go-mod-cache": {}},
        labels=["build"],
    )

    local_resource_in_builder(
        "apocryph-s3-backend-go-compile",
        'go build -v -buildvcs=false -ldflags="-s -w" -o bin/apocryph-s3-backend ./backend',
        "apocryph-s3-go-builder",
        deps=[root_dir + "/backend"],
        allow_parallel=True,
        labels=["build"],
    )

    docker_build_with_restart(
        "comradecoop/s3-aapp/backend",
        root_dir,
        dockerfile="./Dockerfile",
        target="run-backend-copy-local",
        entrypoint=["/usr/local/bin/apocryph-s3-backend"],
        only=[root_dir + "/bin"],
        live_update=[
            sync(root_dir + "/bin", "/usr/local/bin/"),
        ],
    )

    docker_build(
        "comradecoop/s3-aapp/serf",
        root_dir,
        dockerfile=root_dir + "/Dockerfile.serf",
        only=[root_dir + "/Dockerfile.serf"],
    )

def s3_aapp_serve_with_builder():
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
        #TODO: ports=[5173],
        volumes_conf={"pnpm-cache": {}},
        labels=["build"],
    )

    local_resource_in_builder(
        "apocryph-s3-js-serve",
        '',
        "apocryph-s3-js-builder",
        serve_cmd=['sh', '-c', 'cd frontend/ && pnpm install --frozen-lockfile && pnpm run dev'],
        deps=[],
        allow_parallel=True,
        labels=["system"],
    )

def s3_aapp_deploy(
    cluster_names=["one", "two"]
):
    update_settings(k8s_upsert_timeout_secs=160)

    helm_repo("minio-operator-chart", "https://operator.min.io")
    helm_repo("ingress-nginx-chart", "https://kubernetes.github.io/ingress-nginx")

    helm_resource(
        "ingress-nginx",
        "ingress-nginx-chart/ingress-nginx",
        namespace="ingress",
        resource_deps=["ingress-nginx-chart"],
        labels=["system"],
        flags=["--create-namespace"],
    )

    helm_resource(
        "minio-operator",
        "minio-operator-chart/operator",
        namespace="minio-operator",
        labels=["system"],
        resource_deps=["minio-operator-chart"],
        flags=[
            "--create-namespace",
            "--set=operator.replicaCount=1",
        ]
    )

    for name in cluster_names:
        helm_resource(
            "minio-%s-config" % name,
            root_dir + "/charts/config",
            namespace="s3-%s" % name,
            deps=[root_dir + "/charts/config"],
            resource_deps=[],
            labels=["s3-%s" % name],
            flags=[
                "--create-namespace",
            ],
        )
        helm_resource(
            "minio-%s-tenant" % name,
            "minio-operator-chart/tenant",
            namespace="s3-%s" % name,
            resource_deps=["minio-operator", "minio-operator-chart", "minio-%s-config" % name],
            labels=["s3-%s" % name],
            flags=[
                "--create-namespace",
                "--set=tenant.pools[0].name=minio-%s" % name,
                "--set-json=tenant.certificate.requestAutoCert=false",
                "--set=tenant.pools[0].servers=1",
                "--set=tenant.pools[0].volumesPerServer=1",
                "--set=tenant.pools[0].size=5Gi",
                "--set=tenant.configuration.name=minio-config",
            ],
        )
        helm_resource(
            "minio-%s-backend" % name,
            root_dir + "/charts/backend",
            namespace="s3-%s" % name,
            deps=[root_dir + "/charts/backend"],
            resource_deps=["minio-%s-tenant" % name],
            labels=["s3-%s" % name],
            flags=[
                "--create-namespace",
                "--set=ingress.host=%s.s3.localhost" % name,
                "--set=serf.bootstrap=serf-bind.s3-%s.svc.cluster.local" % cluster_names[0],
            ],
            image_deps=["comradecoop/s3-aapp/backend", "comradecoop/s3-aapp/serf"],
            image_keys=["backend.image", "serf.image"],
        )

    namespace_create("eth")
    # TODO: Recreate anvil when we have new contracts code
    k8s_yaml(listdir(apocryph_dir + "/deploy/charts/eth"))
    k8s_resource("anvil", labels=["z_contracts"])


def s3_aapp_deploy_local(
    deployer_key="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
    resource_deps=["anvil", "ingress-nginx"],
):
    if len(resource_deps) == 0:
        local_resource(
            "ingress-nginx-portforward",
            serve_cmd="kubectl port-forward -n keda svc/ingress-nginx-controller 8004:80",
        )
        local_resource(
            "anvil-portforward",
            serve_cmd="kubectl port-forward -n eth svc/eth-rpc 8549:8545",
        )
    else:
        k8s_resource(workload="ingress-nginx", port_forwards=["8004:80"])
        k8s_resource(workload="anvil", port_forwards=["8549:8545"])

        local_resource(  # TODO: Move to container!
            "anvil-deploy-contracts",
            labels=["z_contracts"],
            dir=apocryph_dir + "/contracts",
            cmd="forge script script/Deploy.s.sol --rpc-url http://127.0.0.1:8549 --private-key %s --broadcast || true"
            % (deployer_key,),
            resource_deps=["anvil"],
            deps=[apocryph_dir + "/contracts/src", apocryph_dir + "/contracts/script", apocryph_dir + "/contracts/lib"],
            allow_parallel=True,
        )

s3_aapp_build_with_builder()
s3_aapp_deploy()
s3_aapp_deploy_local()
s3_aapp_serve_with_builder()
