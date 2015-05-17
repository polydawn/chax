#!/bin/bash
set -euo pipefail

apt-get install prosody
update-rc.d prosody disable
service prosody stop
