# -*- mode: Python -*-
# SPDX-License-Identifier: GPL-3.0

load(
    "./deploy/Tiltfile", "apocryph_resource", "apocryph_build_with_builder", "deploy_apocryph_stack"
)

config.define_string_list("include")
cfg = config.parse()

apocryph_build_with_builder()
deploy_apocryph_stack()

if "include" in cfg:
    for f in cfg["include"]:
        load_dynamic(f)
