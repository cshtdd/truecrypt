#/!bin/sh

decrypted_folder="$HOME/decrypted_folder/t/"
encrypted_file="$HOME/Documents/t/t.zip"

echo compressing $decrypted_folder into $encrypted_file 

zip -er $encrypted_file $decrypted_folder
