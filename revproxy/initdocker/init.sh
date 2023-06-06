#!/bin/bash

iptables -t nat -A PREROUTING -p tcp --dport 8080 -j REDIRECT --to-port 11095
iptables -t nat -A OUTPUT -p tcp --dport 8080 -m owner ! --uid-owner 1005 -j DNAT --to 127.0.0.1:11095

# List all iptables rules.
iptables -t nat --list
