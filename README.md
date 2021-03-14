# Infrastructure for trading platform

## Instruction for building
1. go to `tradeserver` folder, run `source platform/docker/.env`
2. run `make`
3. put acc_info.json under `tradeserver` folder
3. run `docker-compose -f platform/docker/compose.yml up -d trade-server`