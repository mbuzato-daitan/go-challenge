#!/usr/bin/env bash
vegeta attack -duration=5s -rate=5 --targets=./targets.txt | vegeta report --type text