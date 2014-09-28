#/!bin/sh

decrypted_folder_parent="$HOME/decrypted_folder/"
encrypted_file="$HOME/Documents/t/t.zip"

unzip -n $encrypted_file -d $decrypted_folder_parent
