#!/bin/bash
curl -X POST http://localhost:8080/clients/ -H "Content-Type: application/json" -d @etc/client.json      