# Olympus moon

Get all the `evmosvalcons` from testnet's active validator set.

## Validator set data

- We are going to use the result from calling `cosmos/staking/v1beta1/validators?pagination.limit=2000`
- This result is store in `validatorset.json`

## Result

The result of running the program is in the file `res.json`, generated running: `./olympus-moon-validatorconverter > res.json`
