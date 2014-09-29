task :default => [:help]

def encrypted_file
    return File.expand_path "~/Documents/t/t.zip"
end

def encrypted_file_parent
    return File.dirname encrypted_file
end

def decrypted_folder_parent
    return File.dirname decrypted_folder_full_path
end

def decrypted_folder_full_path
    return File.expand_path "~/decrypted_folder/#{decrypted_folder_name}"
end

def decrypted_folder_name
    return "t/"
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
    puts "encrypted_file               => #{encrypted_file}"
    puts "encrypted_file_parent        => #{encrypted_file_parent}"
    puts "decrypted_folder_parent      => #{decrypted_folder_parent}"
    puts "decrypted_folder_full_path   => #{decrypted_folder_full_path}"
end

task :encrypt => [:create_decrypted_folder_parent, :create_encrypted_file_folder] do
    Dir.chdir(decrypted_folder_parent){
        sh "zip -er #{encrypted_file} #{decrypted_folder_name}"
    }
end

task :create_decrypted_folder_parent do
    `mkdir #{decrypted_folder_parent}`
end

task :create_encrypted_file_folder do
    `mkdir #{encrypted_file_parent}`
end

task :decrypt => [:ensure_decrypted_folder_doesnt_exist, :create_decrypted_folder_parent] do
    sh "unzip -n #{encrypted_file} -d #{decrypted_folder_parent}"
end

task :ensure_decrypted_folder_doesnt_exist do
    if File.exists? decrypted_folder_full_path
        fail "Decrypted folder already exists at '#{decrypted_folder_full_path}' please manually delete it before continuing"
    end
end
