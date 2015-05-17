#!/bin/bash
set -euo pipefail

prosodyctl --config ./prosody.cfg.lua register testpilot localhost asdf

prosody --config ./prosody.cfg.lua
