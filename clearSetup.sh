#!/bin/bash
rm -rf /home/sergi/.dymint/
rm -rf /home/sergi/.dymint2/
./build/dymint init
./build/dymint init --home=/home/sergi/.dymint2
./build/dymint init --home=/home/sergi/.dymint3
rm /home/sergi/.dymint2/config/genesis.json
rm /home/sergi/.dymint2/config/priv_validator_key.json
rm /home/sergi/.dymint3/config/genesis.json
rm /home/sergi/.dymint3/config/priv_validator_key.json
cp /home/sergi/.dymint/config/genesis.json /home/sergi/.dymint2/config/
cp /home/sergi/.dymint/config/priv_validator_key.json /home/sergi/.dymint2/config/
cp /home/sergi/.dymint/config/genesis.json /home/sergi/.dymint3/config/
cp /home/sergi/.dymint/config/priv_validator_key.json /home/sergi/.dymint3/config/
