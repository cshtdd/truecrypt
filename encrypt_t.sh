#/!bin/sh

decrypted_folder_parent="$HOME/decrypted_folder/"
decrypted_folder="t/"
encrypted_file="$HOME/Documents/t/t.zip"

pushd $decrypted_folder_parent

echo compressing $decrypted_folder_parent$decrypted_folder into $encrypted_file 

zip -er $encrypted_file $decrypted_folder

popd
