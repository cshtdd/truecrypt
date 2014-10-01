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
    puts "decrypt (d)"
    puts "encrypt (e)"
    puts "clean   (c)"
    puts "setup"

    puts ""
    puts "Settings"
    puts "="*16
    puts "encrypted_file               => #{encrypted_file}"
    puts "encrypted_file_parent        => #{encrypted_file_parent}"
    puts "decrypted_folder_parent      => #{decrypted_folder_parent}"
    puts "decrypted_folder_full_path   => #{decrypted_folder_full_path}"
end

task :e => :encrypt
task :encrypt => [:create_decrypted_folder_parent, :create_encrypted_file_folder] do
    Dir.chdir(decrypted_folder_parent){
        sh "zip -er #{encrypted_file} #{decrypted_folder_name}"
    }

    Rake::Task[:clean].invoke
end

task :create_decrypted_folder_parent do
    `mkdir #{decrypted_folder_parent}`
end

task :create_encrypted_file_folder do
    `mkdir #{encrypted_file_parent}`
end

task :d => :decrypt
task :decrypt => [:ensure_decrypted_folder_doesnt_exist, :create_decrypted_folder_parent] do
    sh "unzip -n #{encrypted_file} -d #{decrypted_folder_parent}"
end

task :ensure_decrypted_folder_doesnt_exist do
    if File.exists? decrypted_folder_full_path
        fail "Decrypted folder already exists at '#{decrypted_folder_full_path}' please manually delete it before continuing"
    end
end

task :c => :clean
task :clean do
    rm_rf decrypted_folder_full_path
end

task :setup do
    decrypt_contents = %{
#!/bin/sh
pushd #{basedir}
rake d
read -p "Press any key to continue..."
}
    generate_task_alias "decrypt.sh.command", decrypt_contents

    clean_contents = %{
#!/bin/sh
pushd #{basedir}
rake c
read -p "Press any key to continue..."
}
    generate_task_alias "clean.sh.command", clean_contents
end

def generate_task_alias(filename, file_contents)
    file_path = get_full_path filename
    write_file file_path, file_contents
    chmod "+x", file_path
end

def basedir
   return File.expand_path "."
end

def get_full_path(filename)
    return File.join(basedir, filename)
end

def write_file(file_path, file_content)
    File.open(file_path, "w") do |f|
        f.puts file_content
    end
end
