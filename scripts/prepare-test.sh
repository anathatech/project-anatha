#!/bin/bash
rm -r ~/.anathad
rm -r ~/.anathacli
anathad init vidor --chain-id test

# make folder structure for upgrades
mkdir -p ~/.anathad/upgrade_manager/genesis/bin
mkdir -p ~/.anathad/upgrade_manager/upgrades

# symlink genesis binary to upgrade
cp $(which anathad) ~/.anathad/upgrade_manager/genesis/bin

# MANAGER
# anatha1qaf2gssp652s6np00a5cxdwytdf3vutdumwc0q
echo "dial category ivory again laundry fan doctor walnut real away glory immense throw negative trumpet media bind rebuild embark wolf install diagram toss expand" | anathacli keys add manager --keyring-backend test --recover

# OPERATOR 1
# anatha18klyq6qxyemrgsuyapw80y7xz4gh0lhjk9tj7n
echo "plate jacket brush end vast friend jaguar achieve visa marine drastic strike wage surge deny flush rude supreme half panel common tide turtle ethics" | anathacli keys add operator1 --keyring-backend test --recover

# OPERATOR 2
# anatha1lafvzw2f6nj4sv2k0cpsu4zdxur9l0e98ngzxz
echo "proud electric faint armed market render cube audit firm cactus limb narrow play sweet deputy crisp now valve matter humor lake prize snake flag" | anathacli keys add operator2 --keyring-backend test --recover

# OPERATOR 3
# anatha13fsp5p3mjcxtr833vcln6kjwvyl8tz66tu2j8z
echo "shock very industry icon float fever input analyst candy cake chaos fly toilet rare fiber boy gorilla boring sting fit tennis broccoli spray conduct" | anathacli keys add operator3 --keyring-backend test --recover

# GOVERNOR 1
# anatha170mj6j6veall698u5xwkdhpf2nlza9j6n366v2
echo "elite snow another ocean sad hamster neck alien brisk genre goat creek hire foam artist web curve trust review crisp know glove upgrade base" | anathacli keys add governor1 --keyring-backend test --recover

# GOVERNOR 2
# anatha1e4ay99w2v3qqnl8wrmv3s85mexm6eum2vhucmf
echo "debate candy marine fish chief praise theory cactus hard easily knock before lumber sweet stumble tonight merit echo space special cabin dawn step human" | anathacli keys add governor2 --keyring-backend test --recover

# GOVERNOR 3
# anatha1p470fhtytsym8fpp5j6rak2ec2qrvmwdvpg3ve
echo "flee flip much leader mandate mosquito script solve fox soap december high clap chapter stumble siege embrace first atom calm nothing noodle permit win" | anathacli keys add governor3 --keyring-backend test --recover

# GOVERNOR 4
# anatha1rulgmktamspkx7ecyk6h8wfuwtyylxwldgy4fn
echo "inmate merit banner era furnace save glow misery hollow check measure disease void genius equip joke consider area exhibit aisle subway metal spawn fade" | anathacli keys add governor4 --keyring-backend test --recover

# Validator
# anatha1yp9m0h337z9t0gn7rtxjpe05ttjsyvh6w380yr
# anathavaloper1yp9m0h337z9t0gn7rtxjpe05ttjsyvh6zrk4ld
echo "blind salt icon document best twice wage stadium horn toddler infant flush one slow buddy laptop artist false celery family ahead demand bulk bus" | anathacli keys add validator --keyring-backend test --recover

# Test Owner
# anatha19emxp37355qkw80adr83lq67h8l9jhe2kq4crv
echo "depth vivid text labor kangaroo lend roof sudden achieve basket mechanic brisk rich black deputy uncover birth onion fantasy rigid reunion slow bulk ostrich" | anathacli keys add testOwner --keyring-backend test --recover

# Test Buyer
# anatha1pv9x8vkxypgumth9uzksfs5fgc9nx7c3zcwa3u
echo "clock virtual ice wonder label vacuum visa return local nerve wrap joy wing acoustic plug travel joy slot famous diagram census flag fly foster" | anathacli keys add testBuyer --keyring-backend test --recover

# Test
# anatha18cte6l48hzs4hu50evrhmcvcgtgc85c5cyt0ud
echo "mechanic antenna trigger sugar kidney umbrella result month seek defy spell claw among rough only clarify worry silly ignore fit essence harbor garden useful" | anathacli keys add test --keyring-backend test --recover

#anathad add-genesis-account $(anathacli keys show treasury -a --keyring-backend test) 769700000000000000pin
anathad add-genesis-account $(anathacli keys show manager -a --keyring-backend test) 100000000000000pin
anathad add-genesis-account $(anathacli keys show validator -a --keyring-backend test) 100000000000000pin
anathad add-genesis-account $(anathacli keys show testOwner -a --keyring-backend test) 100000000000000pin
anathad add-genesis-account $(anathacli keys show testBuyer -a --keyring-backend test) 100000000000000pin
anathad gentx --name validator --amount 10000000000pin --keyring-backend test
anathad collect-gentxs
