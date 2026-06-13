#!/bin/sh
set -e

./trama-cli migrate

exec ./trama
