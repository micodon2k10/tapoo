#!bin/bash Shell
# This file is only needed while doing code development locally.

# Export testing environment variables
export TAPOO_DB_NAME=test_db 
export TAPOO_DB_USER_NAME=test_tapoo
export TAPOO_DB_USER_PASSWORD=test1234 
export TAPOO_DB_HOST=localhost:3306

# Delete the maze dependency path if it exists
rm -rf $PWD/vendor/github.com/dmigwi/tapoo

# Create the  maze dependency path
mkdir -p $PWD/vendor/github.com/dmigwi/tapoo

# Copy the contents of the maze package to the dependency path created
ln -sf $PWD/maze $PWD/vendor/github.com/dmigwi/tapoo
