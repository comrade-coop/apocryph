# -*- mode: Python -*-
# SPDX-License-Identifier: GPL-3.0


config.define_string_list("include")
config.define_string("allow-context")
config.define_bool("deploy-stack")
cfg = config.parse()

if "allow-context" in cfg:
    allow_k8s_contexts(cfg["allow-context"])

load(
    "./deploy/Tiltfile",
    "apocryph_resource",
    "apocryph_build_with_builder",
    "deploy_apocryph_stack",
    "deploy_apocryph_local",
)

if cfg.get("deploy-stack", True):
    apocryph_build_with_builder()
    deploy_apocryph_stack()
    deploy_apocryph_local()
else:
    apocryph_build_with_builder(skip_images=True)
    deploy_apocryph_local(resource_deps=[])

for f in cfg.get("include", []):
    load_dynamic(f)
