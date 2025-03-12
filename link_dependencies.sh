#!bin/bash Shell

# Delete the maze dependency path if it exists
rm -rf $PWD/vendor/github.com/dmigwi/tapoo

# Create the  maze dependency path
mkdir -p $PWD/vendor/github.com/dmigwi/tapoo

# Copy the contents of the maze package to the dependency path created
ln -sf $PWD/maze $PWD/vendor/github.com/dmigwi/tapoo