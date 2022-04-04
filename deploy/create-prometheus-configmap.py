#!/usr/bin/env python3

import sys
import os
import subprocess

yaml = os.path.join(
    os.path.abspath(os.path.dirname(sys.argv[0])), 
    "scrapeconfig.yml"
)
command = f"kubectl create configmap prometheus-config --from-file {yaml}"
replaceCommand = f"kubectl create configmap --dry-run=client --from-file {yaml} prometheus-config -o yaml | kubectl replace -f -"
result = subprocess.call(command, shell=True)
result = subprocess.call(replaceCommand, shell=True)
