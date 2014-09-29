task :default => [:help]

def encrypted_file
    return File.expand_path "~/Documents/t/t.zip"
end

def encrypted_file_parent
    return File.dirname encrypted_file
end

def decrypted_folder_parent
    return File.dirname decrypted_folder
end

def decrypted_folder
    return File.expand_path "~/decrypted_folder/t/"
end

task :help do
    puts "TrueCrypt replacement"

    puts ""
    puts "Tasks"
    puts "="*16
    puts "decrypt"
    puts "encrypt"

    puts ""
    puts "Settings"
    puts "="*16
    puts "encrypted_file          => #{encrypted_file}"
    puts "encrypted_file_parent   => #{encrypted_file_parent}"
    puts "decrypted_folder_parent => #{decrypted_folder_parent}"
    puts "decrypted_folder        => #{decrypted_folder}"
end

task :encrypt => [:create_decrypted_folder, :create_encrypted_file_folder] do
    sh "zip -er #{encrypted_file} #{decrypted_folder}"
end

task :create_decrypted_folder do
    `mkdir #{decrypted_folder_parent}`
end

task :create_encrypted_file_folder do
    `mkdir #{encrypted_file_parent}`
end
