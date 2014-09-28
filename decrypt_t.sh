#/!bin/sh

decrypted_folder="$HOME/decrypted_folder/t/"
encrypted_file="$HOME/Documents/t/t.zip"

unzip -n $encrypted_file -d $decrypted_folder
